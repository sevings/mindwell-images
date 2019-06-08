// Code generated by go-swagger; DO NOT EDIT.

package images

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"

	models "github.com/sevings/mindwell-images/models"
)

// DeleteImagesIDHandlerFunc turns a function with the right signature into a delete images ID handler
type DeleteImagesIDHandlerFunc func(DeleteImagesIDParams, *models.UserID) middleware.Responder

// Handle executing the request and returning a response
func (fn DeleteImagesIDHandlerFunc) Handle(params DeleteImagesIDParams, principal *models.UserID) middleware.Responder {
	return fn(params, principal)
}

// DeleteImagesIDHandler interface for that can handle valid delete images ID params
type DeleteImagesIDHandler interface {
	Handle(DeleteImagesIDParams, *models.UserID) middleware.Responder
}

// NewDeleteImagesID creates a new http.Handler for the delete images ID operation
func NewDeleteImagesID(ctx *middleware.Context, handler DeleteImagesIDHandler) *DeleteImagesID {
	return &DeleteImagesID{Context: ctx, Handler: handler}
}

/*DeleteImagesID swagger:route DELETE /images/{id} images deleteImagesId

DeleteImagesID delete images ID API

*/
type DeleteImagesID struct {
	Context *middleware.Context
	Handler DeleteImagesIDHandler
}

func (o *DeleteImagesID) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewDeleteImagesIDParams()

	uprinc, aCtx, err := o.Context.Authorize(r, route)
	if err != nil {
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}
	if aCtx != nil {
		r = aCtx
	}
	var principal *models.UserID
	if uprinc != nil {
		principal = uprinc.(*models.UserID) // this is really a models.UserID, I promise
	}

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params, principal) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
