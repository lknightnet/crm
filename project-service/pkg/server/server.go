package server

import (
	"context"
	"net/http"
	"time"
)

const (
	defaultReadTimeout     = 1 * time.Second
	defaultWriteTimeout    = 1 * time.Second
	defaultShutdownTimeout = 1 * time.Second
	defaultAddr            = ":80"
)

type Server struct {
	httpServer      *http.Server
	shutdownTimeout time.Duration
	notify          chan error
}

func NewServer(handler http.Handler, options ...Options) *Server {
	httpServer := &http.Server{
		Handler:      handler,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
		Addr:         defaultAddr,
	}

	srv := &Server{
		httpServer:      httpServer,
		shutdownTimeout: defaultShutdownTimeout,
		notify:          make(chan error, 1),
	}

	for _, option := range options {
		option(srv)
	}

	srv.start()

	return srv
}

func (s *Server) start() {
	go func() {
		s.notify <- s.httpServer.ListenAndServe()
		close(s.notify)
	}()
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()
	return s.httpServer.Shutdown(ctx)
}
