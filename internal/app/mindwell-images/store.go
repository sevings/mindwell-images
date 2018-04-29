package images

import (
	"image"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/disintegration/imaging"
	"github.com/sevings/mindwell-server/utils"
	goconf "github.com/zpatrick/go-config"
)

type imageStore struct {
	folder   string
	baseURL  string
	savePath string
	saveName string
	image    image.Image
	err      error
}

func newImageStore(cfg *goconf.Config) *imageStore {
	folder, err := cfg.String("images.folder")
	if err != nil {
		log.Println(err)
	}

	baseURL, err := cfg.String("images.base_url")
	if err != nil {
		log.Println(err)
	}

	name := utils.GenerateString(10)
	path := "/" + name[:1] + "/" + name[1:2] + "/"

	return &imageStore{
		folder:   folder,
		baseURL:  baseURL,
		savePath: path,
		saveName: name[2:] + "_",
	}
}

func (is *imageStore) Error() error {
	return is.err
}

func (is *imageStore) Folder() string {
	return is.folder
}

func (is *imageStore) FileName() string {
	return is.savePath + is.saveName
}

func (is *imageStore) SetImage(r io.ReadCloser, name string) {
	is.saveName += name

	defer r.Close()
	is.image, is.err = imaging.Decode(r)
}

func (is *imageStore) Fill(size int) string {
	if is.err != nil {
		return ""
	}

	path := strconv.Itoa(size) + is.savePath
	is.err = os.MkdirAll(is.folder+path, 0777)
	if is.err != nil {
		return ""
	}

	img := imaging.Fill(is.image, size, size, imaging.Center, imaging.Linear)

	fileName := path + is.saveName
	is.err = imaging.Save(img, is.folder+fileName, imaging.JPEGQuality(80))

	return fileName
}
