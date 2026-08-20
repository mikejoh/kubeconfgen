CMDPATH := ./cmd/kubeconfgen
BINARY  := kubeconfgen
BUILDDIR := build

GIT_COMMIT := $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")
BUILD_DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
GO_VERSION := $(shell go env GOVERSION)
OS         := $(shell go env GOOS)
ARCH       := $(shell go env GOARCH)

LDFLAGS := -s -w \
	-X main.GitCommit=$(GIT_COMMIT) \
	-X main.BuildDate=$(BUILD_DATE) \
	-X main.GoVersion=$(GO_VERSION) \
	-X main.OperatingSystem=$(OS) \
	-X main.Architecture=$(ARCH)

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

.PHONY: test testcov dep vet lint clean run install build help all

## test: Run tests.
test:
	go test -v ./...

## testcov: Run tests with a coverage report.
testcov:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

## dep: Download and tidy module dependencies.
dep:
	go mod download
	go mod tidy

## vet: Run go vet.
vet:
	go vet ./...

## lint: Run golangci-lint.
lint:
	golangci-lint run -v --timeout=15m ./...

## clean: Remove build artifacts.
clean:
	rm -rf $(BUILDDIR)
	rm -f coverage.out

## run: Run the application.
run:
	go run $(CMDPATH)

## install: Install the binary into GOPATH/bin.
install:
	go install $(CMDPATH)

## build: Build the binary into $(BUILDDIR)/$(BINARY).
build:
	mkdir -p $(BUILDDIR)
	go build -ldflags '$(LDFLAGS)' -o $(BUILDDIR)/$(BINARY) $(CMDPATH)

## help: Show this help.
.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in $(BINARY):"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' | sed -e 's/^/ /'
	@echo
