package main

import (
	"testing"
	"unicode/utf8"
)

// === Fuzz Testing (Go 1.18+) ===
// Fuzz tests automatically generate inputs to find edge cases
// Run: go test -fuzz=FuzzReverse -fuzztime=10s
// Then run regression: go test ./...  (corpus is saved)

// FuzzReverse checks the property: double reverse = identity
func FuzzReverse(f *testing.F) {
	// Seed corpus: initial inputs to start fuzzing from
	f.Add("")
	f.Add("hello")
	f.Add("Go 🎯")
	f.Add("racecar")

	// Fuzz function receives f *testing.F and input types
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
