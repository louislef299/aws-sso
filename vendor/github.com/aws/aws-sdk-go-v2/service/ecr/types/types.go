// Code generated by smithy-go-codegen DO NOT EDIT.

package types

import (
	smithydocument "github.com/aws/smithy-go/document"
	"time"
)

// This data type is used in the ImageScanFinding data type.
type Attribute struct {

	// The attribute key.
	//
	// This member is required.
	Key *string

	// The value assigned to the attribute key.
	Value *string

	noSmithyDocumentSerde
}

// An object representing authorization data for an Amazon ECR registry.
type AuthorizationData struct {

	// A base64-encoded string that contains authorization data for the specified
	// Amazon ECR registry. When the string is decoded, it is presented in the format
	// user:password for private registry authentication using docker login .
	AuthorizationToken *string

	// The Unix time in seconds and milliseconds when the authorization token expires.
	// Authorization tokens are valid for 12 hours.
	ExpiresAt *time.Time

	// The registry URL to use for this authorization token in a docker login command.
	// The Amazon ECR registry URL format is
	// https://aws_account_id.dkr.ecr.region.amazonaws.com . For example,
	// https://012345678910.dkr.ecr.us-east-1.amazonaws.com ..
	ProxyEndpoint *string

	noSmithyDocumentSerde
}

// The image details of the Amazon ECR container image.
type AwsEcrContainerImageDetails struct {

	// The architecture of the Amazon ECR container image.
	Architecture *string

	// The image author of the Amazon ECR container image.
	Author *string

	// The image hash of the Amazon ECR container image.
	ImageHash *string

	// The image tags attached to the Amazon ECR container image.
	ImageTags []string

	// The platform of the Amazon ECR container image.
	Platform *string

	// The date and time the Amazon ECR container image was pushed.
	PushedAt *time.Time

	// The registry the Amazon ECR container image belongs to.
	Registry *string

	// The name of the repository the Amazon ECR container image resides in.
	RepositoryName *string

	noSmithyDocumentSerde
}

// The CVSS score for a finding.
type CvssScore struct {

	// The base CVSS score used for the finding.
	BaseScore float64

	// The vector string of the CVSS score.
	ScoringVector *string

	// The source of the CVSS score.
	Source *string

	// The version of CVSS used for the score.
	Version *string

	noSmithyDocumentSerde
}

// Details on adjustments Amazon Inspector made to the CVSS score for a finding.
type CvssScoreAdjustment struct {

	// The metric used to adjust the CVSS score.
	Metric *string

	// The reason the CVSS score has been adjustment.
	Reason *string

	noSmithyDocumentSerde
}

// Information about the CVSS score.
type CvssScoreDetails struct {

	// An object that contains details about adjustment Amazon Inspector made to the
	// CVSS score.
	Adjustments []CvssScoreAdjustment

	// The CVSS score.
	Score float64

	// The source for the CVSS score.
	ScoreSource *string

	// The vector for the CVSS score.
	ScoringVector *string

	// The CVSS version used in scoring.
	Version *string

	noSmithyDocumentSerde
}

// An object representing a filter on a DescribeImages operation.
type DescribeImagesFilter struct {

	// The tag status with which to filter your DescribeImages results. You can filter results based
	// on whether they are TAGGED or UNTAGGED .
	TagStatus TagStatus

	noSmithyDocumentSerde
}

