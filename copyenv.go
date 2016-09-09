package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/cloudfoundry/cli/plugin"
)

type CopyEnv struct{}

func fatalIf(err error) {
	if err != nil {
		fmt.Fprintln(os.Stdout, "error: ", err)
		os.Exit(1)
	}
}

func (c *CopyEnv) ExtractAppName(args []string) (string, error) {
	if len(args) < 2 {
		return "", errors.New("missing application name")
	}

	return args[1], nil
}

func (c *CopyEnv) RetrieveAppNameEnv(cliConnection plugin.CliConnection, app_name string) ([]string, error) {

	app, err := cliConnection.GetApp(app_name)

	if err != nil {
		msg := fmt.Sprintf("Failed to retrieve enviroment for \"%s\", is the app name correct?", app_name)
		err = errors.New(msg)
	} else {
		url := fmt.Sprintf("/v2/apps/%s/env", app.Guid)
		output, err := cliConnection.CliCommandWithoutTerminalOutput("curl", url)

		if err != nil {
			msg := fmt.Sprintf("Failed to retrieve enviroment for \"%s\", is the app name correct?", app_name)
			err = errors.New(msg)
		}

		return output, err
	}
	return nil, err
}

func (c *CopyEnv) ExtractCredentialsJSON(envParent string, credKey string, output []string) ([]byte, error) {
	err := errors.New("missing service credentials for application")
	var b []byte

	val := strings.Join(output, "")
	if strings.Contains(val, credKey) {
		var f interface{}
		err = json.Unmarshal([]byte(val), &f)
		if err != nil {
			return nil, err
		}

		envJson := f.(map[string]interface{})
		envParentJson := envJson[envParent].(map[string]interface{})
		b, err = json.Marshal(envParentJson[credKey])
		if err != nil {
			return nil, err
		}
	}

	return b, err
}

func (c *CopyEnv) ExportCredsAsShellVar(cred_key string, creds string) {
	vcap_services := fmt.Sprintf("export %s='%s';", cred_key, creds)
	fmt.Println(vcap_services)
}

func (c *CopyEnv) ExtractAndExportCredentials(env_parent string, cred_key string, app_env []string) {
	creds, err := c.ExtractCredentialsJSON(env_parent, cred_key, app_env)
	fatalIf(err)
	c.ExportCredsAsShellVar(cred_key, string(creds[:]))
}

func (c *CopyEnv) Run(cliConnection plugin.CliConnection, args []string) {
	if len(args) > 0 && args[0] == "CLI-MESSAGE-UNINSTALL" {
		return
	}
	app_name, err := c.ExtractAppName(args)
	fatalIf(err)

	app_env, err := c.RetrieveAppNameEnv(cliConnection, app_name)
	fatalIf(err)

	if len(args) > 2 && args[2] == "--all" {
		c.ExtractAndExportCredentials("application_env_json", "VCAP_APPLICATION", app_env)
		fmt.Println("")
	}

	c.ExtractAndExportCredentials("system_env_json", "VCAP_SERVICES", app_env)
}

func (c *CopyEnv) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "copyenv",
		Version: plugin.VersionType{
			Major: 1,
			Minor: 1,
			Build: 0,
		},
		Commands: []plugin.Command{
			plugin.Command{
				Name:     "copyenv",
				HelpText: "Export application VCAP_SERVICES to local environment variable.",
				UsageDetails: plugin.Usage{
					Usage: "copyenv APP_NAME [--all] - Retrieve and export remote application VCAP_SERVICES to local developer environment.",
					Options: map[string]string{
						"all": "Retrieve both VCAP_SERVICES and VCAP_APPLICATION from remote application",
					},
				},
			},
		},
	}
}

func main() {
	plugin.Start(new(CopyEnv))
}
