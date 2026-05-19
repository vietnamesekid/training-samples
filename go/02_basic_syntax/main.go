// Lesson 2: Basic Syntax — variables, types, constants, iota, type conversion
// Run: go run .
package main

import (
	"fmt"
	"math"
	"unsafe"
)

// === Constants & iota ===

// iota increments within each const block and resets to 0 in a new block
type Weekday int

const (
	Monday    Weekday = iota + 1 // 1
	Tuesday                      // 2
	Wednesday                    // 3
	Thursday                     // 4
	Friday                       // 5
	Saturday                     // 6
	Sunday                       // 7
)

func (d Weekday) String() string {
	names := [...]string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
	if d < Monday || d > Sunday {
		return "Unknown"
	}
	return names[d-1]
}

// iota with bit shift — commonly used for permission flags
type Permission uint

const (
	Read    Permission = 1 << iota // 1   (001)
	Write                          // 2   (010)
	Execute                        // 4   (100)
)

func (p Permission) Has(flags Permission) bool {
	// "does p contain all bits of flags?"
	return (p & flags) == flags
}

func (p Permission) String() string {
	// Direct approach

	// result := ""
	// if p&Read != 0 {
	// 	result += "r"
	// } else {
	// 	result += "-"
	// }
	// if p&Write != 0 {
	// 	result += "w"
	// } else {
	// 	result += "-"
	// }
	// if p&Execute != 0 {
	// 	result += "x"
	// } else {
	// 	result += "-"
	// }
	// return result

	// Using helper method
	result := ""
	if p.Has(Read) {
		result += "r"
	} else {
		result += "-"
	}

	if p.Has(Write) {
		result += "w"
	} else {
		result += "-"
	}

	if p.Has(Execute) {
		result += "x"
	} else {
		result += "-"
	}

	return result
}

// HTTP status codes as a constant block
const (
	StatusOK                  = 200
	StatusCreated             = 201
	StatusBadRequest          = 400
	StatusUnauthorized        = 401
	StatusNotFound            = 404
	StatusInternalServerError = 500
)

