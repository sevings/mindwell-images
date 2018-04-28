// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// ProfileAllOf1Relations profile all of1 relations
// swagger:model profileAllOf1Relations
type ProfileAllOf1Relations struct {

	// from me
	FromMe string `json:"fromMe,omitempty"`

	// to me
	ToMe string `json:"toMe,omitempty"`
}

// Validate validates this profile all of1 relations
func (m *ProfileAllOf1Relations) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateFromMe(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateToMe(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

var profileAllOf1RelationsTypeFromMePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["followed","requested","ignored","none"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		profileAllOf1RelationsTypeFromMePropEnum = append(profileAllOf1RelationsTypeFromMePropEnum, v)
	}
}

const (
	// ProfileAllOf1RelationsFromMeFollowed captures enum value "followed"
	ProfileAllOf1RelationsFromMeFollowed string = "followed"
	// ProfileAllOf1RelationsFromMeRequested captures enum value "requested"
	ProfileAllOf1RelationsFromMeRequested string = "requested"
	// ProfileAllOf1RelationsFromMeIgnored captures enum value "ignored"
	ProfileAllOf1RelationsFromMeIgnored string = "ignored"
	// ProfileAllOf1RelationsFromMeNone captures enum value "none"
	ProfileAllOf1RelationsFromMeNone string = "none"
)

// prop value enum
func (m *ProfileAllOf1Relations) validateFromMeEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, profileAllOf1RelationsTypeFromMePropEnum); err != nil {
		return err
	}
	return nil
}

func (m *ProfileAllOf1Relations) validateFromMe(formats strfmt.Registry) error {

	if swag.IsZero(m.FromMe) { // not required
		return nil
	}

	// value enum
	if err := m.validateFromMeEnum("fromMe", "body", m.FromMe); err != nil {
		return err
	}

	return nil
}

var profileAllOf1RelationsTypeToMePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["followed","requested","ignored","none"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		profileAllOf1RelationsTypeToMePropEnum = append(profileAllOf1RelationsTypeToMePropEnum, v)
	}
}

const (
	// ProfileAllOf1RelationsToMeFollowed captures enum value "followed"
	ProfileAllOf1RelationsToMeFollowed string = "followed"
	// ProfileAllOf1RelationsToMeRequested captures enum value "requested"
	ProfileAllOf1RelationsToMeRequested string = "requested"
	// ProfileAllOf1RelationsToMeIgnored captures enum value "ignored"
	ProfileAllOf1RelationsToMeIgnored string = "ignored"
	// ProfileAllOf1RelationsToMeNone captures enum value "none"
	ProfileAllOf1RelationsToMeNone string = "none"
)

// prop value enum
func (m *ProfileAllOf1Relations) validateToMeEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, profileAllOf1RelationsTypeToMePropEnum); err != nil {
		return err
	}
	return nil
}

func (m *ProfileAllOf1Relations) validateToMe(formats strfmt.Registry) error {

	if swag.IsZero(m.ToMe) { // not required
		return nil
	}

	// value enum
	if err := m.validateToMeEnum("toMe", "body", m.ToMe); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ProfileAllOf1Relations) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ProfileAllOf1Relations) UnmarshalBinary(b []byte) error {
	var res ProfileAllOf1Relations
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