// The encryption configuration for the repository. This determines how the
// contents of your repository are encrypted at rest.
//
// By default, when no encryption configuration is set or the AES256 encryption
// type is used, Amazon ECR uses server-side encryption with Amazon S3-managed
// encryption keys which encrypts your data at rest using an AES256 encryption
// algorithm. This does not require any action on your part.
//
// For more control over the encryption of the contents of your repository, you
// can use server-side encryption with Key Management Service key stored in Key
// Management Service (KMS) to encrypt your images. For more information, see [Amazon ECR encryption at rest]in
// the Amazon Elastic Container Registry User Guide.
//
// [Amazon ECR encryption at rest]: https://docs.aws.amazon.com/AmazonECR/latest/userguide/encryption-at-rest.html
type EncryptionConfiguration struct {

	// The encryption type to use.
	//
	// If you use the KMS encryption type, the contents of the repository will be
	// encrypted using server-side encryption with Key Management Service key stored in
	// KMS. When you use KMS to encrypt your data, you can either use the default
	// Amazon Web Services managed KMS key for Amazon ECR, or specify your own KMS key,
	// which you already created.
	//
	// If you use the KMS_DSSE encryption type, the contents of the repository will be
	// encrypted with two layers of encryption using server-side encryption with the
	// KMS Management Service key stored in KMS. Similar to the KMS encryption type,
	// you can either use the default Amazon Web Services managed KMS key for Amazon
	// ECR, or specify your own KMS key, which you've already created.
	//
	// If you use the AES256 encryption type, Amazon ECR uses server-side encryption
	// with Amazon S3-managed encryption keys which encrypts the images in the
	// repository using an AES256 encryption algorithm.
	//
	// For more information, see [Amazon ECR encryption at rest] in the Amazon Elastic Container Registry User Guide.
	//
	// [Amazon ECR encryption at rest]: https://docs.aws.amazon.com/AmazonECR/latest/userguide/encryption-at-rest.html
	//
	// This member is required.
	EncryptionType EncryptionType

	// If you use the KMS encryption type, specify the KMS key to use for encryption.
	// The alias, key ID, or full ARN of the KMS key can be specified. The key must
	// exist in the same Region as the repository. If no key is specified, the default
	// Amazon Web Services managed KMS key for Amazon ECR will be used.
	KmsKey *string

	noSmithyDocumentSerde
}

// The encryption configuration to associate with the repository creation template.
type EncryptionConfigurationForRepositoryCreationTemplate struct {

	// The encryption type to use.
	//
	// If you use the KMS encryption type, the contents of the repository will be
	// encrypted using server-side encryption with Key Management Service key stored in
	// KMS. When you use KMS to encrypt your data, you can either use the default
	// Amazon Web Services managed KMS key for Amazon ECR, or specify your own KMS key,
	// which you already created. For more information, see [Protecting data using server-side encryption with an KMS key stored in Key Management Service (SSE-KMS)]in the Amazon Simple
	// Storage Service Console Developer Guide.
	//
	// If you use the AES256 encryption type, Amazon ECR uses server-side encryption
	// with Amazon S3-managed encryption keys which encrypts the images in the
	// repository using an AES256 encryption algorithm. For more information, see [Protecting data using server-side encryption with Amazon S3-managed encryption keys (SSE-S3)]in
	// the Amazon Simple Storage Service Console Developer Guide.
	//
	// [Protecting data using server-side encryption with Amazon S3-managed encryption keys (SSE-S3)]: https://docs.aws.amazon.com/AmazonS3/latest/dev/UsingServerSideEncryption.html
	// [Protecting data using server-side encryption with an KMS key stored in Key Management Service (SSE-KMS)]: https://docs.aws.amazon.com/AmazonS3/latest/dev/UsingKMSEncryption.html
	//
	// This member is required.
	EncryptionType EncryptionType

	// If you use the KMS encryption type, specify the KMS key to use for encryption.
	// The full ARN of the KMS key must be specified. The key must exist in the same
	// Region as the repository. If no key is specified, the default Amazon Web
	// Services managed KMS key for Amazon ECR will be used.
	KmsKey *string

	noSmithyDocumentSerde
}