func main() {
	fmt.Println("=== 1. Khai Báo Biến ===")

	// Method 1: explicit var declaration (used at package level or when an explicit type is needed)
	var name string = "Gopher"
	var age int = 10

	// Method 2: short declaration := (ONLY inside functions, most common)
	city := "Hanoi" // Go infers the type: string
	count := 0      // int
	pi := 3.14      // float64
	isReady := true // bool

	// Method 3: declare multiple variables at once
	var (
		host = "localhost"
		port = 8080
	)

	fmt.Printf("name=%s, age=%d, city=%s\n", name, age, city)
	fmt.Printf("count=%d, pi=%f, isReady=%t\n", count, pi, isReady)
	fmt.Printf("host=%s, port=%d\n", host, port)

	// Blank identifier _ — discard values you don't need
	x, _ := 10, 20 // discard the value 20
	fmt.Println("x =", x)

	fmt.Println("\n=== 2. Zero Values — Go does NOT have uninitialized variables ===")
	var (
		zeroInt    int     // 0
		zeroFloat  float64 // 0.0
		zeroBool   bool    // false
		zeroString string  // ""
		zeroPtr    *int    // nil
	)
	// var zeroSlice []int    — nil (different from an empty slice!)
	// var zeroMap map[string]int — nil (will panic on write)
	fmt.Printf("int=%d, float=%f, bool=%t, string=%q, ptr=%v\n",
		zeroInt, zeroFloat, zeroBool, zeroString, zeroPtr)

	fmt.Println("\n=== 3. Bảng Kiểu Dữ Liệu ===")

	// Integers: int, int8, int16, int32, int64
	var i8 int8 = 127          // -128 to 127
	var i16 int16 = 32767      // -32768 to 32767
	var i32 int32 = 2147483647 // ~2.1 billion
	var i64 int64 = math.MaxInt64
	var i int = 42 // 32 or 64-bit depending on platform

	// Unsigned: uint, uint8, uint16, uint32, uint64
	var u8 uint8 = 255 // = byte (0 to 255)
	var u16 uint16 = 65535
	var u32 uint32 = 4294967295

	// Float
	var f32 float32 = 3.14              // ~7 decimal digits
	var f64 float64 = 3.141592653589793 // ~15 decimal digits, default

	// Special types
	var b byte = 'A' // byte = uint8, ASCII value
	var r rune = '🎯' // rune = int32, 1 Unicode code point

	// Bool
	var flag bool = true

	fmt.Printf("int8=%d, int16=%d, int32=%d, int64=%d, int=%d\n", i8, i16, i32, i64, i)
	fmt.Printf("uint8=%d, uint16=%d, uint32=%d\n", u8, u16, u32)
	fmt.Printf("float32=%f, float64=%f\n", f32, f64)
	fmt.Printf("byte=%c(%d), rune=%c(%d), bool=%t\n", b, b, r, r, flag)

	// Size of types (platform-dependent for int/uint)
	fmt.Printf("\nunsafe.Sizeof: int=%d bytes, int64=%d bytes, float64=%d bytes\n",
		unsafe.Sizeof(i), unsafe.Sizeof(i64), unsafe.Sizeof(f64))

	fmt.Println("\n=== 4. Hằng Số & iota ===")
	fmt.Printf("Monday=%d (%s)\n", Monday, Monday)
	fmt.Printf("Friday=%d (%s)\n", Friday, Friday)

	perm := Read | Write
	fmt.Printf("Read|Write = %s (binary: %03b)\n", perm, uint(perm))
	fmt.Printf("Execute = %s (binary: %03b)\n", Execute, uint(Execute))
	fmt.Printf("Có quyền Write: %t\n", perm&Write != 0)
	fmt.Printf("Có quyền Execute: %t\n", perm&Execute != 0)

	fmt.Println("\n=== 5. Type Conversion ===")
	// Go does NOT have implicit conversion — must be explicit
	var intVal int = 42
	var float64Val float64 = float64(intVal) // int → float64
	var int32Val int32 = int32(intVal)       // int → int32
	var byteVal byte = byte(intVal)          // int → byte

	fmt.Printf("int=%d → float64=%f, int32=%d, byte=%c\n",
		intVal, float64Val, int32Val, byteVal)

	// String ↔ rune/byte
	s := "Hello"
	bytes := []byte(s)  // string → []byte (copy)
	runes := []rune(s)  // string → []rune (copy)
	s2 := string(bytes) // []byte → string (copy)
	s3 := string(runes) // []rune → string (copy)
	fmt.Printf("string=%q, bytes=%v, runes=%v\n", s, bytes, runes)
	fmt.Printf("back to string: %q, %q\n", s2, s3)

	// GOTCHA: int → string does NOT do what you might think
	n := 65
	fmt.Printf("string(65) = %q  ← đây là ký tự ASCII, không phải \"65\"!\n", string(rune(n)))
	fmt.Printf("Dùng fmt.Sprintf: %q\n", fmt.Sprintf("%d", n))

	fmt.Println("\n=== 6. fmt Format Verbs ===")
	type Point struct{ X, Y int }
	p := Point{3, 4}
	fmt.Printf("%%v  (default):   %v\n", p)
	fmt.Printf("%%+v (with field names): %+v\n", p)
	fmt.Printf("%%#v (Go syntax): %#v\n", p)
	fmt.Printf("%%T  (type):      %T\n", p)
	fmt.Printf("%%d  (decimal):   %d\n", 255)
	fmt.Printf("%%b  (binary):    %b\n", 255)
	fmt.Printf("%%o  (octal):     %o\n", 255)
	fmt.Printf("%%x  (hex lower): %x\n", 255)
	fmt.Printf("%%X  (hex upper): %X\n", 255)
	fmt.Printf("%%f  (float):     %f\n", 3.14159)
	fmt.Printf("%%.2f (2 decimal): %.2f\n", 3.14159)
	fmt.Printf("%%e  (scientific): %e\n", 3.14159)
	fmt.Printf("%%s  (string):    %s\n", "hello")
	fmt.Printf("%%q  (quoted):    %q\n", "hello")
	fmt.Printf("%%p  (pointer):   %p\n", &n)
	fmt.Printf("%%t  (bool):      %t\n", true)
	fmt.Printf("%%10d (width 10): %10d\n", 42)
	fmt.Printf("%%-10d (left):    %-10d|\n", 42)
	fmt.Printf("%%010d (zero pad): %010d\n", 42)
}
