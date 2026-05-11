package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

func demoGo122() {
	fmt.Println("\n--- 1. Loop Variable Fix ---")
	// TRƯỚC Go 1.22: loop variable được share giữa iterations
	// Đây là lỗi classic #1 của Go:
	//   for i := 0; i < 3; i++ {
	//       go func() { fmt.Println(i) }() // in 3, 3, 3
	//   }
	//
	// Go 1.22+: mỗi iteration có biến riêng
	funcs := make([]func(), 3)
	for i := range 3 {
		i := i // không cần dòng này nữa từ Go 1.22!
		funcs[i] = func() { fmt.Printf("  i = %d\n", i) }
	}
	for _, f := range funcs {
		f() // in 0, 1, 2 (đúng)
	}

	fmt.Println("\n--- 2. for range integer (Go 1.22+) ---")
	// range over integer — không cần range slice/channel
	sum := 0
	for i := range 10 {
		sum += i
	}
	fmt.Printf("  sum(0..9) = %d\n", sum)

	// range over integer with value
	for i := range 5 {
		fmt.Printf("  %d", i)
	}
	fmt.Println()

	fmt.Println("\n--- 3. Enhanced HTTP Mux (Go 1.22+) ---")
	// Mux giờ hỗ trợ method + path pattern: "METHOD /path/{var}"
	mux := http.NewServeMux()

	mux.HandleFunc("GET /users", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "list users")
	})

	mux.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id") // Go 1.22+: lấy path variable
		fmt.Fprintf(w, "get user: %s", id)
	})

	mux.HandleFunc("POST /users", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "create user")
	})

	mux.HandleFunc("DELETE /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		fmt.Fprintf(w, "delete user: %s", id)
	})

	// Test các routes
	tests := []struct {
		method, path string
	}{
		{"GET", "/users"},
		{"GET", "/users/42"},
		{"POST", "/users"},
		{"DELETE", "/users/7"},
	}

	for _, tt := range tests {
		req := httptest.NewRequest(tt.method, tt.path, nil)
		rw := httptest.NewRecorder()
		mux.ServeHTTP(rw, req)
		fmt.Printf("  %s %s → %s", tt.method, tt.path, rw.Body.String())
	}

	// Wildcard route: {path...} match everything
	mux.HandleFunc("GET /files/{path...}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "file: %s", r.PathValue("path"))
	})

	req := httptest.NewRequest("GET", "/files/images/logo.png", nil)
	rw := httptest.NewRecorder()
	mux.ServeHTTP(rw, req)
	fmt.Printf("  GET /files/images/logo.png → %s\n", rw.Body.String())
}
