// Code generated by go-swagger; DO NOT EDIT.

package images

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/sevings/mindwell-images/models"
)

// PostImagesOKCode is the HTTP code returned for type PostImagesOK
const PostImagesOKCode int = 200

/*PostImagesOK Image

swagger:response postImagesOK
*/
type PostImagesOK struct {

	/*
	  In: Body
	*/
	Payload *models.Image `json:"body,omitempty"`
}

// NewPostImagesOK creates PostImagesOK with default headers values
func NewPostImagesOK() *PostImagesOK {

	return &PostImagesOK{}
}

// WithPayload adds the payload to the post images o k response
func (o *PostImagesOK) WithPayload(payload *models.Image) *PostImagesOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post images o k response
func (o *PostImagesOK) SetPayload(payload *models.Image) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostImagesOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PostImagesBadRequestCode is the HTTP code returned for type PostImagesBadRequest
const PostImagesBadRequestCode int = 400

/*PostImagesBadRequest bad request

swagger:response postImagesBadRequest
*/
type PostImagesBadRequest struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPostImagesBadRequest creates PostImagesBadRequest with default headers values
func NewPostImagesBadRequest() *PostImagesBadRequest {

	return &PostImagesBadRequest{}
}

// WithPayload adds the payload to the post images bad request response
func (o *PostImagesBadRequest) WithPayload(payload *models.Error) *PostImagesBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post images bad request response
func (o *PostImagesBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostImagesBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
