package images

import (
	"database/sql"
	"log"
	"time"

	"github.com/sevings/mindwell-images/models"
	"github.com/sevings/mindwell-server/utils"
)

const {
	ActionAvatar = "avatar"
	ActionCover  = "cover"
	ActionAlbum  = "album"
	ActionDelete = "delete"
}

type ImageProcessor struct {
	act string
	ID  int64
	is  *imageStore
	mi  *MindwellImages
}

func (ip *ImageProcessor) Work() {
	defer ip.is.Destroy()

	start := time.Now()
	log.Printf("Working: %s %S\n", ip.act, ip.is.FileName())

	switch ip.act {
	case ActionAvatar:
		ip.saveAvatar()
	case ActionCover:
		ip.saveCover()
	case ActionAlbum:
		ip.saveAlbumPhoto()
	case ActionDelete:
		ip.deleteAlbumPhoto()
	default:
		log.Printf("Unknown ImageProcessor action: %s\n", ip.act)
	}

	elapsed := time.Since(start).Nanoseconds()/1000000
	log.Printf("Done in %d ms\n", elapsed)
}

func (ip *ImageProcessor) saveAvatar() {
	ip.is.Fill(124, "avatars/124")
	ip.is.Fill(92, "avatars/92")
	ip.is.Fill(42, "avatars/42")

	if ip.is.Error() != nil {
		log.Println(ip.is.Error())
		return
	}

	tx := NewAutoTx(ip.mi.DB())
	defer tx.Finish()

	var old string
	tx.Query("select avatar from users where id = $1", ip.ID).Scan(&old)
	tx.Exec("update users set avatar = $2 where id = $1", ip.ID, ip.is.FileName())

	if tx.Error() != nil {
		return 
	}

	ip.is.FolderRemove("avatars/124", old)
	ip.is.FolderRemove("avatars/92", old)
	ip.is.FolderRemove("avatars/42", old)

	if ip.is.Error() != nil {
		log.Print(ip.is.Error())
	}
}

func (ip *ImageProcessor) saveCover() {
	ip.is.FillRect(1920, 640, "covers/1920")
	ip.is.FillRect(318, 122, "covers/318")

	if ip.is.Error() != nil {
		log.Println(ip.is.Error())
		return
	}

	tx := NewAutoTx(ip.mi.DB())
	defer tx.Finish()

	var old string
	tx.Query("select cover from users where id = $1", userID.ID).Scan(&old)
	tx.Exec("update users set cover = $2 where id = $1", userID.ID, ip.is.FileName())

	if tx.Error() != nil {
		return
	}

	ip.is.FolderRemove("covers/1920", old)
	ip.is.FolderRemove("covers/318", old)

	if ip.is.Error() != nil {
		log.Print(ip.is.Error())
	}
}

func (ip *ImageProcessor) saveAlbumPhoto() {
	thumbnail = ip.is.Fill(100, "albums/thumbnails")
	small     = ip.is.FitRect(480, 360, "albums/small")
	medium    = ip.is.FitRect(800, 600, "albums/medium")
	large     = ip.is.FitRect(1280, 960, "albums/large")

	if ip.is.Error() != nil {
		log.Println(ip.is.Error())
		return
	}

	tx := NewAutoTx(ip.mi.DB())
	defer tx.Finish()

	saveImageSize := func(tx *utils.AutoTx, imageID, width, height int64, size string) {
		const q = `
			INSERT INTO image_sizes(image_id, size, width, height)
			VALUES($1, (SELECT id FROM size WHERE type = $2), $3, $4)
		`
	
		tx.Exec(q, imageID, size, width, height)
	}
	
	saveImageSize(tx, ip.ID, thumbnail.Width, thumbnail.Height, "thumbnail")
	saveImageSize(tx, ip.ID, small.Width, small.Height, "small")
	saveImageSize(tx, ip.ID, medium.Width, medium.Height, "medium")
	saveImageSize(tx, ip.ID, large.Width, large.Height, "large")

	tx.Exec("UPDATE images SET processing = false WHERE id = $1", ip.img.ID)
}

func (ip *ImageProcessor) deleteAlbumPhoto() {
	tx := NewAutoTx(ip.mi.DB())
	defer tx.Finish()

	path := tx.QueryString("DELETE FROM images WHERE id = $1 RETURNING path", ip.ID)
	if tx.Error() != nil {
		return
	}

	ip.is.FolderRemove("albums/thumbnails", path)
	ip.is.FolderRemove("albums/small", path)
	ip.is.FolderRemove("albums/medium", path)
	ip.is.FolderRemove("albums/large", path)

	if ip.is.Error() != nil {
		log.Print(ip.is.Error())
	}
}
