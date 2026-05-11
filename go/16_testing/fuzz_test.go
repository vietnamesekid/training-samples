package main

import (
	"testing"
	"unicode/utf8"
)

// === Fuzz Testing (Go 1.18+) ===
// Fuzz test tự động generate inputs để tìm edge cases
// Chạy: go test -fuzz=FuzzReverse -fuzztime=10s
// Sau đó chạy regression: go test ./...  (corpus được save)

// FuzzReverse kiểm tra property: double reverse = identity
func FuzzReverse(f *testing.F) {
	// Seed corpus: các input ban đầu để bắt đầu fuzzing
	f.Add("")
	f.Add("hello")
	f.Add("Go 🎯")
	f.Add("racecar")

	// Fuzz function nhận f *testing.F và input types
	f.Fuzz(func(t *testing.T, s string) {
		// Property 1: double reverse = original
		doubled := Reverse(Reverse(s))
		if doubled != s {
			t.Errorf("Reverse(Reverse(%q)) = %q, want %q", s, doubled, s)
		}

		// Property 2: length preserved (in runes, not bytes)
		original := utf8.RuneCountInString(s)
		reversed := utf8.RuneCountInString(Reverse(s))
		if original != reversed {
			t.Errorf("Reverse(%q) changed rune count: %d → %d", s, original, reversed)
		}
	})
}
