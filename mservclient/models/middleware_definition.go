// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// MiddlewareDefinition middleware definition
//
// swagger:model MiddlewareDefinition
type MiddlewareDefinition struct {

	// name
	Name string `json:"name,omitempty"`

	// path
	Path string `json:"path,omitempty"`

	// raw body only
	RawBodyOnly bool `json:"raw_body_only,omitempty"`

	// require session
	RequireSession bool `json:"require_session,omitempty"`
}

// Validate validates this middleware definition
func (m *MiddlewareDefinition) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *MiddlewareDefinition) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *MiddlewareDefinition) UnmarshalBinary(b []byte) error {
	var res MiddlewareDefinition
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
