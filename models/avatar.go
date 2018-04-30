// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// Avatar avatar
// swagger:model Avatar
type Avatar struct {

	// x100
	X100 string `json:"x100,omitempty"`

	// x400
	X400 string `json:"x400,omitempty"`

	// x800
	X800 string `json:"x800,omitempty"`
}

// Validate validates this avatar
func (m *Avatar) Validate(formats strfmt.Registry) error {
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// MarshalBinary interface implementation
func (m *Avatar) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Avatar) UnmarshalBinary(b []byte) error {
	var res Avatar
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
