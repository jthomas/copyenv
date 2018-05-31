package main

import (
	"errors"

	io_helpers "code.cloudfoundry.org/cli/cf/util/testhelpers/io"
	"code.cloudfoundry.org/cli/plugin/models"
	"code.cloudfoundry.org/cli/plugin/pluginfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cloud Foundry Copyenv Plugin", func() {
	Describe(".Run", func() {
		var fakeCliConnection *pluginfakes.FakeCliConnection
		var callCopyEnvCommandPlugin *CopyEnv

		BeforeEach(func() {
			fakeCliConnection = &pluginfakes.FakeCliConnection{}
			callCopyEnvCommandPlugin = &CopyEnv{}
		})

		It("Extract Application Name From Command Line Args", func() {
			name, err := callCopyEnvCommandPlugin.extractAppName([]string{"copyenv"})
			Ω(err).Should(MatchError("missing application name"))

			name, err = callCopyEnvCommandPlugin.extractAppName([]string{"copyenv", "APP_NAME"})
			Ω(err).ShouldNot(HaveOccurred())
			Ω(name).Should(Equal("APP_NAME"))
		})

		It("Retrieve Application Environment Variables From Name", func() {
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns([]string{"SOME", "OUTPUT", "COMMAND"}, nil)
			fakeCliConnection.GetAppReturns(plugin_models.GetAppModel{Name: "APP_NAME", Guid: "testing"}, nil)
			output, err := callCopyEnvCommandPlugin.retrieveAppNameEnv(fakeCliConnection, "APP_NAME")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(fakeCliConnection.CliCommandWithoutTerminalOutputCallCount()).Should(Equal(1))
			Ω(fakeCliConnection.CliCommandWithoutTerminalOutputArgsForCall(0)).Should(Equal([]string{"curl", "/v2/apps/testing/env"}))
			Ω(output).Should(Equal([]string{"SOME", "OUTPUT", "COMMAND"}))
		})

		It("Return Service Credentials From Application Environment", func() {
			_, err := callCopyEnvCommandPlugin.extractCredentialsJSON("system_env_json", "VCAP_SERVICES", []string{""})
			Ω(err).Should(MatchError("missing service credentials for application"))

			serviceCreds := []string{"{\"system_env_json\": {\"VCAP_SERVICES\":{\"service\": [ { \"credentials\": {} } ]}}}"}
			b, err := callCopyEnvCommandPlugin.extractCredentialsJSON("system_env_json", "VCAP_SERVICES", serviceCreds)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(string(b[:])).Should(Equal("{\"service\":[{\"credentials\":{}}]}"))
		})

		It("Print Service Credentials As Shell Variable", func() {
			output := io_helpers.CaptureOutput(func() {
				callCopyEnvCommandPlugin.exportCredsAsShellVar("VCAP_SERVICES", "testing")
			})
			Ω(output[0]).Should(Equal("export VCAP_SERVICES='testing';"))
		})

		It("Return Error When App Name Is Missing", func() {
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns([]string{}, errors.New(""))
			output, err := callCopyEnvCommandPlugin.retrieveAppNameEnv(fakeCliConnection, "missing_app")
			Ω(output).Should(Equal([]string{}))
			Ω(err).ShouldNot(Equal(nil))
		})

		It("Silently uninstalls", func() {
			callCopyEnvCommandPlugin.Run(fakeCliConnection, []string{"CLI-MESSAGE-UNINSTALL"})
			Ω(fakeCliConnection.CliCommandWithoutTerminalOutputCallCount()).Should(Equal(0))
		})

		Context("when called with --all", func() {
			It("Extracts VCAP_APPLICATION and VCAP_SERVICE", func() {
				services := `{
					"system_env_json": {
						"VCAP_SERVICES": {
							"services": ["services"]
						}
					},
					"application_env_json": {
						"VCAP_APPLICATION": ["application"]
					}
				}`
				fakeCliConnection.CliCommandWithoutTerminalOutputReturns([]string{
					services}, nil)

				output := io_helpers.CaptureOutput(func() {
					callCopyEnvCommandPlugin.Run(fakeCliConnection, []string{"copyenv", "APP_NAME", "--all"})
				})

				Ω(output).Should(ContainElement(
					"export VCAP_APPLICATION='[\"application\"]';",
				))

				Ω(output).Should(ContainElement(
					`export VCAP_SERVICES='{"services":["services"]}';`,
				))
			})
		})
	})
})
