// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// AccessSpec access spec
//
// swagger:model AccessSpec
type AccessSpec struct {

	// methods
	Methods []string `json:"methods"`

	// Url
	URL string `json:"url,omitempty"`
}

// Validate validates this access spec
func (m *AccessSpec) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *AccessSpec) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *AccessSpec) UnmarshalBinary(b []byte) error {
	var res AccessSpec
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
