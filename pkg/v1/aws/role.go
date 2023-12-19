package aws

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sso"
	"github.com/aws/aws-sdk-go-v2/service/sso/types"
	. "github.com/louislef299/aws-sso/internal/envs"
	los "github.com/louislef299/aws-sso/pkg/v1/os"
	"github.com/louislef299/aws-sso/pkg/v1/prompt"
	"github.com/spf13/viper"
)

func contains(l []types.RoleInfo, s string) int {
	for i, v := range l {
		if strings.Compare(*v.RoleName, s) == 0 {
			return i
		}
	}
	return -1
}

func GetRoleCredentials(ctx context.Context, cfg *aws.Config, accountID, roleName, accessToken *string) (*sso.GetRoleCredentialsOutput, error) {
	client := sso.NewFromConfig(*cfg)
	return client.GetRoleCredentials(ctx, &sso.GetRoleCredentialsInput{
		AccountId:   accountID,
		RoleName:    roleName,
		AccessToken: accessToken,
	})
}

func GetAndSaveRoleCredentials(ctx context.Context, cfg *aws.Config, accountID, roleName, accessToken *string, accountName, region string) (string, error) {
	roleCreds, err := GetRoleCredentials(ctx, cfg, accountID, roleName, accessToken)
	if err != nil {
		return "", err
	}
	return saveCredentials(accountName, region, "json", roleCreds)
}

func RetrieveRoleInfo(ctx context.Context, cfg *aws.Config, accountID, accessToken *string) (types.RoleInfo, error) {
	client := sso.NewFromConfig(*cfg)
	roles, err := client.ListAccountRoles(ctx, &sso.ListAccountRolesInput{
		AccountId:   accountID,
		AccessToken: accessToken,
	})
	if err != nil {
		return types.RoleInfo{}, fmt.Errorf("couldn't gather account roles: %v", err)
	}

	r := getConfiguredRole()
	if i := contains(roles.RoleList, r); i >= 0 {
		log.Printf("found role in configuration: %s", r)
		return roles.RoleList[i], nil
	} else if len(roles.RoleList) == 1 {
		log.Printf("only one role available, using role: %s\n", *roles.RoleList[0].RoleName)
		return roles.RoleList[0], nil
	} else {
		log.Println("HINT: if you would like to reuse a specific iam profile, you can set core.defaultRole to your iam profile.")
	}

	var rolesToSelect []string
	for i, info := range roles.RoleList {
		rolesToSelect = append(rolesToSelect, prompt.LINEPREFIX+strconv.Itoa(i)+" "+*info.RoleName)
	}

	label := "Select your role - To choose one role directly just enter #{Int}"
	indexChoice, _ := prompt.Select(label, rolesToSelect, prompt.FuzzySearchWithPrefixAnchor(rolesToSelect))
	roleInfo := roles.RoleList[indexChoice]
	return roleInfo, nil
}

func getConfiguredRole() string {
	r := viper.GetString(SESSION_ROLE)
	if r != "" {
		return r
	}
	return viper.GetString(CORE_DEFAULT_ROLE)
}

func saveCredentials(profile, region, output string, roleCredentials *sso.GetRoleCredentialsOutput) (string, error) {
	// this is where the write to /.aws/credentials happens, going to want to modify this
	custom_profile := los.GetProfile(profile)
	if err := WriteAWSCredentialsFile(custom_profile, roleCredentials); err != nil {
		return "", err
	}
	if err := WriteAWSConfigFile(custom_profile, region, output); err != nil {
		return "", err
	}

	fmt.Printf("If you would like to use these creds with the aws cli, please copy and paste the following command:\n")
	switch runtime.GOOS {
	case "linux", "darwin":
		fmt.Printf("\texport AWS_PROFILE=%s\n", custom_profile)
	case "windows":
		fmt.Printf("\t$env:AWS_PROFILE=%s\n", custom_profile)
	default:
		return "", fmt.Errorf("os not supported")
	}
	viper.Set(CORE_REGION, region)
	return custom_profile, viper.WriteConfig()
}
