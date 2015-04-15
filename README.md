# Cloud Foundry Copy Env CLI Plugin

Cloud Foundry plugin extension to export live application VCAP_SERVICES environment variable into the local developer machine.

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
