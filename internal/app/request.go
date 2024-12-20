package app

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"

	"github.com/AndreiGoStorm/previewer/internal/service"
)

type Request struct {
	Protocol  string
	Hash      string
	Width     int
	Height    int
	URL       string
	Ext       string
	ImageName string
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

	widthTrim := strings.TrimLeft(url, fmt.Sprintf("%s/", parts[0]))
	heightTrim := strings.TrimLeft(widthTrim, fmt.Sprintf("%s/", parts[1]))
	err = req.validateURL(heightTrim)
	if err != nil {
		return err
	}

	err = req.validateExt(req.URL)
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

func (req *Request) validateURL(u string) (err error) {
	if u == "" {
		return fmt.Errorf("loading url is empty")
	}

	req.URL = fmt.Sprintf("%s://%s", req.Protocol, u)
	_, err = url.ParseRequestURI(req.URL)
	if err != nil {
		return fmt.Errorf("wrong url")
	}
	return
}

func (req *Request) validateExt(u string) (err error) {
	req.Ext = strings.ToLower(path.Ext(u))
	if req.Ext == "" {
		return fmt.Errorf("loading image extension is empty")
	}

	extension := strings.TrimLeft(req.Ext, ".")
	if !strings.Contains("jpg,jpeg,png,gif", extension) { //nolint:gocritic
		return fmt.Errorf("loading image has wrong extension: %s", extension)
	}
	req.ImageName = req.Hash + req.Ext
	return
}

func (req *Request) ConvertToServiceImage() *service.Image {
	return &service.Image{
		Width:     req.Width,
		Height:    req.Height,
		URL:       req.URL,
		Ext:       req.Ext,
		ImageName: req.ImageName,
	}
}
