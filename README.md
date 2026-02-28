[![GoDoc reference example](https://img.shields.io/badge/godoc-reference-blue.svg)](https://pkg.go.dev/github.com/louislef299/aws-sso)
[![Go Report Card](https://goreportcard.com/badge/github.com/louislef299/aws-sso)](https://goreportcard.com/report/github.com/louislef299/aws-sso)
[![Releases](https://img.shields.io/github/release-pre/louislef299/aws-sso.svg)](https://github.com/louislef299/aws-sso/releases)
[![Github all releases](https://img.shields.io/github/downloads/louislef299/aws-sso/total.svg)](https://github.com/louislef299/aws-sso/releases/)

# aws-sso

`aws-sso` streamlines local AWS authentication in an idiomatic way, allowing
additional configuration to automatically authenticate to EKS and ECR. This tool
makes AWS authentication easy and repeatable.

## Documentation

For full documentation including installation, configuration, and usage
instructions, visit:
[https://aws-sso.louislefebvre.net](https://aws-sso.louislefebvre.net)

## Quick Install

Install with [homebrew][]:

```bash
brew tap louislef299/aws-sso && brew install --cask aws-sso
```

Or with [krew][]:

```bash
kubectl krew index add louislef299 https://github.com/louislef299/aws-sso.git && \
kubectl krew install louislef299/aws-sso
```

Or manually download from [release assets][]

Or build from source: `just build`

### Verify the Binary

This strategy uses [GnuPG][] and is only required if you installed `aws-sso`
with brew, krew or manually from the release assets. You basically need to
import my public key and then verify the signature of the binary. In the
example, `$BINPATH` will represent the path to the `aws-sso` binary.

```bash
# Import my PGP
curl -s https://louislefebvre.net/public-key.txt | gpg --import

gpg --verify $BINPATH/aws-sso*.sig aws-sso
```

## Basic Usage

1. [Configure your system][]

2. Log in:

    ```bash
    aws-sso login
    ```

3. Check your authentication:

    ```bash
    aws-sso whoami
    ```

## Contributing

Feel free to open up Issues or Feature Requests on GitHub.

[Configure your system]: https://aws-sso.netlify.app/config/
[GnuPG]: https://www.gnupg.org/gph/en/manual/book1.html
[homebrew]: https://brew.sh/
[krew]: https://krew.sigs.k8s.io/
[release assets]: https://github.com/louislef299/aws-sso/releases
