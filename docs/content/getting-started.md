+++
date = '2025-09-04T12:54:49-05:00'
draft = false
title = 'Getting Started'
+++

## Background

`aws-sso` was created after being frustrated with the current authentication
strategy to the cloud. After countless days of running

```bash
aws sso login
aws eks update-kubeconfig --name cluster
aws ecr get-login-password --region us-east-1 |\
docker login --username AWS --password-stdin account.dkr.ecr.us-east-1.amazonaws.com
```

and having to debug and troubleshoot other engineer's local aws configuration to
ensure they had their SSO setup properly as well as having issues authenticating
to different accounts... I decided to fix my problems to make my life easier.

Currently([v1.6][]), `aws-sso` streamlines the authentication process by
authenticating to:

- the AWS account itself
- the specified EKS cluster
- the main ECR of the account

If there are more features you would like to streamline, feel free to [open an
issue][] or make a pull request yourself.

## Installing

Installation can be done with either [homebrew][] or [krew][]:

```bash
brew tap louislef299/aws-sso
brew install --cask aws-sso

kubectl krew index add louislef299 https://github.com/louislef299/aws-sso.git
kubectl krew install louislef299/aws-sso
```

Once installed, you can check your version by running:

```bash
$ aws-sso version
AWS Auth: aws-sso/1.6.2 linux/amd64 built-with/1.25.1
 build-time/2025-09-04T17:04:55Z commit-hash/5f9d4543
```

Since `aws-sso` relies on a configuration file, there is a default config file
that was created if you ran the above command at `$HOME/.aws-sso`. All of the
config settings can be managed with `aws-sso config`, but you can also manually
override the configuration settings in this file if you run into bugs.

[Next, let's configure the damn thing so you can move on with your life!][Local
Configuration]

[v1.6]: https://github.com/louislef299/aws-sso/releases/tag/v1.6.2
[krew]: https://krew.sigs.k8s.io/
[homebrew]: https://brew.sh/
[Local Configuration]: /config
[open an issue]: https://github.com/louislef299/aws-sso/issues/new
