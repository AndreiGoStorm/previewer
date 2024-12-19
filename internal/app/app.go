package app

import (
	"fmt"
	"net/http"

	"github.com/AndreiGoStorm/previewer/internal/cache"
	"github.com/AndreiGoStorm/previewer/internal/logger"
	"github.com/AndreiGoStorm/previewer/internal/service"
)

type App struct {
	logg      *logger.Logger
	lru       cache.Cache
	previewer *service.Previewer
}

func New(logg *logger.Logger, lru cache.Cache, previewer *service.Previewer) *App {
	return &App{
		logg:      logg,
		lru:       lru,
		previewer: previewer,
	}
}

func (a *App) HandleFill(w http.ResponseWriter, r *http.Request) {
	resp := &Response{}
	if r.Method != http.MethodGet {
		err := fmt.Errorf("method %s not not supported on uri %s", r.Method, r.URL.Path)
		resp.Write(w, err, http.StatusMethodNotAllowed)
		return
	}

	req := &Request{}
	req.CreateHash(r.URL.Path)

	ext, ok := a.lru.Get(req.Hash)
	if ok {
		a.logg.Info("Cache Hit " + req.Hash)
		imagePath, err := a.previewer.Storage.GetImagePath(req.GetCachedImageName(ext))
		if err == nil {
			resp.WriteImage(w, r, imagePath)
			return
		}
	}

	if err := req.Validate(r); err != nil {
		a.logg.Warn("app request validate", err)
		resp.Write(w, err, http.StatusUnprocessableEntity)
		return
	}

	if err := a.previewer.Preview(r, req.ConvertToServiceImage()); err != nil {
		a.logg.Warn("app previewer preview", err)
		resp.Write(w, err, http.StatusBadGateway)
		return
	}

	a.lru.Set(req.Hash, req.Ext)
	imagePath, err := a.previewer.Storage.GetImagePath(req.GetImageName())
	if err != nil {
		a.logg.Warn("app previewer GetImagePath", err)
		resp.Write(w, err, http.StatusNotFound)
		return
	}

	resp.WriteImage(w, r, imagePath)
}
