package aws

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sso"
	"github.com/aws/aws-sdk-go-v2/service/ssooidc"
	"github.com/aws/aws-sdk-go-v2/service/ssooidc/types"
	"github.com/briandowns/spinner"
	browser "github.com/louislef299/aws-sso/internal/browser"
	los "github.com/louislef299/aws-sso/pkg/v1/os"
)

var (
	ErrMoreThanOneLocation   = errors.New("too many configuration locations provided")
	ErrStartURLCannotBeEmpty = errors.New("start URL cannot be empty")
)

const (
	grantType  = "urn:ietf:params:oauth:grant-type:device_code"
	clientType = "public"
	clientName = "go-aws-sso"
)

// ClientInfoFileDestination finds local AWS configuration settings. Users can
// optionally input their own home directory location.
func ClientInfoFileDestination(configDir ...string) (string, error) {
	var configLocation string
	if len(configDir) == 0 {
		h, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		configLocation = path.Join(h, AWS_TOKEN_PATH, GetAccessToken())
	} else if len(configDir) > 1 {
		return "", ErrMoreThanOneLocation
	} else {
		cache := path.Join(configDir[0], AWS_TOKEN_PATH)
		err := os.MkdirAll(cache, 0755)
		if err != nil {
			return "", err
		}
		configLocation = path.Join(configDir[0], AWS_TOKEN_PATH, GetAccessToken())
	}

	exists, err := los.IsFileOrFolderExisting(configLocation)
	if err != nil {
		return "", err
	}
	if !exists {
		f, err := os.Create(configLocation)
		if err != nil {
			return "", err
		}
		defer f.Close()
	}
	return configLocation, nil
}

// Attempts to gather current client information. If it doesn't exist, creates new information for the client
func GatherClientInformation(ctx context.Context, cfg *aws.Config, startUrl string, b browser.Browser, refresh bool) (*ClientInformation, error) {
	if startUrl == "" {
		return nil, ErrStartURLCannotBeEmpty
	}
	log.Println("gathering client info")
	infoDest, err := ClientInfoFileDestination()
	if err != nil {
		return nil, err
	}

	clientInfo, err := ReadClientInformation(infoDest)
	if err != nil || clientInfo.StartUrl != startUrl || refresh {
		log.Println("registering client")
		clientInfo, err = RegisterClient(ctx, cfg, startUrl, b)
		if err != nil {
			return nil, err
		}
	} else if clientInfo.IsExpired() {
		log.Println("AccessToken expired. Start retrieving a new AccessToken")
		resp, err := StartDeviceAuthorization(ctx, cfg, startUrl, &ssooidc.RegisterClientOutput{
			ClientId:     &clientInfo.ClientId,
			ClientSecret: &clientInfo.ClientSecret,
		})
		if err != nil {
			return nil, err
		}

		log.Println("user validation code:", *resp.UserCode)
		log.Println("please verify your client request: " + *resp.VerificationUriComplete)
		err = b.OpenURL(ctx, *resp.VerificationUriComplete)
		if err != nil {
			return nil, fmt.Errorf("issue opening url for browser type %s: %v", b.Type(), err)
		}

		clientInfo.DeviceCode = *resp.DeviceCode
	} else {
		return clientInfo, nil
	}
	if err := RetrieveToken(ctx, cfg, clientInfo); err != nil {
		return nil, err
	}
	if err := los.WriteStructToFile(clientInfo, infoDest); err != nil {
		return nil, err
	}
	return clientInfo, nil
}

// Registers a client with AWS OIDC and return the client information
func RegisterClient(ctx context.Context, cfg *aws.Config, startUrl string, b browser.Browser) (*ClientInformation, error) {
	oidc := ssooidc.NewFromConfig(*cfg)
	resp, err := oidc.RegisterClient(ctx, &ssooidc.RegisterClientInput{
		ClientName: aws.String(clientName),
		ClientType: aws.String(clientType),
	})
	if err != nil {
		return nil, err
	}

	output, err := StartDeviceAuthorization(ctx, cfg, startUrl, resp)
	if err != nil {
		return nil, err
	}

	log.Println("user validation code:", *output.UserCode)
	if err = b.OpenURL(ctx, *output.VerificationUriComplete); err != nil {
		return nil, fmt.Errorf("issue opening url for browser type %s: %v", b.Type(), err)
	}

	return &ClientInformation{
		ClientId:                *resp.ClientId,
		ClientSecret:            *resp.ClientSecret,
		ClientSecretExpiresAt:   strconv.FormatInt(resp.ClientSecretExpiresAt, 10),
		DeviceCode:              *output.DeviceCode,
		VerificationUriComplete: *output.VerificationUriComplete,
		StartUrl:                startUrl,
	}, nil
}

func Logout(ctx context.Context, cfg *aws.Config, accessToken string) error {
	client := sso.NewFromConfig(*cfg)
	if _, err := client.Logout(ctx, &sso.LogoutInput{AccessToken: aws.String(accessToken)}); err != nil {
		return err
	}
	return nil
}

func RetrieveToken(ctx context.Context, cfg *aws.Config, clientInfo *ClientInformation) error {
	oidc := ssooidc.NewFromConfig(*cfg)
	s := spinner.New(spinner.CharSets[6], 300*time.Millisecond, spinner.WithWriter(log.Default().Writer()))
	s.Prefix = "Waiting for authorization "
	s.FinalMSG = "Successfully authorized!"
	s.Start()
	defer func() {
		s.Stop()
		fmt.Println()
	}()
	for {
		cto, err := oidc.CreateToken(ctx, &ssooidc.CreateTokenInput{
			ClientId:     &clientInfo.ClientId,
			ClientSecret: &clientInfo.ClientSecret,
			DeviceCode:   &clientInfo.DeviceCode,
			GrantType:    aws.String(grantType),
		})
		if err != nil {
			var authPending *types.AuthorizationPendingException
			if errors.As(err, &authPending) {
				time.Sleep(3 * time.Second)
				continue
			} else {
				return err
			}
		} else {
			clientInfo.AccessToken = *cto.AccessToken
			clientInfo.AccessTokenExpiresAt = time.Now().Add(time.Hour*8 - time.Minute*5)
			return nil
		}
	}
}

func StartDeviceAuthorization(ctx context.Context, cfg *aws.Config, startUrl string, rco *ssooidc.RegisterClientOutput) (*ssooidc.StartDeviceAuthorizationOutput, error) {
	oidc := ssooidc.NewFromConfig(*cfg)
	return oidc.StartDeviceAuthorization(ctx, &ssooidc.StartDeviceAuthorizationInput{
		ClientId:     rco.ClientId,
		ClientSecret: rco.ClientSecret,
		StartUrl:     aws.String(startUrl),
	})
}