// The details of an enhanced image scan. This is returned when enhanced scanning
// is enabled for your private registry.
type EnhancedImageScanFinding struct {

	// The Amazon Web Services account ID associated with the image.
	AwsAccountId *string

	// The description of the finding.
	Description *string

	// If a finding discovered in your environment has an exploit available.
	ExploitAvailable *string

	// The Amazon Resource Number (ARN) of the finding.
	FindingArn *string

	// The date and time that the finding was first observed.
	FirstObservedAt *time.Time

	// Details on whether a fix is available through a version update. This value can
	// be YES , NO , or PARTIAL . A PARTIAL fix means that some, but not all, of the
	// packages identified in the finding have fixes available through updated
	// versions.
	FixAvailable *string

	// The date and time that the finding was last observed.
	LastObservedAt *time.Time

	// An object that contains the details of a package vulnerability finding.
	PackageVulnerabilityDetails *PackageVulnerabilityDetails

	// An object that contains the details about how to remediate a finding.
	Remediation *Remediation

	// Contains information on the resources involved in a finding.
	Resources []Resource

	// The Amazon Inspector score given to the finding.
	Score float64

	// An object that contains details of the Amazon Inspector score.
	ScoreDetails *ScoreDetails

	// The severity of the finding.
	Severity *string

	// The status of the finding.
	Status *string

	// The title of the finding.
	Title *string

	// The type of the finding.
	Type *string

	// The date and time the finding was last updated at.
	UpdatedAt *time.Time

	noSmithyDocumentSerde
}

// An object representing an Amazon ECR image.
type Image struct {

	// An object containing the image tag and image digest associated with an image.
	ImageId *ImageIdentifier

	// The image manifest associated with the image.
	ImageManifest *string

	// The manifest media type of the image.
	ImageManifestMediaType *string

	// The Amazon Web Services account ID associated with the registry containing the
	// image.
	RegistryId *string

	// The name of the repository associated with the image.
	RepositoryName *string

	noSmithyDocumentSerde
}

// An object that describes an image returned by a DescribeImages operation.
type ImageDetail struct {

	// The artifact media type of the image.
	ArtifactMediaType *string

	// The sha256 digest of the image manifest.
	ImageDigest *string

	// The media type of the image manifest.
	ImageManifestMediaType *string

	// The date and time, expressed in standard JavaScript date format, at which the
	// current image was pushed to the repository.
	ImagePushedAt *time.Time

	// A summary of the last completed image scan.
	ImageScanFindingsSummary *ImageScanFindingsSummary

	// The current state of the scan.
	ImageScanStatus *ImageScanStatus

	// The size, in bytes, of the image in the repository.
	//
	// If the image is a manifest list, this will be the max size of all manifests in
	// the list.
	//
	// Starting with Docker version 1.9, the Docker client compresses image layers
	// before pushing them to a V2 Docker registry. The output of the docker images
	// command shows the uncompressed image size. Therefore, Docker might return a
	// larger image than the image sizes returned by DescribeImages.
	ImageSizeInBytes *int64

	// The list of tags associated with this image.
	ImageTags []string

	// The date and time, expressed in standard JavaScript date format, when Amazon
	// ECR recorded the last image pull.
	//
	// Amazon ECR refreshes the last image pull timestamp at least once every 24
	// hours. For example, if you pull an image once a day then the
	// lastRecordedPullTime timestamp will indicate the exact time that the image was
	// last pulled. However, if you pull an image once an hour, because Amazon ECR
	// refreshes the lastRecordedPullTime timestamp at least once every 24 hours, the
	// result may not be the exact time that the image was last pulled.
	LastRecordedPullTime *time.Time

	// The Amazon Web Services account ID associated with the registry to which this
	// image belongs.
	RegistryId *string

	// The name of the repository to which this image belongs.
	RepositoryName *string

	noSmithyDocumentSerde
}

// An object representing an Amazon ECR image failure.
type ImageFailure struct {

	// The code associated with the failure.
	FailureCode ImageFailureCode

	// The reason for the failure.
	FailureReason *string

	// The image ID associated with the failure.
	ImageId *ImageIdentifier

	noSmithyDocumentSerde
}

// An object with identifying information for an image in an Amazon ECR repository.
type ImageIdentifier struct {

	// The sha256 digest of the image manifest.
	ImageDigest *string

	// The tag used for the image.
	ImageTag *string

	noSmithyDocumentSerde
}

// The status of the replication process for an image.
type ImageReplicationStatus struct {

	// The failure code for a replication that has failed.
	FailureCode *string

	// The destination Region for the image replication.
	Region *string

	// The Amazon Web Services account ID associated with the registry to which the
	// image belongs.
	RegistryId *string

	// The image replication status.
	Status ReplicationStatus

	noSmithyDocumentSerde
}

