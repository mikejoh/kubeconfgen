#VERSION := $(shell git describe --tags)
GITCOMMIT := $(shell git rev-parse HEAD)
PROJECTNAME := $(shell basename "$(PWD)")
DATE := $(shell date "+%Y-%m-%d@%H:%M:%S")

# Go related variables.
GOVER=1.13
GOCMD=go
GOBASE := $(shell pwd)
GOBIN := $(GOBASE)/bin

# Environment info
ARCH=amd64
OS=linux

# Use linker flags to provide version/build settings
LDFLAGS=-ldflags "-s -w -X=main.GitCommit=$(GITCOMMIT) -X=main.BuildDate=$(DATE) -X=main.GoVersion=$(GOVER) -X=main.OperatingSystem=$(OS) -X=main.Architecture=$(ARCH)"

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

## exec: Run given command, wrapped with custom GOPATH. e.g; make exec run="go test ./..."
exec:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) $(run)

## clean: Clean build files. Runs `go clean` internally.
clean:
	@-rm $(GOBIN)/$(PROJECTNAME) 2> /dev/null

## build-linux: Build Linux amd64 binary locally.
build-linux:
	@echo " $(LDFLAGS)"
	@echo "  >  Building Linux binary..."
	@GOOS=$(OS) GOARCH=$(ARCH) go build $(LDFLAGS) -o $(GOBIN)/$(PROJECTNAME) $(GOBASE)/cmd/$(PROJECTNAME)/$(wildcard *.go)

## build-linux-docker: Build Linux amd64 binary locally but through a Docker container.
build-linux-docker:
	@echo " > Build Linux binary in a Docker container..."
	@docker run --rm -it -v $(GOBASE):/$(PROJECTNAME) -w="/$(PROJECTNAME)" golang:$(GOVER)-alpine sh -c "GOOS=$(OS) GOARCH=$(ARCH) go build $(LDFLAGS) -o $(GOBIN)/$(PROJECTNAME) ./cmd/$(PROJECTNAME)/$(wildcard *.go)"

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo