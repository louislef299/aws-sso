+++
date = '2025-09-04T12:46:44-05:00'
draft = false
title = 'Local Configuration'
+++

To get started, run `aws-sso config ls`. It should print off something similar
to:

```bash
$ aws-sso config ls
Using config file /Users/louis/.aws-sso
Current config values:
config=/Users/louis
core.browser=chrome
core.defaultregion=us-east-1
core.plugins=[oidc eks ecr]
core.ssoregion=us-east-1
core.url=https://docs.aws.amazon.com/signin/latest/userguide/sign-in-urls-defined.html
ecr.disableecrlogin=false
eks.disableekslogin=false
```

This will list your `[core]` configuration values. These are variables that are
shared amongst the tool regardless of account or plugins in use. You can list
all of your configurable values with `aws-sso config values`.

Here is an example of how to set your [default SSO start url][], which is
required for proper tool functionality:

```bash
$ aws-sso config set core.url https://your_subdomain.awsapps.com/start
Using config file /Users/louis/.aws-sso
set core.url to https://your_subdomain.awsapps.com/start

$ aws-sso config list
Using config file /Users/louis/.aws-sso
Current config values:
config=/Users/louis
core.browser=chrome
core.defaultregion=us-east-1
core.plugins=[oidc eks ecr]
core.ssoregion=us-east-1
core.url=https://your_subdomain.awsapps.com/start
ecr.disableecrlogin=false
eks.disableekslogin=false
```

Once you have set `core.url` and validated that the `core.ssoregion` reflects
your [AWS Organization region][], it's time to setup your first account!

## Accounts

[Accounts][] in `aws-sso` represents the actual AWS Account to target when
signing in. The only required inputs are the actual Account ID and the Account
name(alias), otherwise the tool defaults are used. Core defaults can be
overridden at the Account level, so that they only apply to that specific
account. These include things like SSO URL, region and whether to enable private
browsing or not.

Accounts can be added to `aws-sso` by running the [`account add` command][].
Here is an example:

```bash
$ aws-sso acct add --name prod --id 111111111111
2025/09/04 15:04:31 account.go:54: couldn't find an account URL matching profile , using core default...
2025/09/04 15:04:31 account.go:69: associated account prod to account number 111111111111

# Note: accts is a hidden command
$ aws-sso accts
Account mapping:
dev:
  ID: 000000000000
  Region: us-east-2
  Private: false
  Token: default
  SSO URL: (default) https://your_subdomain.awsapps.com/start
prod:
  ID: 111111111111
  Region: us-east-1
  Private: false
  Token: default
  SSO URL: https://your_subdomain.awsapps.com/start
```

To override values, run the [`account set` command][]:

```bash
$ aws-sso acct set dev --id 222222222222
2025/09/04 15:07:45 account.go:54: couldn't find an account URL matching profile , using core default...
account values have been set for dev
```

If you want to remove accounts, you need to edit `$HOME/.aws-sso` in a text
editor. `aws-sso` uses [viper][] under the hood, so if you would like deleting
to be a part of the service offering, [let their team know][]!

## Logging In

With your SSO start URL configured and your account setup, let's login!

```bash
$ aws-sso login dev
2025/09/04 15:11:47 login.go:89: using token default
2025/09/04 15:11:47 oidc.go:79: using sso region us-east-1 to login
2025/09/04 15:11:47 oidc.go:191: browser set to default(use cookies)
2025/09/04 15:11:47 oidc.go:82: gathering client info
2025/09/04 15:11:47 oidc.go:61: checking config location /Users/louis/.aws/sso/cache/hash-redacted.json
2025/09/04 15:11:47 oidc.go:90: registering client
2025/09/04 15:11:47 oidc.go:147: could not start device authorization
2025/09/04 15:11:47 oidc.go:89: couldn't log into AWS: operation error SSO OIDC: StartDeviceAuthorization, https response error StatusCode: 400, RequestID: 02bd42e4-ea7a-42f4-8328-9b61d95f7761, InvalidRequestException:
```

Since my start URL wasn't properly configured, I'm getting an error. There is
[currently an issue][issue #626] that will address unhelpful error messages, but
you should be able to get an idea of where the error is happening and what is
going on. Feel free to reach out to me over an issue and I can try to help out.

If you successfully logged in, feel free to validate it with `aws-sso whoami`.
Hopefully you find this tool useful!

## Browser Utilization

Setting the `core.browser` config value allows the user to open their browser in
private(incognito) mode for gathering tokens. The currently supported browsers
alongside operating systems as of v1.6 are:

browser           | linux | darwin | windows
 ---------------- | ----- | ------ | -------
brave             | ✅    | ✅     | ✅
chrome            | ✅    | ✅     | ✅
firefox           | ✅    | ✅     | ❌
firefox-developer | ✅    | ✅     | ❌

## Token Usage

If you have multiple SSO sessions you need to manage, you can create additional
tokens to allow authentication caching across multiple AWS Organizations. Let's
say you have your development account in one Organization and a production
account in another. It would be nice to associate the session to a specific
token so that you don't have to login and logout all the time.

In this example, we will keep the `default` token as representing our dev org
and `prod` will be our prod org token. Let's list our existing tokens, add a
prod token, and associate the prod token with our prod account:

```bash
$ aws-sso tokens
Local Tokens:
* default

$ aws-sso token add prod
2025/09/05 09:31:56 successfully added token prod!

# Notice that creating the token associates it with the current session
$ aws-sso tokens        
Local Tokens:
  default
* prod

$ aws-sso acct set prod --token prod
account values have been set for prod

$ aws-sso accts
Account mapping:
dev:
  ID: 000000000000
  Region: us-east-2
  Private: false
  Token: default
  SSO URL: (default) https://your_subdomain.awsapps.com/start
prod:
  ID: 111111111111
  Region: us-east-1
  Private: false
  Token: prod
  SSO URL: https://your_subdomain.awsapps.com/start
```

[Accounts]: https://pkg.go.dev/github.com/louislef299/aws-sso/internal/account
[AWS Organization region]: https://docs.aws.amazon.com/organizations/latest/userguide/region-support.html
[default SSO start url]: https://docs.aws.amazon.com/signin/latest/userguide/sign-in-urls-defined.html
[issue #626]: https://github.com/louislef299/aws-sso/issues/626
[let their team know]: https://forms.gle/R6faU74qPRPAzchZ9
[viper]: https://pkg.go.dev/github.com/spf13/viper
[`account add` command]: {{< ref "cmds/aws-sso_account_add" >}}
[`account set` command]: {{< ref "cmds/aws-sso_account_set" >}}
