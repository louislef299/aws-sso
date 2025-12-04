# https://cheatography.com/linux-china/cheat-sheets/justfile/
[doc('List out available recipes')]
default: 
    @just --list

set shell := ["zsh", "-uc"]

go := require("go")

BINARY_NAME := "knot"
GPG_SIGNING_KEY := shell('git config user.signingkey || echo "not set"')
GOBIN := env('HOME') + '/go/bin'
GOTRACEBACK := "crash"
export GOVERSION := shell("go version | awk '{print $3}'")
GOFLAGS := (
    "-s -w " +
    "-X 'github.com/louislef299/knot/internal/version.Version=local.dev' " +
    "-X 'github.com/louislef299/knot/internal/version.BuildOS=" + `go env GOOS` + "' " +
    "-X 'github.com/louislef299/knot/internal/version.BuildArch=" + `go env GOARCH` + "' " +
    "-X 'github.com/louislef299/knot/internal/version.GoVersion=" + GOVERSION + "' " +
    "-X 'github.com/louislef299/knot/internal/version.BuildTime=" + `date` + "' " +
    "-X 'github.com/louislef299/knot/internal/version.CommitHash=" + `git rev-parse --short HEAD` + "'"
)
GO_LOCATIONS := "cmd internal pkg plugins main.go"

# Run the go program with provided inputs
run *INPUT:
    {{go}} run main.go {{INPUT}}

alias b := build
# Build local.dev version of aws-sso
build:
    @echo "Building {{BINARY_NAME}} binary for your machine..."
    @{{go}} build -mod vendor -ldflags="{{GOFLAGS}}" -o {{BINARY_NAME}}

# Formats Go targets defined by GO_LOCATIONS
fmt:
    #!/usr/bin/env zsh
    for tgt in {{GO_LOCATIONS}}; do \
        gofmt -s -w $tgt &
    done
    wait

# Runs the Go linters
lint:
    @echo "Linting Go files"
    golangci-lint run -c .golangci.yml
    @echo "Checking GoReleaser spec"
    goreleaser check

# Run a verbose Go test with coverage
test:
	@echo "Running tests..."
	@{{go}} test -v -race -cover ./...

# Install a local.dev version in your go/bin
install:
    go install -ldflags="{{GOFLAGS}}"

# Run GoReleaser to build local dist folder
dist: lint test
	@GOVERSION={{GOVERSION}} GPG_SIGNING_KEY={{GPG_SIGNING_KEY}} \
	 goreleaser build --clean --skip=validate --id aws-sso

# Authenticate to GitHub with the gh cli
login:
	@gh auth status || gh auth login --git-protocol https -w -s repo,repo_deployment,workflow

# Validate Git is healthy for release
check-vcs:
    @echo "Ensuring HEAD commit is tagged"
    @git tag --points-at HEAD | grep -q . || \
      (echo "ERROR: HEAD commit is not tagged!(run git tag -s)" && exit 1)
    @echo "Making sure this is the main branch"
    @if [[ "$(git branch --show-current)" != "main" ]]; then exit 1 ; fi;

# Perform a signed release to GitHub
[linux]
[macos]
release: check-vcs lint test login
    @echo "WARNING: the build won't get signed if GPG_SIGNING_KEY isn't set"
    @GITHUB_TOKEN=`gh auth token` GOVERSION={{GOVERSION}} \
      GPG_TTY=`tty` GPG_SIGNING_KEY={{GPG_SIGNING_KEY}} \
      goreleaser release --clean

# Build dist folder, scan binaries & generate SBOM
scan: dist
    @echo "Running Trivy vulnerability scanner:"
    trivy rootfs --scanners vuln,license --license-full --format cyclonedx \
    --output {{BINARY_NAME}}.cyclonedx.json .
    @echo "Generating SBOM {{BINARY_NAME}}.cyclonedx.json"
    trivy sbom {{BINARY_NAME}}.cyclonedx.json

# Generate command documentation
docs:
	{{go}} run main.go docs --dir docs/content/cmds

# Run hugo docs website
serve:
    git submodule update --remote --rebase
    hugo serve -s docs

# Cleanup the filesystem
clean:
	rm -rf {{BINARY_NAME}} {{BINARY_NAME}}.cyclonedx.json dist .hugo_build.lock
