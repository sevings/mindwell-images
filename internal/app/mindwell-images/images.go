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

func NewAvatarUpdater(db *sql.DB, cfg *goconf.Config) func(me.PutUsersMeAvatarParams, *models.UserID) middleware.Responder {
	return func(params me.PutUsersMeAvatarParams, userID *models.UserID) middleware.Responder {
		store := newImageStore(cfg)
		store.ReadImage(params.File.Data, params.File.Header.Size, params.File.Header.Filename)

		avatar := models.Avatar{
			X124: store.Fill(124),
			X42:  store.Fill(42),
		}

		if store.Error() != nil {
			log.Print(store.Error())
			return me.NewPutUsersMeAvatarBadRequest()
		}

		return utils.Transact(db, func(tx *utils.AutoTx) middleware.Responder {
			var old string
			tx.Query("select avatar from users where id = $1", userID).Scan(&old)
			tx.Exec("update users set avatar = $2 where id = $1", userID, store.FileName())
			if tx.Error() != nil {
				return me.NewPutUsersMeAvatarBadRequest()
			}

			store.SizeRemove(124, old)
			store.SizeRemove(42, old)
			if store.Error() != nil {
				log.Print(store.Error())
			}

			return me.NewPutUsersMeAvatarOK().WithPayload(&avatar)
		})
	}
}

func NewCoverUpdater(db *sql.DB, cfg *goconf.Config) func(me.PutUsersMeCoverParams, *models.UserID) middleware.Responder {
	return func(params me.PutUsersMeCoverParams, userID *models.UserID) middleware.Responder {
		store := newImageStore(cfg)
		store.ReadImage(params.File.Data, params.File.Header.Size, params.File.Header.Filename)

		cover := &models.Cover{
			ID:    int64(*userID),
			Cover: store.FillRect(1920, 640, "cover"),
		}

		if store.Error() != nil {
			log.Print(store.Error())
			return me.NewPutUsersMeCoverBadRequest()
		}

		return utils.Transact(db, func(tx *utils.AutoTx) middleware.Responder {
			var old string
			tx.Query("select cover from users where id = $1", userID).Scan(&old)
			tx.Exec("update users set cover = $2 where id = $1", userID, store.FileName())
			if tx.Error() != nil {
				return me.NewPutUsersMeCoverBadRequest()
			}

			store.FolderRemove("cover", old)
			if store.Error() != nil {
				log.Print(store.Error())
			}

			return me.NewPutUsersMeCoverOK().WithPayload(cover)
		})
	}
}
