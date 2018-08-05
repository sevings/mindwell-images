package images

import (
	"io"
	"log"
	"math"
	"os"

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

func (is *imageStore) ReadImage(r io.ReadCloser) {
	defer r.Close()

	blob := make([]byte, 10*1024*1024)
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

	if is.mw.GetNumberImages() == 1 {
		is.saveName += ".jpg"
	} else {
		is.saveName += ".gif"
	}
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

	return is.saveImage(wand, folder)
}

func (is *imageStore) Fit(size uint, folder string) string {
	return is.FitRect(size, size, folder)
}

func (is *imageStore) FitRect(width, height uint, folder string) string {
	if is.err != nil {
		return ""
	}

	wand := is.mw.Clone()
	defer wand.Destroy()

	originHeight := is.mw.GetImageHeight()
	originWidth := is.mw.GetImageWidth()

	if originHeight < height && originWidth < width {
		return is.saveImage(wand, folder)
	}

	ratio := float64(width) / float64(height)
	originRatio := float64(originWidth) / float64(originHeight)

	if ratio > originRatio {
		width = uint(float64(height) * originRatio)
	} else {
		height = uint(float64(width) / originRatio)
	}

	wand.ResetIterator()
	for wand.NextImage() {
		is.err = wand.ResizeImage(width, height, imagick.FILTER_CUBIC, 0.5)
		if is.err != nil {
			return ""
		}
	}

	return is.saveImage(wand, folder)
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

func (is *imageStore) saveImage(wand *imagick.MagickWand, folder string) string {
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
