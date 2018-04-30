package images

import (
	"image"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

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

type storeError string

func (se storeError) Error() string {
	return string(se)
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
		saveName: name[2:],
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

func (is *imageStore) ReadImage(r io.ReadCloser, name string) {
	if strings.HasSuffix(name, ".jpg") ||
		strings.HasSuffix(name, ".jpeg") ||
		strings.HasSuffix(name, ".png") ||
		strings.HasSuffix(name, ".bmp") ||
		strings.HasSuffix(name, ".tiff") ||
		strings.HasSuffix(name, ".tif") {
		is.saveName += ".jpg"
	} else if strings.HasSuffix(name, ".gif") {
		is.saveName += ".gif"
	} else {
		is.err = storeError("Unknown format")
		return
	}

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

	bounds := is.image.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()
	minSize := w
	if h < w {
		minSize = h
	}

	var img image.Image

	if minSize < size {
		img = imaging.CropCenter(is.image, minSize, minSize)
	} else {
		img = imaging.Thumbnail(is.image, size, size, imaging.CatmullRom)
	}

	fileName := path + is.saveName
	is.err = imaging.Save(img, is.folder+fileName, imaging.JPEGQuality(90))

	return is.baseURL + fileName
}

func (is *imageStore) RemoveOld(path string, size int) {
	if is.err != nil {
		return
	}

	is.err = os.Remove(is.folder + strconv.Itoa(size) + "/" + path)
}
