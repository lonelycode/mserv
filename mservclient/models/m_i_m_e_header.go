// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/strfmt"
)

// MIMEHeader A MIMEHeader represents a MIME-style header mapping
// keys to sets of values.
//
// swagger:model MIMEHeader
type MIMEHeader map[string][]string

// Validate validates this m i m e header
func (m MIMEHeader) Validate(formats strfmt.Registry) error {
	return nil
}
