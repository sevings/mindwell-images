package images

import (
	"io"
	"math"
	"os"

	"github.com/sevings/mindwell-images/models"
	"github.com/sevings/mindwell-server/utils"
	"gopkg.in/gographics/imagick.v2/imagick"
)

type imageStore struct {
	savePath  string
	saveName  string
	extension string
	mw        *imagick.MagickWand
	mi        *MindwellImages
	err       error
}

type storeError string

func (se storeError) Error() string {
	return string(se)
}

func newImageStore(mi *MindwellImages) *imageStore {
	name := utils.GenerateString(10)
	path := name[:1] + "/" + name[1:2] + "/"

	return &imageStore{
		savePath: path,
		saveName: name[2:],
		mw:       imagick.NewMagickWand(),
		mi:       mi,
	}
}

func (is *imageStore) Destroy() {
	is.mw.Destroy()
}

func (is *imageStore) Error() error {
	return is.err
}

func (is *imageStore) Folder() string {
	return is.mi.Folder()
}

func (is *imageStore) FileName() string {
	return is.savePath + is.saveName
}

func (is *imageStore) FileExtension() string {
	return is.extension
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
		is.extension = "jpg"
	} else {
		is.extension = "gif"
	}

	is.saveName += "." + is.extension
}

func (is *imageStore) Fill(size uint, folder string) *models.ImageSize {
	return is.FillRect(size, size, folder)
}

func (is *imageStore) FillRect(width, height uint, folder string) *models.ImageSize {
	if is.err != nil {
		return nil
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
				return nil
			}
		}

		is.err = wand.ThumbnailImage(width, height)
		// is.err = wand.AdaptiveResizeImage(width, height)
		if is.err != nil {
			return nil
		}
	}

	return &models.ImageSize{
		Width:  int64(width),
		Height: int64(height),
		URL:    is.saveImage(wand, folder),
	}
}

func (is *imageStore) Fit(size uint, folder string) *models.ImageSize {
	return is.FitRect(size, size, folder)
}

func (is *imageStore) FitRect(width, height uint, folder string) *models.ImageSize {
	if is.err != nil {
		return nil
	}

	wand := is.mw.Clone()
	defer wand.Destroy()

	originHeight := is.mw.GetImageHeight()
	originWidth := is.mw.GetImageWidth()

	if originHeight < height && originWidth < width {
		return &models.ImageSize{
			Width:  int64(originWidth),
			Height: int64(originHeight),
			URL:    is.saveImage(wand, folder),
		}
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
			return nil
		}
	}

	return &models.ImageSize{
		Width:  int64(width),
		Height: int64(height),
		URL:    is.saveImage(wand, folder),
	}
}

func (is *imageStore) FolderRemove(folder, path string) {
	if is.err != nil {
		return
	}

	if len(path) == 0 {
		return
	}

	is.err = os.Remove(is.Folder() + folder + "/" + path)
}

func (is *imageStore) saveImage(wand *imagick.MagickWand, folder string) string {
	is.err = wand.OptimizeImageTransparency()
	if is.err != nil {
		return ""
	}

	is.err = wand.SetCompressionQuality(75)
	if is.err != nil {
		return ""
	}

	path := folder + "/" + is.savePath
	is.err = os.MkdirAll(is.Folder()+path, 0777)
	if is.err != nil {
		return ""
	}

	fileName := path + is.saveName
	is.err = wand.WriteImages(is.Folder()+fileName, true)
	if is.err != nil {
		return ""
	}

	return is.mi.BaseURL() + fileName
}
