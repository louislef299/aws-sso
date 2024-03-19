# Local Configuration

There are certain configurations you can set in your local environment to make
the `aws-sso` cli experience smoother. You can always list all of the
configurable values with the `aws-sso config values` command.

Current configuration values v1.1.2:

```bash
The following values are available for configuration:
core.defaultCluster:		The default cluster to target when logging in, supports go regex expressions(golang.org/s/re2syntax)
core.defaultRole:		The default iam role to use when logging in
core.region:			The default region used when a region is not found in your environment or set with the --region flag
core.url:			The default sso start url used when logging in
core.disableEKSLogin:		Disables automatic detection and login for EKS
core.disableECRLogin:		Disables automatic detection and login for ECR
core.browser:			The default browser to use. Required for advanced features like opening in a private browser
<BOUND_FLAG>session.region:	The region you would like to use at login
<BOUND_FLAG>session.url:	The AWS SSO start url
<BOUND_FLAG>session.role:	The IAM role to use when logging in
```

## Accounts

Account aliases can be added manually to `aws-sso` by running the [`account add`
command][]. Otherwise, the cli will read in pre-existing named profiles in your
local AWS Config file.

## Browser Utilization

Setting the `core.browser` config value allows the user to open their browser in
private(incognito) mode for gathering tokens. The currently supported browsers
alongside operating systems as of v1.1.2 are:

browser | linux | darwin | windows
 ------ | ----- | ------ | -------
brave   | ✅    | ✅     | ✅
chrome  | ✅    | ✅     | ✅
firefox | ✅    | ✅     | ✅

[`account add` command]: ./cmds/aws-sso_account_add.md
