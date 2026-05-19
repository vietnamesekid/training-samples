package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func demoStrings() {
	// String in Go is an immutable sequence of bytes (UTF-8 encoded)
	s := "Xin chào Việt Nam 🇻🇳"

	// GOTCHA: len() returns the number of BYTES, not the number of characters!
	fmt.Printf("Chuỗi: %q\n", s)
	fmt.Printf("len() = %d bytes  ← số bytes, KHÔNG phải ký tự\n", len(s))
	fmt.Printf("utf8.RuneCountInString() = %d ký tự Unicode\n", utf8.RuneCountInString(s))

	// Correct way to iterate over Unicode: use range (returns runes)
	fmt.Println("\nIterate với range (đúng cho Unicode):")
	for i, r := range "Go🎯" {
		fmt.Printf("  index=%d, rune=%c, value=%d\n", i, r, r)
	}

	// Iterate over bytes (only use when raw byte processing is needed)
	fmt.Println("Iterate qua bytes (sai cho multi-byte chars):")
	for i, b := range []byte("Go🎯") {
		fmt.Printf("  index=%d, byte=0x%02x\n", i, b)
	}

	// Strings cannot be mutated — use strings.Builder for efficient concatenation
	fmt.Println("\nstrings.Builder (concat hiệu quả O(n)):")
	var builder strings.Builder
	builder.Grow(50) // pre-allocate to avoid reallocation
	for i := range 5 {
		builder.WriteString(fmt.Sprintf("item%d ", i))
	}
	result := builder.String()
	fmt.Println(" ", result)

	// GOTCHA: string concat with += is O(n²) — do NOT use in large loops
	// bad := ""
	// for i := 0; i < 10000; i++ { bad += "x" } // ← creates a new string every iteration!

	// Converting string ↔ []byte (both create a copy)
	b := []byte("hello")
	b[0] = 'H'
	fmt.Printf("\n[]byte rồi modify: %s\n", string(b))

	// === Common operations with the strings package ===
	fmt.Println("\nstrings package:")
	fmt.Printf("  ToUpper: %s\n", strings.ToUpper("hello"))
	fmt.Printf("  ToLower: %s\n", strings.ToLower("HELLO"))
	fmt.Printf("  Contains: %t\n", strings.Contains("golang", "go"))
	fmt.Printf("  HasPrefix: %t\n", strings.HasPrefix("golang", "go"))
	fmt.Printf("  HasSuffix: %t\n", strings.HasSuffix("golang", "lang"))
	fmt.Printf("  Index: %d\n", strings.Index("golang", "lang"))
	fmt.Printf("  Count: %d\n", strings.Count("cheese", "e"))
	fmt.Printf("  Replace: %s\n", strings.ReplaceAll("foo foo", "foo", "bar"))
	fmt.Printf("  TrimSpace: %q\n", strings.TrimSpace("  hello  "))
	fmt.Printf("  Trim: %q\n", strings.Trim("!!hello!!", "!"))
	fmt.Printf("  Split: %v\n", strings.Split("a,b,c", ","))
	fmt.Printf("  Join: %s\n", strings.Join([]string{"a", "b", "c"}, "-"))
	fmt.Printf("  Fields: %v\n", strings.Fields("  foo   bar  baz  "))
	fmt.Printf("  Repeat: %s\n", strings.Repeat("ab", 3))

	// Go 1.24+: strings.Lines to iterate over each line
	fmt.Println("\nstrings.Lines (Go 1.24+):")
	text := "line1\nline2\nline3"
	for line := range strings.Lines(text) {
		fmt.Printf("  %q\n", line)
	}

	// Go 1.24+: strings.SplitSeq — lazy split (iterator)
	fmt.Println("strings.SplitSeq (Go 1.24+):")
	for part := range strings.SplitSeq("a,b,c,d", ",") {
		fmt.Printf("  %q\n", part)
	}
}
