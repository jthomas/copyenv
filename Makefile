# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOENV=$(GOCMD) env
BUILD_FOLDER=bin
BINARY_NAME=copyenv

CURARCH=$(shell $(GOENV) | grep GOHOSTARCH | cut -d "=" -f2 | tr -d '"')
CUROS=$(shell $(GOENV) | grep GOHOSTOS | cut -d "=" -f2 | tr -d '"')
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

# install: $(PLATFORMS)
# 	ifeq ($(os),$(CUROS))
# 		ifeq ($(arch),$(CURARCH))
# 		 cf install-plugin -f bin/$(BINARY_NAME).'$(ext)'
# 		endif
# 	endif
deps:
	$(GOGET) ./...

.PHONY: deps install clean test $(PLATFORMS)