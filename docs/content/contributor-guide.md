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

The `Authenticate()` method should return a `provider.Credentials` struct
containing:

- **AccessToken**: The primary authentication token (e.g., SSO access token)
- **RefreshToken**: Optional token for refreshing credentials
- **Expiry**: When the access token expires
- **Metadata**: Provider-specific data as `map[string]any`

For the Metadata field, include any information needed for:
- Credential refresh operations
- Native CLI integration (if applicable)
- Additional provider-specific context

Example:

```go
func (p *AWS) Authenticate(ctx context.Context, opts provider.AuthOptions) (*provider.Credentials, error) {
    // Perform authentication flow...
    
    return &provider.Credentials{
        Type:         string(provider.TypeOIDC),
        AccessToken:  clientInfo.AccessToken,
        RefreshToken: "", // If not using refresh tokens
        Expiry:       clientInfo.AccessTokenExpiresAt,
        Metadata: map[string]any{
            // Client information for future operations
            "client_id":     clientInfo.ClientId,
            "client_secret": clientInfo.ClientSecret,
            
            // Native CLI credentials (if implementing NativeCLIIntegration)
            "aws_access_key_id":     roleCredentials.AccessKeyId,
            "aws_secret_access_key": roleCredentials.SecretAccessKey,
            "aws_session_token":     roleCredentials.SessionToken,
            
            // Additional context
            "account_id": accountID,
            "role_name":  roleName,
        },
    }, nil
}
```

### 5. Optional: Implement NativeCLIIntegration

The `provider.NativeCLIIntegration` interface is **optional** and should only be
implemented if your provider needs to integrate with a native CLI tool that
expects credentials in a specific file format.

#### When to Implement NativeCLIIntegration

Implement this interface if:

- Your provider has a native CLI tool (e.g., AWS CLI, Azure CLI, gcloud)
- The CLI expects credentials in a specific file location/format
- Users will want seamless integration with existing CLI workflows
- Credentials need to be written beyond knot's internal cache

#### When NOT to Implement NativeCLIIntegration

Skip this interface if:

- Your provider is a generic OIDC/OAuth provider without a CLI
- Your provider only needs token caching in knot's credential store
- The provider is for a custom/enterprise SSO without external tools
- Credentials are only used within knot itself

#### Implementation Example

```go
// WriteNativeCredentials writes credentials in the format expected by
// the provider's native CLI tool.
func (p *AWS) WriteNativeCredentials(ctx context.Context, creds *provider.Credentials, profile string) error {
    // Extract credentials from metadata
    accessKeyID, ok := creds.Metadata["aws_access_key_id"].(string)
    if !ok {
        return fmt.Errorf("missing aws_access_key_id in credentials")
    }
    
    // Write to native CLI format (e.g., ~/.aws/credentials)
    return writeAWSCredentialsFile(profile, accessKeyID, ...)
}

// CleanNativeCredentials removes credentials written by WriteNativeCredentials
func (p *AWS) CleanNativeCredentials(ctx context.Context, profile string) error {
    // Remove credentials from native CLI files
    return cleanAWSCredentialsFile(profile)
}
```

#### Usage in the Authentication Flow

When a provider implements `NativeCLIIntegration`, knot will automatically
detect and use it:

```go
// Authenticate with provider
creds, err := provider.Authenticate(ctx, opts)

// Save to knot's credential cache
credStore.Save(profile, creds)

// If provider supports native CLI integration, write those too
if nativeCLI, ok := provider.(NativeCLIIntegration); ok {
    nativeCLI.WriteNativeCredentials(ctx, creds, profile)
}
```

This design keeps the core `Provider` interface focused on authentication while
allowing optional integration with native tooling. Providers that only need
token-based authentication don't have to implement file writing logic they'll
never use.

### The Full Lifecycle

```
provider.Activate(ctx, "aws", config)
    ├── 1. GetConfigSchema()   → "What fields do you need?"
    ├── 2. Check required fields are present
    ├── 3. ValidateConfig()    → "Are the values valid?"
    └── 4. Initialize()        → "Store config, do setup"
```

... later ...

```
provider.Authenticate(ctx, opts)
    → Returns credentials with tokens and metadata

credStore.Save(profile, creds)
    → Caches credentials in knot's storage

if provider implements NativeCLIIntegration:
    provider.WriteNativeCredentials(ctx, creds, profile)
        → Writes to native CLI format (e.g., ~/.aws/credentials)
```
