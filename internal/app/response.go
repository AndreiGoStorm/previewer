package app

import (
	"encoding/json"
	"log"
	"net/http"
)

type Response struct {
	Data  interface{} `json:"data"`
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

func (resp *Response) WriteError(w http.ResponseWriter, err error, statusCode int) {
	resp.Error.Message = err.Error()
	w.WriteHeader(statusCode)

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

func (resp *Response) WriteImage(w http.ResponseWriter, r *http.Request, filename string) {
	http.ServeFile(w, r, filename)
}
