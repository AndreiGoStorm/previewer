package app

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"previewer/internal/cache"
	"previewer/internal/logger"
	"previewer/internal/service"
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

type Response struct {
	Data  interface{} `json:"data"`
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

func (a *App) HandleStart(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("handler22Default"))
	w.WriteHeader(http.StatusOK)
}

func (a *App) HandleFill(w http.ResponseWriter, r *http.Request) {
	resp := &Response{}
	if r.Method != http.MethodGet {
		resp.Error.Message = fmt.Sprintf("method %s not not supported on uri %s", r.Method, r.URL.Path)
		w.WriteHeader(http.StatusMethodNotAllowed)
		WriteResponse(w, resp)
		return
	}

	req := &Request{}
	req.CreateHash(r.URL.Path)

	_, ok := a.lru.Get(req.Hash)
	if ok {
		w.Write([]byte("Cache Hit " + req.Hash))
		w.WriteHeader(http.StatusOK)
		return
	}

	if err := req.Validate(r); err != nil {
		a.logg.Warn("app request validate", err)
		resp.Error.Message = err.Error()
		WriteResponse(w, resp)
		return
	}

	if err := a.previewer.Preview(r, req.ConvertToServiceImage()); err != nil {
		a.logg.Warn("app previewer preview", err)
		resp.Error.Message = err.Error()
		WriteResponse(w, resp)
		return
	}

	a.lru.Set(req.Hash, req.Ext)

	w.WriteHeader(http.StatusOK)
}

func WriteResponse(w http.ResponseWriter, resp *Response) {
	resBuf, err := json.Marshal(resp)
	if err != nil {
		log.Printf("response marshal error: %s", err)
	}
	_, err = w.Write(resBuf)
	if err != nil {
		log.Printf("response marshal error: %s", err)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}
