# aws-sso

`aws-sso` makes it easier to authenticate into your EKS cluster. Currently, it
can only be installed with homebrew, but the plugin will is expected to be
brought into krew as it was meant to be a kubectl plugin.

Documentation on the way! This repository needs further testing, but feel free
to use the current implemenation and let me know what can be improved.

## Installation via homebrew

```bash
brew tap louislef299/aws-sso
brew install aws-sso
```

## Installation via krew

```bash
kubectl krew index add louislef299 https://github.com/louislef299/aws-sso.git
kubectl krew install louislef299/aws-sso
```

## Issues

- ✅ An empty `login` returns a strange default profile name. Potentially prompt
  for a profile name?
- ✅ Need to validate the security phrase at login
- ✅ Let the user switch between multiple Tokens
- Delete old access tokens during `token rm`
- On-board additional operating systems and browsers for private flag
- Remove kube config context when logging out with `-c`
- Allow for impersonation with kubeconfig
- Allow assuming another rule with the existing credentials
- SAML integration
