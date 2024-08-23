# aws-sso

`aws-sso` streamlines local AWS authentication in an idiomatic way, allowing
additional configuration to automatically authenticate to EKS and ECR. This tool
is intentionally small and narrowly scoped to attempt to make AWS authentication
as easy and repeatable as possible.

## Installing

Installation can be done with either [homebrew][] or [krew][]:

```bash
brew tap louislef299/aws-sso
brew install aws-sso
```

```bash
kubectl krew index add louislef299 https://github.com/louislef299/aws-sso.git
kubectl krew install louislef299/aws-sso
```

## Quick Start

The first time you authenticate to `aws-sso`, you will need to have gathered the
following information:

- Name and Email
- Environment Alias(for future logins)
- [AWS SSO Start URL][]
- AWS SSO Region

With that information, first set your default SSO region(example sets to
us-east-1):

```bash
aws-sso config set core.ssoregion us-east-1
```

Next, just run `aws-sso login` and the initial configuration will be triggered
and will follow a similar workflow:

```bash
louislef299 ~ % aws-sso login
please enter a prefix alias for this context(ex: env1): env1
2024/01/12 16:25:36 login.go:80: using token default
2024/01/12 16:25:36 account.go:105: couldn't find an account ID matching profile env1, using empty default...
enter your AWS access portal URL: https://your_subdomain.awsapps.com/start
2024/01/12 16:25:46 oidc.go:73: gathering client info
2024/01/12 16:25:46 oidc.go:81: registering client
2024/01/12 16:25:47 oidc.go:132: user validation code: 0000-0000
Successfully authorized!
✔ #1 aws-account-1 000000000000
2024/01/12 16:26:02 account.go:41: Selected account: aws-account-1 - 000000000000
2024/01/12 16:26:02 role.go:64: HINT: if you would like to reuse a specific iam profile, you can set core.defaultRole to your iam profile.
✔ #3 My-Role
2024/01/12 16:26:04 login.go:344: using aws role My-Role
2024/01/12 16:26:04 sts.go:138: saving data to /Users/louislef299/.aws/sso/cache/last-usage.json
If you would like to use these creds with the aws cli, please copy and paste the following command:
  export AWS_PROFILE=env1-aws-sso
2024/01/12 16:26:05 login.go:227: loading up new config env1-aws-sso
```

The above process will create a new [named profile][] in your local AWS
configuration file and cache an access token with last usage information.

Once you have logged in, you can run `aws-sso whoami` to validate a successful
login. The authentication process can be streamlined further with [configuration
settings][].

## Contributing

Feel free to open up Issues or Feature Requests, add documentation or tackle an
item from the Backlog below.

### Backlog

- Make better use of Viper pkg for config precedence
- Delete old access tokens using `token rm`
- Remove kube config context when logging out with `-c`
- Documentation website
- Add sso-section to config file so aws cli understands token hash
- Can modify sso url to add `?env=prod` to generate a new hash for prod

[AWS SSO Start URL]: https://docs.aws.amazon.com/signin/latest/userguide/iam-id-center-sign-in-tutorial.html
[homebrew]: https://brew.sh/
[krew]: https://krew.sigs.k8s.io/
[named profile]: https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html
[configuration settings]: ./docs/CONFIGURATION.md