// Contains information about an image scan finding.
type ImageScanFinding struct {

	// A collection of attributes of the host from which the finding is generated.
	Attributes []Attribute

	// The description of the finding.
	Description *string

	// The name associated with the finding, usually a CVE number.
	Name *string

	// The finding severity.
	Severity FindingSeverity

	// A link containing additional details about the security vulnerability.
	Uri *string

	noSmithyDocumentSerde
}

// The details of an image scan.
type ImageScanFindings struct {

	// Details about the enhanced scan findings from Amazon Inspector.
	EnhancedFindings []EnhancedImageScanFinding

	// The image vulnerability counts, sorted by severity.
	FindingSeverityCounts map[string]int32

	// The findings from the image scan.
	Findings []ImageScanFinding

	// The time of the last completed image scan.
	ImageScanCompletedAt *time.Time

	// The time when the vulnerability data was last scanned.
	VulnerabilitySourceUpdatedAt *time.Time

	noSmithyDocumentSerde
}

// A summary of the last completed image scan.
type ImageScanFindingsSummary struct {

	// The image vulnerability counts, sorted by severity.
	FindingSeverityCounts map[string]int32

	// The time of the last completed image scan.
	ImageScanCompletedAt *time.Time

	// The time when the vulnerability data was last scanned.
	VulnerabilitySourceUpdatedAt *time.Time

	noSmithyDocumentSerde
}

// The image scanning configuration for a repository.
type ImageScanningConfiguration struct {

	// The setting that determines whether images are scanned after being pushed to a
	// repository. If set to true , images will be scanned after being pushed. If this
	// parameter is not specified, it will default to false and images will not be
	// scanned unless a scan is manually started with the [API_StartImageScan]API.
	//
	// [API_StartImageScan]: https://docs.aws.amazon.com/AmazonECR/latest/APIReference/API_StartImageScan.html
	ScanOnPush bool

	noSmithyDocumentSerde
}

// The current status of an image scan.
type ImageScanStatus struct {

	// The description of the image scan status.
	Description *string

	// The current state of an image scan.
	Status ScanStatus

	noSmithyDocumentSerde
}

// An object representing an Amazon ECR image layer.
type Layer struct {

	// The availability status of the image layer.
	LayerAvailability LayerAvailability

	// The sha256 digest of the image layer.
	LayerDigest *string

	// The size, in bytes, of the image layer.
	LayerSize *int64

	// The media type of the layer, such as
	// application/vnd.docker.image.rootfs.diff.tar.gzip or
	// application/vnd.oci.image.layer.v1.tar+gzip .
	MediaType *string

	noSmithyDocumentSerde
}

// An object representing an Amazon ECR image layer failure.
type LayerFailure struct {

	// The failure code associated with the failure.
	FailureCode LayerFailureCode

	// The reason for the failure.
	FailureReason *string

	// The layer digest associated with the failure.
	LayerDigest *string

	noSmithyDocumentSerde
}

// The filter for the lifecycle policy preview.
type LifecyclePolicyPreviewFilter struct {

	// The tag status of the image.
	TagStatus TagStatus

	noSmithyDocumentSerde
}

// The result of the lifecycle policy preview.
type LifecyclePolicyPreviewResult struct {

	// The type of action to be taken.
	Action *LifecyclePolicyRuleAction

	// The priority of the applied rule.
	AppliedRulePriority *int32

	// The sha256 digest of the image manifest.
	ImageDigest *string

	// The date and time, expressed in standard JavaScript date format, at which the
	// current image was pushed to the repository.
	ImagePushedAt *time.Time

	// The list of tags associated with this image.
	ImageTags []string

	noSmithyDocumentSerde
}

// The summary of the lifecycle policy preview request.
type LifecyclePolicyPreviewSummary struct {

	// The number of expiring images.
	ExpiringImageTotalCount *int32

	noSmithyDocumentSerde
}

