package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// === HTTP Handler Testing với httptest ===
// httptest.NewRecorder: capture response
// httptest.NewRequest: tạo fake request

func TestHealthHandler(t *testing.T) {
	// Tạo request và recorder
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	// Gọi handler trực tiếp — không cần start server
	healthHandler(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want %d", res.StatusCode, http.StatusOK)
	}

	contentType := res.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", contentType)
	}

	var body map[string]any
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body["status"] != "ok" {
		t.Errorf("status = %v, want ok", body["status"])
	}
}

func TestUserHandler_CRUD(t *testing.T) {
	store := NewUserStore()
	h := NewUserHandler(store)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /users", h.ListUsers)
	mux.HandleFunc("POST /users", h.CreateUser)
	mux.HandleFunc("GET /users/{id}", h.GetUser)
	mux.HandleFunc("DELETE /users/{id}", h.DeleteUser)

	// Test: List empty store
	t.Run("list empty", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
		}
	})

	// Test: Create user
	t.Run("create user", func(t *testing.T) {
		body := `{"name":"Alice","email":"alice@example.com","age":30}`
		req := httptest.NewRequest(http.MethodPost, "/users",
			strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("status = %d, want %d", w.Code, http.StatusCreated)
		}

		var user User
		if err := json.NewDecoder(w.Body).Decode(&user); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if user.Name != "Alice" {
			t.Errorf("name = %q, want Alice", user.Name)
		}
		if user.ID == 0 {
			t.Error("ID should be assigned")
		}
	})

	// Test: Get nonexistent user
	t.Run("get not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/999", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
		}
	})
}
