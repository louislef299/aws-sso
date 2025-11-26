package aws

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
	"time"

	awsConf "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sso"
	ssotypes "github.com/aws/aws-sdk-go-v2/service/sso/types"
	"github.com/louislef299/knot/internal/browser"
	"github.com/louislef299/knot/internal/envs"
	"github.com/louislef299/knot/internal/region"
	laws "github.com/louislef299/knot/pkg/aws"
	"github.com/louislef299/knot/pkg/provider"
	"github.com/spf13/viper"
	"gopkg.in/ini.v1"
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
	creds *provider.Credentials, opts provider.AuthOptions) (*provider.Credentials, error) {

	// Check if the access token is still valid
	if creds.Expiry.After(time.Now()) {
		log.Printf("AWS SSO token still valid until %s, no refresh needed", creds.Expiry.Format(time.RFC3339))
		return creds, nil
	}

	log.Println("AWS SSO token expired, performing re-authentication")

	// AWS SSO uses OIDC device code flow which doesn't support traditional
	// refresh tokens. When the token expires, we need to perform a full
	// re-authentication with user approval. We'll extract context from the old
	// credentials to minimize user prompts.

	// Extract account_id from credentials to skip account selection prompt
	if accountID, ok := creds.Metadata["account_id"].(string); ok && accountID != "" {
		if opts.Extra == nil {
			opts.Extra = make(map[string]any)
		}
		opts.Extra["account_id"] = accountID
		log.Printf("reusing account_id: %s", accountID)
	}

	// Set refresh flag to force new device authorization
	if opts.Extra == nil {
		opts.Extra = make(map[string]any)
	}
	opts.Extra["refresh"] = true

	// Extract region from credentials if not already specified in opts
	if opts.Region == "" {
		if region, ok := creds.Metadata["region"].(string); ok && region != "" {
			opts.Region = region
			log.Printf("reusing region: %s", region)
		}
	}

	// Perform full re-authentication using the Authenticate flow This will
	// start a new device authorization and require user approval
	newCreds, err := p.Authenticate(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("refresh re-authentication failed: %w", err)
	}

	log.Println("successfully refreshed AWS SSO credentials")
	return newCreds, nil
}

func (p *AWS) Revoke(ctx context.Context, creds *provider.Credentials) error {
	log.Println("revoking AWS SSO credentials")

	// Extract access token from credentials
	accessToken := creds.AccessToken
	if accessToken == "" {
		return fmt.Errorf("no access token available in credentials")
	}

	// Determine SSO region for API calls
	ssoRegion := p.ssoRegion
	if region, ok := creds.Metadata["region"].(string); ok && region != "" {
		// Check for China regions
		if region == "cn-north-1" || region == "cn-northwest-1" {
			ssoRegion = "cn-north-1"
		}
	}

	// Load AWS config for SSO operations
	awsCfg, err := awsConf.LoadDefaultConfig(ctx, awsConf.WithRegion(ssoRegion))
	if err != nil {
		return fmt.Errorf("failed to load AWS config for revocation: %w", err)
	}

	// Call AWS SSO Logout API to revoke the token server-side
	if err := laws.Logout(ctx, &awsCfg, accessToken); err != nil {
		return fmt.Errorf("failed to revoke AWS SSO token: %w", err)
	}

	// Clean up cached client information file This file contains the OIDC
	// client credentials and tokens
	clientInfoPath, err := laws.ClientInfoFileDestination()
	if err != nil {
		return fmt.Errorf("failed to determine client info file path: %w", err)
	}

	// Remove the cache file (idempotent - no error if doesn't exist)
	if err := os.Remove(clientInfoPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove cached client info: %w", err)
	}

	log.Println("successfully revoked AWS SSO credentials")
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
func (p *AWS) WriteCredentials(ctx context.Context, creds *provider.Credentials, profile string) error {
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
// credentials from ~/.aws/credentials and ~/.aws/config for the specified
// profile.
func (p *AWS) CleanCredentials(ctx context.Context, profile string) error {
	log.Printf("cleaning native credentials for profile: %s", profile)

	// Get home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	// Clean ~/.aws/credentials
	credentialsFile := homeDir + "/.aws/credentials"
	if err := cleanProfileFromFile(credentialsFile, profile); err != nil {
		return fmt.Errorf("failed to clean credentials file: %w", err)
	}

	// Clean ~/.aws/config In config files, non-default profiles are prefixed
	// with "profile "
	configFile := homeDir + "/.aws/config"
	configProfileName := profile
	if profile != "default" {
		configProfileName = "profile " + profile
	}
	if err := cleanProfileFromFile(configFile, configProfileName); err != nil {
		return fmt.Errorf("failed to clean config file: %w", err)
	}

	log.Printf("successfully cleaned native credentials for profile: %s", profile)
	return nil
}

// cleanProfileFromFile removes a specific profile section from an INI file.
// This is idempotent - returns nil if the file or section doesn't exist.
func cleanProfileFromFile(filePath, sectionName string) error {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// File doesn't exist, nothing to clean (idempotent)
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to stat file: %w", err)
	}

	// Load the INI file
	cfg, err := ini.Load(filePath)
	if err != nil {
		return fmt.Errorf("failed to load file: %w", err)
	}

	// Delete the section if it exists
	if cfg.HasSection(sectionName) {
		cfg.DeleteSection(sectionName)
		log.Printf("removed section '%s' from %s", sectionName, filePath)
	}

	// Save the file
	if err := cfg.SaveTo(filePath); err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	return nil
}
