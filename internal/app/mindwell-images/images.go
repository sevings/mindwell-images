package images

import (
	"database/sql"
	"log"

	"github.com/go-openapi/runtime/middleware"
	"github.com/sevings/mindwell-images/models"
	"github.com/sevings/mindwell-images/restapi/operations/images"
	"github.com/sevings/mindwell-server/utils"
	goconf "github.com/zpatrick/go-config"
)

func saveImageSize(tx *utils.AutoTx, imageID, width, height int64, size string) {
	const q = `
		INSERT INTO image_sizes(image_id, size, width, height)
		VALUES($1, (SELECT id FROM size WHERE type = $2), $3, $4)
	`

	tx.Exec(q, imageID, size, width, height)
}

func NewImageUploader(db *sql.DB, cfg *goconf.Config) func(images.PostImagesParams, *models.UserID) middleware.Responder {
	return func(params images.PostImagesParams, userID *models.UserID) middleware.Responder {
		store := newImageStore(cfg)
		store.ReadImage(params.File)

		img := &models.Image{
			Author: &models.User{
				ID:   userID.ID,
				Name: userID.Name,
			},
			Type:      store.FileExtension(),
			Thumbnail: store.Fill(100, "albums/thumbnails"),
			Small:     store.FitRect(480, 360, "albums/small"),
			Medium:    store.FitRect(800, 600, "albums/medium"),
			Large:     store.FitRect(1280, 960, "albums/large"),
		}

		if store.Error() != nil {
			log.Print(store.Error())
			return images.NewPostImagesBadRequest()
		}

		return utils.Transact(db, func(tx *utils.AutoTx) middleware.Responder {
			tx.Query("INSERT INTO images(user_id, path) VALUES($1, $2) RETURNING id", userID.ID, store.FileName())
			tx.Scan(&img.ID)

			saveImageSize(tx, img.ID, img.Thumbnail.Width, img.Thumbnail.Height, "thumbnail")
			saveImageSize(tx, img.ID, img.Small.Width, img.Small.Height, "small")
			saveImageSize(tx, img.ID, img.Medium.Width, img.Medium.Height, "medium")
			saveImageSize(tx, img.ID, img.Large.Width, img.Large.Height, "large")

			if tx.Error() != nil {
				return images.NewPostImagesBadRequest()
			}

			return images.NewPostImagesOK().WithPayload(img)
		})
	}
}
