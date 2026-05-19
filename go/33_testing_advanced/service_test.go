package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

// ============================================================
// Hand-written Mocks (no mockery/gomock needed)
// ============================================================

// MockEmailSender records calls for assertions
type MockEmailSender struct {
	mu    sync.Mutex
	calls []EmailCall
	err   error // return this if != nil
}

type EmailCall struct {
	To, Subject, Body string
}

func (m *MockEmailSender) Send(_ context.Context, to, subject, body string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.calls = append(m.calls, EmailCall{To: to, Subject: subject, Body: body})
	return m.err
}

func (m *MockEmailSender) CallCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.calls)
}

func (m *MockEmailSender) LastCall() (EmailCall, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.calls) == 0 {
		return EmailCall{}, false
	}
	return m.calls[len(m.calls)-1], true
}

// MockUserStore with configurable behaviors
type MockUserStore struct {
	mu    sync.Mutex
	users map[string]*User
	saves []string // track saved user IDs
}

func NewMockUserStore(users ...*User) *MockUserStore {
	m := &MockUserStore{users: make(map[string]*User)}
	for _, u := range users {
		m.users[u.ID] = u
	}
	return m
}

func (m *MockUserStore) FindByID(_ context.Context, id string) (*User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	u, ok := m.users[id]
	if !ok {
		return nil, fmt.Errorf("user %s not found", id)
	}
	cp := *u // return copy
	return &cp, nil
}

func (m *MockUserStore) Save(_ context.Context, user *User) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	cp := *user
	m.users[user.ID] = &cp
	m.saves = append(m.saves, user.ID)
	return nil
}

// ============================================================
// t.Helper() — test helper functions
// ============================================================

func requireNoError(t *testing.T, err error) {
	t.Helper() // makes failures report the caller's line, not the helper's line
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func requireEqual[T comparable](t *testing.T, want, got T, msg string) {
	t.Helper()
	if want != got {
		t.Errorf("%s: want %v, got %v", msg, want, got)
	}
}

// ============================================================
// Table-driven tests with subtests
// ============================================================

func TestArithmetic(t *testing.T) {
	tests := []struct {
		name    string
		op      string
		a, b    int
		want    int
		wantErr bool
	}{
		{"add positive", "add", 3, 4, 7, false},
		{"add negative", "add", -3, 4, 1, false},
		{"sub", "sub", 10, 3, 7, false},
		{"mul", "mul", 4, 5, 20, false},
		{"div normal", "div", 10, 2, 5, false},
		{"div by zero", "div", 10, 0, 0, true},
	}

	for _, tt := range tests {
		tt := tt // capture for parallel (pre Go 1.22)
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel() // run subtests in parallel

			var got int
			var err error

			switch tt.op {
			case "add":
				got = Add(tt.a, tt.b)
			case "sub":
				got = Sub(tt.a, tt.b)
			case "mul":
				got = Mul(tt.a, tt.b)
			case "div":
				got, err = Div(tt.a, tt.b)
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
				return
			}
			requireNoError(t, err)
			requireEqual(t, tt.want, got, "result")
		})
	}
}

// ============================================================
// Mock-based service testing
// ============================================================

func TestNotificationService_WelcomeUser(t *testing.T) {
	alice := &User{
		ID: "u1", Name: "Alice", Email: "alice@example.com",
		CreatedAt: time.Now(),
	}

	tests := []struct {
		name        string
		userID      string
		mailerErr   error
		wantErr     bool
		wantMails   int
		wantActive  bool
	}{
		{
			name:       "success",
			userID:     "u1",
			wantMails:  1,
			wantActive: true,
		},
		{
			name:    "user not found",
			userID:  "unknown",
			wantErr: true,
		},
		{
			name:      "mailer error",
			userID:    "u1",
			mailerErr: errors.New("SMTP unavailable"),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewMockUserStore(alice)
			mailer := &MockEmailSender{err: tt.mailerErr}
			svc := NewNotificationService(store, mailer)

			err := svc.WelcomeUser(context.Background(), tt.userID)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got nil")
				}
				return
			}
			requireNoError(t, err)
			requireEqual(t, tt.wantMails, mailer.CallCount(), "email count")

			if call, ok := mailer.LastCall(); ok {
				requireEqual(t, alice.Email, call.To, "email recipient")
				if !strings.Contains(call.Subject, alice.Name) {
					t.Errorf("subject should contain user name, got: %s", call.Subject)
				}
			}

			// Verify user was saved as active
			saved, _ := store.FindByID(context.Background(), "u1")
			requireEqual(t, tt.wantActive, saved.Active, "user.Active")
		})
	}
}

// ============================================================
// Golden file testing — snapshot testing
// ============================================================

func TestFormatUser_Golden(t *testing.T) {
	user := &User{
		ID:    "u123",
		Name:  "Alice Smith",
		Email: "alice@example.com",
		Active: true,
	}

	got := FormatUser(user)

	// Golden file path
	goldenPath := filepath.Join("testdata", "format_user.golden")

	if os.Getenv("UPDATE_GOLDEN") != "" {
		// UPDATE_GOLDEN=1 go test — regenerate golden files
		os.MkdirAll("testdata", 0755)
		os.WriteFile(goldenPath, []byte(got), 0644)
		t.Logf("Updated golden file: %s", goldenPath)
		return
	}

	expected, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("golden file not found: %s (run with UPDATE_GOLDEN=1 to create)", goldenPath)
	}

	if string(expected) != got {
		t.Errorf("output differs from golden file\nWant:\n%s\nGot:\n%s", expected, got)
	}
}

// ============================================================
// t.TempDir() & t.Cleanup()
// ============================================================

func TestTempDirAndCleanup(t *testing.T) {
	// t.TempDir() creates a temp dir that auto-cleans up after the test
	dir := t.TempDir()

	// Create file in temp dir
	filePath := filepath.Join(dir, "test.txt")
	os.WriteFile(filePath, []byte("test content"), 0644)

	// Verify file exists
	data, err := os.ReadFile(filePath)
	requireNoError(t, err)
	requireEqual(t, "test content", string(data), "file content")

	// t.Cleanup() registers a function to run after the test
	t.Cleanup(func() {
		t.Log("cleanup: releasing resources")
		// In practice: close DB connection, stop server, etc.
	})

	// Track side effects
	cleaned := false
	resource := func() func() {
		// Initialize resource
		return func() { cleaned = true } // cleanup function
	}()

	t.Cleanup(resource)

	// Multiple cleanups run in LIFO order (like defer)
	t.Cleanup(func() {
		if !cleaned {
			t.Log("warning: resource not cleaned up yet")
		}
	})
}

// ============================================================
// Reverse tests
// ============================================================

func TestReverse(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"hello", "olleh"},
		{"Go", "oG"},
		{"", ""},
		{"a", "a"},
		{"racecar", "racecar"}, // palindrome
		{"日本語", "語本日"},         // UTF-8
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := Reverse(tt.input)
			if got != tt.want {
				t.Errorf("Reverse(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// ============================================================
// Benchmark
// ============================================================

func BenchmarkReverse(b *testing.B) {
	s := "Hello, World! こんにちは"
	for b.Loop() {
		Reverse(s)
	}
}
