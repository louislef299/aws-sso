# aws-sso

`aws-sso` streamlines local AWS authentication in an idiomatic way, allowing
additional configuration to automatically authenticate to EKS and ECR. This tool
makes AWS authentication easy and repeatable.

## Documentation

For full documentation including installation, configuration, and usage
instructions, visit:
[https://aws-sso.netlify.app/](https://aws-sso.netlify.app/)

## Quick Install

Install with [homebrew][]:

```bash
brew tap louislef299/aws-sso
brew install --cask aws-sso
```

Or with [krew][]:

```bash
kubectl krew index add louislef299 https://github.com/louislef299/aws-sso.git
kubectl krew install louislef299/aws-sso
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
[homebrew]: https://brew.sh/
[krew]: https://krew.sigs.k8s.io/
