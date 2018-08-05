package images

import (
	"database/sql"
	"log"

	"github.com/go-openapi/runtime/middleware"
	"github.com/sevings/mindwell-images/models"
	"github.com/sevings/mindwell-images/restapi/operations/me"
	"github.com/sevings/mindwell-server/utils"
	goconf "github.com/zpatrick/go-config"
)

func NewAvatarUpdater(db *sql.DB, cfg *goconf.Config) func(me.PutMeAvatarParams, *models.UserID) middleware.Responder {
	return func(params me.PutMeAvatarParams, userID *models.UserID) middleware.Responder {
		store := newImageStore(cfg)
		defer store.Destroy()

		store.ReadImage(params.File)

		avatar := models.Avatar{
			X124: store.Fill(124, "avatars/124"),
			X92:  store.Fill(92, "avatars/92"),
			X42:  store.Fill(42, "avatars/42"),
		}

		if store.Error() != nil {
			log.Print(store.Error())
			return me.NewPutMeAvatarBadRequest()
		}

		return utils.Transact(db, func(tx *utils.AutoTx) middleware.Responder {
			var old string
			tx.Query("select avatar from users where id = $1", userID.ID).Scan(&old)
			tx.Exec("update users set avatar = $2 where id = $1", userID.ID, store.FileName())
			if tx.Error() != nil {
				return me.NewPutMeAvatarBadRequest()
			}

			//! \todo remove after 6 months?
			// store.SizeRemove(124, old)
			// store.SizeRemove(92, old)
			// store.SizeRemove(42, old)
			// if store.Error() != nil {
			// 	log.Print(store.Error())
			// }

			return me.NewPutMeAvatarOK().WithPayload(&avatar)
		})
	}
}

func NewCoverUpdater(db *sql.DB, cfg *goconf.Config) func(me.PutMeCoverParams, *models.UserID) middleware.Responder {
	return func(params me.PutMeCoverParams, userID *models.UserID) middleware.Responder {
		store := newImageStore(cfg)
		store.ReadImage(params.File)

		cover := &models.Cover{
			ID:    userID.ID,
			X1920: store.FillRect(1920, 640, "covers/1920"),
			X318:  store.FillRect(318, 122, "covers/318"),
		}

		if store.Error() != nil {
			log.Print(store.Error())
			return me.NewPutMeCoverBadRequest()
		}

		return utils.Transact(db, func(tx *utils.AutoTx) middleware.Responder {
			var old string
			tx.Query("select cover from users where id = $1", userID.ID).Scan(&old)
			tx.Exec("update users set cover = $2 where id = $1", userID.ID, store.FileName())
			if tx.Error() != nil {
				return me.NewPutMeCoverBadRequest()
			}

			// store.FolderRemove("covers", old)
			// if store.Error() != nil {
			// 	log.Print(store.Error())
			// }

			return me.NewPutMeCoverOK().WithPayload(cover)
		})
	}
}
