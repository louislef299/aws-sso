// Code generated by smithy-go-codegen DO NOT EDIT.

package types

type AccessScopeType string

// Enum values for AccessScopeType
const (
	AccessScopeTypeCluster   AccessScopeType = "cluster"
	AccessScopeTypeNamespace AccessScopeType = "namespace"
)

// Values returns all known values for AccessScopeType. Note that this can be
// expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (AccessScopeType) Values() []AccessScopeType {
	return []AccessScopeType{
		"cluster",
		"namespace",
	}
}

type AddonIssueCode string

// Enum values for AddonIssueCode
const (
	AddonIssueCodeAccessDenied                 AddonIssueCode = "AccessDenied"
	AddonIssueCodeInternalFailure              AddonIssueCode = "InternalFailure"
	AddonIssueCodeClusterUnreachable           AddonIssueCode = "ClusterUnreachable"
	AddonIssueCodeInsufficientNumberOfReplicas AddonIssueCode = "InsufficientNumberOfReplicas"
	AddonIssueCodeConfigurationConflict        AddonIssueCode = "ConfigurationConflict"
	AddonIssueCodeAdmissionRequestDenied       AddonIssueCode = "AdmissionRequestDenied"
	AddonIssueCodeUnsupportedAddonModification AddonIssueCode = "UnsupportedAddonModification"
	AddonIssueCodeK8sResourceNotFound          AddonIssueCode = "K8sResourceNotFound"
	AddonIssueCodeAddonSubscriptionNeeded      AddonIssueCode = "AddonSubscriptionNeeded"
	AddonIssueCodeAddonPermissionFailure       AddonIssueCode = "AddonPermissionFailure"
)

// Values returns all known values for AddonIssueCode. Note that this can be
// expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (AddonIssueCode) Values() []AddonIssueCode {
	return []AddonIssueCode{
		"AccessDenied",
		"InternalFailure",
		"ClusterUnreachable",
		"InsufficientNumberOfReplicas",
		"ConfigurationConflict",
		"AdmissionRequestDenied",
		"UnsupportedAddonModification",
		"K8sResourceNotFound",
		"AddonSubscriptionNeeded",
		"AddonPermissionFailure",
	}
}

type AddonStatus string

// Enum values for AddonStatus
const (
	AddonStatusCreating     AddonStatus = "CREATING"
	AddonStatusActive       AddonStatus = "ACTIVE"
	AddonStatusCreateFailed AddonStatus = "CREATE_FAILED"
	AddonStatusUpdating     AddonStatus = "UPDATING"
	AddonStatusDeleting     AddonStatus = "DELETING"
	AddonStatusDeleteFailed AddonStatus = "DELETE_FAILED"
	AddonStatusDegraded     AddonStatus = "DEGRADED"
	AddonStatusUpdateFailed AddonStatus = "UPDATE_FAILED"
)

// Values returns all known values for AddonStatus. Note that this can be expanded
// in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (AddonStatus) Values() []AddonStatus {
	return []AddonStatus{
		"CREATING",
		"ACTIVE",
		"CREATE_FAILED",
		"UPDATING",
		"DELETING",
		"DELETE_FAILED",
		"DEGRADED",
		"UPDATE_FAILED",
	}
}

type AMITypes string

// Enum values for AMITypes
const (
	AMITypesAl2X8664                AMITypes = "AL2_x86_64"
	AMITypesAl2X8664Gpu             AMITypes = "AL2_x86_64_GPU"
	AMITypesAl2Arm64                AMITypes = "AL2_ARM_64"
	AMITypesCustom                  AMITypes = "CUSTOM"
	AMITypesBottlerocketArm64       AMITypes = "BOTTLEROCKET_ARM_64"
	AMITypesBottlerocketX8664       AMITypes = "BOTTLEROCKET_x86_64"
	AMITypesBottlerocketArm64Nvidia AMITypes = "BOTTLEROCKET_ARM_64_NVIDIA"
	AMITypesBottlerocketX8664Nvidia AMITypes = "BOTTLEROCKET_x86_64_NVIDIA"
	AMITypesWindowsCore2019X8664    AMITypes = "WINDOWS_CORE_2019_x86_64"
	AMITypesWindowsFull2019X8664    AMITypes = "WINDOWS_FULL_2019_x86_64"
	AMITypesWindowsCore2022X8664    AMITypes = "WINDOWS_CORE_2022_x86_64"
	AMITypesWindowsFull2022X8664    AMITypes = "WINDOWS_FULL_2022_x86_64"
	AMITypesAl2023X8664Standard     AMITypes = "AL2023_x86_64_STANDARD"
	AMITypesAl2023Arm64Standard     AMITypes = "AL2023_ARM_64_STANDARD"
	AMITypesAl2023X8664Neuron       AMITypes = "AL2023_x86_64_NEURON"
	AMITypesAl2023X8664Nvidia       AMITypes = "AL2023_x86_64_NVIDIA"
)