// The type of action to be taken.
type LifecyclePolicyRuleAction struct {

	// The type of action to be taken.
	Type ImageActionType

	noSmithyDocumentSerde
}

// An object representing a filter on a ListImages operation.
type ListImagesFilter struct {

	// The tag status with which to filter your ListImages results. You can filter results based
	// on whether they are TAGGED or UNTAGGED .
	TagStatus TagStatus

	noSmithyDocumentSerde
}

// Information about a package vulnerability finding.
type PackageVulnerabilityDetails struct {

	// An object that contains details about the CVSS score of a finding.
	Cvss []CvssScore

	// One or more URLs that contain details about this vulnerability type.
	ReferenceUrls []string

	// One or more vulnerabilities related to the one identified in this finding.
	RelatedVulnerabilities []string

	// The source of the vulnerability information.
	Source *string

	// A URL to the source of the vulnerability information.
	SourceUrl *string

	// The date and time that this vulnerability was first added to the vendor's
	// database.
	VendorCreatedAt *time.Time

	// The severity the vendor has given to this vulnerability type.
	VendorSeverity *string

	// The date and time the vendor last updated this vulnerability in their database.
	VendorUpdatedAt *time.Time

	// The ID given to this vulnerability.
	VulnerabilityId *string

	// The packages impacted by this vulnerability.
	VulnerablePackages []VulnerablePackage

	noSmithyDocumentSerde
}

// The details of a pull through cache rule.
type PullThroughCacheRule struct {

	// The date and time the pull through cache was created.
	CreatedAt *time.Time

	// The ARN of the Secrets Manager secret associated with the pull through cache
	// rule.
	CredentialArn *string

	// The ARN of the IAM role associated with the pull through cache rule.
	CustomRoleArn *string

	// The Amazon ECR repository prefix associated with the pull through cache rule.
	EcrRepositoryPrefix *string

	// The Amazon Web Services account ID associated with the registry the pull
	// through cache rule is associated with.
	RegistryId *string

	// The date and time, in JavaScript date format, when the pull through cache rule
	// was last updated.
	UpdatedAt *time.Time

	// The name of the upstream source registry associated with the pull through cache
	// rule.
	UpstreamRegistry UpstreamRegistry

	// The upstream registry URL associated with the pull through cache rule.
	UpstreamRegistryUrl *string

	// The upstream repository prefix associated with the pull through cache rule.
	UpstreamRepositoryPrefix *string

	noSmithyDocumentSerde
}

// Details about the recommended course of action to remediate the finding.
type Recommendation struct {

	// The recommended course of action to remediate the finding.
	Text *string

	// The URL address to the CVE remediation recommendations.
	Url *string

	noSmithyDocumentSerde
}

// The scanning configuration for a private registry.
type RegistryScanningConfiguration struct {

	// The scanning rules associated with the registry.
	Rules []RegistryScanningRule

	// The type of scanning configured for the registry.
	ScanType ScanType

	noSmithyDocumentSerde
}

// The details of a scanning rule for a private registry.
type RegistryScanningRule struct {

	// The repository filters associated with the scanning configuration for a private
	// registry.
	//
	// This member is required.
	RepositoryFilters []ScanningRepositoryFilter

	// The frequency that scans are performed at for a private registry. When the
	// ENHANCED scan type is specified, the supported scan frequencies are
	// CONTINUOUS_SCAN and SCAN_ON_PUSH . When the BASIC scan type is specified, the
	// SCAN_ON_PUSH scan frequency is supported. If scan on push is not specified, then
	// the MANUAL scan frequency is set by default.
	//
	// This member is required.
	ScanFrequency ScanFrequency

	noSmithyDocumentSerde
}

// Information on how to remediate a finding.
type Remediation struct {

	// An object that contains information about the recommended course of action to
	// remediate the finding.
	Recommendation *Recommendation

	noSmithyDocumentSerde
}

// The replication configuration for a registry.
type ReplicationConfiguration struct {

	// An array of objects representing the replication destinations and repository
	// filters for a replication configuration.
	//
	// This member is required.
	Rules []ReplicationRule

	noSmithyDocumentSerde
}

