// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// MW m w
//
// swagger:model MW
type MW struct {

	// API ID
	APIID string `json:"APIID,omitempty"`

	// active
	Active bool `json:"Active,omitempty"`

	// added
	// Format: date-time
	Added strfmt.DateTime `json:"Added,omitempty"`

	// bundle ref
	BundleRef string `json:"BundleRef,omitempty"`

	// download only
	DownloadOnly bool `json:"DownloadOnly,omitempty"`

	// manifest
	Manifest *BundleManifest `json:"Manifest,omitempty"`

	// org ID
	OrgID string `json:"OrgID,omitempty"`

	// plugins
	Plugins []*Plugin `json:"Plugins"`

	// UID
	UID string `json:"UID,omitempty"`
}

// Validate validates this m w
func (m *MW) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAdded(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateManifest(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePlugins(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *MW) validateAdded(formats strfmt.Registry) error {

	if swag.IsZero(m.Added) { // not required
		return nil
	}

	if err := validate.FormatOf("Added", "body", "date-time", m.Added.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *MW) validateManifest(formats strfmt.Registry) error {

	if swag.IsZero(m.Manifest) { // not required
		return nil
	}

	if m.Manifest != nil {
		if err := m.Manifest.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("Manifest")
			}
			return err
		}
	}

	return nil
}

func (m *MW) validatePlugins(formats strfmt.Registry) error {

	if swag.IsZero(m.Plugins) { // not required
		return nil
	}

	for i := 0; i < len(m.Plugins); i++ {
		if swag.IsZero(m.Plugins[i]) { // not required
			continue
		}

		if m.Plugins[i] != nil {
			if err := m.Plugins[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("Plugins" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *MW) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *MW) UnmarshalBinary(b []byte) error {
	var res MW
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}