// Values returns all known values for AMITypes. Note that this can be expanded in
// the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (AMITypes) Values() []AMITypes {
	return []AMITypes{
		"AL2_x86_64",
		"AL2_x86_64_GPU",
		"AL2_ARM_64",
		"CUSTOM",
		"BOTTLEROCKET_ARM_64",
		"BOTTLEROCKET_x86_64",
		"BOTTLEROCKET_ARM_64_NVIDIA",
		"BOTTLEROCKET_x86_64_NVIDIA",
		"WINDOWS_CORE_2019_x86_64",
		"WINDOWS_FULL_2019_x86_64",
		"WINDOWS_CORE_2022_x86_64",
		"WINDOWS_FULL_2022_x86_64",
		"AL2023_x86_64_STANDARD",
		"AL2023_ARM_64_STANDARD",
		"AL2023_x86_64_NEURON",
		"AL2023_x86_64_NVIDIA",
	}
}

type AuthenticationMode string

// Enum values for AuthenticationMode
const (
	AuthenticationModeApi             AuthenticationMode = "API"
	AuthenticationModeApiAndConfigMap AuthenticationMode = "API_AND_CONFIG_MAP"
	AuthenticationModeConfigMap       AuthenticationMode = "CONFIG_MAP"
)

// Values returns all known values for AuthenticationMode. Note that this can be
// expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (AuthenticationMode) Values() []AuthenticationMode {
	return []AuthenticationMode{
		"API",
		"API_AND_CONFIG_MAP",
		"CONFIG_MAP",
	}
}

type CapacityTypes string

// Enum values for CapacityTypes
const (
	CapacityTypesOnDemand      CapacityTypes = "ON_DEMAND"
	CapacityTypesSpot          CapacityTypes = "SPOT"
	CapacityTypesCapacityBlock CapacityTypes = "CAPACITY_BLOCK"
)

// Values returns all known values for CapacityTypes. Note that this can be
// expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (CapacityTypes) Values() []CapacityTypes {
	return []CapacityTypes{
		"ON_DEMAND",
		"SPOT",
		"CAPACITY_BLOCK",
	}
}

type Category string

// Enum values for Category
const (
	CategoryUpgradeReadiness Category = "UPGRADE_READINESS"
)

// Values returns all known values for Category. Note that this can be expanded in
// the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (Category) Values() []Category {
	return []Category{
		"UPGRADE_READINESS",
	}
}

type ClusterIssueCode string

// Enum values for ClusterIssueCode
const (
	ClusterIssueCodeAccessDenied                ClusterIssueCode = "AccessDenied"
	ClusterIssueCodeClusterUnreachable          ClusterIssueCode = "ClusterUnreachable"
	ClusterIssueCodeConfigurationConflict       ClusterIssueCode = "ConfigurationConflict"
	ClusterIssueCodeInternalFailure             ClusterIssueCode = "InternalFailure"
	ClusterIssueCodeResourceLimitExceeded       ClusterIssueCode = "ResourceLimitExceeded"
	ClusterIssueCodeResourceNotFound            ClusterIssueCode = "ResourceNotFound"
	ClusterIssueCodeIamRoleNotFound             ClusterIssueCode = "IamRoleNotFound"
	ClusterIssueCodeVpcNotFound                 ClusterIssueCode = "VpcNotFound"
	ClusterIssueCodeInsufficientFreeAddresses   ClusterIssueCode = "InsufficientFreeAddresses"
	ClusterIssueCodeEc2ServiceNotSubscribed     ClusterIssueCode = "Ec2ServiceNotSubscribed"
	ClusterIssueCodeEc2SubnetNotFound           ClusterIssueCode = "Ec2SubnetNotFound"
	ClusterIssueCodeEc2SecurityGroupNotFound    ClusterIssueCode = "Ec2SecurityGroupNotFound"
	ClusterIssueCodeKmsGrantRevoked             ClusterIssueCode = "KmsGrantRevoked"
	ClusterIssueCodeKmsKeyNotFound              ClusterIssueCode = "KmsKeyNotFound"
	ClusterIssueCodeKmsKeyMarkedForDeletion     ClusterIssueCode = "KmsKeyMarkedForDeletion"
	ClusterIssueCodeKmsKeyDisabled              ClusterIssueCode = "KmsKeyDisabled"
	ClusterIssueCodeStsRegionalEndpointDisabled ClusterIssueCode = "StsRegionalEndpointDisabled"
	ClusterIssueCodeUnsupportedVersion          ClusterIssueCode = "UnsupportedVersion"
	ClusterIssueCodeOther                       ClusterIssueCode = "Other"
)

