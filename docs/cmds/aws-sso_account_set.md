## aws-sso account set

Set the AWS account profile values.

```
aws-sso account set [flags]
```

### Examples

```
  aws-sso account set env1 --number 000000000 --region us-west-2
```

### Options

```
  -h, --help            help for set
      --number string   The account number of the account associated to the account name
  -r, --region string   The default region to associate to the account
```

### Options inherited from parent commands

```
      --as string               Username to impersonate for the operation. User could be a regular user or a service account in a namespace.
      --as-group stringArray    Group to impersonate for the operation, this flag can be repeated to specify multiple groups.
      --commandTimeout string   the default timeout for network commands executed (default "3s")
```

### SEE ALSO

* [aws-sso account](aws-sso_account.md)	 - Manage AWS account aliases

###### Auto generated by spf13/cobra on 21-Jul-2024
