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

### 6. Implement Refresh

The `Refresh()` method handles credential renewal. The approach varies by
provider type and the authentication mechanism used.

#### Refresh Method Signature

```go
Refresh(ctx context.Context, creds *Credentials, opts AuthOptions) (*Credentials, error)
```

The `opts` parameter provides context for providers that require re-authentication,
allowing them to preserve user preferences (private browser, skip defaults, etc.).

#### Implementation Patterns

**For providers with refresh tokens (OAuth2, OIDC with offline_access):**

These providers can refresh credentials silently without user interaction:

```go
func (p *GenericOIDC) Refresh(ctx context.Context, creds *provider.Credentials, opts provider.AuthOptions) (*provider.Credentials, error) {
    // Check if still valid
    if creds.Expiry.After(time.Now()) {
        return creds, nil
    }
    
    // Extract refresh token
    refreshToken := creds.RefreshToken
    if refreshToken == "" {
        return nil, fmt.Errorf("no refresh token available")
    }
    
    // Use OAuth2 token refresh
    newToken, err := p.oauthClient.RefreshToken(ctx, refreshToken)
    if err != nil {
        return nil, fmt.Errorf("token refresh failed: %w", err)
    }
    
    // Return updated credentials
    return &provider.Credentials{
        Type:         creds.Type,
        AccessToken:  newToken.AccessToken,
        RefreshToken: newToken.RefreshToken, // May be rotated
        Expiry:       newToken.Expiry,
        Metadata:     creds.Metadata, // Preserve existing metadata
    }, nil
}
```

**For providers requiring re-authentication (AWS SSO, device flow):**

These providers need user interaction, so they reuse `Authenticate()` with
preserved context:

```go
func (p *AWS) Refresh(ctx context.Context, creds *provider.Credentials, opts provider.AuthOptions) (*provider.Credentials, error) {
    // Check if still valid
    if creds.Expiry.After(time.Now()) {
        log.Printf("token still valid until %s, no refresh needed", creds.Expiry)
        return creds, nil
    }
    
    log.Println("token expired, performing re-authentication")
    
    // Extract context from old credentials to minimize user prompts
    if opts.Extra == nil {
        opts.Extra = make(map[string]any)
    }
    
    // Reuse account_id to skip account selection
    if accountID, ok := creds.Metadata["account_id"].(string); ok {
        opts.Extra["account_id"] = accountID
    }
    
    // Set refresh flag
    opts.Extra["refresh"] = true
    
    // Reuse region if not specified
    if opts.Region == "" {
        if region, ok := creds.Metadata["region"].(string); ok {
            opts.Region = region
        }
    }
    
    // Delegate to Authenticate() which handles the device flow
    return p.Authenticate(ctx, opts)
}
```

#### Key Principles

1. **Always check if still valid first** - If credentials haven't expired, return
   them immediately to avoid unnecessary work
2. **Preserve context** - Extract account IDs, regions, or other context from old
   credentials to minimize user prompts
3. **Use opts for preferences** - Respect browser preferences, defaults, etc. from
   the opts parameter
4. **Be clear about user interaction** - Log when user interaction is required
5. **Return detailed errors** - Help users understand why refresh failed

### 7. Implement Revoke

The `Revoke()` method handles server-side credential revocation. It should
**only** handle revoking tokens with the identity provider, not local file
cleanup (that's handled by `CleanCredentials()` in NativeCLIIntegration).

#### Revoke Method Signature

```go
Revoke(ctx context.Context, creds *Credentials) error
```

#### Implementation Guidelines

**Revoke should:**
- Call the provider's token revocation endpoint (if available)
- Clean up provider-specific cached state (e.g., cached OIDC client info)
- Be idempotent (revoking already-revoked credentials should not error)
- Return wrapped errors with helpful context

**Revoke should NOT:**
- Clean native CLI credential files (use `CleanCredentials()` for that)
- Require user interaction
- Fail if the token is already expired/invalid

#### Separation of Concerns

- **`Revoke()`**: Server-side token revocation + provider-specific cache cleanup
- **`CleanCredentials()`** (NativeCLIIntegration): Local native CLI file cleanup

The caller orchestrates: `provider.Revoke()` → `provider.CleanCredentials()`

#### Example Implementation (AWS SSO)

```go
func (p *AWS) Revoke(ctx context.Context, creds *provider.Credentials) error {
    log.Println("revoking AWS SSO credentials")
    
    // Extract access token
    accessToken := creds.AccessToken
    if accessToken == "" {
        return fmt.Errorf("no access token available in credentials")
    }
    
    // Load AWS config
    awsCfg, err := awsConfig.LoadDefaultConfig(ctx, awsConfig.WithRegion(p.ssoRegion))
    if err != nil {
        return fmt.Errorf("failed to load AWS config for revocation: %w", err)
    }
    
    // Call AWS SSO Logout API (server-side revocation)
    if err := awsSDK.Logout(ctx, &awsCfg, accessToken); err != nil {
        return fmt.Errorf("failed to revoke AWS SSO token: %w", err)
    }
    
    // Clean up cached OIDC client info (provider-specific state)
    clientInfoPath, err := getClientInfoPath()
    if err != nil {
        return fmt.Errorf("failed to determine client info path: %w", err)
    }
    
    // Remove cache file (idempotent)
    if err := os.Remove(clientInfoPath); err != nil && !os.IsNotExist(err) {
        return fmt.Errorf("failed to remove cached client info: %w", err)
    }
    
    log.Println("successfully revoked AWS SSO credentials")
    return nil
}
```

#### Example CleanCredentials (NativeCLIIntegration)

```go
func (p *AWS) CleanCredentials(ctx context.Context, profile string) error {
    log.Printf("cleaning native credentials for profile: %s", profile)
    
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return fmt.Errorf("failed to get home directory: %w", err)
    }
    
    // Clean ~/.aws/credentials
    credentialsFile := homeDir + "/.aws/credentials"
    if err := removeProfileFromINI(credentialsFile, profile); err != nil {
        return fmt.Errorf("failed to clean credentials file: %w", err)
    }
    
    // Clean ~/.aws/config (profiles are prefixed with "profile ")
    configFile := homeDir + "/.aws/config"
    configProfile := profile
    if profile != "default" {
        configProfile = "profile " + profile
    }
    if err := removeProfileFromINI(configFile, configProfile); err != nil {
        return fmt.Errorf("failed to clean config file: %w", err)
    }
    
    return nil
}
```

#### Key Principles

1. **Separation of concerns** - Revoke handles server-side, CleanCredentials handles local files
2. **Idempotent** - Multiple calls should not fail
3. **Graceful failure** - Expired tokens should not cause revocation to fail
4. **Wrapped errors** - Return helpful error messages with context
5. **Log progress** - Help users understand what's happening

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