// Values returns all known values for ClusterIssueCode. Note that this can be
// expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (ClusterIssueCode) Values() []ClusterIssueCode {
	return []ClusterIssueCode{
		"AccessDenied",
		"ClusterUnreachable",
		"ConfigurationConflict",
		"InternalFailure",
		"ResourceLimitExceeded",
		"ResourceNotFound",
		"IamRoleNotFound",
		"VpcNotFound",
		"InsufficientFreeAddresses",
		"Ec2ServiceNotSubscribed",
		"Ec2SubnetNotFound",
		"Ec2SecurityGroupNotFound",
		"KmsGrantRevoked",
		"KmsKeyNotFound",
		"KmsKeyMarkedForDeletion",
		"KmsKeyDisabled",
		"StsRegionalEndpointDisabled",
		"UnsupportedVersion",
		"Other",
	}
}

type ClusterStatus string

// Enum values for ClusterStatus
const (
	ClusterStatusCreating ClusterStatus = "CREATING"
	ClusterStatusActive   ClusterStatus = "ACTIVE"
	ClusterStatusDeleting ClusterStatus = "DELETING"
	ClusterStatusFailed   ClusterStatus = "FAILED"
	ClusterStatusUpdating ClusterStatus = "UPDATING"
	ClusterStatusPending  ClusterStatus = "PENDING"
)

// Values returns all known values for ClusterStatus. Note that this can be
// expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (ClusterStatus) Values() []ClusterStatus {
	return []ClusterStatus{
		"CREATING",
		"ACTIVE",
		"DELETING",
		"FAILED",
		"UPDATING",
		"PENDING",
	}
}

type ConfigStatus string

// Enum values for ConfigStatus
const (
	ConfigStatusCreating ConfigStatus = "CREATING"
	ConfigStatusDeleting ConfigStatus = "DELETING"
	ConfigStatusActive   ConfigStatus = "ACTIVE"
)

// Values returns all known values for ConfigStatus. Note that this can be
// expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (ConfigStatus) Values() []ConfigStatus {
	return []ConfigStatus{
		"CREATING",
		"DELETING",
		"ACTIVE",
	}
}

type ConnectorConfigProvider string

// Enum values for ConnectorConfigProvider
const (
	ConnectorConfigProviderEksAnywhere ConnectorConfigProvider = "EKS_ANYWHERE"
	ConnectorConfigProviderAnthos      ConnectorConfigProvider = "ANTHOS"
	ConnectorConfigProviderGke         ConnectorConfigProvider = "GKE"
	ConnectorConfigProviderAks         ConnectorConfigProvider = "AKS"
	ConnectorConfigProviderOpenshift   ConnectorConfigProvider = "OPENSHIFT"
	ConnectorConfigProviderTanzu       ConnectorConfigProvider = "TANZU"
	ConnectorConfigProviderRancher     ConnectorConfigProvider = "RANCHER"
	ConnectorConfigProviderEc2         ConnectorConfigProvider = "EC2"
	ConnectorConfigProviderOther       ConnectorConfigProvider = "OTHER"
)

// Values returns all known values for ConnectorConfigProvider. Note that this can
// be expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (ConnectorConfigProvider) Values() []ConnectorConfigProvider {
	return []ConnectorConfigProvider{
		"EKS_ANYWHERE",
		"ANTHOS",
		"GKE",
		"AKS",
		"OPENSHIFT",
		"TANZU",
		"RANCHER",
		"EC2",
		"OTHER",
	}
}

type EksAnywhereSubscriptionLicenseType string

// Enum values for EksAnywhereSubscriptionLicenseType
const (
	EksAnywhereSubscriptionLicenseTypeCluster EksAnywhereSubscriptionLicenseType = "Cluster"
)

// Values returns all known values for EksAnywhereSubscriptionLicenseType. Note
// that this can be expanded in the future, and so it is only as up to date as the
// client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (EksAnywhereSubscriptionLicenseType) Values() []EksAnywhereSubscriptionLicenseType {
	return []EksAnywhereSubscriptionLicenseType{
		"Cluster",
	}
}

type EksAnywhereSubscriptionStatus string

