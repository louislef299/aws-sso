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

// Updates the Kubernetes version or AMI version of an Amazon EKS managed node
// group.
//
// You can update a node group using a launch template only if the node group was
// originally deployed with a launch template. If you need to update a custom AMI
// in a node group that was deployed with a launch template, then update your
// custom AMI, specify the new ID in a new version of the launch template, and then
// update the node group to the new version of the launch template.
//
// If you update without a launch template, then you can update to the latest
// available AMI version of a node group's current Kubernetes version by not
// specifying a Kubernetes version in the request. You can update to the latest AMI
// version of your cluster's current Kubernetes version by specifying your
// cluster's Kubernetes version in the request. For information about Linux
// versions, see [Amazon EKS optimized Amazon Linux AMI versions]in the Amazon EKS User Guide. For information about Windows
// versions, see [Amazon EKS optimized Windows AMI versions]in the Amazon EKS User Guide.
//
// You cannot roll back a node group to an earlier Kubernetes version or AMI
// version.
//
// When a node in a managed node group is terminated due to a scaling action or
// update, every Pod on that node is drained first. Amazon EKS attempts to drain
// the nodes gracefully and will fail if it is unable to do so. You can force the
// update if Amazon EKS is unable to drain the nodes as a result of a Pod
// disruption budget issue.
//
// [Amazon EKS optimized Amazon Linux AMI versions]: https://docs.aws.amazon.com/eks/latest/userguide/eks-linux-ami-versions.html
// [Amazon EKS optimized Windows AMI versions]: https://docs.aws.amazon.com/eks/latest/userguide/eks-ami-versions-windows.html
func (c *Client) UpdateNodegroupVersion(ctx context.Context, params *UpdateNodegroupVersionInput, optFns ...func(*Options)) (*UpdateNodegroupVersionOutput, error) {
	if params == nil {
		params = &UpdateNodegroupVersionInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "UpdateNodegroupVersion", params, optFns, c.addOperationUpdateNodegroupVersionMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*UpdateNodegroupVersionOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type UpdateNodegroupVersionInput struct {

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

	// Force the update if any Pod on the existing node group can't be drained due to
	// a Pod disruption budget issue. If an update fails because all Pods can't be
	// drained, you can force the update after it fails to terminate the old node
	// whether or not any Pod is running on the node.
	Force bool

	// An object representing a node group's launch template specification. You can
	// only update a node group using a launch template if the node group was
	// originally deployed with a launch template.
	LaunchTemplate *types.LaunchTemplateSpecification

	// The AMI version of the Amazon EKS optimized AMI to use for the update. By
	// default, the latest available AMI version for the node group's Kubernetes
	// version is used. For information about Linux versions, see [Amazon EKS optimized Amazon Linux AMI versions]in the Amazon EKS
	// User Guide. Amazon EKS managed node groups support the November 2022 and later
	// releases of the Windows AMIs. For information about Windows versions, see [Amazon EKS optimized Windows AMI versions]in
	// the Amazon EKS User Guide.
	//
	// If you specify launchTemplate , and your launch template uses a custom AMI, then
	// don't specify releaseVersion , or the node group update will fail. For more
	// information about using launch templates with Amazon EKS, see [Customizing managed nodes with launch templates]in the Amazon EKS
	// User Guide.
	//
	// [Customizing managed nodes with launch templates]: https://docs.aws.amazon.com/eks/latest/userguide/launch-templates.html
	// [Amazon EKS optimized Amazon Linux AMI versions]: https://docs.aws.amazon.com/eks/latest/userguide/eks-linux-ami-versions.html
	// [Amazon EKS optimized Windows AMI versions]: https://docs.aws.amazon.com/eks/latest/userguide/eks-ami-versions-windows.html
	ReleaseVersion *string

	// The Kubernetes version to update to. If no version is specified, then the
	// Kubernetes version of the node group does not change. You can specify the
	// Kubernetes version of the cluster to update the node group to the latest AMI
	// version of the cluster's Kubernetes version. If you specify launchTemplate , and
	// your launch template uses a custom AMI, then don't specify version , or the node
	// group update will fail. For more information about using launch templates with
	// Amazon EKS, see [Customizing managed nodes with launch templates]in the Amazon EKS User Guide.
	//
	// [Customizing managed nodes with launch templates]: https://docs.aws.amazon.com/eks/latest/userguide/launch-templates.html
	Version *string

	noSmithyDocumentSerde
}

type UpdateNodegroupVersionOutput struct {

	// An object representing an asynchronous update.
	Update *types.Update

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationUpdateNodegroupVersionMiddlewares(stack *middleware.Stack, options Options) (err error) {
	if err := stack.Serialize.Add(&setOperationInputMiddleware{}, middleware.After); err != nil {
		return err
	}
	err = stack.Serialize.Add(&awsRestjson1_serializeOpUpdateNodegroupVersion{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsRestjson1_deserializeOpUpdateNodegroupVersion{}, middleware.After)
	if err != nil {
		return err
	}
	if err := addProtocolFinalizerMiddlewares(stack, options, "UpdateNodegroupVersion"); err != nil {
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
	if err = addIdempotencyToken_opUpdateNodegroupVersionMiddleware(stack, options); err != nil {
		return err
	}
	if err = addOpUpdateNodegroupVersionValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opUpdateNodegroupVersion(options.Region), middleware.Before); err != nil {
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

type idempotencyToken_initializeOpUpdateNodegroupVersion struct {
	tokenProvider IdempotencyTokenProvider
}

func (*idempotencyToken_initializeOpUpdateNodegroupVersion) ID() string {
	return "OperationIdempotencyTokenAutoFill"
}

func (m *idempotencyToken_initializeOpUpdateNodegroupVersion) HandleInitialize(ctx context.Context, in middleware.InitializeInput, next middleware.InitializeHandler) (
	out middleware.InitializeOutput, metadata middleware.Metadata, err error,
) {
	if m.tokenProvider == nil {
		return next.HandleInitialize(ctx, in)
	}

	input, ok := in.Parameters.(*UpdateNodegroupVersionInput)
	if !ok {
		return out, metadata, fmt.Errorf("expected middleware input to be of type *UpdateNodegroupVersionInput ")
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
func addIdempotencyToken_opUpdateNodegroupVersionMiddleware(stack *middleware.Stack, cfg Options) error {
	return stack.Initialize.Add(&idempotencyToken_initializeOpUpdateNodegroupVersion{tokenProvider: cfg.IdempotencyTokenProvider}, middleware.Before)
}

func newServiceMetadataMiddleware_opUpdateNodegroupVersion(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		OperationName: "UpdateNodegroupVersion",
	}
}
