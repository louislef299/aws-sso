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

// Updates an Amazon EKS cluster configuration. Your cluster continues to function
// during the update. The response output includes an update ID that you can use to
// track the status of your cluster update with DescribeUpdate "/>.
//
// You can use this API operation to enable or disable exporting the Kubernetes
// control plane logs for your cluster to CloudWatch Logs. By default, cluster
// control plane logs aren't exported to CloudWatch Logs. For more information, see
// [Amazon EKS Cluster control plane logs]in the Amazon EKS User Guide .
//
// CloudWatch Logs ingestion, archive storage, and data scanning rates apply to
// exported control plane logs. For more information, see [CloudWatch Pricing].
//
// You can also use this API operation to enable or disable public and private
// access to your cluster's Kubernetes API server endpoint. By default, public
// access is enabled, and private access is disabled. For more information, see [Amazon EKS cluster endpoint access control]in
// the Amazon EKS User Guide .
//
// You can also use this API operation to choose different subnets and security
// groups for the cluster. You must specify at least two subnets that are in
// different Availability Zones. You can't change which VPC the subnets are from,
// the subnets must be in the same VPC as the subnets that the cluster was created
// with. For more information about the VPC requirements, see [https://docs.aws.amazon.com/eks/latest/userguide/network_reqs.html]in the Amazon EKS
// User Guide .
//
// Cluster updates are asynchronous, and they should finish within a few minutes.
// During an update, the cluster status moves to UPDATING (this status transition
// is eventually consistent). When the update is complete (either Failed or
// Successful ), the cluster status moves to Active .
//
// [Amazon EKS Cluster control plane logs]: https://docs.aws.amazon.com/eks/latest/userguide/control-plane-logs.html
// [CloudWatch Pricing]: http://aws.amazon.com/cloudwatch/pricing/
// [https://docs.aws.amazon.com/eks/latest/userguide/network_reqs.html]: https://docs.aws.amazon.com/eks/latest/userguide/network_reqs.html
// [Amazon EKS cluster endpoint access control]: https://docs.aws.amazon.com/eks/latest/userguide/cluster-endpoint.html
func (c *Client) UpdateClusterConfig(ctx context.Context, params *UpdateClusterConfigInput, optFns ...func(*Options)) (*UpdateClusterConfigOutput, error) {
	if params == nil {
		params = &UpdateClusterConfigInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "UpdateClusterConfig", params, optFns, c.addOperationUpdateClusterConfigMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*UpdateClusterConfigOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type UpdateClusterConfigInput struct {

	// The name of the Amazon EKS cluster to update.
	//
	// This member is required.
	Name *string

	// The access configuration for the cluster.
	AccessConfig *types.UpdateAccessConfigRequest

	// A unique, case-sensitive identifier that you provide to ensure the idempotency
	// of the request.
	ClientRequestToken *string

	// Enable or disable exporting the Kubernetes control plane logs for your cluster
	// to CloudWatch Logs. By default, cluster control plane logs aren't exported to
	// CloudWatch Logs. For more information, see [Amazon EKS cluster control plane logs]in the Amazon EKS User Guide .
	//
	// CloudWatch Logs ingestion, archive storage, and data scanning rates apply to
	// exported control plane logs. For more information, see [CloudWatch Pricing].
	//
	// [CloudWatch Pricing]: http://aws.amazon.com/cloudwatch/pricing/
	// [Amazon EKS cluster control plane logs]: https://docs.aws.amazon.com/eks/latest/userguide/control-plane-logs.html
	Logging *types.Logging

	// An object representing the VPC configuration to use for an Amazon EKS cluster.
	ResourcesVpcConfig *types.VpcConfigRequest

	noSmithyDocumentSerde
}

type UpdateClusterConfigOutput struct {

	// An object representing an asynchronous update.
	Update *types.Update

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationUpdateClusterConfigMiddlewares(stack *middleware.Stack, options Options) (err error) {
	if err := stack.Serialize.Add(&setOperationInputMiddleware{}, middleware.After); err != nil {
		return err
	}
	err = stack.Serialize.Add(&awsRestjson1_serializeOpUpdateClusterConfig{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsRestjson1_deserializeOpUpdateClusterConfig{}, middleware.After)
	if err != nil {
		return err
	}
	if err := addProtocolFinalizerMiddlewares(stack, options, "UpdateClusterConfig"); err != nil {
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
	if err = addIdempotencyToken_opUpdateClusterConfigMiddleware(stack, options); err != nil {
		return err
	}
	if err = addOpUpdateClusterConfigValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opUpdateClusterConfig(options.Region), middleware.Before); err != nil {
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

type idempotencyToken_initializeOpUpdateClusterConfig struct {
	tokenProvider IdempotencyTokenProvider
}

func (*idempotencyToken_initializeOpUpdateClusterConfig) ID() string {
	return "OperationIdempotencyTokenAutoFill"
}

func (m *idempotencyToken_initializeOpUpdateClusterConfig) HandleInitialize(ctx context.Context, in middleware.InitializeInput, next middleware.InitializeHandler) (
	out middleware.InitializeOutput, metadata middleware.Metadata, err error,
) {
	if m.tokenProvider == nil {
		return next.HandleInitialize(ctx, in)
	}

	input, ok := in.Parameters.(*UpdateClusterConfigInput)
	if !ok {
		return out, metadata, fmt.Errorf("expected middleware input to be of type *UpdateClusterConfigInput ")
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
func addIdempotencyToken_opUpdateClusterConfigMiddleware(stack *middleware.Stack, cfg Options) error {
	return stack.Initialize.Add(&idempotencyToken_initializeOpUpdateClusterConfig{tokenProvider: cfg.IdempotencyTokenProvider}, middleware.Before)
}

func newServiceMetadataMiddleware_opUpdateClusterConfig(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		OperationName: "UpdateClusterConfig",
	}
}
