// Code generated by private/model/cli/gen-api/main.go. DO NOT EDIT.

package lambda

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/awsutil"
	"github.com/aws/aws-sdk-go-v2/private/protocol"
)

type CreateAliasInput struct {
	_ struct{} `type:"structure"`

	// A description of the alias.
	Description *string `type:"string"`

	// The name of the Lambda function.
	//
	// Name formats
	//
	//    * Function name - MyFunction.
	//
	//    * Function ARN - arn:aws:lambda:us-west-2:123456789012:function:MyFunction.
	//
	//    * Partial ARN - 123456789012:function:MyFunction.
	//
	// The length constraint applies only to the full ARN. If you specify only the
	// function name, it is limited to 64 characters in length.
	//
	// FunctionName is a required field
	FunctionName *string `location:"uri" locationName:"FunctionName" min:"1" type:"string" required:"true"`

	// The function version that the alias invokes.
	//
	// FunctionVersion is a required field
	FunctionVersion *string `min:"1" type:"string" required:"true"`

	// The name of the alias.
	//
	// Name is a required field
	Name *string `min:"1" type:"string" required:"true"`

	// The routing configuration (https://docs.aws.amazon.com/lambda/latest/dg/lambda-traffic-shifting-using-aliases.html)
	// of the alias.
	RoutingConfig *AliasRoutingConfiguration `type:"structure"`
}