// An array of objects representing the destination for a replication rule.
type ReplicationDestination struct {

	// The Region to replicate to.
	//
	// This member is required.
	Region *string

	// The Amazon Web Services account ID of the Amazon ECR private registry to
	// replicate to. When configuring cross-Region replication within your own
	// registry, specify your own account ID.
	//
	// This member is required.
	RegistryId *string

	noSmithyDocumentSerde
}

// An array of objects representing the replication destinations and repository
// filters for a replication configuration.
type ReplicationRule struct {

	// An array of objects representing the destination for a replication rule.
	//
	// This member is required.
	Destinations []ReplicationDestination

	// An array of objects representing the filters for a replication rule. Specifying
	// a repository filter for a replication rule provides a method for controlling
	// which repositories in a private registry are replicated.
	RepositoryFilters []RepositoryFilter

	noSmithyDocumentSerde
}

// An object representing a repository.
type Repository struct {

	// The date and time, in JavaScript date format, when the repository was created.
	CreatedAt *time.Time

	// The encryption configuration for the repository. This determines how the
	// contents of your repository are encrypted at rest.
	EncryptionConfiguration *EncryptionConfiguration

	// The image scanning configuration for a repository.
	ImageScanningConfiguration *ImageScanningConfiguration

	// The tag mutability setting for the repository.
	ImageTagMutability ImageTagMutability

	// The Amazon Web Services account ID associated with the registry that contains
	// the repository.
	RegistryId *string

	// The Amazon Resource Name (ARN) that identifies the repository. The ARN contains
	// the arn:aws:ecr namespace, followed by the region of the repository, Amazon Web
	// Services account ID of the repository owner, repository namespace, and
	// repository name. For example,
	// arn:aws:ecr:region:012345678910:repository-namespace/repository-name .
	RepositoryArn *string

	// The name of the repository.
	RepositoryName *string

	// The URI for the repository. You can use this URI for container image push and
	// pull operations.
	RepositoryUri *string

	noSmithyDocumentSerde
}

// The details of the repository creation template associated with the request.
type RepositoryCreationTemplate struct {

	// A list of enumerable Strings representing the repository creation scenarios
	// that this template will apply towards. The two supported scenarios are
	// PULL_THROUGH_CACHE and REPLICATION
	AppliedFor []RCTAppliedFor

	// The date and time, in JavaScript date format, when the repository creation
	// template was created.
	CreatedAt *time.Time

	// The ARN of the role to be assumed by Amazon ECR. Amazon ECR will assume your
	// supplied role when the customRoleArn is specified. When this field isn't
	// specified, Amazon ECR will use the service-linked role for the repository
	// creation template.
	CustomRoleArn *string

	// The description associated with the repository creation template.
	Description *string

	// The encryption configuration associated with the repository creation template.
	EncryptionConfiguration *EncryptionConfigurationForRepositoryCreationTemplate

	// The tag mutability setting for the repository. If this parameter is omitted,
	// the default setting of MUTABLE will be used which will allow image tags to be
	// overwritten. If IMMUTABLE is specified, all image tags within the repository
	// will be immutable which will prevent them from being overwritten.
	ImageTagMutability ImageTagMutability

	// The lifecycle policy to use for repositories created using the template.
	LifecyclePolicy *string

	// The repository namespace prefix associated with the repository creation
	// template.
	Prefix *string

	// The repository policy to apply to repositories created using the template. A
	// repository policy is a permissions policy associated with a repository to
	// control access permissions.
	RepositoryPolicy *string

	// The metadata to apply to the repository to help you categorize and organize.
	// Each tag consists of a key and an optional value, both of which you define. Tag
	// keys can have a maximum character length of 128 characters, and tag values can
	// have a maximum length of 256 characters.
	ResourceTags []Tag

	// The date and time, in JavaScript date format, when the repository creation
	// template was last updated.
	UpdatedAt *time.Time

	noSmithyDocumentSerde
}

