package provider

import (
	"context"
	"time"
)

type Type string

const (
	TypeOIDC  Type = "oidc"
	TypeSAML  Type = "saml"
	TypeOAuth Type = "oauth2"
)

type Provider interface {
	Name() string
	Type() Type

	// Lifecycle
	Initialize(ctx context.Context, config map[string]any) error
	Authenticate(ctx context.Context, opts AuthOptions) (*Credentials, error)
	Refresh(ctx context.Context, creds *Credentials) (*Credentials, error)
	Revoke(ctx context.Context, creds *Credentials) error

	// Metadata
	GetConfigSchema() ConfigSchema
	ValidateConfig(config map[string]any) error
}

type Credentials struct {
	Type         string
	AccessToken  string
	RefreshToken string
	Expiry       time.Time
	Metadata     map[string]any // Provider-specific data
}

type AuthOptions struct {
	Profile      string
	Region       string
	Private      bool
	SkipDefaults bool
	Extra        map[string]any
}

type ConfigSchema struct {
	Fields []ConfigField
}

type ConfigField struct {
	Name        string
	Type        string // "string", "bool", "int", "duration"
	Required    bool
	Default     any
	Description string
}
