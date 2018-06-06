# Cloud Foundry CLI Copy Env Plugin

[![Build Status](https://travis-ci.org/vaughnh/copyenv.svg?branch=master)](https://travis-ci.org/vaughnh/copyenv)

Cloud Foundry CLI plugin to export application VCAP_SERVICES and VCAP_APPLICATION onto the local machine.

Applications running on Cloud Foundry rely on the VCAP_SERVICES environment variable to provide service credentials. Application specific configuration environment is provided in VCAP_APPLICATION.

When running applications locally for development and testing, it's useful to have the same VCAP_SERVICES values available in the local environment to simulate running on the host platform.

This plugin will export the remote application environment variables, available using cf env, into a format that makes it simple to expose those same values locally. This plugin will only export VCAP_SERVICES by default. To include VCAP_APPLICATION as well use the --all option.

## Install

```
$ go get github.com/jthomas/copyenv
$ cf install-plugin $GOPATH/bin/copyenv
```
or
```
$ make install-plugin
```

## Usage

```
$ cf copyenv APP_NAME

export VCAP_SERVICES='...'
```

The plugin output needs to be evaluated in your shell to set up the
local environment variables.

Using eval: 
```
$ eval `cf copyenv APP_NAME` 
```

Using a temporary file:
```
$ cf copyenv APP_NAME > temp.json
$ source temp.json
```

You can obtain additional information with flags.
```
$ cf copyenv APP_NAME -all

export VCAP_SERVICES='...'
export VCAP_APPLICATION='...'
```


## Uninstall

```
$ cf uninstall-plugin copyenv
```
