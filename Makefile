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

all: deps test build release

release: $(PLATFORMS)

$(PLATFORMS):
	GOOS=$(os) GOARCH=$(arch) go build -o bin/$(BINARY_NAME).'$(ext)'

build:

test: 
	$(GOTEST) -v ./...
clean:
	$(GOCLEAN)

install-plugin: release
ifeq ($(CUROS),darwin)
	cf install-plugin -f bin/$(BINARY_NAME).osx
endif
ifeq ($(CUROS),linux)
	cf install-plugin -f bin/$(BINARY_NAME).linux$(ARCHWIDTH)
endif
ifeq ($(CUROS),windows)
	cf install-plugin -f bin/$(BINARY_NAME).win$(ARCHWIDTH)
endif

deps:
	$(GOGET) ./...



.PHONY: deps clean test install-plugin $(PLATFORMS)