package images

import (
	"io"
	"log"
	"math"
	"os"
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

func (is *imageStore) Destroy() {
	is.mw.Destroy()
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
	if strings.HasSuffix(name, ".gif") {
		is.saveName += ".gif"
	} else {
		is.saveName += ".jpg"
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

func (is *imageStore) Fill(size uint, folder string) string {
	return is.FillRect(size, size, folder)
}

func (is *imageStore) FillRect(width, height uint, folder string) string {
	if is.err != nil {
		return ""
	}

	originWidth := is.mw.GetImageWidth()
	originHeight := is.mw.GetImageHeight()

	ratio := float64(width) / float64(height)
	originRatio := float64(originWidth) / float64(originHeight)

	crop := math.Abs(ratio-originRatio) > 0.01

	cropWidth, cropHeight := originWidth, originHeight

	if ratio < originRatio {
		cropWidth = uint(float64(originHeight) * ratio)
	} else {
		cropHeight = uint(float64(originWidth) / ratio)
	}

	if width > originWidth || height > originHeight {
		width, height = cropWidth, cropHeight
	}

	x := int(originWidth-cropWidth) / 2
	y := int(originHeight-cropHeight) / 2

	wand := is.mw.Clone()
	defer wand.Destroy()

	wand.ResetIterator()
	for wand.NextImage() {
		if crop {
			is.err = wand.CropImage(cropWidth, cropHeight, x, y)
			if is.err != nil {
				return ""
			}
		}

		is.err = wand.ThumbnailImage(width, height)
		// is.err = wand.AdaptiveResizeImage(width, height)
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

	path := folder + "/" + is.savePath
	is.err = os.MkdirAll(is.folder+path, 0777)
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

func (is *imageStore) FolderRemove(folder, path string) {
	if is.err != nil {
		return
	}

	if len(path) == 0 {
		return
	}

	is.err = os.Remove(is.folder + folder + "/" + path)
}
