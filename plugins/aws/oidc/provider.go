package oidc

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ssooidc"
	"github.com/louislef299/knot/pkg/provider"
)

func init() {
	provider.Register(&AWSProvider{})
}

type AWSProvider struct {
	startURL  string
	ssoRegion string
	client    *ssooidc.Client
}

func (p *AWSProvider) Name() string        { return "aws" }
func (p *AWSProvider) Type() provider.Type { return provider.TypeOIDC }

func (p *AWSProvider) Initialize(ctx context.Context, config map[string]any) error {
	p.startURL = config["url"].(string)
	p.ssoRegion = config["sso_region"].(string)

	cfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(p.ssoRegion))
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %w", err)
	}
	p.client = ssooidc.NewFromConfig(cfg)
	return nil
}

func (p *AWSProvider) Authenticate(ctx context.Context, opts provider.AuthOptions) (*provider.Credentials, error) {
	// 1. Register OIDC client (or load from cache)
	clientInfo, err := p.getOrRegisterClient(ctx)
	if err != nil {
		return nil, err
	}

	// 2. Start device authorization
	deviceAuth, err := p.client.StartDeviceAuthorization(ctx,
		&ssooidc.StartDeviceAuthorizationInput{
			ClientId:     &clientInfo.ClientId,
			ClientSecret: &clientInfo.ClientSecret,
			StartUrl:     &p.startURL,
		})
	if err != nil {
		return nil, err
	}

	// 3. Open browser, poll for token
	token, err := p.pollForToken(ctx, clientInfo, deviceAuth)
	if err != nil {
		return nil, err
	}

	// 4. Get role credentials
	roleCreds, err := p.getRoleCredentials(ctx, token, opts)
	if err != nil {
		return nil, err
	}

	return &provider.Credentials{
		AccessToken:  *roleCreds.AccessKeyId,
		RefreshToken: *token.RefreshToken,
		ExpiresAt:    time.Unix(roleCreds.Expiration, 0),
		Metadata: map[string]any{
			"secret_access_key": *roleCreds.SecretAccessKey,
			"session_token":     *roleCreds.SessionToken,
			"account_id":        opts.Extra["account_id"],
			"role_name":         opts.Extra["role_name"],
		},
	}, nil
}

func (p *AWSProvider) GetConfigSchema() provider.ConfigSchema {
	return provider.ConfigSchema{
		Fields: []provider.ConfigField{
			{
				Name: "url", Type: "string", Required: true, Description: "AWS SSO start URL",
			},
			{
				Name: "sso_region", Type: "string", Required: true, Default: "us-east-1",
			},
			{
				Name: "default_region", Type: "string", Required: false, Default: "us-east-1",
			},
		},
	}
}

func (p *AWSProvider) ValidateConfig(config map[string]any) error {
	if _, ok := config["url"]; !ok {
		return errors.New("url is required")
	}
	return nil
}
