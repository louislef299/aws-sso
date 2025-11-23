package aws

import (
	"context"
	"errors"
	"slices"
	"strings"

	"github.com/louislef299/knot/internal/region"
	"github.com/louislef299/knot/pkg/provider"
)

type AWS struct{}

const (
	SSO_START_URL = "sso_start_url"
	SSO_REGION    = "sso_region"
	DEFAULT_ROLE  = "default_role"
)

var (
	ErrConfigStartUrl = errors.New("sso_start_url must be an HTTPS URL")
	ErrInvalidRegion  = errors.New("the provided sso_region is invalid")
)

func (p *AWS) Name() string {
	return "aws"
}

func (p *AWS) Type() provider.Type {
	return provider.TypeOIDC
}

// Initialize stores the provider configuration (passed via Activate) so that
// subsequent Authenticate calls can use it without needing config passed again.
func (p *AWS) Initialize(ctx context.Context,
	config map[string]any) error {
	return nil
}

func (p *AWS) Authenticate(ctx context.Context,
	opts provider.AuthOptions) (*provider.Credentials, error) {
	return nil, nil
}

func (p *AWS) Refresh(ctx context.Context,
	creds *provider.Credentials) (*provider.Credentials, error) {
	return nil, nil
}

func (p *AWS) Revoke(ctx context.Context, creds *provider.Credentials) error {
	return nil
}

func (p *AWS) GetConfigSchema() provider.ConfigSchema {
	return provider.ConfigSchema{
		Fields: []provider.ConfigField{
			{
				Name:        SSO_START_URL,
				Type:        "string",
				Required:    true,
				Description: "The AWS SSO start URL (e.g., https://my-sso.awsapps.com/start)",
			},
			{
				Name:        SSO_REGION,
				Type:        "string",
				Required:    true,
				Description: "The AWS region where SSO is configured (e.g., us-east-1)",
			},
			{
				Name:        DEFAULT_ROLE,
				Type:        "string",
				Required:    false,
				Description: "Default IAM role to assume if not specified at login time",
			},
		},
	}
}

func (p *AWS) ValidateConfig(config map[string]any) error {
	// Validate Start URL
	startURL, err := provider.ConfigGet[string](config, SSO_START_URL)
	if err != nil {
		return err
	}
	if !strings.HasPrefix(startURL, "https://") {
		return ErrConfigStartUrl
	}

	// Validate Region
	r, err := provider.ConfigGet[string](config, SSO_REGION)
	if err != nil {
		return err
	}
	validRegion := slices.Contains(region.AwsRegions, r)
	if !validRegion {
		return ErrInvalidRegion
	}

	// TODO: Validate the optional default_role
	return nil
}