// Enum values for EksAnywhereSubscriptionStatus
const (
	EksAnywhereSubscriptionStatusCreating EksAnywhereSubscriptionStatus = "CREATING"
	EksAnywhereSubscriptionStatusActive   EksAnywhereSubscriptionStatus = "ACTIVE"
	EksAnywhereSubscriptionStatusUpdating EksAnywhereSubscriptionStatus = "UPDATING"
	EksAnywhereSubscriptionStatusExpiring EksAnywhereSubscriptionStatus = "EXPIRING"
	EksAnywhereSubscriptionStatusExpired  EksAnywhereSubscriptionStatus = "EXPIRED"
	EksAnywhereSubscriptionStatusDeleting EksAnywhereSubscriptionStatus = "DELETING"
)

// Values returns all known values for EksAnywhereSubscriptionStatus. Note that
// this can be expanded in the future, and so it is only as up to date as the
// client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (EksAnywhereSubscriptionStatus) Values() []EksAnywhereSubscriptionStatus {
	return []EksAnywhereSubscriptionStatus{
		"CREATING",
		"ACTIVE",
		"UPDATING",
		"EXPIRING",
		"EXPIRED",
		"DELETING",
	}
}

type EksAnywhereSubscriptionTermUnit string

// Enum values for EksAnywhereSubscriptionTermUnit
const (
	EksAnywhereSubscriptionTermUnitMonths EksAnywhereSubscriptionTermUnit = "MONTHS"
)

// Values returns all known values for EksAnywhereSubscriptionTermUnit. Note that
// this can be expanded in the future, and so it is only as up to date as the
// client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (EksAnywhereSubscriptionTermUnit) Values() []EksAnywhereSubscriptionTermUnit {
	return []EksAnywhereSubscriptionTermUnit{
		"MONTHS",
	}
}

type ErrorCode string

// Enum values for ErrorCode
const (
	ErrorCodeSubnetNotFound               ErrorCode = "SubnetNotFound"
	ErrorCodeSecurityGroupNotFound        ErrorCode = "SecurityGroupNotFound"
	ErrorCodeEniLimitReached              ErrorCode = "EniLimitReached"
	ErrorCodeIpNotAvailable               ErrorCode = "IpNotAvailable"
	ErrorCodeAccessDenied                 ErrorCode = "AccessDenied"
	ErrorCodeOperationNotPermitted        ErrorCode = "OperationNotPermitted"
	ErrorCodeVpcIdNotFound                ErrorCode = "VpcIdNotFound"
	ErrorCodeUnknown                      ErrorCode = "Unknown"
	ErrorCodeNodeCreationFailure          ErrorCode = "NodeCreationFailure"
	ErrorCodePodEvictionFailure           ErrorCode = "PodEvictionFailure"
	ErrorCodeInsufficientFreeAddresses    ErrorCode = "InsufficientFreeAddresses"
	ErrorCodeClusterUnreachable           ErrorCode = "ClusterUnreachable"
	ErrorCodeInsufficientNumberOfReplicas ErrorCode = "InsufficientNumberOfReplicas"
	ErrorCodeConfigurationConflict        ErrorCode = "ConfigurationConflict"
	ErrorCodeAdmissionRequestDenied       ErrorCode = "AdmissionRequestDenied"
	ErrorCodeUnsupportedAddonModification ErrorCode = "UnsupportedAddonModification"
	ErrorCodeK8sResourceNotFound          ErrorCode = "K8sResourceNotFound"
)

// Values returns all known values for ErrorCode. Note that this can be expanded
// in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (ErrorCode) Values() []ErrorCode {
	return []ErrorCode{
		"SubnetNotFound",
		"SecurityGroupNotFound",
		"EniLimitReached",
		"IpNotAvailable",
		"AccessDenied",
		"OperationNotPermitted",
		"VpcIdNotFound",
		"Unknown",
		"NodeCreationFailure",
		"PodEvictionFailure",
		"InsufficientFreeAddresses",
		"ClusterUnreachable",
		"InsufficientNumberOfReplicas",
		"ConfigurationConflict",
		"AdmissionRequestDenied",
		"UnsupportedAddonModification",
		"K8sResourceNotFound",
	}
}

type FargateProfileIssueCode string

// Enum values for FargateProfileIssueCode
const (
	FargateProfileIssueCodePodExecutionRoleAlreadyInUse FargateProfileIssueCode = "PodExecutionRoleAlreadyInUse"
	FargateProfileIssueCodeAccessDenied                 FargateProfileIssueCode = "AccessDenied"
	FargateProfileIssueCodeClusterUnreachable           FargateProfileIssueCode = "ClusterUnreachable"
	FargateProfileIssueCodeInternalFailure              FargateProfileIssueCode = "InternalFailure"
)

