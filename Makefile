# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOENV=$(GOCMD) env
sed=sed
BUILD_FOLDER=bin
BINARY_NAME=copyenv

CURARCH=$(shell $(GOENV) | $(sed) -n 's@GOHOSTARCH="\([^"]*\).*@\1@p')
CUROS=$(shell $(GOENV) |  $(sed) -n 's@GOHOSTOS="\([^"]*\).*@\1@p')
ARCHWIDTH=$(shell echo $(CURARCH) | $(sed) s'@.*\(..\)@\1@')

PLATFORMS := linux/amd64/linux64 linux/386/linux32 windows/amd64/win64 windows/386/win32 darwin/amd64/osx

temp = $(subst /, ,$@)
os = $(word 1, $(temp))
arch = $(word 2, $(temp))
ext = $(word 3, $(temp))

release: $(PLATFORMS)

$(PLATFORMS):
	GOOS=$(os) GOARCH=$(arch) go build -o bin/$(BINARY_NAME).'$(ext)'

all: test build
build:

test: 
	$(GOTEST) -v ./...
clean:
	$(GOCLEAN)

# install: release
# ifeq ($(CUROS),darwin)
# 	INSTALL_TARGET=bin/$(BINARY_NAME).osx
# endif
# ifeq ($(CUROS),linux)
# 	INSTALL_TARGET=bin/$(BINARY_NAME).linux$(ARCHWIDTH)
# endif
# ifeq ($(CUROS),windows)
# 	INSTALL_TARGET=bin/$(BINARY_NAME).win$(ARCHWIDTH)
# endif
# 	cf install-plugin -f '$(INSTALL_TARGET)'

deps:
	$(GOGET) ./...

.PHONY: deps install clean test $(PLATFORMS)