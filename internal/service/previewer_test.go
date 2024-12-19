package service

import (
	"os"
	"path"
	"path/filepath"
	"previewer/internal/logger"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPreviewer(t *testing.T) {
	var (
		testImageJpeg = "images/image.jpeg"
		testImagePng  = "images/image.png"
	)

	logg := logger.New("INFO")

	t.Run("resize image jpeg", func(t *testing.T) {
		im := &Image{
			Width:           500,
			Height:          500,
			Ext:             path.Ext(testImageJpeg),
			ImageName:       filepath.Base(testImageJpeg),
			LoadedImageName: "image_for_resize_jpeg.jpeg",
		}
		previewer := New(logg)

		fileBytes, err := os.ReadFile(testImageJpeg)
		require.NoError(t, err)
		imageFile, err := os.Create(filepath.Join(previewer.Storage.Dir, im.LoadedImageName))
		require.NoError(t, err)
		defer imageFile.Close()
		_, err = imageFile.Write(fileBytes)
		require.NoError(t, err)

		err = previewer.Resize(im)
		require.NoError(t, err)

		names, err := previewer.Storage.ReadDirNames()
		require.NoError(t, err)

		require.Len(t, names, 2)

		err = os.RemoveAll(previewer.Storage.Dir)
		require.NoError(t, err)
	})

	t.Run("resize image png", func(t *testing.T) {
		im := &Image{
			Width:           300,
			Height:          300,
			Ext:             path.Ext(testImagePng),
			ImageName:       filepath.Base(testImagePng),
			LoadedImageName: "image_for_resize_png.png",
		}
		previewer := New(logg)

		fileBytes, err := os.ReadFile(testImagePng)
		require.NoError(t, err)
		imageFile, err := os.Create(filepath.Join(previewer.Storage.Dir, im.LoadedImageName))
		require.NoError(t, err)
		defer imageFile.Close()
		_, err = imageFile.Write(fileBytes)
		require.NoError(t, err)

		err = previewer.Resize(im)
		require.NoError(t, err)

		names, err := previewer.Storage.ReadDirNames()
		require.NoError(t, err)

		require.Len(t, names, 2)

		err = os.RemoveAll(previewer.Storage.Dir)
		require.NoError(t, err)
	})
}
