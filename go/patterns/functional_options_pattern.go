package main

import (
	"time"
)

type Server struct {
	host    string
	port    int
	timeout time.Duration
}

type ServerOption func(*Server)

func WithHost(host string) ServerOption {
	return func(s *Server) {
		s.host = host
	}
}

func WithPort(port int) ServerOption {
	return func(s *Server) {
		s.port = port
	}
}

func WithTimeout(timeout time.Duration) ServerOption {
	return func(s *Server) {
		s.timeout = timeout
	}
}

func NewServer(host string, opts ...ServerOption) *Server {
	s := &Server{
		host:    host,
		port:    8080,
		timeout: 30 * time.Second,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}
