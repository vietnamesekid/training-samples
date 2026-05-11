package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

// === User types & store (minimal copy for testing demo) ===

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

var ErrUserNotFound = errors.New("user not found")

type UserStore struct {
	mu     sync.RWMutex
	users  map[int]User
	nextID atomic.Int64
}

func NewUserStore() *UserStore {
	s := &UserStore{users: make(map[int]User)}
	s.nextID.Store(1)
	return s
}

func (s *UserStore) Create(u User) User {
	s.mu.Lock()
	defer s.mu.Unlock()
	u.ID = int(s.nextID.Add(1)) - 1
	s.users[u.ID] = u
	return u
}

func (s *UserStore) GetByID(id int) (User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, ok := s.users[id]
	if !ok {
		return User{}, ErrUserNotFound
	}
	return u, nil
}

func (s *UserStore) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.users[id]; !ok {
		return ErrUserNotFound
	}
	delete(s.users, id)
	return nil
}

func (s *UserStore) List() []User {
	s.mu.RLock()
	defer s.mu.RUnlock()
	users := make([]User, 0, len(s.users))
	for _, u := range s.users {
		users = append(users, u)
	}
	return users
}

type UserHandler struct{ store *UserStore }

func NewUserHandler(store *UserStore) *UserHandler { return &UserHandler{store: store} }

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users := h.store.List()
	writeJSON(w, http.StatusOK, map[string]any{"users": users, "count": len(users)})
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	var id int
	fmt.Sscanf(r.PathValue("id"), "%d", &id)
	u, err := h.store.GetByID(id)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, u)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}
	if req.Name == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "name required"})
		return
	}
	user := h.store.Create(User{Name: req.Name, Email: req.Email, Age: req.Age})
	writeJSON(w, http.StatusCreated, user)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	var id int
	fmt.Sscanf(r.PathValue("id"), "%d", &id)
	if err := h.store.Delete(id); err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, `{"status":"ok","time":"`+time.Now().Format(time.RFC3339)+`"}`)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
