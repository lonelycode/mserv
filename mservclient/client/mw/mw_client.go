// Code generated by go-swagger; DO NOT EDIT.

package mw

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

// New creates a new mw API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) ClientService {
	return &Client{transport: transport, formats: formats}
}

/*
Client for mw API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

// ClientService is the interface for Client methods
type ClientService interface {
	MwAdd(params *MwAddParams, authInfo runtime.ClientAuthInfoWriter) (*MwAddOK, error)

	MwDelete(params *MwDeleteParams, authInfo runtime.ClientAuthInfoWriter) (*MwDeleteOK, error)

	MwFetch(params *MwFetchParams, authInfo runtime.ClientAuthInfoWriter) (*MwFetchOK, error)

	MwFetchBundle(params *MwFetchBundleParams, authInfo runtime.ClientAuthInfoWriter) (*MwFetchBundleOK, error)

	MwListAll(params *MwListAllParams, authInfo runtime.ClientAuthInfoWriter) (*MwListAllOK, error)

	MwUpdate(params *MwUpdateParams, authInfo runtime.ClientAuthInfoWriter) (*MwUpdateOK, error)

	SetTransport(transport runtime.ClientTransport)
}

/*
  MwAdd adds a new middleware if store only field is true it will only be available for download

  Expects a file bundle in `uploadfile` form field.
*/
func (a *Client) MwAdd(params *MwAddParams, authInfo runtime.ClientAuthInfoWriter) (*MwAddOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewMwAddParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "mwAdd",
		Method:             "POST",
		PathPattern:        "/api/mw",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &MwAddReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*MwAddOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for mwAdd: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
  MwDelete deletes a middleware specified by id
*/
func (a *Client) MwDelete(params *MwDeleteParams, authInfo runtime.ClientAuthInfoWriter) (*MwDeleteOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewMwDeleteParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "mwDelete",
		Method:             "DELETE",
		PathPattern:        "/api/mw/{id}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &MwDeleteReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*MwDeleteOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for mwDelete: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
  MwFetch fetches a middleware specified by id
*/
func (a *Client) MwFetch(params *MwFetchParams, authInfo runtime.ClientAuthInfoWriter) (*MwFetchOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewMwFetchParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "mwFetch",
		Method:             "GET",
		PathPattern:        "/api/mw/{id}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &MwFetchReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*MwFetchOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for mwFetch: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
  MwFetchBundle fetches a middleware bundle file specified by id
*/
func (a *Client) MwFetchBundle(params *MwFetchBundleParams, authInfo runtime.ClientAuthInfoWriter) (*MwFetchBundleOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewMwFetchBundleParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "mwFetchBundle",
		Method:             "GET",
		PathPattern:        "/api/mw/bundle/{id}",
		ProducesMediaTypes: []string{"application/octet-stream"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &MwFetchBundleReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*MwFetchBundleOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for mwFetchBundle: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
  MwListAll lists all middleware
*/
func (a *Client) MwListAll(params *MwListAllParams, authInfo runtime.ClientAuthInfoWriter) (*MwListAllOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewMwListAllParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "mwListAll",
		Method:             "GET",
		PathPattern:        "/api/mw/master/all",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &MwListAllReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*MwListAllOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for mwListAll: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
  MwUpdate updates a middleware specified by id

  Expects a file bundle in `uploadfile` form field.
*/
func (a *Client) MwUpdate(params *MwUpdateParams, authInfo runtime.ClientAuthInfoWriter) (*MwUpdateOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewMwUpdateParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "mwUpdate",
		Method:             "PUT",
		PathPattern:        "/api/mw/{id}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &MwUpdateReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*MwUpdateOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for mwUpdate: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}