// Values returns all known values for FargateProfileIssueCode. Note that this can
// be expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (FargateProfileIssueCode) Values() []FargateProfileIssueCode {
	return []FargateProfileIssueCode{
		"PodExecutionRoleAlreadyInUse",
		"AccessDenied",
		"ClusterUnreachable",
		"InternalFailure",
	}
}

type FargateProfileStatus string

// Enum values for FargateProfileStatus
const (
	FargateProfileStatusCreating     FargateProfileStatus = "CREATING"
	FargateProfileStatusActive       FargateProfileStatus = "ACTIVE"
	FargateProfileStatusDeleting     FargateProfileStatus = "DELETING"
	FargateProfileStatusCreateFailed FargateProfileStatus = "CREATE_FAILED"
	FargateProfileStatusDeleteFailed FargateProfileStatus = "DELETE_FAILED"
)

// Values returns all known values for FargateProfileStatus. Note that this can be
// expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (FargateProfileStatus) Values() []FargateProfileStatus {
	return []FargateProfileStatus{
		"CREATING",
		"ACTIVE",
		"DELETING",
		"CREATE_FAILED",
		"DELETE_FAILED",
	}
}

type InsightStatusValue string

// Enum values for InsightStatusValue
const (
	InsightStatusValuePassing InsightStatusValue = "PASSING"
	InsightStatusValueWarning InsightStatusValue = "WARNING"
	InsightStatusValueError   InsightStatusValue = "ERROR"
	InsightStatusValueUnknown InsightStatusValue = "UNKNOWN"
)

// Values returns all known values for InsightStatusValue. Note that this can be
// expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (InsightStatusValue) Values() []InsightStatusValue {
	return []InsightStatusValue{
		"PASSING",
		"WARNING",
		"ERROR",
		"UNKNOWN",
	}
}

type IpFamily string

// Enum values for IpFamily
const (
	IpFamilyIpv4 IpFamily = "ipv4"
	IpFamilyIpv6 IpFamily = "ipv6"
)

// Values returns all known values for IpFamily. Note that this can be expanded in
// the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (IpFamily) Values() []IpFamily {
	return []IpFamily{
		"ipv4",
		"ipv6",
	}
}

type LogType string

// Enum values for LogType
const (
	LogTypeApi               LogType = "api"
	LogTypeAudit             LogType = "audit"
	LogTypeAuthenticator     LogType = "authenticator"
	LogTypeControllerManager LogType = "controllerManager"
	LogTypeScheduler         LogType = "scheduler"
)

// Values returns all known values for LogType. Note that this can be expanded in
// the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (LogType) Values() []LogType {
	return []LogType{
		"api",
		"audit",
		"authenticator",
		"controllerManager",
		"scheduler",
	}
}

type NodegroupIssueCode string

