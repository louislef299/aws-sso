package aws

import (
	"context"
	"errors"
	"fmt"
	"log"
	"slices"
	"strings"

	awsConf "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sso"
	ssotypes "github.com/aws/aws-sdk-go-v2/service/sso/types"
	"github.com/louislef299/knot/internal/browser"
	"github.com/louislef299/knot/internal/envs"
	"github.com/louislef299/knot/internal/region"
	laws "github.com/louislef299/knot/pkg/aws"
	"github.com/louislef299/knot/pkg/provider"
	"github.com/spf13/viper"
)

type AWS struct {
	startURL  string
	ssoRegion string
}

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

func (p *AWS) Initialize(ctx context.Context, config map[string]any) (err error) {
	p.startURL, err = provider.ConfigGet[string](config, SSO_START_URL)
	if err != nil {
		return
	}
	p.ssoRegion, err = provider.ConfigGet[string](config, SSO_REGION)
	return
}

func (p *AWS) Authenticate(ctx context.Context,
	opts provider.AuthOptions) (*provider.Credentials, error) {

	// Determine the SSO region to use for authentication
	ssoRegion := p.ssoRegion
	// Special handling for China regions
	if opts.Region == "cn-north-1" || opts.Region == "cn-northwest-1" {
		ssoRegion = "cn-north-1"
	}
	log.Printf("using sso region %s to login", ssoRegion)

	// Load AWS config for SSO operations
	awsCfg, err := awsConf.LoadDefaultConfig(ctx, awsConf.WithRegion(ssoRegion))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Determine if we should refresh the token
	refresh := false
	if v, ok := opts.Extra["refresh"].(bool); ok {
		refresh = v
	}

	// Get browser for device authorization
	browserType := viper.GetString(envs.CORE_BROWSER)
	b := browser.GetBrowser(browserType, opts.Private)
	if opts.Private {
		log.Println("browser set to open incognito (no cookies)")
	} else {
		log.Println("browser set to default (use cookies)")
	}

	// Perform OIDC device flow and get access token
	clientInfo, err := laws.GatherClientInformation(ctx, &awsCfg, p.startURL, b, refresh)
	if err != nil {
		return nil, fmt.Errorf("failed to gather client information: %w", err)
	}

	// Get account ID from options if provided
	accountID := ""
	if v, ok := opts.Extra["account_id"].(string); ok {
		accountID = v
	}

	// Retrieve account information (with user selection if accountID is empty)
	account, err := laws.RetrieveAccountInformation(ctx, &awsCfg, &clientInfo.AccessToken, &accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account information: %w", err)
	}
	accountID = *account.AccountId
	log.Printf("using AWS account: %s (%s)", *account.AccountName, accountID)

	// Retrieve role information (with user selection if not in skipDefaults
	// mode)
	role, err := laws.RetrieveRoleInfo(ctx, &awsCfg, accountID, clientInfo.AccessToken, opts.SkipDefaults)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve role information: %w", err)
	}
	log.Printf("using AWS role: %s", *role.RoleName)

	// Save usage information for defaults
	if err := laws.SaveUsageInformation(account, &role); err != nil {
		log.Printf("warning: couldn't save usage information: %v", err)
	}

	// Get role credentials (the actual AWS access key/secret/token)
	roleCreds, err := laws.GetRoleCredentials(ctx, &awsCfg, account.AccountId, role.RoleName, &clientInfo.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get role credentials: %w", err)
	}

	// Build and return credentials with all necessary information
	return &provider.Credentials{
		Type:         string(provider.TypeOIDC),
		AccessToken:  clientInfo.AccessToken,
		RefreshToken: "", // AWS SSO uses device code flow, not refresh tokens
		Expiry:       clientInfo.AccessTokenExpiresAt,
		Metadata: map[string]any{
			// OIDC client information (for future device flow renewals)
			"client_id":                 clientInfo.ClientId,
			"client_secret":             clientInfo.ClientSecret,
			"client_secret_expires_at":  clientInfo.ClientSecretExpiresAt,
			"device_code":               clientInfo.DeviceCode,
			"verification_uri_complete": clientInfo.VerificationUriComplete,
			"start_url":                 clientInfo.StartUrl,

			// Account and role information
			"account_id":   accountID,
			"account_name": *account.AccountName,
			"role_name":    *role.RoleName,

			// AWS role credentials (for native CLI integration)
			"aws_access_key_id":      *roleCreds.RoleCredentials.AccessKeyId,
			"aws_secret_access_key":  *roleCreds.RoleCredentials.SecretAccessKey,
			"aws_session_token":      *roleCreds.RoleCredentials.SessionToken,
			"aws_credentials_expiry": roleCreds.RoleCredentials.Expiration,

			// Additional context
			"region": opts.Region,
		},
	}, nil
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

// WriteCredentials implements provider.NativeCLIIntegration by writing AWS
// credentials to ~/.aws/credentials and ~/.aws/config in the format expected by
// the AWS CLI.
func (p *AWS) WriteNativeCredentials(ctx context.Context, creds *provider.Credentials, profile string) error {
	// Extract AWS credentials from metadata
	accessKeyID, ok := creds.Metadata["aws_access_key_id"].(string)
	if !ok {
		return fmt.Errorf("missing or invalid aws_access_key_id in credentials")
	}
	secretAccessKey, ok := creds.Metadata["aws_secret_access_key"].(string)
	if !ok {
		return fmt.Errorf("missing or invalid aws_secret_access_key in credentials")
	}
	sessionToken, ok := creds.Metadata["aws_session_token"].(string)
	if !ok {
		return fmt.Errorf("missing or invalid aws_session_token in credentials")
	}

	// Get region from metadata (fallback to configured ssoRegion if not
	// present)
	credRegion, ok := creds.Metadata["region"].(string)
	if !ok || credRegion == "" {
		credRegion = p.ssoRegion
	}

	// Build an sso.GetRoleCredentialsOutput with the credentials This allows us
	// to reuse the existing WriteAWSCredentialsFile function
	roleCredsOutput := &sso.GetRoleCredentialsOutput{
		RoleCredentials: &ssotypes.RoleCredentials{
			AccessKeyId:     &accessKeyID,
			SecretAccessKey: &secretAccessKey,
			SessionToken:    &sessionToken,
		},
	}

	// Write to ~/.aws/credentials
	if err := laws.WriteAWSCredentialsFile(profile, roleCredsOutput); err != nil {
		return fmt.Errorf("failed to write AWS credentials file: %w", err)
	}

	// Write to ~/.aws/config (with output format defaulting to "json")
	if err := laws.WriteAWSConfigFile(profile, credRegion, "json"); err != nil {
		return fmt.Errorf("failed to write AWS config file: %w", err)
	}

	log.Printf("successfully wrote AWS credentials for profile: %s", profile)
	return nil
}

// CleanCredentials implements provider.NativeCLIIntegration by removing AWS
// credentials from ~/.aws/credentials and ~/.aws/config.
func (p *AWS) CleanCredentials(ctx context.Context, profile string) error {
	// The existing clean logic in plugins/aws/oidc/oidc.go handles this For
	// now, return nil as cleanup will be handled separately TODO: Extract the
	// clean logic from oidc.go into a shared function
	log.Printf("cleaning native credentials for profile: %s", profile)
	return nil
}
