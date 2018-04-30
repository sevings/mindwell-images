package images

import (
	"database/sql"
	"log"
	"os"

	"github.com/go-openapi/runtime/middleware"
	"github.com/sevings/mindwell-images/models"
	"github.com/sevings/mindwell-images/restapi/operations/me"
	"github.com/sevings/mindwell-server/utils"
	goconf "github.com/zpatrick/go-config"
)

func NewAvatarUpdater(db *sql.DB, cfg *goconf.Config) func(me.PutUsersMeAvatarParams, *models.UserID) middleware.Responder {
	return func(params me.PutUsersMeAvatarParams, userID *models.UserID) middleware.Responder {
		store := newImageStore(cfg)
		store.ReadImage(params.File.Data, params.File.Header.Filename)

		av800 := store.Fill(800)
		av400 := store.Fill(400)

		if store.err != nil {
			log.Print(store.Error())
			return me.NewPutUsersMeAvatarBadRequest()
		}

		avatar := models.Avatar{
			Nr800: av800,
			Nr400: av400,
		}

		return utils.Transact(db, func(tx *utils.AutoTx) middleware.Responder {

			var old string
			tx.Query("select avatar from users where id = $1", userID).Scan(&old)
			tx.Exec("update users set avatar = $2 where id = $1", userID, store.FileName())
			if tx.Error() != nil {
				return me.NewPutUsersMeAvatarBadRequest()
			}

			err := os.Remove(store.Folder() + old)
			if err != nil {
				log.Print(err)
			}

			return me.NewPutUsersMeAvatarOK().WithPayload(&avatar)
		})
	}
}
