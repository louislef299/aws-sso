// Package provider defines the interface that all authentication providers must
// implement to be compatible with the knot provider plugin system.
//
// This package replaces the legacy dlogin.ILogin interface with a more
// comprehensive abstraction that supports credential lifecycle management,
// configuration validation, and multiple authentication protocols (OIDC, SAML,
// OAuth2).
//
// Each provider is responsible for registering itself (typically in an init
// function) and implementing the full authentication lifecycle:
//
//   - Initialize: One-time setup with provider configuration (endpoints, client IDs)
//   - Authenticate: Perform authentication and acquire credentials
//   - Refresh: Proactively refresh credentials before expiry
//   - Revoke: Clean up credentials on logout or explicit revocation
//
// Usage:
//
//	provider.Register(myProvider)           // Register at init time
//	p, _ := provider.Get("my-provider")     // Retrieve by name
//	p.Initialize(ctx, config)               // Configure the provider
//	creds, _ := p.Authenticate(ctx, opts)   // Authenticate
//	creds, _ = p.Refresh(ctx, creds)        // Refresh before expiry
//	p.Revoke(ctx, creds)                    // Revoke on logout
package provider

import (
	"context"
	"time"
)

// Type represents the authentication protocol used by a provider.
// This helps consumers understand the underlying mechanism and handle
// provider-specific behaviors appropriately.
type Type string

const (
	// TypeOIDC indicates an OpenID Connect provider. OIDC providers return
	// ID tokens in addition to access tokens, enabling identity verification.
	TypeOIDC Type = "oidc"

	// TypeSAML indicates a SAML 2.0 provider. SAML providers exchange XML
	// assertions and are commonly used in enterprise SSO scenarios.
	TypeSAML Type = "saml"

	// TypeOAuth indicates an OAuth 2.0 provider. OAuth providers issue access
	// tokens for API authorization without identity claims.
	TypeOAuth Type = "oauth2"
)

// Provider defines the interface that all authentication providers must
// implement. This interface replaces dlogin.ILogin and provides a complete
// credential lifecycle: initialization, authentication, refresh, and revocation.
//
// Implementations should be safe for concurrent use after Initialize is called.
// Providers register themselves via provider.Register, typically in an init func.
//
// Lifecycle methods should be called in order:
//  1. Initialize (once per provider instance)
//  2. Authenticate (as needed)
//  3. Refresh (proactively, before credentials expire)
//  4. Revoke (on logout or explicit revocation request)
type Provider interface {
	// Name returns the unique identifier for this provider. This name is used
	// for registration and lookup via provider.Get. It should be lowercase,
	// hyphen-separated (e.g., "aws-oidc", "okta-oidc").
	Name() string

	// Type returns the authentication protocol this provider implements.
	// Callers may use this to apply protocol-specific handling.
	Type() Type

	// Initialize configures the provider with static configuration that remains
	// constant across authentication attempts. This includes provider endpoints,
	// client IDs, secrets, and other setup parameters.
	//
	// Initialize must be called before Authenticate. The config map keys should
	// match the fields defined in GetConfigSchema. Invalid configuration should
	// cause Initialize to return an error.
	//
	// This method should be idempotent - calling it multiple times with the same
	// config should have the same effect as calling it once.
	Initialize(ctx context.Context, config map[string]any) error

	// Authenticate performs the authentication flow and returns credentials.
	// The AuthOptions contain per-authentication parameters that may vary
	// between calls (profile, region, MFA codes via Extra, etc.).
	//
	// Implementations should respect context cancellation for interactive flows
	// (e.g., browser-based authentication).
	//
	// Returns an error if authentication fails or is cancelled.
	Authenticate(ctx context.Context, opts AuthOptions) (*Credentials, error)

	// Refresh obtains new credentials using an existing credential's refresh
	// token or equivalent mechanism. Callers should invoke Refresh proactively
	// before the credential's Expiry time, not after credentials have expired.
	//
	// If the provider does not support refresh (e.g., SAML assertions), or if
	// the credential cannot be refreshed, return an error. The caller should
	// then fall back to a full Authenticate flow.
	//
	// The returned Credentials may have a new RefreshToken; callers should
	// persist the updated credentials.
	Refresh(ctx context.Context, creds *Credentials) (*Credentials, error)

	// Revoke invalidates the given credentials. This serves dual purposes:
	//   - Logout cleanup: Called during logout to clean up local state and
	//     optionally revoke tokens with the identity provider
	//   - Explicit revocation: Called when the user explicitly requests
	//     token revocation (e.g., OAuth token revocation endpoint)
	//
	// Implementations should handle both scenarios gracefully. If the provider
	// does not support server-side revocation, it should still return nil
	// (cleanup-only behavior is acceptable).
	//
	// Revoke should be idempotent - revoking already-revoked credentials
	// should not return an error.
	Revoke(ctx context.Context, creds *Credentials) error

	// GetConfigSchema returns the schema describing the configuration fields
	// this provider accepts in Initialize. This enables runtime validation,
	// CLI flag generation, and documentation generation.
	GetConfigSchema() ConfigSchema

	// ValidateConfig checks whether the given configuration is valid for this
	// provider without actually initializing. Use this for early validation
	// (e.g., config file parsing) before calling Initialize.
	//
	// Returns nil if the config is valid, or an error describing what is wrong.
	ValidateConfig(config map[string]any) error
}

