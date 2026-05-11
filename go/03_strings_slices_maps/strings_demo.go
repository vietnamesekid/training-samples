package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func demoStrings() {
	// String trong Go là immutable sequence of bytes (UTF-8 encoded)
	s := "Xin chào Việt Nam 🇻🇳"

	// GOTCHA: len() trả về số BYTES, không phải số ký tự!
	fmt.Printf("Chuỗi: %q\n", s)
	fmt.Printf("len() = %d bytes  ← số bytes, KHÔNG phải ký tự\n", len(s))
	fmt.Printf("utf8.RuneCountInString() = %d ký tự Unicode\n", utf8.RuneCountInString(s))

	// Iterate đúng cách qua Unicode: dùng range (trả về rune)
	fmt.Println("\nIterate với range (đúng cho Unicode):")
	for i, r := range "Go🎯" {
		fmt.Printf("  index=%d, rune=%c, value=%d\n", i, r, r)
	}

	// Iterate qua bytes (chỉ dùng khi cần xử lý raw bytes)
	fmt.Println("Iterate qua bytes (sai cho multi-byte chars):")
	for i, b := range []byte("Go🎯") {
		fmt.Printf("  index=%d, byte=0x%02x\n", i, b)
	}

	// String không thể mutate — dùng strings.Builder để concat hiệu quả
	fmt.Println("\nstrings.Builder (concat hiệu quả O(n)):")
	var builder strings.Builder
	builder.Grow(50) // pre-allocate để tránh reallocation
	for i := range 5 {
		builder.WriteString(fmt.Sprintf("item%d ", i))
	}
	result := builder.String()
	fmt.Println(" ", result)

	// GOTCHA: string concat với += là O(n²) — KHÔNG dùng trong loop lớn
	// bad := ""
	// for i := 0; i < 10000; i++ { bad += "x" } // ← tạo string mới mỗi lần!

	// Chuyển đổi string ↔ []byte (đều tạo copy)
	b := []byte("hello")
	b[0] = 'H'
	fmt.Printf("\n[]byte rồi modify: %s\n", string(b))

	// === Thao tác phổ biến với strings package ===
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

	// Go 1.24+: strings.Lines để iterate qua từng dòng
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