// Enum values for NodegroupIssueCode
const (
	NodegroupIssueCodeAutoScalingGroupNotFound                 NodegroupIssueCode = "AutoScalingGroupNotFound"
	NodegroupIssueCodeAutoScalingGroupInvalidConfiguration     NodegroupIssueCode = "AutoScalingGroupInvalidConfiguration"
	NodegroupIssueCodeEc2SecurityGroupNotFound                 NodegroupIssueCode = "Ec2SecurityGroupNotFound"
	NodegroupIssueCodeEc2SecurityGroupDeletionFailure          NodegroupIssueCode = "Ec2SecurityGroupDeletionFailure"
	NodegroupIssueCodeEc2LaunchTemplateNotFound                NodegroupIssueCode = "Ec2LaunchTemplateNotFound"
	NodegroupIssueCodeEc2LaunchTemplateVersionMismatch         NodegroupIssueCode = "Ec2LaunchTemplateVersionMismatch"
	NodegroupIssueCodeEc2SubnetNotFound                        NodegroupIssueCode = "Ec2SubnetNotFound"
	NodegroupIssueCodeEc2SubnetInvalidConfiguration            NodegroupIssueCode = "Ec2SubnetInvalidConfiguration"
	NodegroupIssueCodeIamInstanceProfileNotFound               NodegroupIssueCode = "IamInstanceProfileNotFound"
	NodegroupIssueCodeEc2SubnetMissingIpv6Assignment           NodegroupIssueCode = "Ec2SubnetMissingIpv6Assignment"
	NodegroupIssueCodeIamLimitExceeded                         NodegroupIssueCode = "IamLimitExceeded"
	NodegroupIssueCodeIamNodeRoleNotFound                      NodegroupIssueCode = "IamNodeRoleNotFound"
	NodegroupIssueCodeNodeCreationFailure                      NodegroupIssueCode = "NodeCreationFailure"
	NodegroupIssueCodeAsgInstanceLaunchFailures                NodegroupIssueCode = "AsgInstanceLaunchFailures"
	NodegroupIssueCodeInstanceLimitExceeded                    NodegroupIssueCode = "InstanceLimitExceeded"
	NodegroupIssueCodeInsufficientFreeAddresses                NodegroupIssueCode = "InsufficientFreeAddresses"
	NodegroupIssueCodeAccessDenied                             NodegroupIssueCode = "AccessDenied"
	NodegroupIssueCodeInternalFailure                          NodegroupIssueCode = "InternalFailure"
	NodegroupIssueCodeClusterUnreachable                       NodegroupIssueCode = "ClusterUnreachable"
	NodegroupIssueCodeAmiIdNotFound                            NodegroupIssueCode = "AmiIdNotFound"
	NodegroupIssueCodeAutoScalingGroupOptInRequired            NodegroupIssueCode = "AutoScalingGroupOptInRequired"
	NodegroupIssueCodeAutoScalingGroupRateLimitExceeded        NodegroupIssueCode = "AutoScalingGroupRateLimitExceeded"
	NodegroupIssueCodeEc2LaunchTemplateDeletionFailure         NodegroupIssueCode = "Ec2LaunchTemplateDeletionFailure"
	NodegroupIssueCodeEc2LaunchTemplateInvalidConfiguration    NodegroupIssueCode = "Ec2LaunchTemplateInvalidConfiguration"
	NodegroupIssueCodeEc2LaunchTemplateMaxLimitExceeded        NodegroupIssueCode = "Ec2LaunchTemplateMaxLimitExceeded"
	NodegroupIssueCodeEc2SubnetListTooLong                     NodegroupIssueCode = "Ec2SubnetListTooLong"
	NodegroupIssueCodeIamThrottling                            NodegroupIssueCode = "IamThrottling"
	NodegroupIssueCodeNodeTerminationFailure                   NodegroupIssueCode = "NodeTerminationFailure"
	NodegroupIssueCodePodEvictionFailure                       NodegroupIssueCode = "PodEvictionFailure"
	NodegroupIssueCodeSourceEc2LaunchTemplateNotFound          NodegroupIssueCode = "SourceEc2LaunchTemplateNotFound"
	NodegroupIssueCodeLimitExceeded                            NodegroupIssueCode = "LimitExceeded"
	NodegroupIssueCodeUnknown                                  NodegroupIssueCode = "Unknown"
	NodegroupIssueCodeAutoScalingGroupInstanceRefreshActive    NodegroupIssueCode = "AutoScalingGroupInstanceRefreshActive"
	NodegroupIssueCodeKubernetesLabelInvalid                   NodegroupIssueCode = "KubernetesLabelInvalid"
	NodegroupIssueCodeEc2LaunchTemplateVersionMaxLimitExceeded NodegroupIssueCode = "Ec2LaunchTemplateVersionMaxLimitExceeded"
)

// Values returns all known values for NodegroupIssueCode. Note that this can be
// expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (NodegroupIssueCode) Values() []NodegroupIssueCode {
	return []NodegroupIssueCode{
		"AutoScalingGroupNotFound",
		"AutoScalingGroupInvalidConfiguration",
		"Ec2SecurityGroupNotFound",
		"Ec2SecurityGroupDeletionFailure",
		"Ec2LaunchTemplateNotFound",
		"Ec2LaunchTemplateVersionMismatch",
		"Ec2SubnetNotFound",
		"Ec2SubnetInvalidConfiguration",
		"IamInstanceProfileNotFound",
		"Ec2SubnetMissingIpv6Assignment",
		"IamLimitExceeded",
		"IamNodeRoleNotFound",
		"NodeCreationFailure",
		"AsgInstanceLaunchFailures",
		"InstanceLimitExceeded",
		"InsufficientFreeAddresses",
		"AccessDenied",
		"InternalFailure",
		"ClusterUnreachable",
		"AmiIdNotFound",
		"AutoScalingGroupOptInRequired",
		"AutoScalingGroupRateLimitExceeded",
		"Ec2LaunchTemplateDeletionFailure",
		"Ec2LaunchTemplateInvalidConfiguration",
		"Ec2LaunchTemplateMaxLimitExceeded",
		"Ec2SubnetListTooLong",
		"IamThrottling",
		"NodeTerminationFailure",
		"PodEvictionFailure",
		"SourceEc2LaunchTemplateNotFound",
		"LimitExceeded",
		"Unknown",
		"AutoScalingGroupInstanceRefreshActive",
		"KubernetesLabelInvalid",
		"Ec2LaunchTemplateVersionMaxLimitExceeded",
	}
}

