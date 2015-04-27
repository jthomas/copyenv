package main_test

import (
	. "."
	"github.com/cloudfoundry/cli/plugin/fakes"
	io_helpers "github.com/cloudfoundry/cli/testhelpers/io"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cloud Foundry Copyenv Command", func() {
	Describe(".Run", func() {
		var fakeCliConnection *fakes.FakeCliConnection
		var callCopyEnvCommandPlugin *CopyEnv

		BeforeEach(func() {
			fakeCliConnection = &fakes.FakeCliConnection{}
			callCopyEnvCommandPlugin = &CopyEnv{}
		})

		It("Extract Application Name From Command Line Args", func() {
			name, err := callCopyEnvCommandPlugin.ExtractAppName([]string{"copyenv"})
			Ω(err).Should(MatchError("missing application name"))

			name, err = callCopyEnvCommandPlugin.ExtractAppName([]string{"copyenv", "APP_NAME"})
			Ω(err).ShouldNot(HaveOccurred())
			Ω(name).Should(Equal("APP_NAME"))
		})

		It("Retrieve Application Environment Variables From Name", func() {
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns([]string{"SOME", "OUTPUT", "COMMAND"}, nil)
			output, err := callCopyEnvCommandPlugin.RetrieveAppNameEnv(fakeCliConnection, "APP_NAME")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(fakeCliConnection.CliCommandWithoutTerminalOutputCallCount()).Should(Equal(1))
			Ω(fakeCliConnection.CliCommandWithoutTerminalOutputArgsForCall(0)).Should(Equal([]string{"env", "APP_NAME"}))
			Ω(output).Should(Equal([]string{"SOME", "OUTPUT", "COMMAND"}))
		})

		It("Return Service Credentials From Appplication Environment", func() {
			_, err := callCopyEnvCommandPlugin.ExtractServiceCredentialsJSON([]string{""})
			Ω(err).Should(MatchError("missing service credentials for application"))

			service_creds := []string{"{\"VCAP_SERVICES\":{\"service\": [ { \"credentials\": {} } ]}}"}
			b, err := callCopyEnvCommandPlugin.ExtractServiceCredentialsJSON(service_creds)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(string(b[:])).Should(Equal("{\"service\":[{\"credentials\":{}}]}"))
		})

		It("Print Service Credentials As Shell Variable", func() {
			output := io_helpers.CaptureOutput(func() {
				callCopyEnvCommandPlugin.ExportCredsAsShellVar("testing")
			})
			Ω(output[0]).Should(Equal("export VCAP_SERVICES='testing';"))
		})
	})
})
