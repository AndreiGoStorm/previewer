package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Loader struct {
	client       http.Client
	dir          string
	timeDuration time.Duration
}

func NewLoader(dir string) *Loader {
	return &Loader{
		client: http.Client{
			Transport: &http.Transport{
				MaxIdleConns:    100,
				IdleConnTimeout: 90 * time.Second,
			},
		},
		dir:          dir,
		timeDuration: 5 * time.Second,
	}
}

func (l *Loader) LoadImage(r *http.Request, im *Image) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), l.timeDuration)
	defer cancel()
	imageRequest, err := http.NewRequestWithContext(ctx, http.MethodGet, im.URL, nil)
	if err != nil {
		return err
	}
	for k, v := range r.Header {
		imageRequest.Header[k] = v
	}

	if err = l.upload(imageRequest, im); err != nil {
		return err
	}

	return
}

func (l *Loader) upload(imageRequest *http.Request, im *Image) error {
	respImage, err := l.client.Do(imageRequest)
	if err != nil {
		return err
	}
	defer respImage.Body.Close()

	if respImage.StatusCode != http.StatusOK {
		return fmt.Errorf("loader upload StatusCode: %s", respImage.Status)
	}

	data, err := io.ReadAll(respImage.Body)
	if err != nil {
		return err
	}

	im.LoadedImageName = l.generateRandomFileName(im.Ext)
	file, err := os.Create(filepath.Join(l.dir, im.LoadedImageName))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (l *Loader) generateRandomFileName(ext string) string {
	randBytes := make([]byte, 16)
	rand.Read(randBytes)
	return hex.EncodeToString(randBytes) + ext
}
