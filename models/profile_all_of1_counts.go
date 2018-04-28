// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// ProfileAllOf1Counts profile all of1 counts
// swagger:model profileAllOf1Counts
type ProfileAllOf1Counts struct {

	// comments
	Comments int64 `json:"comments,omitempty"`

	// entries
	Entries int64 `json:"entries,omitempty"`

	// favorites
	Favorites int64 `json:"favorites,omitempty"`

	// followers
	Followers int64 `json:"followers,omitempty"`

	// followings
	Followings int64 `json:"followings,omitempty"`

	// ignored
	Ignored int64 `json:"ignored,omitempty"`

	// invited
	Invited int64 `json:"invited,omitempty"`

	// tags
	Tags int64 `json:"tags,omitempty"`
}

// Validate validates this profile all of1 counts
func (m *ProfileAllOf1Counts) Validate(formats strfmt.Registry) error {
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// MarshalBinary interface implementation
func (m *ProfileAllOf1Counts) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ProfileAllOf1Counts) UnmarshalBinary(b []byte) error {
	var res ProfileAllOf1Counts
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
