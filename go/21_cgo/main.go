// Bài 21: CGo — gọi code C từ Go
// Chạy: go run .
// Build static binary: CGO_ENABLED=0 go build .  (không cần C)
// Build với CGo: go build .
//
// Điều kiện: cần C compiler (gcc hoặc clang)
// macOS: xcode-select --install
// Linux: apt-get install build-essential
package main

/*
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <math.h>

// Hàm C đơn giản
int add(int a, int b) {
    return a + b;
}

double square_root(double x) {
    return sqrt(x);
}

// Hàm nhận và trả về C string
char* greet(const char* name) {
    static char buf[256];
    snprintf(buf, sizeof(buf), "Hello, %s from C!", name);
    return buf;
}

// Struct trong C
typedef struct {
    int x;
    int y;
} Point;

Point make_point(int x, int y) {
    Point p;
    p.x = x;
    p.y = y;
    return p;
}

double distance(Point a, Point b) {
    int dx = a.x - b.x;
    int dy = a.y - b.y;
    return sqrt(dx*dx + dy*dy);
}
*/
import "C" // QUAN TRỌNG: không có blank line giữa comment và import "C"!

import (
	"fmt"
	"unsafe"
)

func main() {
	fmt.Println("=== 1. Gọi Hàm C Đơn Giản ===")

	// Gọi C.add — tự động convert Go int → C int
	result := C.add(3, 4)
	fmt.Printf("  C.add(3, 4) = %d\n", int(result))

	// Gọi C.square_root với C.double
	sqrt := C.square_root(2.0)
	fmt.Printf("  C.square_root(2.0) = %f\n", float64(sqrt))

	fmt.Println("\n=== 2. Go String ↔ C String ===")
	// QUAN TRỌNG: C.CString allocate memory — phải gọi C.free!
	cName := C.CString("Gopher") // Go string → C char*
	defer C.free(unsafe.Pointer(cName)) // PHẢI free để tránh memory leak

	greeting := C.greet(cName)
	// C.GoString: C char* → Go string (copy)
	fmt.Printf("  %s\n", C.GoString(greeting))

	fmt.Println("\n=== 3. C Struct ===")
	p1 := C.make_point(0, 0)
	p2 := C.make_point(3, 4)
	fmt.Printf("  p1 = {%d, %d}\n", int(p1.x), int(p1.y))
	fmt.Printf("  p2 = {%d, %d}\n", int(p2.x), int(p2.y))
	dist := C.distance(p1, p2)
	fmt.Printf("  distance(p1, p2) = %.1f\n", float64(dist))

	fmt.Println("\n=== 4. Conversion Types ===")
	fmt.Println("  Go int     → C.int:    C.int(goInt)")
	fmt.Println("  C.int      → Go int:   int(cInt)")
	fmt.Println("  Go string  → C char*:  C.CString(s) — PHẢI free!")
	fmt.Println("  C char*    → Go string: C.GoString(cs) — copy")
	fmt.Println("  Go []byte  → C void*:  unsafe.Pointer(&slice[0])")
	fmt.Println("  C void*    → Go ptr:   unsafe.Pointer(cPtr)")

	fmt.Println("\n=== 5. CGo Overhead Warning ===")
	fmt.Println("  CGo call overhead: ~20-100x so với pure Go function call")
	fmt.Println("  Nguyên nhân: context switch, argument marshaling")
	fmt.Println()
	fmt.Println("  Khi nào dùng CGo:")
	fmt.Println("  ✓ Wrapper cho C libraries (OpenSSL, SQLite, BLAS)")
	fmt.Println("  ✓ System calls không có Go bindings")
	fmt.Println("  ✓ Legacy C code integration")
	fmt.Println()
	fmt.Println("  Khi nào KHÔNG dùng CGo:")
	fmt.Println("  ✗ Hot paths — overhead quá lớn")
	fmt.Println("  ✗ Cross-compilation phức tạp hơn nhiều")
	fmt.Println("  ✗ CGO_ENABLED=0 không hoạt động")

	fmt.Println("\n=== 6. CGO_ENABLED ===")
	fmt.Println("  CGO_ENABLED=1 (default): build với CGo support")
	fmt.Println("  CGO_ENABLED=0: pure Go, dễ cross-compile, static binary")
	fmt.Println("  go build -ldflags='-extldflags=-static': fully static binary")
}
