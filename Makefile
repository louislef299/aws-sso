.DEFAULT_GOAL := default
.PHONY: clean assembly docs

BINARY_NAME= aws-sso

COMMIT_HASH ?= $(shell git rev-parse --short HEAD)
GOBIN = ${HOME}/go/bin
GOTRACEBACK = 'crash'
GOVERSION= $(shell go version | awk '{print $$3}')
GOFLAGS= -s -w -X 'github.com/louislef299/aws-sso/pkg/v1/version.Version=$(shell cat version.txt)' \
-X 'github.com/louislef299/aws-sso/pkg/v1/version.BuildOS=$(shell go env GOOS)' \
-X 'github.com/louislef299/aws-sso/pkg/v1/version.BuildArch=$(shell go env GOARCH)' \
-X 'github.com/louislef299/aws-sso/pkg/v1/version.GoVersion=$(GOVERSION)' \
-X 'github.com/louislef299/aws-sso/pkg/v1/version.BuildTime=$(shell date)' \
-X 'github.com/louislef299/aws-sso/pkg/v1/version.CommitHash=$(COMMIT_HASH)'

default: lint test clean $(BINARY_NAME)
	@echo "Run './$(BINARY_NAME) -h' to get started"

docs:
	@echo "Generating command documentation in docs/cmd"
	@go run main.go docs --dir docs/cmds

local: lint test $(BINARY_NAME)
	@echo "Installing $(BINARY_NAME) on your machine..."
	@go install -ldflags="$(GOFLAGS)"

$(BINARY_NAME):
	@echo "Building $(BINARY_NAME) binary for your machine..."
	@go build -mod vendor -ldflags="$(GOFLAGS)" -o $(BINARY_NAME)

test:
	@echo "Running tests..."
	@go test -v -race -cover ./...

lint: releaser-lint
	@echo "Linting Go program files"
	@golangci-lint run

update:
	go mod tidy
	go mod vendor

login:
	@gh auth status || gh auth login --git-protocol https -w -s repo,repo_deployment,workflow

releaser-lint: .goreleaser.yaml
	@echo "Checking goreleaser spec"
	@goreleaser check

release: lint test login
	@GITHUB_TOKEN=$(shell gh auth token) GOVERSION=$(GOVERSION) \
	 goreleaser release --clean

build: lint test
	@GOVERSION=$(GOVERSION) goreleaser build --clean --skip=validate

clean:
	@rm -rf $(BINARY_NAME) dist