type NodegroupStatus string

// Enum values for NodegroupStatus
const (
	NodegroupStatusCreating     NodegroupStatus = "CREATING"
	NodegroupStatusActive       NodegroupStatus = "ACTIVE"
	NodegroupStatusUpdating     NodegroupStatus = "UPDATING"
	NodegroupStatusDeleting     NodegroupStatus = "DELETING"
	NodegroupStatusCreateFailed NodegroupStatus = "CREATE_FAILED"
	NodegroupStatusDeleteFailed NodegroupStatus = "DELETE_FAILED"
	NodegroupStatusDegraded     NodegroupStatus = "DEGRADED"
)

// Values returns all known values for NodegroupStatus. Note that this can be
// expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (NodegroupStatus) Values() []NodegroupStatus {
	return []NodegroupStatus{
		"CREATING",
		"ACTIVE",
		"UPDATING",
		"DELETING",
		"CREATE_FAILED",
		"DELETE_FAILED",
		"DEGRADED",
	}
}

type ResolveConflicts string

// Enum values for ResolveConflicts
const (
	ResolveConflictsOverwrite ResolveConflicts = "OVERWRITE"
	ResolveConflictsNone      ResolveConflicts = "NONE"
	ResolveConflictsPreserve  ResolveConflicts = "PRESERVE"
)

// Values returns all known values for ResolveConflicts. Note that this can be
// expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (ResolveConflicts) Values() []ResolveConflicts {
	return []ResolveConflicts{
		"OVERWRITE",
		"NONE",
		"PRESERVE",
	}
}

type SupportType string

// Enum values for SupportType
const (
	SupportTypeStandard SupportType = "STANDARD"
	SupportTypeExtended SupportType = "EXTENDED"
)

// Values returns all known values for SupportType. Note that this can be expanded
// in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (SupportType) Values() []SupportType {
	return []SupportType{
		"STANDARD",
		"EXTENDED",
	}
}

type TaintEffect string

// Enum values for TaintEffect
const (
	TaintEffectNoSchedule       TaintEffect = "NO_SCHEDULE"
	TaintEffectNoExecute        TaintEffect = "NO_EXECUTE"
	TaintEffectPreferNoSchedule TaintEffect = "PREFER_NO_SCHEDULE"
)

// Values returns all known values for TaintEffect. Note that this can be expanded
// in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (TaintEffect) Values() []TaintEffect {
	return []TaintEffect{
		"NO_SCHEDULE",
		"NO_EXECUTE",
		"PREFER_NO_SCHEDULE",
	}
}

type UpdateParamType string

// Enum values for UpdateParamType
const (
	UpdateParamTypeVersion                  UpdateParamType = "Version"
	UpdateParamTypePlatformVersion          UpdateParamType = "PlatformVersion"
	UpdateParamTypeEndpointPrivateAccess    UpdateParamType = "EndpointPrivateAccess"
	UpdateParamTypeEndpointPublicAccess     UpdateParamType = "EndpointPublicAccess"
	UpdateParamTypeClusterLogging           UpdateParamType = "ClusterLogging"
	UpdateParamTypeDesiredSize              UpdateParamType = "DesiredSize"
	UpdateParamTypeLabelsToAdd              UpdateParamType = "LabelsToAdd"
	UpdateParamTypeLabelsToRemove           UpdateParamType = "LabelsToRemove"
	UpdateParamTypeTaintsToAdd              UpdateParamType = "TaintsToAdd"
	UpdateParamTypeTaintsToRemove           UpdateParamType = "TaintsToRemove"
	UpdateParamTypeMaxSize                  UpdateParamType = "MaxSize"
	UpdateParamTypeMinSize                  UpdateParamType = "MinSize"
	UpdateParamTypeReleaseVersion           UpdateParamType = "ReleaseVersion"
	UpdateParamTypePublicAccessCidrs        UpdateParamType = "PublicAccessCidrs"
	UpdateParamTypeLaunchTemplateName       UpdateParamType = "LaunchTemplateName"
	UpdateParamTypeLaunchTemplateVersion    UpdateParamType = "LaunchTemplateVersion"
	UpdateParamTypeIdentityProviderConfig   UpdateParamType = "IdentityProviderConfig"
	UpdateParamTypeEncryptionConfig         UpdateParamType = "EncryptionConfig"
	UpdateParamTypeAddonVersion             UpdateParamType = "AddonVersion"
	UpdateParamTypeServiceAccountRoleArn    UpdateParamType = "ServiceAccountRoleArn"
	UpdateParamTypeResolveConflicts         UpdateParamType = "ResolveConflicts"
	UpdateParamTypeMaxUnavailable           UpdateParamType = "MaxUnavailable"
	UpdateParamTypeMaxUnavailablePercentage UpdateParamType = "MaxUnavailablePercentage"
	UpdateParamTypeConfigurationValues      UpdateParamType = "ConfigurationValues"
	UpdateParamTypeSecurityGroups           UpdateParamType = "SecurityGroups"
	UpdateParamTypeSubnets                  UpdateParamType = "Subnets"
	UpdateParamTypeAuthenticationMode       UpdateParamType = "AuthenticationMode"
	UpdateParamTypePodIdentityAssociations  UpdateParamType = "PodIdentityAssociations"
	UpdateParamTypeUpgradePolicy            UpdateParamType = "UpgradePolicy"
	UpdateParamTypeZonalShiftConfig         UpdateParamType = "ZonalShiftConfig"
)

