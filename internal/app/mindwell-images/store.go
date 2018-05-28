package images

import (
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/sevings/mindwell-server/utils"
	goconf "github.com/zpatrick/go-config"
	"gopkg.in/gographics/imagick.v2/imagick"
)

type imageStore struct {
	folder   string
	baseURL  string
	savePath string
	saveName string
	mw       *imagick.MagickWand
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
	path := name[:1] + "/" + name[1:2] + "/"

	return &imageStore{
		folder:   folder,
		baseURL:  baseURL,
		savePath: path,
		saveName: name[2:],
		mw:       imagick.NewMagickWand(),
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

func (is *imageStore) ReadImage(r io.ReadCloser, size int64, name string) {
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

	blob := make([]byte, size)
	_, is.err = r.Read(blob)
	if is.err != nil {
		return
	}

	is.err = is.mw.ReadImageBlob(blob)
	if is.err != nil {
		return
	}

	wand := is.mw.CoalesceImages()
	is.mw.Destroy()
	is.mw = wand
}

func (is *imageStore) Fill(size uint) string {
	if is.err != nil {
		return ""
	}

	path := strconv.Itoa(int(size)) + "/" + is.savePath
	is.err = os.MkdirAll(is.folder+path, 0777)
	if is.err != nil {
		return ""
	}

	wand := is.mw.Clone()
	defer wand.Destroy()

	w := wand.GetImageWidth()
	h := wand.GetImageHeight()
	if w < size {
		size = w
	}
	if h < size {
		size = h
	}

	cropSize := w
	if h < cropSize {
		cropSize = h
	}

	x := int(w-cropSize) / 2
	y := int(h-cropSize) / 2

	wand.ResetIterator()
	for wand.NextImage() {
		is.err = wand.CropImage(cropSize, cropSize, x, y)
		if is.err != nil {
			return ""
		}

		is.err = wand.ThumbnailImage(size, size)
		// is.err = wand.AdaptiveResizeImage(size, size)
		if is.err != nil {
			return ""
		}
	}

	is.err = wand.OptimizeImageTransparency()
	if is.err != nil {
		return ""
	}

	is.err = wand.SetCompressionQuality(85)
	if is.err != nil {
		return ""
	}

	fileName := path + is.saveName
	is.err = wand.WriteImages(is.folder+fileName, true)
	if is.err != nil {
		return ""
	}

	return is.baseURL + fileName
}

func (is *imageStore) RemoveOld(path string, size int) {
	if is.err != nil {
		return
	}

	is.err = os.Remove(is.folder + strconv.Itoa(size) + "/" + path)
}
