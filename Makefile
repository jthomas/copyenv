# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOENV=$(GOCMD) env
SED=sed
BUILD_FOLDER=bin
BINARY_NAME=copyenv
SHASUM=shasum
CF=cf

CURARCH=$(shell $(GOENV) | $(SED) -n 's@GOHOSTARCH="\([^"]*\).*@\1@p')
CUROS=$(shell $(GOENV) |  $(SED) -n 's@GOHOSTOS="\([^"]*\).*@\1@p')
ARCHWIDTH=$(shell echo $(CURARCH) | $(SED) s'@.*\(..\)@\1@')

# Parameters for the cross platform release binaries
PLATFORMS := linux/amd64/linux64 linux/386/linux32 windows/amd64/win64 windows/386/win32 darwin/amd64/osx

BINARIES = $(wildcard bin/copyenv.*)

TEMP = $(subst /, ,$@)
OS = $(word 1, $(TEMP))
ARCH = $(word 2, $(TEMP))
EXT = $(word 3, $(TEMP))

# DEFAULT
all: deps test release

release: $(PLATFORMS)

$(PLATFORMS):
	GOOS=$(OS) GOARCH=$(ARCH) go build -o bin/$(BINARY_NAME).'$(EXT)'
	$(SHASUM) -a 1 bin/$(BINARY_NAME).'$(EXT)'

test: deps
	$(GOTEST) -v ./...
clean:
	$(GOCLEAN)

install-plugin: release
ifeq ($(CUROS),darwin)
	$(CF) install-plugin -f bin/$(BINARY_NAME).osx
endif
ifeq ($(CUROS),linux)
	$(CF) install-plugin -f bin/$(BINARY_NAME).linux$(ARCHWIDTH)
endif
ifeq ($(CUROS),windows)
	$(CF) install-plugin -f bin/$(BINARY_NAME).win$(ARCHWIDTH)
endif

deps:
	$(GOGET) -t ./...


.PHONY: deps clean test install-plugin $(PLATFORMS)