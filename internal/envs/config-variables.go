/*
Copyright © 2023 Louis Lefebvre <louislefebvre1999@gmail.com>
*/
package envs

// These variables are meant to represent config variables used by Viper or used
// across packages
const (
	// Core variables
	CORE_BROWSER           = "core.browser"
	CORE_DEFAULT_CLUSTER   = "core.defaultCluster"
	CORE_DEFAULT_ROLE      = "core.defaultRole"
	CORE_DEFAULT_REGION    = "core.defaultRegion"
	CORE_SSO_REGION        = "core.ssoRegion"
	CORE_URL               = "core.url"
	CORE_PLUGINS           = "core.plugins"
	CORE_DISABLE_EKS_LOGIN = "core.disableEKSLogin"

	// Session variables
	SESSION_HEADER      = "session"
	SESSION_PROFILE     = "session.profile"
	SESSION_CLUSTER     = "session.cluster"
	SESSION_REGION      = "session.region"
	SESSION_URL         = "session.url"
	SESSION_ROLE        = "session.role"
	SESSION_TOKEN       = "session.token"
	SESSION_LAST_VCHECK = "session.lastVersionCheck"

	// Token variables
	TOKEN_HEADER       = "token"
	DEFAULT_TOKEN_NAME = "default"
	TOKEN_LOCK         = "token.lock"

	// Profile variables
	PROFILE_HEADER = "profile"
)
