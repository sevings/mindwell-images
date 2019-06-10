package images

import (
	"database/sql"
	"log"

	"github.com/sevings/mindwell-server/utils"

	"github.com/zpatrick/go-config"
)

type MindwellImages struct {
	cfg     *config.Config
	db      *sql.DB
	folder  string
	baseURL string
}

func NewMindwellImages(cfg *config.Config) *MindwellImages {
	mi := &MindwellImages{
		cfg: cfg,
		db:  utils.OpenDatabase(cfg),
	}

	mi.baseURL = mi.ConfigString("images.base_url")
	mi.folder = mi.ConfigString("images.folder")

	return mi
}

func (mi *MindwellImages) ConfigString(key string) string {
	value, err := mi.cfg.String(key)
	if err != nil {
		log.Println(err)
	}

	return value
}

func (mi *MindwellImages) Folder() string {
	return mi.folder
}

func (mi *MindwellImages) BaseURL() string {
	return mi.baseURL
}

func (mi *MindwellImages) DB() *sql.DB {
	return mi.db
}
