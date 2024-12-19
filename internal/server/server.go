package server

import (
	"context"
	"net"
	"net/http"
	"previewer/internal/app"
	"previewer/internal/config"
	"previewer/internal/logger"
	"strconv"
	"time"
)

const (
	httpTimeout = 5 * time.Second
)

type Server struct {
	server *http.Server
	notify chan error
	logg   *logger.Logger
}

func New(conf config.HTTP, l *logger.Logger) *Server {
	httpServer := &http.Server{
		ReadTimeout:  httpTimeout,
		WriteTimeout: httpTimeout,
		Addr:         net.JoinHostPort(conf.Host, strconv.Itoa(conf.Port)),
	}

	return &Server{
		server: httpServer,
		notify: make(chan error, 1),
		logg:   l,
	}
}

func (s *Server) Start(application *app.App) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", application.HandleStart)
	mux.HandleFunc("/fill/", application.HandleFill)
	s.server.Handler = mux

	go func() {
		s.logg.Info("server started on " + s.server.Addr)
		s.notify <- s.server.ListenAndServe()
		close(s.notify)
	}()
}

func (s *Server) Notify() <-chan error {
	s.logg.Info("server notify")
	return s.notify
}

func (s *Server) Stop() error {
	s.logg.Info("server stopped")
	ctx, cancel := context.WithTimeout(context.Background(), httpTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}
