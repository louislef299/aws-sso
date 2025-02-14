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

// List the EKS Pod Identity associations in a cluster. You can filter the list by
// the namespace that the association is in or the service account that the
// association uses.
func (c *Client) ListPodIdentityAssociations(ctx context.Context, params *ListPodIdentityAssociationsInput, optFns ...func(*Options)) (*ListPodIdentityAssociationsOutput, error) {
	if params == nil {
		params = &ListPodIdentityAssociationsInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "ListPodIdentityAssociations", params, optFns, c.addOperationListPodIdentityAssociationsMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*ListPodIdentityAssociationsOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type ListPodIdentityAssociationsInput struct {

	// The name of the cluster that the associations are in.
	//
	// This member is required.
	ClusterName *string

	// The maximum number of EKS Pod Identity association results returned by
	// ListPodIdentityAssociations in paginated output. When you use this parameter,
	// ListPodIdentityAssociations returns only maxResults results in a single page
	// along with a nextToken response element. You can see the remaining results of
	// the initial request by sending another ListPodIdentityAssociations request with
	// the returned nextToken value. This value can be between 1 and 100. If you don't
	// use this parameter, ListPodIdentityAssociations returns up to 100 results and a
	// nextToken value if applicable.
	MaxResults *int32

	// The name of the Kubernetes namespace inside the cluster that the associations
	// are in.
	Namespace *string

	// The nextToken value returned from a previous paginated ListUpdates request
	// where maxResults was used and the results exceeded the value of that parameter.
	// Pagination continues from the end of the previous results that returned the
	// nextToken value.
	//
	// This token should be treated as an opaque identifier that is used only to
	// retrieve the next items in a list and not for other programmatic purposes.
	NextToken *string

	// The name of the Kubernetes service account that the associations use.
	ServiceAccount *string

	noSmithyDocumentSerde
}

type ListPodIdentityAssociationsOutput struct {

	// The list of summarized descriptions of the associations that are in the cluster
	// and match any filters that you provided.
	//
	// Each summary is simplified by removing these fields compared to the full [PodIdentityAssociation]
	// PodIdentityAssociation :
	//
	//   - The IAM role: roleArn
	//
	//   - The timestamp that the association was created at: createdAt
	//
	//   - The most recent timestamp that the association was modified at:. modifiedAt
	//
	//   - The tags on the association: tags
	//
	// [PodIdentityAssociation]: https://docs.aws.amazon.com/eks/latest/APIReference/API_PodIdentityAssociation.html
	Associations []types.PodIdentityAssociationSummary

	// The nextToken value to include in a future ListPodIdentityAssociations request.
	// When the results of a ListPodIdentityAssociations request exceed maxResults ,
	// you can use this value to retrieve the next page of results. This value is null
	// when there are no more results to return.
	//
	// This token should be treated as an opaque identifier that is used only to
	// retrieve the next items in a list and not for other programmatic purposes.
	NextToken *string

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationListPodIdentityAssociationsMiddlewares(stack *middleware.Stack, options Options) (err error) {
	if err := stack.Serialize.Add(&setOperationInputMiddleware{}, middleware.After); err != nil {
		return err
	}
	err = stack.Serialize.Add(&awsRestjson1_serializeOpListPodIdentityAssociations{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsRestjson1_deserializeOpListPodIdentityAssociations{}, middleware.After)
	if err != nil {
		return err
	}
	if err := addProtocolFinalizerMiddlewares(stack, options, "ListPodIdentityAssociations"); err != nil {
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
	if err = addOpListPodIdentityAssociationsValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opListPodIdentityAssociations(options.Region), middleware.Before); err != nil {
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

// ListPodIdentityAssociationsPaginatorOptions is the paginator options for
// ListPodIdentityAssociations
type ListPodIdentityAssociationsPaginatorOptions struct {
	// The maximum number of EKS Pod Identity association results returned by
	// ListPodIdentityAssociations in paginated output. When you use this parameter,
	// ListPodIdentityAssociations returns only maxResults results in a single page
	// along with a nextToken response element. You can see the remaining results of
	// the initial request by sending another ListPodIdentityAssociations request with
	// the returned nextToken value. This value can be between 1 and 100. If you don't
	// use this parameter, ListPodIdentityAssociations returns up to 100 results and a
	// nextToken value if applicable.
	Limit int32

	// Set to true if pagination should stop if the service returns a pagination token
	// that matches the most recent token provided to the service.
	StopOnDuplicateToken bool
}

// ListPodIdentityAssociationsPaginator is a paginator for
// ListPodIdentityAssociations
type ListPodIdentityAssociationsPaginator struct {
	options   ListPodIdentityAssociationsPaginatorOptions
	client    ListPodIdentityAssociationsAPIClient
	params    *ListPodIdentityAssociationsInput
	nextToken *string
	firstPage bool
}

// NewListPodIdentityAssociationsPaginator returns a new
// ListPodIdentityAssociationsPaginator
func NewListPodIdentityAssociationsPaginator(client ListPodIdentityAssociationsAPIClient, params *ListPodIdentityAssociationsInput, optFns ...func(*ListPodIdentityAssociationsPaginatorOptions)) *ListPodIdentityAssociationsPaginator {
	if params == nil {
		params = &ListPodIdentityAssociationsInput{}
	}

	options := ListPodIdentityAssociationsPaginatorOptions{}
	if params.MaxResults != nil {
		options.Limit = *params.MaxResults
	}

	for _, fn := range optFns {
		fn(&options)
	}

	return &ListPodIdentityAssociationsPaginator{
		options:   options,
		client:    client,
		params:    params,
		firstPage: true,
		nextToken: params.NextToken,
	}
}

// HasMorePages returns a boolean indicating whether more pages are available
func (p *ListPodIdentityAssociationsPaginator) HasMorePages() bool {
	return p.firstPage || (p.nextToken != nil && len(*p.nextToken) != 0)
}

// NextPage retrieves the next ListPodIdentityAssociations page.
func (p *ListPodIdentityAssociationsPaginator) NextPage(ctx context.Context, optFns ...func(*Options)) (*ListPodIdentityAssociationsOutput, error) {
	if !p.HasMorePages() {
		return nil, fmt.Errorf("no more pages available")
	}

	params := *p.params
	params.NextToken = p.nextToken

	var limit *int32
	if p.options.Limit > 0 {
		limit = &p.options.Limit
	}
	params.MaxResults = limit

	optFns = append([]func(*Options){
		addIsPaginatorUserAgent,
	}, optFns...)
	result, err := p.client.ListPodIdentityAssociations(ctx, &params, optFns...)
	if err != nil {
		return nil, err
	}
	p.firstPage = false

	prevToken := p.nextToken
	p.nextToken = result.NextToken

	if p.options.StopOnDuplicateToken &&
		prevToken != nil &&
		p.nextToken != nil &&
		*prevToken == *p.nextToken {
		p.nextToken = nil
	}

	return result, nil
}

// ListPodIdentityAssociationsAPIClient is a client that implements the
// ListPodIdentityAssociations operation.
type ListPodIdentityAssociationsAPIClient interface {
	ListPodIdentityAssociations(context.Context, *ListPodIdentityAssociationsInput, ...func(*Options)) (*ListPodIdentityAssociationsOutput, error)
}

var _ ListPodIdentityAssociationsAPIClient = (*Client)(nil)

func newServiceMetadataMiddleware_opListPodIdentityAssociations(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		OperationName: "ListPodIdentityAssociations",
	}
}
