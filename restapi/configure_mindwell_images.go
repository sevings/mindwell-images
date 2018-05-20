package restapi

import (
	"crypto/tls"
	"math/rand"
	"net/http"
	"time"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	graceful "github.com/tylerb/graceful"
	"gopkg.in/gographics/imagick.v2/imagick"

	"github.com/sevings/mindwell-server/utils"

	"github.com/sevings/mindwell-images/internal/app/mindwell-images"
	"github.com/sevings/mindwell-images/models"
	"github.com/sevings/mindwell-images/restapi/operations"
	"github.com/sevings/mindwell-images/restapi/operations/me"
)

// This file is safe to edit. Once it exists it will not be overwritten

//go:generate swagger generate server --target .. --name mindwell-images --spec ../../mindwell-server/web/swagger.yaml --operation PutUsersMeAvatar --principal models.UserID

func configureFlags(api *operations.MindwellImagesAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.MindwellImagesAPI) http.Handler {
	rand.Seed(time.Now().UTC().UnixNano())

	config := utils.LoadConfig("configs/images")
	db := utils.OpenDatabase(config)

	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.UrlformConsumer = runtime.DiscardConsumer

	api.MultipartformConsumer = runtime.DiscardConsumer

	api.JSONProducer = runtime.JSONProducer()

	// Applies when the "X-User-Key" header is set
	keyAuth := utils.NewKeyAuth(db)
	api.APIKeyHeaderAuth = func(apiKey string) (*models.UserID, error) {
		id, err := keyAuth(apiKey)
		userID := models.UserID(*id)
		return &userID, err
	}

	// Set your custom authorizer if needed. Default one is security.Authorized()
	// Expected interface runtime.Authorizer
	//
	// Example:
	// api.APIAuthorizer = security.Authorized()

	api.MePutUsersMeAvatarHandler = me.PutUsersMeAvatarHandlerFunc(images.NewAvatarUpdater(db, config))

	api.ServerShutdown = func() {
		imagick.Terminate()
	}

	imagick.Initialize()

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *graceful.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
