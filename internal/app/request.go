package app

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"path"
	"previewer/internal/service"
	"strconv"
	"strings"
)

type Request struct {
	Hash   string
	Width  int
	Height int
	Url    string
	Ext    string
}

func (req *Request) CreateHash(url string) {
	h1 := sha256.New()
	h1.Write([]byte(url))
	hash := h1.Sum(nil)
	req.Hash = hex.EncodeToString(hash)
}

func (req *Request) Validate(r *http.Request) (err error) {
	url := strings.TrimLeft(r.URL.Path, r.Pattern)
	parts := strings.Split(url, "/")
	if len(parts) < 2 {
		return fmt.Errorf("wrong loading url: %s", url)
	}

	err = req.validateWidth(parts[0])
	if err != nil {
		return err
	}

	err = req.validateHeight(parts[1])
	if err != nil {
		return err
	}

	err = req.validateUrl(strings.TrimLeft(url, fmt.Sprintf("%s/%s/", parts[0], parts[1])))
	if err != nil {
		return err
	}

	return nil
}

func (req *Request) validateWidth(width string) (err error) {
	req.Width, err = strconv.Atoi(width)
	if err != nil {
		return fmt.Errorf("wrong width: %s", width)
	}
	if req.Width <= 0 || req.Width >= 10000 {
		return fmt.Errorf("wrong width: %d", req.Width)
	}
	return
}

func (req *Request) validateHeight(height string) (err error) {
	req.Height, err = strconv.Atoi(height)
	if err != nil {
		return fmt.Errorf("wrong height: %s", height)
	}
	if req.Height <= 0 || req.Height >= 10000 {
		return fmt.Errorf("wrong height: %d", req.Height)
	}
	return
}

func (req *Request) validateUrl(url string) (err error) {
	if url == "" {
		return fmt.Errorf("loading url is empty")
	}
	req.Url = fmt.Sprintf("https://%s", url)

	req.Ext = strings.ToLower(path.Ext(url))
	if !strings.Contains("jpg,jpeg,png,gif", strings.TrimLeft(req.Ext, ".")) {
		return fmt.Errorf("loading image has wrong extension: %s", req.Ext)
	}
	return
}

func (req *Request) ConvertToServiceImage() *service.Image {
	im := &service.Image{
		Width:     req.Width,
		Height:    req.Height,
		Url:       req.Url,
		Ext:       req.Ext,
		ImageName: req.Hash + req.Ext,
	}

	return im
}
