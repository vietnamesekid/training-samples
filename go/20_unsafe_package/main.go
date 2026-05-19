// Lesson 20: unsafe Package — low-level memory operations
// Run: go run .
// WARNING: unsafe breaks Go's type safety
// Only use when: performance critical, interop with C, zero-copy operations
package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

// === Struct Layout & Padding ===

// BAD: padding waste — Go does not auto-reorder fields
type BadLayout struct {
	A bool    // 1 byte
	// 7 bytes padding (align B to 8-byte boundary)
	B float64 // 8 bytes
	C bool    // 1 byte
	// 7 bytes padding
	// Total: 24 bytes
}

// GOOD: sorted by size descending
type GoodLayout struct {
	B float64 // 8 bytes
	C bool    // 1 byte
	A bool    // 1 byte
	// 6 bytes padding
	// Total: 16 bytes
}

// === Visualize struct layout ===

type StructInfo struct {
	FieldName string
	Offset    uintptr
	Size      uintptr
	TypeName  string
}

func inspectStruct(v any) []StructInfo {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	var info []StructInfo
	for i := range t.NumField() {
		f := t.Field(i)
		info = append(info, StructInfo{
			FieldName: f.Name,
			Offset:    f.Offset,
			Size:      f.Type.Size(),
			TypeName:  f.Type.String(),
		})
	}
	return info
}

// === Zero-copy string ↔ []byte (Go 1.20+) ===
// unsafe.String and unsafe.Slice allow conversion without copying

func stringToBytes(s string) []byte {
	// Go 1.20+: unsafe.StringData returns pointer to underlying bytes
	if len(s) == 0 {
		return nil
	}
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

func bytesToString(b []byte) string {
	// Go 1.20+: unsafe.String creates a string from pointer + length
	if len(b) == 0 {
		return ""
	}
	return unsafe.String(&b[0], len(b))
}

// === unsafe.Pointer type reinterpretation ===
// DANGEROUS: only do this if you fully understand the memory layout

func float64Bits(f float64) uint64 {
	// Reinterpret float64 as uint64 (bit manipulation)
	return *(*uint64)(unsafe.Pointer(&f))
}

func main() {
	fmt.Println("=== 1. unsafe.Sizeof, Alignof, Offsetof ===")

	var i8 int8
	var i16 int16
	var i32 int32
	var i64 int64
	var f64 float64
	var b bool
	var s string
	var sl []int

	fmt.Printf("  int8:    size=%d, align=%d\n", unsafe.Sizeof(i8), unsafe.Alignof(i8))
	fmt.Printf("  int16:   size=%d, align=%d\n", unsafe.Sizeof(i16), unsafe.Alignof(i16))
	fmt.Printf("  int32:   size=%d, align=%d\n", unsafe.Sizeof(i32), unsafe.Alignof(i32))
	fmt.Printf("  int64:   size=%d, align=%d\n", unsafe.Sizeof(i64), unsafe.Alignof(i64))
	fmt.Printf("  float64: size=%d, align=%d\n", unsafe.Sizeof(f64), unsafe.Alignof(f64))
	fmt.Printf("  bool:    size=%d, align=%d\n", unsafe.Sizeof(b), unsafe.Alignof(b))
	fmt.Printf("  string:  size=%d (ptr+len)\n", unsafe.Sizeof(s))
	fmt.Printf("  []int:   size=%d (ptr+len+cap)\n", unsafe.Sizeof(sl))

	fmt.Println("\n=== 2. Struct Padding ===")
	fmt.Printf("  BadLayout:  %d bytes\n", unsafe.Sizeof(BadLayout{}))
	fmt.Printf("  GoodLayout: %d bytes\n", unsafe.Sizeof(GoodLayout{}))
	fmt.Printf("  Savings: %d bytes per instance\n",
		unsafe.Sizeof(BadLayout{})-unsafe.Sizeof(GoodLayout{}))

	fmt.Println("\n  BadLayout fields:")
	for _, info := range inspectStruct(BadLayout{}) {
		fmt.Printf("    %-10s offset=%2d size=%d type=%s\n",
			info.FieldName, info.Offset, info.Size, info.TypeName)
	}
	fmt.Println("  GoodLayout fields:")
	for _, info := range inspectStruct(GoodLayout{}) {
		fmt.Printf("    %-10s offset=%2d size=%d type=%s\n",
			info.FieldName, info.Offset, info.Size, info.TypeName)
	}

	fmt.Println("\n=== 3. unsafe.Offsetof ===")
	type Example struct {
		A int32
		B int64
		C bool
	}
	var ex Example
	fmt.Printf("  Example.A offset: %d\n", unsafe.Offsetof(ex.A))
	fmt.Printf("  Example.B offset: %d\n", unsafe.Offsetof(ex.B))
	fmt.Printf("  Example.C offset: %d\n", unsafe.Offsetof(ex.C))
	fmt.Printf("  Example total size: %d bytes\n", unsafe.Sizeof(ex))

	fmt.Println("\n=== 4. Zero-copy string ↔ []byte (Go 1.20+) ===")
	original := "Hello, Go! 🎯"
	bytes := stringToBytes(original)
	fmt.Printf("  Original string: %q (len=%d bytes)\n", original, len(original))
	fmt.Printf("  As bytes (first 5): %v\n", bytes[:5])

	// WARNING: string is immutable — do not modify bytes!
	// bytes[0] = 'h'  // ← undefined behavior / segfault!

	bSlice := []byte("mutable bytes")
	str := bytesToString(bSlice)
	fmt.Printf("  []byte to string: %q\n", str)

	fmt.Println("\n=== 5. Float64 bit manipulation ===")
	f := 1.0
	bits := float64Bits(f)
	fmt.Printf("  1.0 as bits: 0x%016X\n", bits)
	fmt.Printf("  IEEE 754: sign=0, exponent=01111111111, mantissa=0...0\n")

	fmt.Println("\n=== 6. NGUYÊN TẮC sử dụng unsafe ===")
	fmt.Println("  ✓ Tối ưu zero-copy string/byte conversion (1.20+)")
	fmt.Println("  ✓ Kiểm tra struct layout/alignment")
	fmt.Println("  ✓ Interop với C qua CGo")
	fmt.Println("  ✗ Không dùng cho general type conversion")
	fmt.Println("  ✗ Không bao giờ modify immutable data (string)")
	fmt.Println("  ✗ Pointer arithmetic không được hỗ trợ trực tiếp")
	fmt.Println()
	fmt.Println("  Go 1.20+ unsafe additions:")
	fmt.Println("  - unsafe.String(ptr, len): *byte → string (no copy)")
	fmt.Println("  - unsafe.StringData(s): string → *byte")
	fmt.Println("  - unsafe.Slice(ptr, len): *T → []T")
	fmt.Println("  - unsafe.SliceData(s): []T → *T")
}
