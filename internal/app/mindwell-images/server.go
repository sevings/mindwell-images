package images

import (
	"database/sql"
	"log"

	"github.com/sevings/mindwell-images/models"
	"github.com/sevings/mindwell-server/utils"

	"github.com/zpatrick/go-config"
)

type MindwellImages struct {
	cfg     *config.Config
	db      *sql.DB
	acts    chan ImageProcessor
	stop    chan bool
	folder  string
	baseURL string
}

func NewMindwellImages(cfg *config.Config) *MindwellImages {
	mi := &MindwellImages{
		cfg:  cfg,
		db:   utils.OpenDatabase(cfg),
		acts: make(chan ImageProcessor, 50),
		stop: make(chan bool),
	}

	mi.baseURL = mi.ConfigString("images.base_url")
	mi.folder = mi.ConfigString("images.folder")

	go func() {
		for act := range mi.acts {
			act.Work()
		}
		mi.stop <- true
	}

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

func (mi *MindwellImages) QueueAction(is *imageStore, ID int64, action string) {
	log.Printf("Queued: %s %s\n", action, is.FileName())
	mi.acts <- ImageProcessor{is: is, ID: ID, act: action, mi: mi}
}

func (mi *MindwellImages) Shutdown() {
	close(mi.acts)
	<- mi.stop
}
