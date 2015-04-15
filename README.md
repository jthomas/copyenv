# Cloud Foundry CLI Copy Env Plugin

Cloud Foundry CLI plugin to export application VCAP_SERVICES onto the local machine.

Applications running on Cloud Foundry rely on the VCAP_SERVICES environment variable to provide service credentials. 

When running applications locally for development and testing, it's useful to have the same VCAP_SERVICES values available in the local environment to simulate running on the host platform.


This plugin will export the remote application environment variables, available using cf env, into a format that makes it simple to expose those same values locally. 

## Install

```
$ go get github.com/jthomas/copyenv
$ cf install-plugin $GOPATH/bin/copyenv
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

## Uninstall

```
$ cf uninstall-plugin copyenv
```
