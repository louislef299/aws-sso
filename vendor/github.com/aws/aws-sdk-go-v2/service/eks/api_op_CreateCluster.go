// Code generated by smithy-go-codegen DO NOT EDIT.

package eks

import (
	"context"
	"fmt"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// Creates an Amazon EKS control plane.
//
// The Amazon EKS control plane consists of control plane instances that run the
// Kubernetes software, such as etcd and the API server. The control plane runs in
// an account managed by Amazon Web Services, and the Kubernetes API is exposed by
// the Amazon EKS API server endpoint. Each Amazon EKS cluster control plane is
// single tenant and unique. It runs on its own set of Amazon EC2 instances.
//
// The cluster control plane is provisioned across multiple Availability Zones and
// fronted by an Elastic Load Balancing Network Load Balancer. Amazon EKS also
// provisions elastic network interfaces in your VPC subnets to provide
// connectivity from the control plane instances to the nodes (for example, to
// support kubectl exec , logs , and proxy data flows).
//
// Amazon EKS nodes run in your Amazon Web Services account and connect to your
// cluster's control plane over the Kubernetes API server endpoint and a
// certificate file that is created for your cluster.
//
// You can use the endpointPublicAccess and endpointPrivateAccess parameters to
// enable or disable public and private access to your cluster's Kubernetes API
// server endpoint. By default, public access is enabled, and private access is
// disabled. For more information, see [Amazon EKS Cluster Endpoint Access Control]in the Amazon EKS User Guide .
//
// You can use the logging parameter to enable or disable exporting the Kubernetes
// control plane logs for your cluster to CloudWatch Logs. By default, cluster
// control plane logs aren't exported to CloudWatch Logs. For more information, see
// [Amazon EKS Cluster Control Plane Logs]in the Amazon EKS User Guide .
//
// CloudWatch Logs ingestion, archive storage, and data scanning rates apply to
// exported control plane logs. For more information, see [CloudWatch Pricing].
//
// In most cases, it takes several minutes to create a cluster. After you create
// an Amazon EKS cluster, you must configure your Kubernetes tooling to communicate
// with the API server and launch nodes into your cluster. For more information,
// see [Allowing users to access your cluster]and [Launching Amazon EKS nodes] in the Amazon EKS User Guide.
//
// [Allowing users to access your cluster]: https://docs.aws.amazon.com/eks/latest/userguide/cluster-auth.html
// [CloudWatch Pricing]: http://aws.amazon.com/cloudwatch/pricing/
// [Amazon EKS Cluster Control Plane Logs]: https://docs.aws.amazon.com/eks/latest/userguide/control-plane-logs.html
// [Amazon EKS Cluster Endpoint Access Control]: https://docs.aws.amazon.com/eks/latest/userguide/cluster-endpoint.html
// [Launching Amazon EKS nodes]: https://docs.aws.amazon.com/eks/latest/userguide/launch-workers.html
func (c *Client) CreateCluster(ctx context.Context, params *CreateClusterInput, optFns ...func(*Options)) (*CreateClusterOutput, error) {
	if params == nil {
		params = &CreateClusterInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "CreateCluster", params, optFns, c.addOperationCreateClusterMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*CreateClusterOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type CreateClusterInput struct {

	// The unique name to give to your cluster.
	//
	// This member is required.
	Name *string

	// The VPC configuration that's used by the cluster control plane. Amazon EKS VPC
	// resources have specific requirements to work properly with Kubernetes. For more
	// information, see [Cluster VPC Considerations]and [Cluster Security Group Considerations] in the Amazon EKS User Guide. You must specify at least
	// two subnets. You can specify up to five security groups. However, we recommend
	// that you use a dedicated security group for your cluster control plane.
	//
	// [Cluster Security Group Considerations]: https://docs.aws.amazon.com/eks/latest/userguide/sec-group-reqs.html
	// [Cluster VPC Considerations]: https://docs.aws.amazon.com/eks/latest/userguide/network_reqs.html
	//
	// This member is required.
	ResourcesVpcConfig *types.VpcConfigRequest

	// The Amazon Resource Name (ARN) of the IAM role that provides permissions for
	// the Kubernetes control plane to make calls to Amazon Web Services API operations
	// on your behalf. For more information, see [Amazon EKS Service IAM Role]in the Amazon EKS User Guide .
	//
	// [Amazon EKS Service IAM Role]: https://docs.aws.amazon.com/eks/latest/userguide/service_IAM_role.html
	//
	// This member is required.
	RoleArn *string

	// The access configuration for the cluster.
	AccessConfig *types.CreateAccessConfigRequest

	// If you set this value to False when creating a cluster, the default networking
	// add-ons will not be installed.
	//
	// The default networking addons include vpc-cni, coredns, and kube-proxy.
	//
	// Use this option when you plan to install third-party alternative add-ons or
	// self-manage the default networking add-ons.
	BootstrapSelfManagedAddons *bool

	// A unique, case-sensitive identifier that you provide to ensure the idempotency
	// of the request.
	ClientRequestToken *string

	// The encryption configuration for the cluster.
	EncryptionConfig []types.EncryptionConfig

	// The Kubernetes network configuration for the cluster.
	KubernetesNetworkConfig *types.KubernetesNetworkConfigRequest

	// Enable or disable exporting the Kubernetes control plane logs for your cluster
	// to CloudWatch Logs. By default, cluster control plane logs aren't exported to
	// CloudWatch Logs. For more information, see [Amazon EKS Cluster control plane logs]in the Amazon EKS User Guide .
	//
	// CloudWatch Logs ingestion, archive storage, and data scanning rates apply to
	// exported control plane logs. For more information, see [CloudWatch Pricing].
	//
	// [Amazon EKS Cluster control plane logs]: https://docs.aws.amazon.com/eks/latest/userguide/control-plane-logs.html
	// [CloudWatch Pricing]: http://aws.amazon.com/cloudwatch/pricing/
	Logging *types.Logging

	// An object representing the configuration of your local Amazon EKS cluster on an
	// Amazon Web Services Outpost. Before creating a local cluster on an Outpost,
	// review [Local clusters for Amazon EKS on Amazon Web Services Outposts]in the Amazon EKS User Guide. This object isn't available for creating
	// Amazon EKS clusters on the Amazon Web Services cloud.
	//
	// [Local clusters for Amazon EKS on Amazon Web Services Outposts]: https://docs.aws.amazon.com/eks/latest/userguide/eks-outposts-local-cluster-overview.html
	OutpostConfig *types.OutpostConfigRequest

	// Metadata that assists with categorization and organization. Each tag consists
	// of a key and an optional value. You define both. Tags don't propagate to any
	// other cluster or Amazon Web Services resources.
	Tags map[string]string

	// New clusters, by default, have extended support enabled. You can disable
	// extended support when creating a cluster by setting this value to STANDARD .
	UpgradePolicy *types.UpgradePolicyRequest

	// The desired Kubernetes version for your cluster. If you don't specify a value
	// here, the default version available in Amazon EKS is used.
	//
	// The default version might not be the latest version available.
	Version *string

	noSmithyDocumentSerde
}

type CreateClusterOutput struct {

	// The full description of your new cluster.
	Cluster *types.Cluster

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationCreateClusterMiddlewares(stack *middleware.Stack, options Options) (err error) {
	if err := stack.Serialize.Add(&setOperationInputMiddleware{}, middleware.After); err != nil {
		return err
	}
	err = stack.Serialize.Add(&awsRestjson1_serializeOpCreateCluster{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsRestjson1_deserializeOpCreateCluster{}, middleware.After)
	if err != nil {
		return err
	}
	if err := addProtocolFinalizerMiddlewares(stack, options, "CreateCluster"); err != nil {
		return fmt.Errorf("add protocol finalizers: %v", err)
	}

	if err = addlegacyEndpointContextSetter(stack, options); err != nil {
		return err
	}
	if err = addSetLoggerMiddleware(stack, options); err != nil {
		return err
	}
	if err = addClientRequestID(stack); err != nil {
		return err
	}
	if err = addComputeContentLength(stack); err != nil {
		return err
	}
	if err = addResolveEndpointMiddleware(stack, options); err != nil {
		return err
	}
	if err = addComputePayloadSHA256(stack); err != nil {
		return err
	}
	if err = addRetry(stack, options); err != nil {
		return err
	}
	if err = addRawResponseToMetadata(stack); err != nil {
		return err
	}
	if err = addRecordResponseTiming(stack); err != nil {
		return err
	}
	if err = addClientUserAgent(stack, options); err != nil {
		return err
	}
	if err = smithyhttp.AddErrorCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = addSetLegacyContextSigningOptionsMiddleware(stack); err != nil {
		return err
	}
	if err = addTimeOffsetBuild(stack, c); err != nil {
		return err
	}
	if err = addUserAgentRetryMode(stack, options); err != nil {
		return err
	}
	if err = addIdempotencyToken_opCreateClusterMiddleware(stack, options); err != nil {
		return err
	}
	if err = addOpCreateClusterValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opCreateCluster(options.Region), middleware.Before); err != nil {
		return err
	}
	if err = addRecursionDetection(stack); err != nil {
		return err
	}
	if err = addRequestIDRetrieverMiddleware(stack); err != nil {
		return err
	}
	if err = addResponseErrorMiddleware(stack); err != nil {
		return err
	}
	if err = addRequestResponseLogging(stack, options); err != nil {
		return err
	}
	if err = addDisableHTTPSMiddleware(stack, options); err != nil {
		return err
	}
	return nil
}

type idempotencyToken_initializeOpCreateCluster struct {
	tokenProvider IdempotencyTokenProvider
}

func (*idempotencyToken_initializeOpCreateCluster) ID() string {
	return "OperationIdempotencyTokenAutoFill"
}

func (m *idempotencyToken_initializeOpCreateCluster) HandleInitialize(ctx context.Context, in middleware.InitializeInput, next middleware.InitializeHandler) (
	out middleware.InitializeOutput, metadata middleware.Metadata, err error,
) {
	if m.tokenProvider == nil {
		return next.HandleInitialize(ctx, in)
	}

	input, ok := in.Parameters.(*CreateClusterInput)
	if !ok {
		return out, metadata, fmt.Errorf("expected middleware input to be of type *CreateClusterInput ")
	}

	if input.ClientRequestToken == nil {
		t, err := m.tokenProvider.GetIdempotencyToken()
		if err != nil {
			return out, metadata, err
		}
		input.ClientRequestToken = &t
	}
	return next.HandleInitialize(ctx, in)
}
func addIdempotencyToken_opCreateClusterMiddleware(stack *middleware.Stack, cfg Options) error {
	return stack.Initialize.Add(&idempotencyToken_initializeOpCreateCluster{tokenProvider: cfg.IdempotencyTokenProvider}, middleware.Before)
}

func newServiceMetadataMiddleware_opCreateCluster(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		OperationName: "CreateCluster",
	}
}
