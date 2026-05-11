package main

import (
	"fmt"
	"time"
)

// === Functional Options Pattern ===
// Pattern này giải quyết vấn đề: "Làm sao có constructor flexible với nhiều optional params?"
//
// Vấn đề:
//   NewServer("localhost", 8080, 30*time.Second, true, false, 100) ← không rõ ràng
//
// Solution: functional options

type HTTPServer struct {
	host        string
	port        int
	timeout     time.Duration
	maxConns    int
	tlsEnabled  bool
	readTimeout time.Duration
}

// ServerOption là một function modify *HTTPServer
type ServerOption func(*HTTPServer)

// Các option functions — mỗi function set một field
func WithPort(port int) ServerOption {
	return func(s *HTTPServer) {
		s.port = port
	}
}

func WithTimeout(d time.Duration) ServerOption {
	return func(s *HTTPServer) {
		s.timeout = d
		s.readTimeout = d
	}
}

func WithMaxConns(n int) ServerOption {
	return func(s *HTTPServer) {
		s.maxConns = n
	}
}

func WithTLS(enabled bool) ServerOption {
	return func(s *HTTPServer) {
		s.tlsEnabled = enabled
		if enabled && s.port == 8080 {
			s.port = 443 // default HTTPS port
		}
	}
}

// NewHTTPServer với default values + optional overrides
func NewHTTPServer(host string, opts ...ServerOption) *HTTPServer {
	// Default configuration
	s := &HTTPServer{
		host:        host,
		port:        8080,
		timeout:     30 * time.Second,
		maxConns:    1000,
		tlsEnabled:  false,
		readTimeout: 10 * time.Second,
	}

	// Apply options theo thứ tự
	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *HTTPServer) String() string {
	scheme := "http"
	if s.tlsEnabled {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s:%d (timeout=%s, maxConns=%d)",
		scheme, s.host, s.port, s.timeout, s.maxConns)
}

func demoFunctionalOptions() {
	// Minimal — chỉ dùng defaults
	s1 := NewHTTPServer("localhost")
	fmt.Printf("  Default server: %s\n", s1)

	// Custom port và timeout
	s2 := NewHTTPServer("example.com",
		WithPort(9090),
		WithTimeout(60*time.Second),
	)
	fmt.Printf("  Custom server: %s\n", s2)

	// Production với TLS
	s3 := NewHTTPServer("api.example.com",
		WithTLS(true),
		WithMaxConns(5000),
		WithTimeout(120*time.Second),
	)
	fmt.Printf("  Production server: %s\n", s3)

	fmt.Println("\n  Ưu điểm của Functional Options:")
	fmt.Println("  - Self-documenting: NewServer(host, WithPort(9090), WithTLS(true))")
	fmt.Println("  - Backward compatible: thêm option mới không break existing code")
	fmt.Println("  - Default values rõ ràng trong constructor")
	fmt.Println("  - Dễ validate trong mỗi option function")
}