// Values returns all known values for UpdateParamType. Note that this can be
// expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (UpdateParamType) Values() []UpdateParamType {
	return []UpdateParamType{
		"Version",
		"PlatformVersion",
		"EndpointPrivateAccess",
		"EndpointPublicAccess",
		"ClusterLogging",
		"DesiredSize",
		"LabelsToAdd",
		"LabelsToRemove",
		"TaintsToAdd",
		"TaintsToRemove",
		"MaxSize",
		"MinSize",
		"ReleaseVersion",
		"PublicAccessCidrs",
		"LaunchTemplateName",
		"LaunchTemplateVersion",
		"IdentityProviderConfig",
		"EncryptionConfig",
		"AddonVersion",
		"ServiceAccountRoleArn",
		"ResolveConflicts",
		"MaxUnavailable",
		"MaxUnavailablePercentage",
		"ConfigurationValues",
		"SecurityGroups",
		"Subnets",
		"AuthenticationMode",
		"PodIdentityAssociations",
		"UpgradePolicy",
		"ZonalShiftConfig",
	}
}

type UpdateStatus string

// Enum values for UpdateStatus
const (
	UpdateStatusInProgress UpdateStatus = "InProgress"
	UpdateStatusFailed     UpdateStatus = "Failed"
	UpdateStatusCancelled  UpdateStatus = "Cancelled"
	UpdateStatusSuccessful UpdateStatus = "Successful"
)

// Values returns all known values for UpdateStatus. Note that this can be
// expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (UpdateStatus) Values() []UpdateStatus {
	return []UpdateStatus{
		"InProgress",
		"Failed",
		"Cancelled",
		"Successful",
	}
}

type UpdateType string

// Enum values for UpdateType
const (
	UpdateTypeVersionUpdate                      UpdateType = "VersionUpdate"
	UpdateTypeEndpointAccessUpdate               UpdateType = "EndpointAccessUpdate"
	UpdateTypeLoggingUpdate                      UpdateType = "LoggingUpdate"
	UpdateTypeConfigUpdate                       UpdateType = "ConfigUpdate"
	UpdateTypeAssociateIdentityProviderConfig    UpdateType = "AssociateIdentityProviderConfig"
	UpdateTypeDisassociateIdentityProviderConfig UpdateType = "DisassociateIdentityProviderConfig"
	UpdateTypeAssociateEncryptionConfig          UpdateType = "AssociateEncryptionConfig"
	UpdateTypeAddonUpdate                        UpdateType = "AddonUpdate"
	UpdateTypeVpcConfigUpdate                    UpdateType = "VpcConfigUpdate"
	UpdateTypeAccessConfigUpdate                 UpdateType = "AccessConfigUpdate"
	UpdateTypeUpgradePolicyUpdate                UpdateType = "UpgradePolicyUpdate"
	UpdateTypeZonalShiftConfigUpdate             UpdateType = "ZonalShiftConfigUpdate"
)

// Values returns all known values for UpdateType. Note that this can be expanded
// in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (UpdateType) Values() []UpdateType {
	return []UpdateType{
		"VersionUpdate",
		"EndpointAccessUpdate",
		"LoggingUpdate",
		"ConfigUpdate",
		"AssociateIdentityProviderConfig",
		"DisassociateIdentityProviderConfig",
		"AssociateEncryptionConfig",
		"AddonUpdate",
		"VpcConfigUpdate",
		"AccessConfigUpdate",
		"UpgradePolicyUpdate",
		"ZonalShiftConfigUpdate",
	}
}
