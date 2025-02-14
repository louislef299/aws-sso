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

// Updates an Amazon EKS managed node group configuration. Your node group
// continues to function during the update. The response output includes an update
// ID that you can use to track the status of your node group update with the [DescribeUpdate]
// DescribeUpdate API operation. You can update the Kubernetes labels and taints
// for a node group and the scaling and version update configuration.
//
// [DescribeUpdate]: https://docs.aws.amazon.com/eks/latest/APIReference/API_DescribeUpdate.html
func (c *Client) UpdateNodegroupConfig(ctx context.Context, params *UpdateNodegroupConfigInput, optFns ...func(*Options)) (*UpdateNodegroupConfigOutput, error) {
	if params == nil {
		params = &UpdateNodegroupConfigInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "UpdateNodegroupConfig", params, optFns, c.addOperationUpdateNodegroupConfigMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*UpdateNodegroupConfigOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type UpdateNodegroupConfigInput struct {

	// The name of your cluster.
	//
	// This member is required.
	ClusterName *string

	// The name of the managed node group to update.
	//
	// This member is required.
	NodegroupName *string

	// A unique, case-sensitive identifier that you provide to ensure the idempotency
	// of the request.
	ClientRequestToken *string

	// The Kubernetes labels to apply to the nodes in the node group after the update.
	Labels *types.UpdateLabelsPayload

	// The node auto repair configuration for the node group.
	NodeRepairConfig *types.NodeRepairConfig

	// The scaling configuration details for the Auto Scaling group after the update.
	ScalingConfig *types.NodegroupScalingConfig

	// The Kubernetes taints to be applied to the nodes in the node group after the
	// update. For more information, see [Node taints on managed node groups].
	//
	// [Node taints on managed node groups]: https://docs.aws.amazon.com/eks/latest/userguide/node-taints-managed-node-groups.html
	Taints *types.UpdateTaintsPayload

	// The node group update configuration.
	UpdateConfig *types.NodegroupUpdateConfig

	noSmithyDocumentSerde
}

type UpdateNodegroupConfigOutput struct {

	// An object representing an asynchronous update.
	Update *types.Update

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationUpdateNodegroupConfigMiddlewares(stack *middleware.Stack, options Options) (err error) {
	if err := stack.Serialize.Add(&setOperationInputMiddleware{}, middleware.After); err != nil {
		return err
	}
	err = stack.Serialize.Add(&awsRestjson1_serializeOpUpdateNodegroupConfig{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsRestjson1_deserializeOpUpdateNodegroupConfig{}, middleware.After)
	if err != nil {
		return err
	}
	if err := addProtocolFinalizerMiddlewares(stack, options, "UpdateNodegroupConfig"); err != nil {
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
	if err = addSpanRetryLoop(stack, options); err != nil {
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
	if err = addIdempotencyToken_opUpdateNodegroupConfigMiddleware(stack, options); err != nil {
		return err
	}
	if err = addOpUpdateNodegroupConfigValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opUpdateNodegroupConfig(options.Region), middleware.Before); err != nil {
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
	if err = addSpanInitializeStart(stack); err != nil {
		return err
	}
	if err = addSpanInitializeEnd(stack); err != nil {
		return err
	}
	if err = addSpanBuildRequestStart(stack); err != nil {
		return err
	}
	if err = addSpanBuildRequestEnd(stack); err != nil {
		return err
	}
	return nil
}

type idempotencyToken_initializeOpUpdateNodegroupConfig struct {
	tokenProvider IdempotencyTokenProvider
}

func (*idempotencyToken_initializeOpUpdateNodegroupConfig) ID() string {
	return "OperationIdempotencyTokenAutoFill"
}

func (m *idempotencyToken_initializeOpUpdateNodegroupConfig) HandleInitialize(ctx context.Context, in middleware.InitializeInput, next middleware.InitializeHandler) (
	out middleware.InitializeOutput, metadata middleware.Metadata, err error,
) {
	if m.tokenProvider == nil {
		return next.HandleInitialize(ctx, in)
	}

	input, ok := in.Parameters.(*UpdateNodegroupConfigInput)
	if !ok {
		return out, metadata, fmt.Errorf("expected middleware input to be of type *UpdateNodegroupConfigInput ")
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
func addIdempotencyToken_opUpdateNodegroupConfigMiddleware(stack *middleware.Stack, cfg Options) error {
	return stack.Initialize.Add(&idempotencyToken_initializeOpUpdateNodegroupConfig{tokenProvider: cfg.IdempotencyTokenProvider}, middleware.Before)
}

func newServiceMetadataMiddleware_opUpdateNodegroupConfig(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		OperationName: "UpdateNodegroupConfig",
	}
}
