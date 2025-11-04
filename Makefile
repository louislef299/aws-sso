.DEFAULT_GOAL := default
.PHONY: clean assembly docs

BINARY_NAME= aws-sso

COMMIT_HASH ?= $(shell git rev-parse --short HEAD)
GOBIN = ${HOME}/go/bin
GOTRACEBACK = 'crash'
GOVERSION= $(shell go version | awk '{print $$3}')
GOFLAGS= -s -w -X 'github.com/louislef299/aws-sso/pkg/version.Version=$(shell cat version.txt)' \
-X 'github.com/louislef299/aws-sso/pkg/version.BuildOS=$(shell go env GOOS)' \
-X 'github.com/louislef299/aws-sso/pkg/version.BuildArch=$(shell go env GOARCH)' \
-X 'github.com/louislef299/aws-sso/pkg/version.GoVersion=$(GOVERSION)' \
-X 'github.com/louislef299/aws-sso/pkg/version.BuildTime=$(shell date)' \
-X 'github.com/louislef299/aws-sso/pkg/version.CommitHash=$(COMMIT_HASH)'

default: lint test clean $(BINARY_NAME)
	@echo "Run './$(BINARY_NAME) -h' to get started"

docs:
	@echo "Generating command documentation in docs/cmd"
	@go run main.go docs --dir docs/content/cmds

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
	@golangci-lint run -c .golangci.yml

update:
	go mod tidy
	go mod vendor

login:
	@gh auth status || gh auth login --git-protocol https -w -s repo,repo_deployment,workflow

releaser-lint: .goreleaser.yaml
	@echo "Checking goreleaser spec"
	@goreleaser check

release: lint test login
	@echo "the build won't get signed if GPGKEYID isn't set"
	@GITHUB_TOKEN=$(shell gh auth token) GOVERSION=$(GOVERSION) \
	  GPG_TTY=$(shell tty) goreleaser release --clean

build: lint test
	@echo "the build won't get signed if GPGKEYID isn't set"
	@GOVERSION=$(GOVERSION) GPG_TTY=$(shell tty) \
	  goreleaser build --clean --skip=validate

scan: $(BINARY_NAME) license-scan
	trivy rootfs --scanners vuln --format cyclonedx --output $(BINARY_NAME).cyclonedx.json .
	trivy sbom $(BINARY_NAME).cyclonedx.json

license-scan:
	trivy fs --scanners license --license-full  .

clean:
	@rm -rf $(BINARY_NAME) $(BINARY_NAME).cyclonedx.json dist .hugo_build.lock
