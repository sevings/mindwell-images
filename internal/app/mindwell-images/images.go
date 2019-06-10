package images

import (
	"log"

	"github.com/go-openapi/runtime/middleware"
	"github.com/sevings/mindwell-images/models"
	"github.com/sevings/mindwell-images/restapi/operations/images"
	"github.com/sevings/mindwell-server/utils"
)

func saveImageSize(tx *utils.AutoTx, imageID, width, height int64, size string) {
	const q = `
		INSERT INTO image_sizes(image_id, size, width, height)
		VALUES($1, (SELECT id FROM size WHERE type = $2), $3, $4)
	`

	tx.Exec(q, imageID, size, width, height)
}

func NewImageUploader(mi *MindwellImages) func(images.PostImagesParams, *models.UserID) middleware.Responder {
	return func(params images.PostImagesParams, userID *models.UserID) middleware.Responder {
		store := newImageStore(mi)
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

		return utils.Transact(mi.DB(), func(tx *utils.AutoTx) middleware.Responder {
			tx.Query("INSERT INTO images(user_id, path, extension, processing) VALUES($1, $2, $3, false) RETURNING id",
				userID.ID, store.FileName(), store.FileExtension())
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

func NewImageLoader(mi *MindwellImages) func(images.GetImagesIDParams, *models.UserID) middleware.Responder {
	return func(params images.GetImagesIDParams, userID *models.UserID) middleware.Responder {
		return utils.Transact(mi.DB(), func(tx *utils.AutoTx) middleware.Responder {
			var authorID int64
			var path, extension string
			var processing bool

			tx.Query("SELECT user_id, path, extension, processing FROM images WHERE id = $1", params.ID).
				Scan(&authorID, &path, &extension, &processing)
			if authorID == 0 {
				return images.NewGetImagesIDNotFound()
			}

			if authorID != userID.ID {
				entryID := tx.QueryInt64("SELECT entry_id FROM entry_images WHERE image_id = $1", params.ID)
				if entryID == 0 {
					return images.NewGetImagesIDForbidden()
				}

				allowed := utils.CanViewEntry(tx, userID.ID, entryID)
				if !allowed {
					return images.NewGetImagesIDForbidden()
				}
			}

			img := &models.Image{
				ID: params.ID,
				Author: &models.User{
					ID: authorID,
				},
				Type:       extension,
				Processing: processing,
			}

			if processing {
				return images.NewGetImagesIDOK().WithPayload(img)
			}

			var width, height int64
			var size string
			tx.Query(`
				SELECT width, height, (SELECT type FROM size WHERE size.id = image_sizes.size)
				FROM image_sizes
				WHERE image_sizes.image_id = $1
			`, params.ID)

			for tx.Scan(&width, &height, &size) {
				switch size {
				case "thumbnail":
					img.Thumbnail = &models.ImageSize{
						Height: height,
						Width:  width,
						URL:    mi.BaseURL() + "albums/thumbnails/" + path,
					}
				case "small":
					img.Small = &models.ImageSize{
						Height: height,
						Width:  width,
						URL:    mi.BaseURL() + "albums/small/" + path,
					}
				case "medium":
					img.Medium = &models.ImageSize{
						Height: height,
						Width:  width,
						URL:    mi.BaseURL() + "albums/medium/" + path,
					}
				case "large":
					img.Large = &models.ImageSize{
						Height: height,
						Width:  width,
						URL:    mi.BaseURL() + "albums/large/" + path,
					}
				}
			}

			return images.NewGetImagesIDOK().WithPayload(img)
		})
	}
}

func NewImageDeleter(mi *MindwellImages) func(images.DeleteImagesIDParams, *models.UserID) middleware.Responder {
	return func(params images.DeleteImagesIDParams, userID *models.UserID) middleware.Responder {
		return utils.Transact(mi.DB(), func(tx *utils.AutoTx) middleware.Responder {
			authorID := tx.QueryInt64("SELECT user_id FROM images WHERE id = $1", params.ID)
			if authorID == 0 {
				return images.NewDeleteImagesIDNotFound()
			}

			if authorID != userID.ID {
				return images.NewDeleteImagesIDForbidden()
			}

			tx.Exec("DELETE FROM images WHERE id = $1", params.ID)
			return images.NewDeleteImagesIDNoContent()
		})
	}
}