// The filter settings used with image replication. Specifying a repository filter
// to a replication rule provides a method for controlling which repositories in a
// private registry are replicated. If no filters are added, the contents of all
// repositories are replicated.
type RepositoryFilter struct {

	// The repository filter details. When the PREFIX_MATCH filter type is specified,
	// this value is required and should be the repository name prefix to configure
	// replication for.
	//
	// This member is required.
	Filter *string

	// The repository filter type. The only supported value is PREFIX_MATCH , which is
	// a repository name prefix specified with the filter parameter.
	//
	// This member is required.
	FilterType RepositoryFilterType

	noSmithyDocumentSerde
}

// The details of the scanning configuration for a repository.
type RepositoryScanningConfiguration struct {

	// The scan filters applied to the repository.
	AppliedScanFilters []ScanningRepositoryFilter

	// The ARN of the repository.
	RepositoryArn *string

	// The name of the repository.
	RepositoryName *string

	// The scan frequency for the repository.
	ScanFrequency ScanFrequency

	// Whether or not scan on push is configured for the repository.
	ScanOnPush bool

	noSmithyDocumentSerde
}

// The details about any failures associated with the scanning configuration of a
// repository.
type RepositoryScanningConfigurationFailure struct {

	// The failure code.
	FailureCode ScanningConfigurationFailureCode

	// The reason for the failure.
	FailureReason *string

	// The name of the repository.
	RepositoryName *string

	noSmithyDocumentSerde
}

// Details about the resource involved in a finding.
type Resource struct {

	// An object that contains details about the resource involved in a finding.
	Details *ResourceDetails

	// The ID of the resource.
	Id *string

	// The tags attached to the resource.
	Tags map[string]string

	// The type of resource.
	Type *string

	noSmithyDocumentSerde
}

// Contains details about the resource involved in the finding.
type ResourceDetails struct {

	// An object that contains details about the Amazon ECR container image involved
	// in the finding.
	AwsEcrContainerImage *AwsEcrContainerImageDetails

	noSmithyDocumentSerde
}

// The details of a scanning repository filter. For more information on how to use
// filters, see [Using filters]in the Amazon Elastic Container Registry User Guide.
//
// [Using filters]: https://docs.aws.amazon.com/AmazonECR/latest/userguide/image-scanning.html#image-scanning-filters
type ScanningRepositoryFilter struct {

	// The filter to use when scanning.
	//
	// This member is required.
	Filter *string

	// The type associated with the filter.
	//
	// This member is required.
	FilterType ScanningRepositoryFilterType

	noSmithyDocumentSerde
}

// Information about the Amazon Inspector score given to a finding.
type ScoreDetails struct {

	// An object that contains details about the CVSS score given to a finding.
	Cvss *CvssScoreDetails

	noSmithyDocumentSerde
}

// The metadata to apply to a resource to help you categorize and organize them.
// Each tag consists of a key and a value, both of which you define. Tag keys can
// have a maximum character length of 128 characters, and tag values can have a
// maximum length of 256 characters.
type Tag struct {

	// One part of a key-value pair that make up a tag. A key is a general label that
	// acts like a category for more specific tag values.
	//
	// This member is required.
	Key *string

	// A value acts as a descriptor within a tag category (key).
	//
	// This member is required.
	Value *string

	noSmithyDocumentSerde
}

// Information on the vulnerable package identified by a finding.
type VulnerablePackage struct {

	// The architecture of the vulnerable package.
	Arch *string

	// The epoch of the vulnerable package.
	Epoch *int32

	// The file path of the vulnerable package.
	FilePath *string

	// The version of the package that contains the vulnerability fix.
	FixedInVersion *string

	// The name of the vulnerable package.
	Name *string

	// The package manager of the vulnerable package.
	PackageManager *string

	// The release of the vulnerable package.
	Release *string

	// The source layer hash of the vulnerable package.
	SourceLayerHash *string

	// The version of the vulnerable package.
	Version *string

	noSmithyDocumentSerde
}

type noSmithyDocumentSerde = smithydocument.NoSerde
