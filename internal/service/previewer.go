package service

import (
	"net/http"

	"github.com/disintegration/imaging"
	"previewer/internal/logger"
)

type Previewer struct {
	logg    *logger.Logger
	dir     string
	loader  *Loader
	Storage *Storage
}

func New(logg *logger.Logger) *Previewer {
	pr := &Previewer{
		logg: logg,
	}
	pr.Storage = NewStorage(pr.logg)
	pr.loader = NewLoader(pr.Storage.Dir)

	return pr
}

func (pr *Previewer) Preview(r *http.Request, im *Image) (err error) {
	if err := pr.loader.LoadImage(r, im); err != nil {
		pr.logg.Error("previewer Preview LoadImage: %w", err)
		return err
	}

	if err = pr.Resize(im); err != nil {
		pr.logg.Error("previewer Preview Resize: %w", err)
		return err
	}

	if err = pr.Storage.DeleteFile(im.LoadedImageName); err != nil {
		pr.logg.Error("previewer Preview DeleteImage: %w", err)
		return err
	}

	return nil
}

func (pr *Previewer) Resize(im *Image) error {
	path := pr.Storage.getStorageFullPath(im.LoadedImageName)
	image, err := imaging.Open(path)
	if err != nil {
		return err
	}

	resized := imaging.Resize(image, im.Width, im.Height, imaging.Lanczos)

	path = pr.Storage.getStorageFullPath(im.ImageName)
	if err = imaging.Save(resized, path); err != nil {
		return err
	}

	return nil
}
