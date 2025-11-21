# knot

🚧 **THIS BRANCH IS UNDER CONSTRUCTION** 🚧

Use at your own risk, but follow along with issue [#817][]

---

`knot` streamlines local AWS authentication in an idiomatic way, allowing
additional configuration to automatically authenticate to EKS and ECR. This tool
makes AWS authentication easy and repeatable.

## Documentation

For full documentation including installation, configuration, and usage
instructions, visit:
[https://aws-sso.netlify.app/](https://aws-sso.netlify.app/)

## Quick Install

Install with [homebrew][]:

```bash
brew tap louislef299/aws-sso && brew install --cask aws-sso
```

Or with [krew][]:

```bash
kubectl krew index add louislef299 https://github.com/louislef299/knot.git && \
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

[#817]: https://github.com/louislef299/aws-sso/issues/817
[Configure your system]: https://aws-sso.netlify.app/config/
[GnuPG]: https://www.gnupg.org/gph/en/manual/book1.html
[homebrew]: https://brew.sh/
[krew]: https://krew.sigs.k8s.io/
[release assets]: https://github.com/louislef299/knot/releases
