+++
date = '2025-09-05T15:55:03-05:00'
draft = false
title = 'Plugin Architecture'
+++

The `aws-sso` tool currently is *attempting* to use a plugin architecture for
each step of the authentication process. This gives developers(mostly just me)
the flexibility to try out newer features in isolation without breaking the
current tool. As of v1.6, there are three main plugins:

1. OIDC Plugin (Always Enabled - Core Auth)
    - Handles AWS SSO authentication
    - Manages token retrieval

2. EKS Plugin (Optional - Can be disabled)
    - Configures kubectl credentials
    - Creates kubeconfig entries
    - Can be disabled with `eks.disableEKSLogin` config value or
      `--disableEKSLogin`

3. ECR Plugin (Optional - Can be disabled)  
    - Configures Docker credentials
    - Allows container image pulls
    - Can be disabled with `ecr.disableECRLogin` config value or
      `--disableECRLogin`

```bash
                        +-------------------+
                        |                   |
                        |    aws-sso CLI    |
                        |                   |
                        +-------------------+
                                |
                                v
        +-----------------------------------------------------+
        |                                                     |
        |                  Login Command                      |
        |                                                     |
        +-----------------------------------------------------+
                                |
                                v
    +-------------------------------------------------------------+
    |                       Plugin Registry                       |
    |      (dlogin package - manages plugin registration)         |
    +-------------------------------------------------------------+
                |               |                 |
                v               v                 v
+------------------+    +------------------+    +------------------+
|                  |    |                  |    |                  |
|   OIDC Plugin    |    |    EKS Plugin    |    |    ECR Plugin    |
|   (Required)     |    |    (Optional)    |    |    (Optional)    |
|                  |    |                  |    |                  |
+------------------+    +------------------+    +------------------+
        |                       |                       |
        v                       v                       v
+------------------+    +------------------+    +------------------+
|                  |    |                  |    |                  |
| AWS SSO/OIDC     |    | Kubernetes       |    | Docker Registry  |
| Authentication   |    | Authentication   |    | Authentication   |
|                  |    |                  |    |                  |
+------------------+    +------------------+    +------------------+
        |
        v
+------------------+
|                  |
|   AWS Console    |
|   Access         |
|                  |
+------------------+
```

## Provider Interface (v2 Architecture)

> **Status**: In development on `knot` branch

The `dlogin.ILogin` interface is being replaced by a more comprehensive
`provider.Provider` interface that supports the full credential lifecycle and
multiple authentication protocols (OIDC, SAML, OAuth2).

### Design Decisions

| Decision | Rationale |
|----------|-----------|
| Provider replaces ILogin | Unified abstraction with credential lifecycle (refresh, revoke) |
| Initialize receives static config | Endpoints, client IDs go in Initialize; per-auth options (MFA, profile) go in AuthOptions |
| Proactive refresh | Callers should refresh before credential expiry, not after |
| Revoke serves dual purpose | Handles both logout filesystem cleanup and explicit OAuth token revocation |
| Activate calls Initialize | Single entry point ensures providers are never left in broken state |

### Provider Interface

```go
type Provider interface {
    Name() string                                                    // Unique identifier
    Type() Type                                                      // "oidc", "saml", "oauth2"

    // Lifecycle
    Initialize(ctx context.Context, config map[string]any) error     // One-time setup
    Authenticate(ctx context.Context, opts AuthOptions) (*Credentials, error)
    Refresh(ctx context.Context, creds *Credentials) (*Credentials, error)
    Revoke(ctx context.Context, creds *Credentials) error

    // Metadata
    GetConfigSchema() ConfigSchema                                   // For CLI flag generation
    ValidateConfig(config map[string]any) error                      // Early validation
}
```

### Registry State Model

The provider registry manages four states to ensure providers are never left
in a broken state and failures are clearly visible to users:

```
┌──────────────────┐
│    Registered    │◀───────────────────────────────────┐
└──────────────────┘                                    │
        │                                               │
        │ Activate()                                    │ Deactivate()
        ▼                                               │
   ┌─────────┐     success      ┌─────────────────────┐ │
   │ Validate│─────────────────▶│Active + Initialized │─┤
   │   +     │                  └─────────────────────┘ │
   │  Init   │                                          │
   └─────────┘                                          │
        │                                               │
        │ failure                                       │
        ▼                                               │
┌──────────────────┐                                    │
│     Invalid      │────────────────────────────────────┘
│  (stores error)  │
└──────────────────┘
        │       ▲
        │       │ still fails
        └───────┘
     Activate() auto-retries
```

| State | Description |
|-------|-------------|
| **Registered** | Provider is known (via init import) but never activated |
| **Invalid** | Activation failed; error stored via `GetInvalidReason()` |
| **Active + Initialized** | Config validated, Initialize() called, ready for auth |

Key behaviors:
- `Activate(ctx, name, config)` validates AND initializes atomically
- On failure, provider enters **Invalid** state (not silent rollback)
- Calling `Activate` on Invalid provider **auto-retries**
- `Deactivate` clears any state back to Registered

### Configuration Split

```
┌─────────────────────────────────────────────────────────────────┐
│                     Initialize(config)                          │
│  Static provider setup - doesn't change per authentication      │
│  Examples: sso_start_url, client_id, sso_region                 │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                   Authenticate(AuthOptions)                     │
│  Per-authentication parameters - may vary each call             │
│  Examples: profile, region, MFA codes, session preferences      │
└─────────────────────────────────────────────────────────────────┘
```

### Example Usage

```go
// At init time (automatic via import)
func init() {
    provider.Register(&AWSProvider{})
}

// At runtime
missing, _ := provider.MissingConfig("aws-oidc", userConfig)
if len(missing) > 0 {
    // Interactive: prompt for missing fields
    // Non-interactive: print "run: knot config set provider.field <value>"
}

// Activate validates config AND calls Initialize
err := provider.Activate(ctx, "aws-oidc", userConfig)
if err != nil {
    // Provider is now Invalid - error is stored and visible
    if provider.IsInvalid("aws-oidc") {
        reason, _ := provider.GetInvalidReason("aws-oidc")
        fmt.Printf("Provider failed: %v\n", reason)
    }
    // Fix config and retry - Activate auto-retries on Invalid
    err = provider.Activate(ctx, "aws-oidc", fixedConfig)
}

// Now safe to use - Initialize was already called
for _, p := range provider.ActiveProviders() {
    creds, _ := p.Authenticate(ctx, opts)
    // Later: p.Refresh(ctx, creds) before expiry
    // On logout: p.Revoke(ctx, creds)
}

// List provider states for CLI
for _, p := range provider.All() {
    status := provider.Status(p.Name())  // "registered", "active", "invalid"
    fmt.Printf("%s: %s\n", p.Name(), status)
}
```

### Migration from dlogin

| dlogin.ILogin | provider.Provider |
|---------------|-------------------|
| `Init(cmd)` | `GetConfigSchema()` + CLI generates flags |
| `Login(ctx, config)` | `Authenticate(ctx, opts)` |
| `Logout(ctx, config)` | `Revoke(ctx, creds)` |
| (none) | `Initialize(ctx, config)` - one-time setup |
| (none) | `Refresh(ctx, creds)` - proactive refresh |
| (none) | `ValidateConfig(config)` - early validation |
