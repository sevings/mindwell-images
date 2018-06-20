// Code generated by go-swagger; DO NOT EDIT.

package me

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"mime/multipart"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
)

// NewPutUsersMeCoverParams creates a new PutUsersMeCoverParams object
// with the default values initialized.
func NewPutUsersMeCoverParams() PutUsersMeCoverParams {
	var ()
	return PutUsersMeCoverParams{}
}

// PutUsersMeCoverParams contains all the bound params for the put users me cover operation
// typically these are obtained from a http.Request
//
// swagger:parameters PutUsersMeCover
type PutUsersMeCoverParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  In: formData
	*/
	File *runtime.File
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls
func (o *PutUsersMeCoverParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error
	o.HTTPRequest = r

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		if err != http.ErrNotMultipart {
			return err
		} else if err := r.ParseForm(); err != nil {
			return err
		}
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil && err != http.ErrMissingFile {
		res = append(res, errors.New(400, "reading file %q failed: %v", "file", err))
	} else if err == http.ErrMissingFile {
		// no-op for missing but optional file parameter
	} else if err := o.bindFile(file, fileHeader); err != nil {
		res = append(res, err)
	} else {
		o.File = &runtime.File{Data: file, Header: fileHeader}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PutUsersMeCoverParams) bindFile(file multipart.File, header *multipart.FileHeader) error {

	return nil
}