// Credentials represents the authentication tokens and metadata returned by
// a successful Authenticate or Refresh call. Callers should persist credentials
// and use the Expiry field to determine when to proactively call Refresh.
type Credentials struct {
	// Type indicates the credential format, typically matching the provider's
	// Type (e.g., "oidc", "saml", "oauth2"). This helps callers handle
	// credentials appropriately.
	Type string

	// AccessToken is the primary token used to access protected resources.
	// For OIDC/OAuth, this is the access token. For SAML, this may be the
	// assertion or a derived token.
	AccessToken string

	// RefreshToken is used to obtain new credentials without re-authentication.
	// May be empty if the provider does not support refresh (e.g., SAML).
	// Callers should persist this and pass it back via Refresh.
	RefreshToken string

	// Expiry indicates when the AccessToken expires. Callers should call
	// Refresh proactively before this time. A zero value indicates the
	// credential does not expire (uncommon).
	Expiry time.Time

	// Metadata contains provider-specific data that doesn't fit the standard
	// fields. Examples: ID tokens (OIDC), session IDs, role ARNs (AWS).
	// The keys and structure are provider-defined.
	Metadata map[string]any
}

// AuthOptions contains per-authentication parameters passed to Authenticate.
// Unlike the config in Initialize (which is static provider setup), these
// options may vary between authentication attempts.
type AuthOptions struct {
	// Profile specifies the named profile to use for this authentication.
	// The interpretation is provider-specific (e.g., AWS profile name).
	Profile string

	// Region specifies the geographic region for the authentication.
	// The interpretation is provider-specific (e.g., AWS region).
	Region string

	// Private indicates whether to use private/incognito browser windows
	// for interactive authentication flows. Useful for multi-account scenarios.
	Private bool

	// SkipDefaults indicates whether to skip loading default configuration
	// values. When true, only explicitly provided options are used.
	SkipDefaults bool

	// Extra contains provider-specific options that don't fit the standard
	// fields. Examples: MFA codes, session duration preferences, specific
	// scopes to request. Keys and values are provider-defined.
	Extra map[string]any
}

// ConfigSchema describes the configuration fields a provider accepts in
// Initialize. This enables tooling to validate configuration, generate CLI
// flags, and produce documentation without provider-specific knowledge.
type ConfigSchema struct {
	// Fields lists all configuration fields the provider accepts.
	Fields []ConfigField
}

// ConfigField describes a single configuration field for a provider.
type ConfigField struct {
	// Name is the configuration key (e.g., "sso_start_url", "client_id").
	// Should be snake_case for consistency with config files.
	Name string

	// Type indicates the expected value type. Supported types:
	//   - "string": text value
	//   - "bool": true/false
	//   - "int": integer number
	//   - "duration": time duration (e.g., "1h30m")
	Type string

	// Required indicates whether this field must be provided. If true and
	// the field is missing, ValidateConfig and Initialize should return errors.
	Required bool

	// Default is the value used when the field is not provided. Only
	// meaningful when Required is false. The type should match Type.
	Default any

	// Description explains what this field configures. Used for help text
	// and documentation generation. Should be a complete sentence.
	Description string
}
