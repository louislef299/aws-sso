+++
date = '2025-11-23T12:00:06-06:00'
draft = false
title = 'Contributor Guide'
+++

blah blah blah this will eventually be an introduction, don't worry about it for
now. it should be short yet helpful, so I'll make it once I'm done migrating AWS
over to the Provider arch.

## Implementing a Provider

When creating a new provider, I'd recommend implementing the `provider.Provider`
interface in this order:

### 1. Start with GetConfigSchema

`GetConfigSchema()` defines what configuration fields your provider needs.
Implement this first because:

- It establishes the contract for what config your provider accepts
- The registry uses it to check for missing required fields before validation
- It powers CLI help like `knot plugin show <provider>` to display required
  config
- `ValidateConfig` and `Initialize` both depend on knowing what fields exist

The registry calls your schema via `provider.MissingConfig()` →
`getMissingConfigLocked()` → `p.GetConfigSchema()`.

Example implementation:

```go
func (p *AWS) GetConfigSchema() provider.ConfigSchema {
    return provider.ConfigSchema{
        Fields: []provider.ConfigField{
            {
                Name:        "sso_start_url",
                Type:        "string",
                Required:    true,
                Description: "The AWS SSO start URL (e.g., https://my-sso.awsapps.com/start)",
            },
            {
                Name:        "sso_region",
                Type:        "string",
                Required:    true,
                Description: "The AWS region where SSO is configured (e.g., us-east-1)",
            },
            {
                Name:        "default_role",
                Type:        "string",
                Required:    false,
                Description: "Default IAM role to assume if not specified at login time",
            },
        },
    }
}
```

### 2. Then ValidateConfig

`ValidateConfig()` validates that config values are correct beyond just being
present. By the time this is called, the registry has already verified required
fields exist using your schema.

Use this for format validation, range checks, or cross-field validation:

```go
func (p *AWS) ValidateConfig(config map[string]any) error {
    startURL, _ := config["sso_start_url"].(string)
    if !strings.HasPrefix(startURL, "https://") {
        return fmt.Errorf("sso_start_url must be an HTTPS URL")
    }
    return nil
}
```

### 3. Then Initialize

`Initialize()` is called by `provider.Activate()` after schema checks and
validation pass. Its job is to store the validated config so `Authenticate()`
can use it later:

```go
type AWS struct {
    startURL  string
    ssoRegion string
}

func (p *AWS) Initialize(ctx context.Context, config map[string]any) error {
    p.startURL = config["sso_start_url"].(string)
    p.ssoRegion = config["sso_region"].(string)
    return nil
}
```

### 4. Finally Authenticate

With config stored, `Authenticate()` can perform the actual authentication flow
using the values set during `Initialize()`.

### The Full Lifecycle

```
provider.Activate(ctx, "aws", config)
    ├── 1. GetConfigSchema()   → "What fields do you need?"
    ├── 2. Check required fields are present
    ├── 3. ValidateConfig()    → "Are the values valid?"
    └── 4. Initialize()        → "Store config, do setup"
```

... later ...

provider.Authenticate(ctx, opts)  → Uses stored config from Initialize
