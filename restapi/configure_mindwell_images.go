package restapi

import (
	"crypto/tls"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	"github.com/unrolled/logger"
	"gopkg.in/gographics/imagick.v2/imagick"

	imagesImpl "github.com/sevings/mindwell-images/internal/app/mindwell-images"
	"github.com/sevings/mindwell-images/models"
	"github.com/sevings/mindwell-images/restapi/operations"
	"github.com/sevings/mindwell-images/restapi/operations/images"
	"github.com/sevings/mindwell-images/restapi/operations/me"
	"github.com/sevings/mindwell-server/utils"
)

// This file is safe to edit. Once it exists it will not be overwritten

//go:generate swagger generate server --target .. --name mindwell-images --spec ../../mindwell-server/web/swagger.yaml --operation PutMeAvatar --operation PutMeCover --principal models.UserID --model Avatar --model Cover --model UserID --model Error

func configureFlags(api *operations.MindwellImagesAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.MindwellImagesAPI) http.Handler {
	rand.Seed(time.Now().UTC().UnixNano())

	config := utils.LoadConfig("configs/images")
	mi := imagesImpl.NewMindwellImages(config)

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

	apiSecret := mi.ConfigString("server.api_secret")

	// Applies when the "X-User-Key" header is set
	keyAuth := utils.NewKeyAuth(mi.DB(), []byte(apiSecret))
	api.APIKeyHeaderAuth = func(apiKey string) (*models.UserID, error) {
		id, err := keyAuth(apiKey)
		userID := models.UserID{
			ID:             id.ID,
			Name:           id.Name,
			FollowersCount: id.FollowersCount,
			IsInvited:      id.IsInvited,
			NegKarma:       id.NegKarma,
			Ban: &models.UserIDBan{
				Comment: id.Ban.Comment,
				Invite:  id.Ban.Invite,
				Live:    id.Ban.Live,
				Vote:    id.Ban.Vote,
			},
		}
		return &userID, err
	}

	// Set your custom authorizer if needed. Default one is security.Authorized()
	// Expected interface runtime.Authorizer
	//
	// Example:
	// api.APIAuthorizer = security.Authorized()

	api.MePutMeAvatarHandler = me.PutMeAvatarHandlerFunc(imagesImpl.NewAvatarUpdater(mi))
	api.MePutMeCoverHandler = me.PutMeCoverHandlerFunc(imagesImpl.NewCoverUpdater(mi))

	api.ImagesPostImagesHandler = images.PostImagesHandlerFunc(imagesImpl.NewImageUploader(mi))
	api.ImagesGetImagesIDHandler = images.GetImagesIDHandlerFunc(imagesImpl.NewImageLoader(mi))
	api.ImagesDeleteImagesIDHandler = images.DeleteImagesIDHandlerFunc(imagesImpl.NewImageDeleter(mi))

	api.ServerShutdown = func() {
		mi.Shutdown()
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
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	lmt := tollbooth.NewLimiter(1, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})
	lmt.SetIPLookups([]string{"X-Forwarded-For"})
	lmt.SetMessage(`{"message":"You have reached maximum request limit."}`)
	lmt.SetMessageContentType("application/json")

	return tollbooth.LimitHandler(lmt, handler)
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	log := logger.New(logger.Options{
		RemoteAddressHeaders: []string{"X-Forwarded-For"},
		Out:                  os.Stdout,
	})

	return log.Handler(handler)
}