// String returns the string representation
func (s CreateAliasInput) String() string {
	return awsutil.Prettify(s)
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *CreateAliasInput) Validate() error {
	invalidParams := aws.ErrInvalidParams{Context: "CreateAliasInput"}

	if s.FunctionName == nil {
		invalidParams.Add(aws.NewErrParamRequired("FunctionName"))
	}
	if s.FunctionName != nil && len(*s.FunctionName) < 1 {
		invalidParams.Add(aws.NewErrParamMinLen("FunctionName", 1))
	}

	if s.FunctionVersion == nil {
		invalidParams.Add(aws.NewErrParamRequired("FunctionVersion"))
	}
	if s.FunctionVersion != nil && len(*s.FunctionVersion) < 1 {
		invalidParams.Add(aws.NewErrParamMinLen("FunctionVersion", 1))
	}

	if s.Name == nil {
		invalidParams.Add(aws.NewErrParamRequired("Name"))
	}
	if s.Name != nil && len(*s.Name) < 1 {
		invalidParams.Add(aws.NewErrParamMinLen("Name", 1))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s CreateAliasInput) MarshalFields(e protocol.FieldEncoder) error {
	e.SetValue(protocol.HeaderTarget, "Content-Type", protocol.StringValue("application/json"), protocol.Metadata{})

	if s.Description != nil {
		v := *s.Description

		metadata := protocol.Metadata{}
		e.SetValue(protocol.BodyTarget, "Description", protocol.QuotedValue{ValueMarshaler: protocol.StringValue(v)}, metadata)
	}
	if s.FunctionVersion != nil {
		v := *s.FunctionVersion

		metadata := protocol.Metadata{}
		e.SetValue(protocol.BodyTarget, "FunctionVersion", protocol.QuotedValue{ValueMarshaler: protocol.StringValue(v)}, metadata)
	}
	if s.Name != nil {
		v := *s.Name

		metadata := protocol.Metadata{}
		e.SetValue(protocol.BodyTarget, "Name", protocol.QuotedValue{ValueMarshaler: protocol.StringValue(v)}, metadata)
	}
	if s.RoutingConfig != nil {
		v := s.RoutingConfig

		metadata := protocol.Metadata{}
		e.SetFields(protocol.BodyTarget, "RoutingConfig", v, metadata)
	}
	if s.FunctionName != nil {
		v := *s.FunctionName

		metadata := protocol.Metadata{}
		e.SetValue(protocol.PathTarget, "FunctionName", protocol.QuotedValue{ValueMarshaler: protocol.StringValue(v)}, metadata)
	}
	return nil
}

// Provides configuration information about a Lambda function alias (https://docs.aws.amazon.com/lambda/latest/dg/versioning-aliases.html).
type CreateAliasOutput struct {
	_ struct{} `type:"structure"`

	// The Amazon Resource Name (ARN) of the alias.
	AliasArn *string `type:"string"`

	// A description of the alias.
	Description *string `type:"string"`

	// The function version that the alias invokes.
	FunctionVersion *string `min:"1" type:"string"`

	// The name of the alias.
	Name *string `min:"1" type:"string"`

	// A unique identifier that changes when you update the alias.
	RevisionId *string `type:"string"`

	// The routing configuration (https://docs.aws.amazon.com/lambda/latest/dg/lambda-traffic-shifting-using-aliases.html)
	// of the alias.
	RoutingConfig *AliasRoutingConfiguration `type:"structure"`
}

// String returns the string representation
func (s CreateAliasOutput) String() string {
	return awsutil.Prettify(s)
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s CreateAliasOutput) MarshalFields(e protocol.FieldEncoder) error {
	if s.AliasArn != nil {
		v := *s.AliasArn

		metadata := protocol.Metadata{}
		e.SetValue(protocol.BodyTarget, "AliasArn", protocol.QuotedValue{ValueMarshaler: protocol.StringValue(v)}, metadata)
	}
	if s.Description != nil {
		v := *s.Description

		metadata := protocol.Metadata{}
		e.SetValue(protocol.BodyTarget, "Description", protocol.QuotedValue{ValueMarshaler: protocol.StringValue(v)}, metadata)
	}
	if s.FunctionVersion != nil {
		v := *s.FunctionVersion

		metadata := protocol.Metadata{}
		e.SetValue(protocol.BodyTarget, "FunctionVersion", protocol.QuotedValue{ValueMarshaler: protocol.StringValue(v)}, metadata)
	}
	if s.Name != nil {
		v := *s.Name

		metadata := protocol.Metadata{}
		e.SetValue(protocol.BodyTarget, "Name", protocol.QuotedValue{ValueMarshaler: protocol.StringValue(v)}, metadata)
	}
	if s.RevisionId != nil {
		v := *s.RevisionId

		metadata := protocol.Metadata{}
		e.SetValue(protocol.BodyTarget, "RevisionId", protocol.QuotedValue{ValueMarshaler: protocol.StringValue(v)}, metadata)
	}
	if s.RoutingConfig != nil {
		v := s.RoutingConfig

		metadata := protocol.Metadata{}
		e.SetFields(protocol.BodyTarget, "RoutingConfig", v, metadata)
	}
	return nil
}

const opCreateAlias = "CreateAlias"

// CreateAliasRequest returns a request value for making API operation for
// AWS Lambda.
//
// Creates an alias (https://docs.aws.amazon.com/lambda/latest/dg/versioning-aliases.html)
// for a Lambda function version. Use aliases to provide clients with a function
// identifier that you can update to invoke a different version.
//
// You can also map an alias to split invocation requests between two versions.
// Use the RoutingConfig parameter to specify a second version and the percentage
// of invocation requests that it receives.
//
//    // Example sending a request using CreateAliasRequest.
//    req := client.CreateAliasRequest(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
//
// Please also see https://docs.aws.amazon.com/goto/WebAPI/lambda-2015-03-31/CreateAlias
func (c *Client) CreateAliasRequest(input *CreateAliasInput) CreateAliasRequest {
	op := &aws.Operation{
		Name:       opCreateAlias,
		HTTPMethod: "POST",
		HTTPPath:   "/2015-03-31/functions/{FunctionName}/aliases",
	}

	if input == nil {
		input = &CreateAliasInput{}
	}

	req := c.newRequest(op, input, &CreateAliasOutput{})
	return CreateAliasRequest{Request: req, Input: input, Copy: c.CreateAliasRequest}
}

// CreateAliasRequest is the request type for the
// CreateAlias API operation.
type CreateAliasRequest struct {
	*aws.Request
	Input *CreateAliasInput
	Copy  func(*CreateAliasInput) CreateAliasRequest
}

// Send marshals and sends the CreateAlias API request.
func (r CreateAliasRequest) Send(ctx context.Context) (*CreateAliasResponse, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &CreateAliasResponse{
		CreateAliasOutput: r.Request.Data.(*CreateAliasOutput),
		response:          &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// CreateAliasResponse is the response type for the
// CreateAlias API operation.
type CreateAliasResponse struct {
	*CreateAliasOutput

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// CreateAlias request.
func (r *CreateAliasResponse) SDKResponseMetdata() *aws.Response {
	return r.response
}
