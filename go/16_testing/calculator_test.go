package main

import (
	"errors"
	"testing"
)

// === Basic Tests ===

func TestCalculator_Add(t *testing.T) {
	c := NewCalculator()
	result := c.Add(2, 3)
	if result != 5 {
		t.Errorf("Add(2,3) = %v, want 5", result)
	}
}

// === Table-Driven Tests — pattern phổ biến nhất trong Go ===

func TestCalculator_Operations(t *testing.T) {
	c := NewCalculator()

	tests := []struct {
		name    string
		op      func() float64
		want    float64
	}{
		{"add positive", func() float64 { return c.Add(2, 3) }, 5},
		{"add negative", func() float64 { return c.Add(-1, -2) }, -3},
		{"add zero", func() float64 { return c.Add(5, 0) }, 5},
		{"sub", func() float64 { return c.Sub(10, 3) }, 7},
		{"mul", func() float64 { return c.Mul(4, 5) }, 20},
		{"mul by zero", func() float64 { return c.Mul(99, 0) }, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.op()
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculator_Div(t *testing.T) {
	c := NewCalculator()

	tests := []struct {
		name    string
		a, b    float64
		want    float64
		wantErr error
	}{
		{"normal division", 10, 2, 5, nil},
		{"decimal", 7, 2, 3.5, nil},
		{"divide by zero", 5, 0, 0, ErrDivisionByZero},
		{"negative", -10, 2, -5, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.Div(tt.a, tt.b)

			// Check error
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("Div(%v,%v) error = %v, want %v", tt.a, tt.b, err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Div(%v,%v) unexpected error: %v", tt.a, tt.b, err)
			}

			// Check result
			if got != tt.want {
				t.Errorf("Div(%v,%v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

// === t.Helper() — helper functions cho tests ===

func assertEqual(t *testing.T, got, want float64) {
	t.Helper() // làm cho error messages chỉ vào caller, không phải helper
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestCalculator_History(t *testing.T) {
	c := NewCalculator()
	c.Add(1, 2)
	c.Sub(5, 3)

	history := c.History()
	if len(history) != 2 {
		t.Fatalf("history len = %d, want 2", len(history))
	}

	c.Clear()
	if len(c.History()) != 0 {
		t.Error("History() should be empty after Clear()")
	}
}

// === Parallel Tests — chạy độc lập, nhanh hơn ===

func TestReverse(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"empty", "", ""},
		{"single char", "a", "a"},
		{"ascii", "hello", "olleh"},
		{"unicode", "🎯Go", "oG🎯"},
		{"palindrome", "racecar", "racecar"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel() // các subtests chạy song song
			got := Reverse(tt.input)
			if got != tt.want {
				t.Errorf("Reverse(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// === TestMain — setup/teardown toàn bộ test suite ===

func TestMain(m *testing.M) {
	// Setup: chạy trước tất cả tests
	// Ví dụ: setup DB, start server, load fixtures
	// log.Println("Setting up test suite...")

	code := m.Run() // chạy tất cả tests

	// Teardown: chạy sau tất cả tests
	// Ví dụ: cleanup DB, stop server
	// log.Println("Tearing down test suite...")

	// Không dùng os.Exit trực tiếp — dùng m.Run() return code
	if code != 0 {
		// Test failures
	}
}
