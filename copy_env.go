package main

import (
        "os"
	"fmt"
	"strings"
        "encoding/json"
	"github.com/cloudfoundry/cli/plugin"
)

type CopyEnv struct{}

func (c *CopyEnv) Run(cliConnection plugin.CliConnection, args []string) {
	if len(args) < 2 {
              fmt.Println("ERROR: Missing application name")
	      os.Exit(1)
	}

        args[0] = "env"
        output, err := cliConnection.CliCommandWithoutTerminalOutput(args...)

	if err != nil {
                for _, val := range output {
                        fmt.Println(val)
                }
	        os.Exit(1)
	}

	for _, val := range output {
                if (strings.Contains(val, "VCAP_SERVICES")) {
                  var f interface{}
                  err := json.Unmarshal([]byte(val), &f)
                  if err != nil {
                           fmt.Println(err)
                           os.Exit(1)
                  }

                  m := f.(map[string]interface{})
                  b, err := json.Marshal(m["VCAP_SERVICES"])
                  if err != nil {
                            fmt.Println(err)
                            os.Exit(1)
                  }

                  vcap_services := "export VCAP_SERVICES='" + string(b[:]) + "';"
                  fmt.Println(vcap_services)
                }
	}
}

func (c *CopyEnv) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "copyenv",
		Version: plugin.VersionType{
			Major: 1,
			Minor: 0,
			Build: 0,
		},
		Commands: []plugin.Command{
			plugin.Command{
				Name:     "copyenv",
				HelpText: "Export application VCAP_SERVICES to local environment variable.",
				UsageDetails: plugin.Usage{
					Usage: "copyenv APP_NAME - Retrieve and export remote application VCAP_SERVICES to local developer environment.",
				},
			},
		},
	}
}

func main() {
	plugin.Start(new(CopyEnv))
}
