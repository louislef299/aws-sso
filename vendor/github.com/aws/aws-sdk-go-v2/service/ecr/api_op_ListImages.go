// Code generated by smithy-go-codegen DO NOT EDIT.

package ecr

import (
	"context"
	"fmt"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// Lists all the image IDs for the specified repository.
//
// You can filter images based on whether or not they are tagged by using the
// tagStatus filter and specifying either TAGGED , UNTAGGED or ANY . For example,
// you can filter your results to return only UNTAGGED images and then pipe that
// result to a BatchDeleteImageoperation to delete them. Or, you can filter your results to return
// only TAGGED images to list all of the tags in your repository.
func (c *Client) ListImages(ctx context.Context, params *ListImagesInput, optFns ...func(*Options)) (*ListImagesOutput, error) {
	if params == nil {
		params = &ListImagesInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "ListImages", params, optFns, c.addOperationListImagesMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*ListImagesOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type ListImagesInput struct {

	// The repository with image IDs to be listed.
	//
	// This member is required.
	RepositoryName *string

	// The filter key and value with which to filter your ListImages results.
	Filter *types.ListImagesFilter

	// The maximum number of image results returned by ListImages in paginated output.
	// When this parameter is used, ListImages only returns maxResults results in a
	// single page along with a nextToken response element. The remaining results of
	// the initial request can be seen by sending another ListImages request with the
	// returned nextToken value. This value can be between 1 and 1000. If this
	// parameter is not used, then ListImages returns up to 100 results and a nextToken
	// value, if applicable.
	MaxResults *int32

	// The nextToken value returned from a previous paginated ListImages request where
	// maxResults was used and the results exceeded the value of that parameter.
	// Pagination continues from the end of the previous results that returned the
	// nextToken value. This value is null when there are no more results to return.
	//
	// This token should be treated as an opaque identifier that is only used to
	// retrieve the next items in a list and not for other programmatic purposes.
	NextToken *string

	// The Amazon Web Services account ID associated with the registry that contains
	// the repository in which to list images. If you do not specify a registry, the
	// default registry is assumed.
	RegistryId *string

	noSmithyDocumentSerde
}

type ListImagesOutput struct {

	// The list of image IDs for the requested repository.
	ImageIds []types.ImageIdentifier

	// The nextToken value to include in a future ListImages request. When the results
	// of a ListImages request exceed maxResults , this value can be used to retrieve
	// the next page of results. This value is null when there are no more results to
	// return.
	NextToken *string

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationListImagesMiddlewares(stack *middleware.Stack, options Options) (err error) {
	if err := stack.Serialize.Add(&setOperationInputMiddleware{}, middleware.After); err != nil {
		return err
	}
	err = stack.Serialize.Add(&awsAwsjson11_serializeOpListImages{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsAwsjson11_deserializeOpListImages{}, middleware.After)
	if err != nil {
		return err
	}
	if err := addProtocolFinalizerMiddlewares(stack, options, "ListImages"); err != nil {
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
	if err = addOpListImagesValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opListImages(options.Region), middleware.Before); err != nil {
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

// ListImagesAPIClient is a client that implements the ListImages operation.
type ListImagesAPIClient interface {
	ListImages(context.Context, *ListImagesInput, ...func(*Options)) (*ListImagesOutput, error)
}

var _ ListImagesAPIClient = (*Client)(nil)

// ListImagesPaginatorOptions is the paginator options for ListImages
type ListImagesPaginatorOptions struct {
	// The maximum number of image results returned by ListImages in paginated output.
	// When this parameter is used, ListImages only returns maxResults results in a
	// single page along with a nextToken response element. The remaining results of
	// the initial request can be seen by sending another ListImages request with the
	// returned nextToken value. This value can be between 1 and 1000. If this
	// parameter is not used, then ListImages returns up to 100 results and a nextToken
	// value, if applicable.
	Limit int32

	// Set to true if pagination should stop if the service returns a pagination token
	// that matches the most recent token provided to the service.
	StopOnDuplicateToken bool
}

// ListImagesPaginator is a paginator for ListImages
type ListImagesPaginator struct {
	options   ListImagesPaginatorOptions
	client    ListImagesAPIClient
	params    *ListImagesInput
	nextToken *string
	firstPage bool
}

// NewListImagesPaginator returns a new ListImagesPaginator
func NewListImagesPaginator(client ListImagesAPIClient, params *ListImagesInput, optFns ...func(*ListImagesPaginatorOptions)) *ListImagesPaginator {
	if params == nil {
		params = &ListImagesInput{}
	}

	options := ListImagesPaginatorOptions{}
	if params.MaxResults != nil {
		options.Limit = *params.MaxResults
	}

	for _, fn := range optFns {
		fn(&options)
	}

	return &ListImagesPaginator{
		options:   options,
		client:    client,
		params:    params,
		firstPage: true,
		nextToken: params.NextToken,
	}
}

// HasMorePages returns a boolean indicating whether more pages are available
func (p *ListImagesPaginator) HasMorePages() bool {
	return p.firstPage || (p.nextToken != nil && len(*p.nextToken) != 0)
}

// NextPage retrieves the next ListImages page.
func (p *ListImagesPaginator) NextPage(ctx context.Context, optFns ...func(*Options)) (*ListImagesOutput, error) {
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

	result, err := p.client.ListImages(ctx, &params, optFns...)
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

func newServiceMetadataMiddleware_opListImages(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		OperationName: "ListImages",
	}
}
