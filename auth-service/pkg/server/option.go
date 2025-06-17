package server

import (
	"net"
	"time"
)

type Options func(s *Server)

func Port(addr string) Options {
	return func(s *Server) {
		s.httpServer.Addr = net.JoinHostPort("", addr)
	}
}

func ReadTimeout(timeout time.Duration) Options {
	return func(s *Server) {
		s.httpServer.ReadTimeout = timeout
	}
}
func WriteTimeout(timeout time.Duration) Options {
	return func(s *Server) {
		s.httpServer.WriteTimeout = timeout
	}
}

func ShutdownTimeout(timeout time.Duration) Options {
	return func(s *Server) {
		s.shutdownTimeout = timeout
	}
}
