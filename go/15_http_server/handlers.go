package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
)

// === Domain Types ===

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

type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var ErrUserNotFound = errors.New("user not found")

// === In-Memory Store ===

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

func (s *UserStore) Update(id int, u User) (User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.users[id]; !ok {
		return User{}, ErrUserNotFound
	}
	u.ID = id
	s.users[id] = u
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

// === HTTP Handler ===

type UserHandler struct {
	store *UserStore
}

func NewUserHandler(store *UserStore) *UserHandler {
	return &UserHandler{store: store}
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users := h.store.List()
	writeJSON(w, http.StatusOK, map[string]any{
		"users": users,
		"count": len(users),
	})
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	// Go 1.22+: r.PathValue() để lấy path parameter {id}
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user ID", fmt.Sprintf("id must be integer, got: %s", idStr))
		return
	}

	user, err := h.store.GetByID(id)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			writeError(w, http.StatusNotFound, "user not found", fmt.Sprintf("user with id=%d does not exist", id))
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, user)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body", err.Error())
		return
	}

	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "validation error", "name is required")
		return
	}

	user := h.store.Create(User{
		Name:  req.Name,
		Email: req.Email,
		Age:   req.Age,
	})

	writeJSON(w, http.StatusCreated, user)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user ID", idStr)
		return
	}

	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body", err.Error())
		return
	}

	user, err := h.store.Update(id, User{Name: req.Name, Email: req.Email, Age: req.Age})
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			writeError(w, http.StatusNotFound, "user not found", "")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, user)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user ID", idStr)
		return
	}

	if err := h.store.Delete(id); err != nil {
		if errors.Is(err, ErrUserNotFound) {
			writeError(w, http.StatusNotFound, "user not found", "")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error", err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// === Helpers ===

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		// log the error but can't write again
		_ = err
	}
}

func writeError(w http.ResponseWriter, status int, errMsg, detail string) {
	writeJSON(w, status, ErrorResponse{
		Error:   http.StatusText(status),
		Code:    status,
		Message: fmt.Sprintf("%s: %s", errMsg, detail),
	})
}
