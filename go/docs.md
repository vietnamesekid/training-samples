# 🎯 MASTER GO — Tài Liệu Học Toàn Diện
> **Phiên bản tham chiếu:** Go 1.25 (August 2025) — phiên bản mới nhất  
> **Mục tiêu:** Zero → Go Expert  
> **Ngôn ngữ:** Tiếng Việt, code và thuật ngữ kỹ thuật giữ nguyên tiếng Anh

---

## 📋 MỤC LỤC

1. [Triết Lý Go & Lịch Sử](#1-triết-lý-go--lịch-sử)
2. [Cài Đặt & Toolchain](#2-cài-đặt--toolchain)
3. [Syntax Cơ Bản](#3-syntax-cơ-bản)
4. [Strings, Arrays, Slices, Maps](#4-strings-arrays-slices-maps)
5. [Structs & Methods](#5-structs--methods)
6. [Control Flow & Functions](#6-control-flow--functions)
7. [Interfaces](#7-interfaces)
8. [Error Handling](#8-error-handling)
9. [Packages & Modules](#9-packages--modules)
10. [Pointers](#10-pointers)
11. [Goroutines & Channels](#11-goroutines--channels)
12. [Sync Primitives](#12-sync-primitives)
13. [Context](#13-context)
14. [Race Conditions](#14-race-conditions)
15. [Advanced Concurrency Patterns](#15-advanced-concurrency-patterns)
16. [HTTP Server](#16-http-server)
17. [Testing](#17-testing)
18. [Profiling & Performance](#18-profiling--performance)
19. [Logging với slog](#19-logging-với-slog)
20. [Generics](#20-generics)
21. [Reflection](#21-reflection)
22. [unsafe Package](#22-unsafe-package)
23. [CGo](#23-cgo)
24. [Design Patterns trong Go](#24-design-patterns-trong-go)
25. [Microservices & gRPC](#25-microservices--grpc)
26. [Database Patterns](#26-database-patterns)
27. [Go Runtime Internals — GMP Model](#27-go-runtime-internals--gmp-model)
28. [Memory Management & GC](#28-memory-management--gc)
29. [Performance Optimization](#29-performance-optimization)
30. [Tools Ecosystem](#30-tools-ecosystem)
31. [Go 1.22 → 1.25 — Tính Năng Mới Quan Trọng](#31-go-122--125--tính-năng-mới-quan-trọng)
32. [Production Checklist](#32-production-checklist)
33. [Những Sai Lầm Phổ Biến Nhất](#33-những-sai-lầm-phổ-biến-nhất)
34. [Bài Tập Thực Hành](#34-bài-tập-thực-hành)
35. [Tài Nguyên Học Tập](#35-tài-nguyên-học-tập)
36. [Milestones & Checkpoints](#36-milestones--checkpoints)

---

## 1. Triết Lý Go & Lịch Sử

Go được tạo ra năm 2007 bởi **Robert Griesemer, Rob Pike, và Ken Thompson** tại Google, ra mắt công khai năm 2009. Vấn đề họ muốn giải quyết:
- C++ compile quá chậm, quá phức tạp
- Java/Python không đủ performance cho hệ thống lớn
- Cần ngôn ngữ mà engineer mới có thể productive ngay

### Triết Lý Thiết Kế Core

| Nguyên tắc | Ý nghĩa |
|---|---|
| **Simplicity over cleverness** | Một cách làm đúng, không nhiều cách |
| **Explicit over implicit** | Không có magic, không có hidden behavior |
| **Composition over inheritance** | Không có class hierarchy, dùng interfaces |
| **Concurrency as first-class citizen** | Goroutines và channels built-in |
| **Fast compilation** | Compile cả project lớn trong vài giây |

> **WHY:** Go được thiết kế để giải quyết các vấn đề thực tế ở Google quy mô lớn: hàng ngàn kỹ sư, hàng triệu dòng code, hệ thống distributed phức tạp. Triết lý "boring language" là intentional — code Go dễ đọc hơn là dễ viết.

---

## 2. Cài Đặt & Toolchain

```bash
# Tải Go tại: https://go.dev/dl/
# Kiểm tra cài đặt
go version          # go version go1.25.x linux/amd64
go env              # xem tất cả environment variables
go env GOPATH       # $HOME/go
go env GOMODCACHE   # nơi cache modules
go env GOPROXY      # proxy để tải dependencies

# Cấu trúc workspace (dự án hiện đại dùng Go Modules từ Go 1.11)
mkdir myproject && cd myproject
go mod init github.com/yourname/myproject
# Tạo ra file go.mod

# Các lệnh thiết yếu
go run main.go          # compile và chạy ngay
go build ./...          # compile tất cả packages
go test ./...           # chạy tất cả tests
go fmt ./...            # format code (BẮT BUỘC)
go vet ./...            # linting cơ bản
go doc fmt.Println      # xem documentation
go mod tidy             # dọn dẹp dependencies
go clean -cache         # xóa build cache
```

### Cấu Trúc Dự Án Chuẩn

```
myapp/
├── cmd/
│   ├── server/
│   │   └── main.go          # entry point cho server
│   └── cli/
│       └── main.go          # entry point cho CLI tool
├── internal/                # package chỉ dùng trong module này
│   ├── user/
│   │   ├── user.go
│   │   ├── user_test.go
│   │   ├── repository.go
│   │   └── service.go
│   ├── order/
│   └── config/
├── pkg/                     # package có thể import bởi người khác
│   ├── validator/
│   └── logger/
├── api/                     # Protocol definitions (proto, OpenAPI)
├── docs/
├── scripts/
├── go.mod
├── go.sum
└── Makefile
```

---

## 3. Syntax Cơ Bản

### Biến & Kiểu Dữ Liệu

```go
package main

import "fmt"

func main() {
    // Khai báo tường minh
    var name string = "Gopher"
    var age int = 10

    // Short declaration — CHỈ dùng trong function
    city := "Hanoi"        // Go tự suy luận kiểu (type inference)
    count := 0             // int
    pi := 3.14             // float64
    isReady := true        // bool

    // Zero values — QUAN TRỌNG: Go không có uninitialized variables
    var x int              // 0
    var s string           // ""
    var b bool             // false
    var p *int             // nil
    var sl []int           // nil (không phải empty slice!)
    var m map[string]int   // nil (không phải empty map!)

    // Hằng số
    const MaxRetries = 3
    const (
        StatusOK    = 200
        StatusError = 500
    )

    // iota — tạo hằng số tuần tự
    type Weekday int
    const (
        Monday Weekday = iota + 1  // 1
        Tuesday                     // 2
        Wednesday                   // 3
        Thursday                    // 4
        Friday                      // 5
    )

    // iota với bit shift — phổ biến cho flags
    type Permission uint
    const (
        Read    Permission = 1 << iota  // 1 (001)
        Write                           // 2 (010)
        Execute                         // 4 (100)
    )

    fmt.Println(name, age, city, count, pi, isReady)
    fmt.Println(x, s, b, p, sl, m)
}
```

### Bảng Kiểu Dữ Liệu Đầy Đủ

| Category | Types | Ghi chú |
|---|---|---|
| Integer | `int`, `int8`, `int16`, `int32`, `int64` | `int` = 32 hoặc 64-bit tùy platform |
| Unsigned | `uint`, `uint8`, `uint16`, `uint32`, `uint64`, `uintptr` | |
| Float | `float32`, `float64` | Mặc định `float64` |
| Complex | `complex64`, `complex128` | |
| String | `string` | Immutable, UTF-8 |
| Boolean | `bool` | |
| Alias | `byte` (= uint8), `rune` (= int32) | rune = 1 Unicode code point |

> **GOTCHA:** `int` có size phụ thuộc platform (32-bit hoặc 64-bit). Khi cần chính xác, dùng `int64`.

---

## 4. Strings, Arrays, Slices, Maps

### Strings — Hiểu Sâu

```go
import (
    "fmt"
    "strings"
    "unicode/utf8"
)

s := "Xin chào Việt Nam 🇻🇳"

// len() trả về số BYTES, không phải số ký tự!
fmt.Println(len(s))                          // số bytes (lớn hơn số ký tự)
fmt.Println(utf8.RuneCountInString(s))       // số ký tự Unicode thực sự

// Iterate đúng cách với rune (Unicode code points)
for i, r := range s {
    fmt.Printf("index=%d, rune=%c, value=%d\n", i, r, r)
}

// String không thể mutate — dùng strings.Builder
var builder strings.Builder
builder.Grow(estimatedSize)    // pre-allocate để tránh reallocations
for i := 0; i < 5; i++ {
    builder.WriteString("item")
}
result := builder.String()

// Chuyển đổi (đều tạo copy)
b := []byte(s)     // string → []byte
s2 := string(b)    // []byte → string

// Thao tác phổ biến
strings.ToUpper("hello")
strings.Contains("golang", "go")
strings.Split("a,b,c", ",")          // ["a", "b", "c"]
strings.Join([]string{"a","b"}, "-") // "a-b"
strings.TrimSpace("  hello  ")
strings.HasPrefix("golang", "go")
strings.HasSuffix("golang", "lang")
strings.ReplaceAll("foo foo", "foo", "bar")
strings.Count("cheese", "e")          // 3
fmt.Sprintf("Hello, %s! You are %d years old.", "Alice", 30)
```

### Arrays — Fixed Size, Value Type

```go
// Array có fixed size, là VALUE TYPE (copy khi gán)
arr := [3]int{1, 2, 3}
arr2 := [...]int{4, 5, 6}     // compiler đếm size

a1 := [3]int{1, 2, 3}
a2 := a1                       // COPY toàn bộ array
a2[0] = 99                     // a1[0] vẫn là 1

// So sánh được nếu element type comparable
fmt.Println([3]int{1,2,3} == [3]int{1,2,3})  // true
```

### Slices — Dynamic, Reference Type

```go
// Slice internal: struct{ptr *T, len int, cap int}
// LÀ REFERENCE TYPE — share underlying array!

s := []int{1, 2, 3}
s = append(s, 4, 5)            // thêm elements

// Tạo slice với make (KHUYẾN KHÍCH pre-allocate)
s3 := make([]int, 5)           // len=5, cap=5, zero values
s4 := make([]int, 0, 10)       // len=0, cap=10 — pre-allocate

// Slicing — CHÚ Ý: share memory!
a := []int{1, 2, 3, 4, 5}
b := a[1:3]      // b = [2, 3], nhưng SHARE underlying array với a
b[0] = 99        // a[1] cũng thành 99! ← GOTCHA phổ biến

// Fix: dùng full slice expression để giới hạn capacity
b2 := a[1:3:3]   // capacity bị giới hạn, append sẽ tạo array mới

// Copy — deep copy
dst := make([]int, len(a))
n := copy(dst, a)  // trả về số elements copied

// Xóa element tại index i (không giữ thứ tự)
a = append(a[:i], a[i+1:]...)

// Xóa element tại index i (giữ thứ tự) — dùng slices.Delete Go 1.21+
import "slices"
a = slices.Delete(a, i, i+1)

// Thêm vào giữa
a = append(a[:i], append([]int{newVal}, a[i:]...)...)

// 2D slice
matrix := make([][]int, rows)
for i := range matrix {
    matrix[i] = make([]int, cols)
}
```

### Maps

```go
// Map không đảm bảo thứ tự khi iterate
m := map[string]int{
    "alice": 30,
    "bob":   25,
}

// Tạo với make
m2 := make(map[string]int)
m2["key"] = 100

// CRUD
m["charlie"] = 35
delete(m, "bob")

// Kiểm tra key tồn tại — LUÔN dùng 2-value form
val, ok := m["alice"]
if !ok {
    // key không tồn tại
}

// Iterate — thứ tự RANDOM mỗi lần
for k, v := range m {
    fmt.Printf("%s: %d\n", k, v)
}

// Iterate theo thứ tự: sort keys trước
import "sort"
keys := make([]string, 0, len(m))
for k := range m {
    keys = append(keys, k)
}
sort.Strings(keys)
for _, k := range keys {
    fmt.Println(k, m[k])
}

// Map với struct values — không thể set field trực tiếp
type Point struct{ X, Y int }
points := map[string]Point{"a": {1, 2}}
// points["a"].X = 10  // ERROR: cannot assign to struct field in map
// Fix: copy, modify, assign back
p := points["a"]
p.X = 10
points["a"] = p
// Hoặc: dùng *Point
```

---

## 5. Structs & Methods

```go
type Person struct {
    Name    string
    Age     int
    Address Address
    score   int     // unexported — chỉ access trong package
}

type Address struct {
    Street string
    City   string
}

// Constructor pattern (Go không có constructor)
func NewPerson(name string, age int) *Person {
    if name == "" {
        panic("name cannot be empty")
    }
    return &Person{Name: name, Age: age}
}

// Value receiver — không mutate, nhận COPY
func (p Person) String() string {
    return fmt.Sprintf("%s (%d tuổi)", p.Name, p.Age)
}

// Pointer receiver — có thể mutate, nhận POINTER
func (p *Person) Birthday() {
    p.Age++
}

// NGUYÊN TẮC: Nếu bất kỳ method nào cần pointer receiver,
// hãy dùng pointer receiver cho TẤT CẢ methods của type đó.

// Embedded struct — composition
type Employee struct {
    Person           // promoted: Employee.Birthday() = Person.Birthday()
    Company string
    Salary  float64
}

func main() {
    p := NewPerson("Alice", 30)
    p.Birthday()
    fmt.Println(p)   // gọi String() tự động qua fmt.Stringer interface

    e := Employee{
        Person:  Person{Name: "Bob", Age: 25},
        Company: "Acme Corp",
    }
    e.Birthday()     // gọi Person.Birthday() qua promotion
    fmt.Println(e.Name)  // access Person.Name trực tiếp

    // Struct literal — luôn đặt field names (tránh vỡ khi thêm field)
    addr := Address{Street: "123 Main St", City: "Hanoi"}
    _ = addr

    // Anonymous struct
    point := struct{ X, Y int }{X: 10, Y: 20}
    fmt.Println(point)

    // Struct tags
    type JSONUser struct {
        ID    int    `json:"id"`
        Name  string `json:"name"`
        Email string `json:"email,omitempty"`
        Pass  string `json:"-"`           // bỏ qua khi marshal
    }
}
```

---

## 6. Control Flow & Functions

### Control Flow

```go
// If — không cần parentheses
if x > 0 {
    // ...
} else if x < 0 {
    // ...
}

// If với initialization statement — biến chỉ sống trong block
if err := doSomething(); err != nil {
    return err
}

// For — Go chỉ có FOR, thay thế for/while/do-while
for i := 0; i < 10; i++ { }           // C-style
for condition { }                       // while-style
for { }                                 // infinite loop
for i, v := range slice { }            // range slice
for k, v := range myMap { }            // range map
for i := range slice { }               // chỉ index
for _, v := range slice { }            // chỉ value
for v := range channel { }             // receive từ channel đến khi close
for i := range 5 { }                   // Go 1.22+: range integer

// Switch — không cần break (implicit break)
switch os := runtime.GOOS; os {
case "linux":
    fmt.Println("Linux")
case "darwin":
    fmt.Println("macOS")
default:
    fmt.Printf("Other: %s\n", os)
}

// Dùng fallthrough để tiếp tục case tiếp theo
switch x {
case 1:
    fmt.Println("one")
    fallthrough
case 2:
    fmt.Println("one or two")
}

// Switch không có expression — như if-else chain
switch {
case x < 0:
    fmt.Println("negative")
case x == 0:
    fmt.Println("zero")
default:
    fmt.Println("positive")
}

// Type switch
switch v := i.(type) {
case int:
    fmt.Printf("int: %d\n", v)
case string:
    fmt.Printf("string: %s\n", v)
case fmt.Stringer:
    fmt.Printf("Stringer: %s\n", v.String())
default:
    fmt.Printf("unknown: %T\n", v)
}

// Defer — LIFO order, arguments evaluated IMMEDIATELY
func readFile(name string) error {
    f, err := os.Open(name)
    if err != nil {
        return err
    }
    defer f.Close()    // đảm bảo close kể cả khi panic

    defer fmt.Println("first to defer = last to run")
    defer fmt.Println("second to defer = second to last")
    // ...
    return nil
}

// Defer với closure capture — WATCH OUT!
for i := 0; i < 3; i++ {
    defer fmt.Println(i)       // i được evaluate ngay: prints 0,1,2 (LIFO: 2,1,0)
    // vs:
    defer func() { fmt.Println(i) }() // capture by reference: prints 3,3,3 sau loop
}

// Labeled break/continue
outer:
for i := 0; i < 5; i++ {
    for j := 0; j < 5; j++ {
        if i+j == 6 {
            break outer  // break khỏi cả 2 vòng lặp
        }
    }
}
```

### Functions

```go
// Function cơ bản
func add(a, b int) int {
    return a + b
}

// Multiple return values — ĐẶC TRƯNG CỦA GO
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, fmt.Errorf("division by zero")
    }
    return a / b, nil
}

// Named return values — dùng cẩn thận (chỉ cho function ngắn)
func minMax(arr []int) (min, max int) {
    min, max = arr[0], arr[0]
    for _, v := range arr[1:] {
        if v < min { min = v }
        if v > max { max = v }
    }
    return  // naked return — TRÁNH dùng cho function dài
}

// Variadic function
func sum(nums ...int) int {
    total := 0
    for _, n := range nums {
        total += n
    }
    return total
}
sum(1, 2, 3)
nums := []int{1, 2, 3}
sum(nums...)    // unpack slice

// Function là first-class value
type MathFunc func(int, int) int

func apply(f MathFunc, a, b int) int { return f(a, b) }

multiply := func(a, b int) int { return a * b }
apply(multiply, 3, 4)  // 12
apply(func(a, b int) int { return a + b }, 1, 2)  // 3

// Closure — capture biến từ outer scope
func counter(start int) func() int {
    count := start
    return func() int {
        count++
        return count
    }
}
c := counter(10)
fmt.Println(c())  // 11
fmt.Println(c())  // 12

// Panic và Recover
func safeDivide(a, b int) (result int, err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("recovered from panic: %v", r)
        }
    }()
    return a / b, nil
}

// NGUYÊN TẮC: Chỉ dùng panic cho lỗi programming không thể recover được.
// Dùng errors cho runtime errors.
func mustPositive(n int) int {
    if n <= 0 {
        panic(fmt.Sprintf("expected positive, got %d", n))
    }
    return n
}
```

---

## 7. Interfaces

> **Đây là concept quan trọng nhất trong Go.** Hiểu interfaces sâu = hiểu Go.

### Implicit Interface — Duck Typing

```go
// Interface = tập hợp method signatures
// Go dùng IMPLICIT INTERFACE — không cần khai báo "implements"
type Animal interface {
    Sound() string
    Move() string
}

type Dog struct{ Name string }
type Bird struct{ Name string }

func (d Dog) Sound() string { return "Woof" }
func (d Dog) Move() string  { return "Run" }
func (b Bird) Sound() string { return "Tweet" }
func (b Bird) Move() string  { return "Fly" }

// Dog và Bird tự động implement Animal
func describe(a Animal) {
    fmt.Printf("Sound: %s, Move: %s\n", a.Sound(), a.Move())
}

describe(Dog{Name: "Rex"})
describe(Bird{Name: "Tweety"})
```

### Interface Composition

```go
// Standard library interfaces quan trọng nhất
type Reader interface {
    Read(p []byte) (n int, err error)
}
type Writer interface {
    Write(p []byte) (n int, err error)
}
type Closer interface {
    Close() error
}

// Interface embed interface khác
type ReadWriter interface {
    Reader
    Writer
}
type ReadWriteCloser interface {
    Reader
    Writer
    Closer
}

// fmt.Stringer — implement để customize fmt.Println
type Stringer interface {
    String() string
}

// error — interface có 1 method
type error interface {
    Error() string
}
```

### Empty Interface & Type Assertions

```go
// any (= interface{} từ Go 1.18+) — chứa bất kỳ giá trị nào
func printAnything(v any) {
    // Type assertion với safety check
    if s, ok := v.(string); ok {
        fmt.Println("String:", s)
        return
    }

    // Type switch
    switch val := v.(type) {
    case int:
        fmt.Printf("int: %d\n", val)
    case []int:
        fmt.Printf("[]int: %v\n", val)
    case fmt.Stringer:
        fmt.Printf("Stringer: %s\n", val.String())
    default:
        fmt.Printf("unknown: %T = %v\n", val, val)
    }
}

// NGUY HIỂM: assertion không có safety check → panic nếu sai type
s := v.(string)  // panic nếu v không phải string!
```

### Interface Internals — Nil Gotcha

```go
// Interface value = (type, value)
// Hai components: pointer đến type descriptor + pointer đến data

var err error = nil                // (nil, nil)       → err == nil → TRUE
var p *MyError = nil
var err2 error = p                  // (*MyError, nil)  → err2 == nil → FALSE!

// Đây là Go gotcha NGUY HIỂM NHẤT về interfaces
func getError() error {
    var p *MyError = nil
    if someCondition {
        return p    // WRONG! trả về non-nil interface chứa nil pointer
    }
    return nil      // RIGHT! trả về nil interface
}

// FIX:
func getErrorFixed() error {
    if someCondition {
        return &MyError{...}  // return concrete error
    }
    return nil
}
```

### Functional Options Pattern

```go
type Server struct {
    host    string
    port    int
    timeout time.Duration
}

type Option func(*Server)

func WithPort(port int) Option {
    return func(s *Server) { s.port = port }
}

func WithTimeout(d time.Duration) Option {
    return func(s *Server) { s.timeout = d }
}

func NewServer(host string, opts ...Option) *Server {
    s := &Server{host: host, port: 8080, timeout: 30 * time.Second}
    for _, opt := range opts {
        opt(s)
    }
    return s
}

// Usage
server := NewServer("localhost",
    WithPort(9090),
    WithTimeout(60*time.Second),
)
```

### Interface Best Practices

```go
// 1. Accept interfaces, return concrete types
func NewBufferedReader(r io.Reader) *bufio.Reader { ... }

// 2. Small interfaces tốt hơn big interfaces (ISP)
// GOOD:
type Saver interface { Save(ctx context.Context) error }
type Loader interface { Load(ctx context.Context, id int) error }
// BAD:
type Repository interface {
    Save(ctx context.Context) error
    Load(ctx context.Context, id int) error
    Delete(ctx context.Context, id int) error
    List(ctx context.Context) ([]any, error)
    // ... 10 more methods
}

// 3. Define interfaces tại nơi DÙNG, không phải nơi implement
// Package user KHÔNG define UserRepository interface
// Package service DEFINE interface nó cần:
package service
type UserStore interface {
    GetUser(id int) (*User, error)
}
```

---

## 8. Error Handling

### Error Là Values — Không Có Exceptions

```go
// error interface cơ bản
type error interface {
    Error() string
}

// Cách tạo errors
err1 := errors.New("something went wrong")
err2 := fmt.Errorf("failed to connect to %s: %w", host, err1)  // %w = wrap

// Custom error types
type ValidationError struct {
    Field   string
    Message string
}
func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation failed for %s: %s", e.Field, e.Message)
}

// Error với Unwrap để chain
type DBError struct {
    Op  string
    Err error
}
func (e *DBError) Error() string { return fmt.Sprintf("db.%s: %v", e.Op, e.Err) }
func (e *DBError) Unwrap() error { return e.Err }

// errors.Is — kiểm tra error trong chain (Go 1.13+)
var ErrNotFound = errors.New("not found")
err := fmt.Errorf("findUser: %w", ErrNotFound)
if errors.Is(err, ErrNotFound) {       // traverse toàn bộ chain
    fmt.Println("User not found")
}

// errors.As — extract concrete error type từ chain (Go 1.13+)
var valErr *ValidationError
if errors.As(err, &valErr) {
    fmt.Println("Field:", valErr.Field)
}

// Sentinel errors — dùng khi cần compare
var (
    ErrTimeout     = errors.New("timeout")
    ErrPermission  = errors.New("permission denied")
    ErrUnavailable = errors.New("service unavailable")
)
```

### Error Handling Patterns

```go
// Pattern 1: Early return (Go idiomatic style)
func process(data []byte) (Result, error) {
    if len(data) == 0 {
        return Result{}, errors.New("empty data")
    }
    parsed, err := parse(data)
    if err != nil {
        return Result{}, fmt.Errorf("process: parse: %w", err)
    }
    validated, err := validate(parsed)
    if err != nil {
        return Result{}, fmt.Errorf("process: validate: %w", err)
    }
    return validated, nil
}

// Pattern 2: Error wrapping với context đầy đủ
// GOOD: "process: parse: json: cannot unmarshal string into Go value"
// BAD: "invalid input" (không có context)

// Pattern 3: Error type cho HTTP mapping
func toHTTPStatus(err error) int {
    switch {
    case errors.Is(err, ErrNotFound):
        return http.StatusNotFound
    case errors.Is(err, ErrPermission):
        return http.StatusForbidden
    case errors.Is(err, ErrTimeout):
        return http.StatusGatewayTimeout
    default:
        return http.StatusInternalServerError
    }
}

// Pattern 4: Accumulate errors (validation)
type MultiError []error

func (m MultiError) Error() string {
    msgs := make([]string, len(m))
    for i, e := range m {
        msgs[i] = e.Error()
    }
    return strings.Join(msgs, "; ")
}

// Pattern 5: errors.Join (Go 1.20+)
err = errors.Join(err1, err2, err3)
```

---

## 9. Packages & Modules

### Package Organization

```go
// Package naming: lowercase, single word, no underscores
package user    // GOOD
package userService  // BAD — dùng user/service thay

// Exported = viết hoa chữ cái đầu
// Unexported = viết thường (chỉ access trong package)
type User struct {         // exported
    ID   int
    Name string
    age  int               // unexported
}

func NewUser(name string) *User { ... }  // exported constructor
func (u *User) GetAge() int { return u.age }
```

### go.mod — Quản Lý Dependencies

```
module github.com/yourname/myapp

go 1.25

require (
    github.com/gin-gonic/gin v1.9.1
    go.uber.org/zap v1.27.0
)
```

```bash
# Thêm dependency
go get github.com/some/package@v1.2.3
go get github.com/some/package@latest

# Upgrade
go get -u ./...           # upgrade all
go get -u=patch ./...     # chỉ patch versions

# Tool dependencies (Go 1.24+)
go get -tool golang.org/x/tools/cmd/stringer
go tool stringer -type=Status  # run tool

# Xem dependency graph
go mod graph

# Vendor mode (offline, deterministic builds)
go mod vendor
go build -mod=vendor ./...

# Kiểm tra vulnerabilities
govulncheck ./...
```

### init() Function

```go
// init() chạy tự động khi package được import
// Có thể có nhiều init() trong một package (chạy theo thứ tự)
// init() chạy SAU package-level variables được init

var db *sql.DB

func init() {
    var err error
    db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
    if err != nil {
        log.Fatal("cannot connect to db:", err)
    }
}

// Import side effect (chạy init() mà không dùng package)
import _ "net/http/pprof"     // register pprof handlers
import _ "github.com/lib/pq"  // register postgres driver
```

---

## 10. Pointers

```go
// Pointer cơ bản
x := 42
p := &x          // p là *int, chứa địa chỉ của x
*p = 100         // dereference — thay đổi x
fmt.Println(x)   // 100

// Khi nào NÊN dùng pointer:
// 1. Mutate giá trị trong function
// 2. Large structs (tránh expensive copy)
// 3. Optional values (nil = "không có")
// 4. Shared mutable state (cẩn thận races!)

// Khi nào KHÔNG nên dùng pointer:
// 1. Small values (int, bool, float)
// 2. Immutable data
// 3. Khi không cần mutate

// nil pointer → panic khi dereference
var ptr *int     // nil
// *ptr = 5     // PANIC: nil pointer dereference

// new() — zero value + pointer
p2 := new(int)   // *int → 0
p3 := new([]string)  // *[]string → nil slice

// Pointer to struct
cfg := &Config{Port: 8080}
cfg.Debug = true    // Go tự dereference: (*cfg).Debug = true

// Không có pointer arithmetic trong Go (unlike C)
// Đây là design decision để memory safety

// Hàm nhận và trả về pointer
func newIntPointer(v int) *int {
    return &v  // SAFE: v escapes to heap, không bị destroy sau return
}
```

---

## 11. Goroutines & Channels

### Goroutines

```go
// Goroutine = lightweight thread, managed by Go runtime
// Cost: ~2KB stack ban đầu (tự grow, max ~1GB)
// OS Thread: 1-8MB stack — goroutines rẻ hơn rất nhiều!
// Go có thể chạy hàng triệu goroutines đồng thời

go func() {
    fmt.Println("Hello from goroutine")
}()

// QUAN TRỌNG: main() kết thúc → tất cả goroutines bị kill
// Cần synchronization!

// WaitGroup — đợi nhiều goroutines hoàn thành
var wg sync.WaitGroup
for i := 0; i < 5; i++ {
    wg.Add(1)
    go func(id int) {
        defer wg.Done()
        fmt.Printf("Worker %d done\n", id)
    }(i)  // TRUYỀN i làm argument, KHÔNG capture trực tiếp (Go < 1.22)
}
wg.Wait()

// Go 1.25+: WaitGroup.Go — đơn giản hơn nhiều!
var wg sync.WaitGroup
for i := 0; i < 5; i++ {
    id := i  // capture cho Go < 1.22
    wg.Go(func() {
        fmt.Printf("Worker %d done\n", id)
    })
}
wg.Wait()

// Go 1.22+: Loop variable mỗi iteration là independent
// for i := range 5 { go func() { fmt.Println(i) }() } — OK từ 1.22

// Goroutine leak — NGUY HIỂM
func leakyFunc() {
    ch := make(chan int)
    go func() {
        val := <-ch    // blocked mãi mãi nếu không ai gửi → LEAK
        fmt.Println(val)
    }()
    // return mà không send vào ch hay close ch → goroutine bị leak
}
// FIX: luôn đảm bảo goroutines có exit condition
```

### Channels

```go
// Channel = typed conduit — "communicate by sharing, don't share to communicate"

// Unbuffered channel — synchronous
ch := make(chan int)
go func() { ch <- 42 }()    // send (block cho đến khi có receiver)
val := <-ch                   // receive (block cho đến khi có sender)

// Buffered channel — asynchronous
buffered := make(chan string, 3)  // buffer size = 3
buffered <- "a"    // không block
buffered <- "b"    // không block
buffered <- "c"    // không block
// buffered <- "d" // BLOCK — buffer đầy

// Đóng channel
close(ch)
val, ok := <-ch    // ok = false khi channel đóng và rỗng
for v := range ch { // range tự dừng khi channel đóng
    fmt.Println(v)
}

// NGUYÊN TẮC: chỉ sender mới nên close channel
// Close channel đã close → PANIC

// Directional channels — restrict access
func producer(out chan<- int) {    // chỉ send
    for i := 0; i < 5; i++ {
        out <- i
    }
    close(out)
}

func consumer(in <-chan int) {     // chỉ receive
    for v := range in {
        fmt.Println(v)
    }
}

// Select — multiplex channels (như switch nhưng cho channels)
select {
case msg1 := <-ch1:
    fmt.Println("ch1:", msg1)
case msg2 := <-ch2:
    fmt.Println("ch2:", msg2)
case <-time.After(3 * time.Second):
    fmt.Println("timeout!")
    return
default:
    // non-blocking: chạy nếu không có case nào ready
    fmt.Println("no message ready")
}

// Channel patterns quan trọng

// Done channel — signal cancellation
done := make(chan struct{})  // struct{} = zero bytes!
go func() {
    defer close(done)
    // do work...
}()
// Caller:
<-done    // wait for completion

// Pipeline
func generate(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for _, n := range nums {
            out <- n
        }
    }()
    return out
}

func square(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for n := range in {
            out <- n * n
        }
    }()
    return out
}

// Usage: square(square(generate(2, 3, 4)))
// Output: 16, 81, 256
```

---

## 12. Sync Primitives

```go
import "sync"
import "sync/atomic"

// Mutex — mutual exclusion
type SafeCounter struct {
    mu    sync.Mutex
    count int
}

func (c *SafeCounter) Inc() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.count++
}

func (c *SafeCounter) Value() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.count
}

// RWMutex — nhiều readers HOẶC một writer
type Cache struct {
    mu    sync.RWMutex
    items map[string]string
}

func (c *Cache) Get(key string) (string, bool) {
    c.mu.RLock()         // nhiều goroutines có thể RLock cùng lúc
    defer c.mu.RUnlock()
    v, ok := c.items[key]
    return v, ok
}

func (c *Cache) Set(key, val string) {
    c.mu.Lock()          // exclusive — không có reader/writer nào khác
    defer c.mu.Unlock()
    c.items[key] = val
}

// Once — chạy function đúng 1 lần (thread-safe singleton)
var (
    once     sync.Once
    instance *Database
)

func GetDB() *Database {
    once.Do(func() {
        instance = connectDB()  // chỉ chạy 1 lần, kể cả khi nhiều goroutines gọi
    })
    return instance
}

// Pool — reuse objects, giảm GC pressure
var bufPool = sync.Pool{
    New: func() any {
        return bytes.NewBuffer(make([]byte, 0, 4096))
    },
}

func processData(data []byte) string {
    buf := bufPool.Get().(*bytes.Buffer)
    defer func() {
        buf.Reset()
        bufPool.Put(buf)   // trả về pool để reuse
    }()
    buf.Write(data)
    return buf.String()
}

// sync.Map — concurrent map (phù hợp khi read nhiều, write ít)
var m sync.Map
m.Store("key", "value")
v, ok := m.Load("key")
m.Delete("key")
m.LoadOrStore("key", "default")
m.Range(func(k, v any) bool {
    fmt.Println(k, v)
    return true  // false để stop iteration
})

// Atomic operations — lock-free, nhanh nhất
var counter atomic.Int64     // Go 1.19+ typed atomics
counter.Add(1)
val := counter.Load()
counter.Store(0)
counter.CompareAndSwap(0, 1)  // CAS

// Legacy style (Go < 1.19)
var count int64
atomic.AddInt64(&count, 1)
atomic.LoadInt64(&count)
atomic.StoreInt64(&count, 0)
atomic.CompareAndSwapInt64(&count, 0, 1)

// sync.Cond — complex synchronization
type BlockingQueue struct {
    mu    sync.Mutex
    cond  *sync.Cond
    items []int
    max   int
}

func NewBlockingQueue(max int) *BlockingQueue {
    q := &BlockingQueue{max: max}
    q.cond = sync.NewCond(&q.mu)
    return q
}

func (q *BlockingQueue) Push(item int) {
    q.mu.Lock()
    defer q.mu.Unlock()
    for len(q.items) >= q.max {
        q.cond.Wait()  // release lock, wait for signal
    }
    q.items = append(q.items, item)
    q.cond.Signal()    // notify one waiter
}

func (q *BlockingQueue) Pop() int {
    q.mu.Lock()
    defer q.mu.Unlock()
    for len(q.items) == 0 {
        q.cond.Wait()
    }
    item := q.items[0]
    q.items = q.items[1:]
    q.cond.Signal()
    return item
}
```

---

## 13. Context

```go
import "context"

// Context propagates: cancellation, deadlines, values qua call chain
// RULE: Context là ARGUMENT ĐẦU TIÊN, không bao giờ nil

// Root contexts
ctx := context.Background()    // root, không bao giờ cancel
ctx2 := context.TODO()         // placeholder khi chưa quyết định

// Cancellation
ctx, cancel := context.WithCancel(context.Background())
defer cancel()  // LUÔN defer cancel để tránh goroutine/resource leak

go func() {
    select {
    case <-time.After(5 * time.Second):
        fmt.Println("done normally")
    case <-ctx.Done():
        fmt.Println("cancelled:", ctx.Err())  // context.Canceled
    }
}()
time.Sleep(time.Second)
cancel()  // cancel tất cả goroutines đang listen ctx.Done()

// Timeout
ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// Deadline
ctx, cancel = context.WithDeadline(
    context.Background(),
    time.Now().Add(5*time.Second),
)
defer cancel()

// Values — chỉ dùng cho request-scoped data
type contextKey string
const (
    userIDKey    contextKey = "userID"
    requestIDKey contextKey = "requestID"
)

// Set value
ctx = context.WithValue(ctx, userIDKey, 12345)

// Get value — LUÔN check type assertion
if userID, ok := ctx.Value(userIDKey).(int); ok {
    fmt.Println("User:", userID)
}

// HTTP handler với context
func handler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    result, err := fetchData(ctx)  // propagate context
    if err != nil {
        if errors.Is(err, context.Canceled) {
            return  // client disconnected — không cần response
        }
        if errors.Is(err, context.DeadlineExceeded) {
            http.Error(w, "timeout", http.StatusGatewayTimeout)
            return
        }
        http.Error(w, "internal error", 500)
        return
    }
    json.NewEncoder(w).Encode(result)
}

func fetchData(ctx context.Context) ([]byte, error) {
    req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
    resp, err := http.DefaultClient.Do(req)
    // Context cancellation tự động cancel HTTP request
    if err != nil {
        return nil, fmt.Errorf("fetchData: %w", err)
    }
    defer resp.Body.Close()
    return io.ReadAll(resp.Body)
}
```

---

## 14. Race Conditions

```go
// Data race: concurrent access, ít nhất một write, không có synchronization
// LUÔN chạy với -race flag!
// go test -race ./...
// go run -race main.go
// go build -race -o app ./...

// DATA RACE — sẽ bị detect bởi race detector
var count int
go func() { count++ }()    // DATA RACE!
go func() { count++ }()

// Fix 1: Mutex
var mu sync.Mutex
var count int
mu.Lock()
count++
mu.Unlock()

// Fix 2: Atomic
var count atomic.Int64
count.Add(1)

// Fix 3: Channel
ch := make(chan int, 1)
ch <- 0
go func() {
    v := <-ch
    ch <- v + 1
}()

// RACE CONDITION (logic error, khó detect hơn data race)
// Check-then-act KHÔNG atomic:
if _, ok := m[key]; !ok {   // check
    m[key] = computeValue()  // act ← có thể bị race giữa check và act
}
// Fix: sync.Map.LoadOrStore hoặc mutex bảo vệ cả hai operations

// GOTCHA: map concurrent read+write → PANIC (không phải data race)
m := make(map[string]int)
go func() { m["key"] = 1 }()
go func() { _ = m["key"] }()   // panic: concurrent map read and map write

// Fix: sync.Map hoặc sync.RWMutex
```

---

## 15. Advanced Concurrency Patterns

### Worker Pool

```go
type Job struct{ ID int }
type Result struct{ JobID int; Output string }

func workerPool(ctx context.Context, numWorkers int, jobs <-chan Job) <-chan Result {
    results := make(chan Result)

    var wg sync.WaitGroup
    for i := 0; i < numWorkers; i++ {
        wg.Go(func() {     // Go 1.25+
            for {
                select {
                case job, ok := <-jobs:
                    if !ok {
                        return
                    }
                    results <- processJob(job)
                case <-ctx.Done():
                    return
                }
            }
        })
    }

    go func() {
        wg.Wait()
        close(results)
    }()

    return results
}
```

### errgroup — Goroutines với Error Handling

```go
import "golang.org/x/sync/errgroup"

func fetchAll(ctx context.Context, urls []string) ([][]byte, error) {
    g, ctx := errgroup.WithContext(ctx)
    results := make([][]byte, len(urls))

    for i, url := range urls {
        i, url := i, url   // capture (không cần từ Go 1.22+)
        g.Go(func() error {
            data, err := fetch(ctx, url)
            if err != nil {
                return fmt.Errorf("fetch %s: %w", url, err)
            }
            results[i] = data
            return nil
        })
    }

    // Wait đến khi tất cả goroutines xong
    // Nếu bất kỳ goroutine nào trả về error → ctx bị cancel
    if err := g.Wait(); err != nil {
        return nil, err
    }
    return results, nil
}

// errgroup với limit concurrency
g := new(errgroup.Group)
g.SetLimit(10)  // max 10 concurrent goroutines
for _, url := range urls {
    url := url
    g.Go(func() error { return process(url) })
}
```

### Semaphore — Giới Hạn Concurrent Operations

```go
// Dùng buffered channel làm semaphore
sem := make(chan struct{}, maxConcurrent)

for _, item := range items {
    sem <- struct{}{}     // acquire
    go func(item Item) {
        defer func() { <-sem }()  // release
        process(item)
    }(item)
}

// Drain semaphore (wait for all)
for i := 0; i < cap(sem); i++ {
    sem <- struct{}{}
}
```

### Rate Limiter

```go
// Token bucket với time.Tick
func rateLimited(handler func()) func() {
    limiter := time.NewTicker(100 * time.Millisecond)  // 10 req/sec
    return func() {
        <-limiter.C
        handler()
    }
}

// Hoặc dùng golang.org/x/time/rate
import "golang.org/x/time/rate"

limiter := rate.NewLimiter(rate.Limit(10), 20)  // 10/sec, burst 20

func handleRequest(ctx context.Context) error {
    if err := limiter.Wait(ctx); err != nil {
        return fmt.Errorf("rate limit: %w", err)
    }
    return processRequest()
}
```

### Fan-out / Fan-in

```go
// Fan-out: một input channel → nhiều workers
func fanOut(input <-chan int, n int) []<-chan int {
    outputs := make([]<-chan int, n)
    for i := range outputs {
        outputs[i] = worker(input)
    }
    return outputs
}

// Fan-in: nhiều input channels → một output channel
func fanIn(ctx context.Context, inputs ...<-chan int) <-chan int {
    merged := make(chan int)
    var wg sync.WaitGroup

    for _, in := range inputs {
        in := in
        wg.Go(func() {
            for {
                select {
                case v, ok := <-in:
                    if !ok {
                        return
                    }
                    merged <- v
                case <-ctx.Done():
                    return
                }
            }
        })
    }

    go func() {
        wg.Wait()
        close(merged)
    }()

    return merged
}
```

---

## 16. HTTP Server

```go
package main

import (
    "encoding/json"
    "errors"
    "log/slog"
    "net/http"
    "time"
)

func main() {
    mux := http.NewServeMux()

    // Go 1.22+: Method + Path parameters
    mux.HandleFunc("GET /users/{id}", getUser)
    mux.HandleFunc("POST /users", createUser)
    mux.HandleFunc("PUT /users/{id}", updateUser)
    mux.HandleFunc("DELETE /users/{id}", deleteUser)
    mux.HandleFunc("GET /users", listUsers)

    // Static files
    mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

    // Middleware chain
    handler := chain(mux,
        loggingMiddleware,
        recoveryMiddleware,
        corsMiddleware,
    )

    srv := &http.Server{
        Addr:              ":8080",
        Handler:           handler,
        ReadTimeout:       15 * time.Second,
        ReadHeaderTimeout: 5 * time.Second,
        WriteTimeout:      15 * time.Second,
        IdleTimeout:       60 * time.Second,
        MaxHeaderBytes:    1 << 20,  // 1MB
    }

    slog.Info("server starting", "addr", srv.Addr)
    if err := srv.ListenAndServe(); err != http.ErrServerClosed {
        slog.Error("server failed", "error", err)
    }
}

// Handler
func getUser(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")    // Go 1.22+

    user, err := userService.GetByID(r.Context(), id)
    if err != nil {
        writeError(w, err)
        return
    }

    writeJSON(w, http.StatusOK, user)
}

// Helpers
func writeJSON(w http.ResponseWriter, status int, v any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, err error) {
    switch {
    case errors.Is(err, ErrNotFound):
        http.Error(w, "not found", http.StatusNotFound)
    case errors.Is(err, ErrForbidden):
        http.Error(w, "forbidden", http.StatusForbidden)
    default:
        slog.Error("internal error", "error", err)
        http.Error(w, "internal server error", http.StatusInternalServerError)
    }
}

// Middleware
type Middleware func(http.Handler) http.Handler

func chain(h http.Handler, middlewares ...Middleware) http.Handler {
    for i := len(middlewares) - 1; i >= 0; i-- {
        h = middlewares[i](h)
    }
    return h
}

func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
        next.ServeHTTP(rw, r)
        slog.Info("request",
            "method", r.Method,
            "path", r.URL.Path,
            "status", rw.status,
            "duration", time.Since(start),
        )
    })
}

type responseWriter struct {
    http.ResponseWriter
    status int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.status = code
    rw.ResponseWriter.WriteHeader(code)
}

func recoveryMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if rec := recover(); rec != nil {
                slog.Error("panic recovered", "panic", rec)
                http.Error(w, "internal server error", http.StatusInternalServerError)
            }
        }()
        next.ServeHTTP(w, r)
    })
}

// HTTP Client — LUÔN có timeout và reuse
var httpClient = &http.Client{
    Timeout: 30 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 10,
        IdleConnTimeout:     90 * time.Second,
    },
}

// KHÔNG dùng http.DefaultClient trong production — không có timeout!
```

---

## 17. Testing

```go
package user_test   // black-box test

import (
    "context"
    "testing"
    "github.com/yourname/myapp/internal/user"
)

// Unit test cơ bản
func TestGetUser(t *testing.T) {
    repo := user.NewMockRepository()
    svc := user.NewService(repo)

    u, err := svc.GetUser(context.Background(), 1)

    if err != nil {
        t.Fatalf("expected no error, got %v", err)
    }
    if u.ID != 1 {
        t.Errorf("expected ID=1, got %d", u.ID)
    }
}

// Table-driven tests — GO WAY (phải biết!)
func TestAdd(t *testing.T) {
    tests := []struct {
        name     string
        a, b     int
        expected int
    }{
        {"positive numbers", 1, 2, 3},
        {"negative numbers", -1, -2, -3},
        {"zero", 0, 0, 0},
        {"mixed", -1, 2, 1},
        {"large numbers", 1<<31 - 1, 1, 1 << 31},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := add(tt.a, tt.b)
            if got != tt.expected {
                t.Errorf("add(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.expected)
            }
        })
    }
}

// Test helpers
func assertNoError(t *testing.T, err error) {
    t.Helper()  // makes error point to caller, not here
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
}

// Benchmarks
func BenchmarkSort(b *testing.B) {
    data := generateData(1000)
    b.ResetTimer()             // reset sau setup
    b.ReportAllocs()           // báo cáo allocations

    for i := 0; i < b.N; i++ {
        sort.Ints(data)
    }
}

// Go 1.24+: testing.B.Loop (đơn giản hơn)
func BenchmarkSortNew(b *testing.B) {
    data := generateData(1000)
    for b.Loop() {
        sort.Ints(data)
    }
}

// go test -bench=. -benchmem -benchtime=10s -count=5

// Fuzz testing (Go 1.18+)
func FuzzReverse(f *testing.F) {
    f.Add("hello")          // seed corpus
    f.Add("")
    f.Add("🇻🇳")

    f.Fuzz(func(t *testing.T, s string) {
        reversed := reverse(s)
        if reverse(reversed) != s {
            t.Errorf("double reverse of %q = %q", s, reverse(reversed))
        }
    })
}
// go test -fuzz=FuzzReverse -fuzztime=60s

// Test với timeout
func TestSlowOperation(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping in short mode")
    }
    // ...
}
// go test -short ./...

// Parallel tests
func TestParallel(t *testing.T) {
    t.Parallel()  // cho phép chạy song song với test khác
    // ...
}

// Test fixtures với TestMain
func TestMain(m *testing.M) {
    // Setup
    db := setupTestDB()

    // Run tests
    code := m.Run()

    // Teardown
    db.Close()
    os.Exit(code)
}

// httptest — test HTTP handlers
func TestGetUserHandler(t *testing.T) {
    req := httptest.NewRequest("GET", "/users/1", nil)
    rr := httptest.NewRecorder()

    handler := http.HandlerFunc(getUser)
    handler.ServeHTTP(rr, req)

    if rr.Code != http.StatusOK {
        t.Errorf("expected 200, got %d", rr.Code)
    }
}

// Testify (library phổ biến)
import "github.com/stretchr/testify/assert"
import "github.com/stretchr/testify/require"

assert.Equal(t, expected, actual, "message")
assert.NoError(t, err)
assert.ErrorIs(t, err, ErrNotFound)
require.NotNil(t, user)  // require → t.FailNow() nếu fail
```

---

## 18. Profiling & Performance

```go
// pprof — built-in profiler — PHẢI BIẾT!
import _ "net/http/pprof"  // register /debug/pprof endpoints

// Trong server:
go func() {
    log.Println(http.ListenAndServe("localhost:6060", nil))
}()

// Commands:
// go tool pprof http://localhost:6060/debug/pprof/heap          # memory
// go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30  # CPU
// go tool pprof http://localhost:6060/debug/pprof/goroutine     # goroutines
// go tool pprof http://localhost:6060/debug/pprof/allocs        # allocations
// go tool pprof http://localhost:6060/debug/pprof/mutex         # mutex contention
// go tool pprof http://localhost:6060/debug/pprof/block         # goroutine blocking

// Trong tests:
// go test -bench=BenchmarkHeavy -memprofile=mem.pprof -cpuprofile=cpu.pprof
// go tool pprof -http=:8081 mem.pprof  # web UI

// Execution tracer — timeline view
import "runtime/trace"

f, _ := os.Create("trace.out")
trace.Start(f)
// ... code to trace ...
trace.Stop()
f.Close()
// go tool trace trace.out

// Escape analysis — xem variable nào lên heap
// go build -gcflags="-m -m" ./...
// go build -gcflags="-m=2" 2>&1 | grep "escapes to heap"

// Memory stats
var m runtime.MemStats
runtime.ReadMemStats(&m)
fmt.Printf("Alloc: %v MB\n", m.Alloc/1024/1024)
fmt.Printf("TotalAlloc: %v MB\n", m.TotalAlloc/1024/1024)
fmt.Printf("NumGC: %v\n", m.NumGC)

// Force GC (for testing)
runtime.GC()
```

---

## 19. Logging với slog

```go
import "log/slog"    // Go 1.21+

// Default logger (text format)
slog.Info("server started", "port", 8080)
slog.Warn("high latency", "duration", 2.5)
slog.Error("connection failed", "error", err, "host", host)
slog.Debug("request received", "path", r.URL.Path)

// Custom JSON logger
logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
    Level:     slog.LevelDebug,
    AddSource: true,  // thêm file:line
}))
slog.SetDefault(logger)  // set làm default

// Structured logging
logger.Info("user created",
    slog.Int("userID", 123),
    slog.String("email", "user@example.com"),
    slog.Duration("elapsed", time.Since(start)),
    slog.Any("metadata", map[string]any{"role": "admin"}),
)

// Logger với persistent attributes
requestLogger := logger.With(
    slog.String("requestID", requestID),
    slog.String("userIP", r.RemoteAddr),
)
requestLogger.Info("handling request")
requestLogger.Info("request complete", slog.Int("status", 200))

// Group attributes
logger.Info("http request",
    slog.Group("request",
        slog.String("method", "GET"),
        slog.String("path", "/users/1"),
    ),
    slog.Group("response",
        slog.Int("status", 200),
        slog.Duration("duration", 45*time.Millisecond),
    ),
)

// Context-aware logging
logger.InfoContext(ctx, "processing job")

// Log levels
// Debug < Info < Warn < Error
// Production: Info hoặc Warn
// Development: Debug
```

---

## 20. Generics

```go
// Generics (Go 1.18+) — Type Parameters

// Generic function
func Map[T, U any](slice []T, f func(T) U) []U {
    result := make([]U, len(slice))
    for i, v := range slice {
        result[i] = f(v)
    }
    return result
}

func Filter[T any](slice []T, pred func(T) bool) []T {
    var result []T
    for _, v := range slice {
        if pred(v) {
            result = append(result, v)
        }
    }
    return result
}

func Reduce[T, U any](slice []T, init U, f func(U, T) U) U {
    acc := init
    for _, v := range slice {
        acc = f(acc, v)
    }
    return acc
}

// Usage
nums := []int{1, 2, 3, 4, 5}
doubled := Map(nums, func(n int) int { return n * 2 })
evens := Filter(nums, func(n int) bool { return n%2 == 0 })
sum := Reduce(nums, 0, func(acc, n int) int { return acc + n })

// Constraints — giới hạn type parameters
type Number interface {
    ~int | ~int64 | ~float64  // ~ = underlying type
}

func Sum[T Number](nums []T) T {
    var total T
    for _, n := range nums {
        total += n
    }
    return total
}

// Built-in constraints (golang.org/x/exp/constraints)
type Ordered interface {
    Integer | Float | ~string
}

// Generic data structures
type Stack[T any] struct {
    items []T
}

func (s *Stack[T]) Push(item T) {
    s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() (T, bool) {
    if len(s.items) == 0 {
        var zero T
        return zero, false
    }
    top := s.items[len(s.items)-1]
    s.items = s.items[:len(s.items)-1]
    return top, true
}

func (s *Stack[T]) Len() int { return len(s.items) }

// Generic type alias (Go 1.24+)
type Node[T any] struct {
    Value T
    Next  *Node[T]
}
type IntNode = Node[int]         // OK từ trước
type GenericNode[T any] = Node[T] // Go 1.24+: generic type alias

// slices package (Go 1.21+) — generic functions cho slices
import "slices"
slices.Sort(nums)
slices.Contains(nums, 3)
slices.Index(nums, 3)
slices.Reverse(nums)
idx, found := slices.BinarySearch(sorted, 3)

// maps package (Go 1.21+)
import "maps"
maps.Keys(m)         // iterator over keys (Go 1.23+: returns iter.Seq)
maps.Values(m)       // iterator over values
maps.Clone(m)        // shallow copy
maps.DeleteFunc(m, func(k, v) bool { return v == 0 })

// Khi KHÔNG nên dùng generics:
// - Khi interface{} đủ dùng
// - Khi logic khác nhau cho mỗi type (cần type switch)
// - Khi làm code phức tạp hơn không cần thiết
// "If in doubt, leave it out" — Rob Pike
```

---

## 21. Reflection

```go
import "reflect"

// reflect.TypeOf và reflect.ValueOf
var x float64 = 3.14
t := reflect.TypeOf(x)    // *reflect.rtype
v := reflect.ValueOf(x)   // reflect.Value

fmt.Println(t.Kind())     // float64
fmt.Println(t.String())   // "float64"
fmt.Println(v.Float())    // 3.14

// Modify value (cần pointer)
p := reflect.ValueOf(&x)
p.Elem().SetFloat(2.71)   // x = 2.71

// Inspect struct với tags
type User struct {
    Name  string `json:"name" validate:"required"`
    Email string `json:"email" validate:"email"`
    Age   int    `json:"age,omitempty"`
}

u := User{Name: "Alice", Email: "alice@example.com", Age: 30}
t = reflect.TypeOf(u)
v = reflect.ValueOf(u)

for i := 0; i < t.NumField(); i++ {
    field := t.Field(i)
    value := v.Field(i)
    jsonTag := field.Tag.Get("json")
    validateTag := field.Tag.Get("validate")
    fmt.Printf("%s: value=%v, json=%s, validate=%s\n",
        field.Name, value, jsonTag, validateTag)
}

// Tạo struct dynamically
structType := reflect.StructOf([]reflect.StructField{
    {Name: "Name", Type: reflect.TypeOf(""), Tag: `json:"name"`},
    {Name: "Age", Type: reflect.TypeOf(0), Tag: `json:"age"`},
})
instance := reflect.New(structType).Elem()
instance.FieldByName("Name").SetString("Alice")

// reflect.DeepEqual — so sánh deep
reflect.DeepEqual([]int{1,2,3}, []int{1,2,3})  // true

// Khi dùng reflection:
// - JSON/YAML/XML marshaling
// - ORM libraries
// - Dependency injection frameworks
// - Template engines

// NHƯỢC ĐIỂM: chậm (~10-100x), không có compile-time safety
```

---

## 22. unsafe Package

```go
import "unsafe"

// unsafe.Sizeof — size tại compile time (không allocate)
fmt.Println(unsafe.Sizeof(int64(0)))      // 8
fmt.Println(unsafe.Sizeof(struct{}{}))    // 0
fmt.Println(unsafe.Sizeof([1000]byte{}))  // 1000

// unsafe.Alignof — alignment requirement
fmt.Println(unsafe.Alignof(float64(0)))   // 8

// unsafe.Offsetof — field offset trong struct
type S struct { A byte; B float64 }
fmt.Println(unsafe.Offsetof(S{}.B))       // 8 (4 bytes padding sau A)

// unsafe.Pointer — convert giữa pointer types
var x int64 = 0x0102030405060708
p := unsafe.Pointer(&x)
b := (*[8]byte)(p)   // view int64 as byte array
fmt.Println(*b)

// String ↔ []byte không copy (Go 1.20+)
func unsafeStringToBytes(s string) []byte {
    return unsafe.Slice(unsafe.StringData(s), len(s))
}
func unsafeBytesToString(b []byte) string {
    return unsafe.String(unsafe.SliceData(b), len(b))
}
// CẢNH BÁO: chỉ dùng read-only! Mutation → undefined behavior

// Struct alignment optimization
type BadStruct struct {   // 24 bytes (với padding)
    a bool    // 1 byte + 7 padding
    b int64   // 8 bytes
    c bool    // 1 byte + 7 padding
}

type GoodStruct struct {  // 16 bytes (optimized)
    b int64   // 8 bytes
    a bool    // 1 byte
    c bool    // 1 byte + 6 padding
}

// Khi NÊN dùng unsafe:
// - CGo interfacing
// - Performance-critical hot paths (zero-copy)
// - Binary protocol parsing
// - Implementing sync primitives cấp thấp

// KHÔNG DÙNG khi có alternative an toàn!
```

---

## 23. CGo

```go
package main

/*
#include <stdio.h>
#include <stdlib.h>

void hello_from_c() {
    printf("Hello from C!\n");
}

int add(int a, int b) {
    return a + b;
}

typedef struct {
    int x;
    int y;
} Point;
*/
import "C"   // PHẢI ngay sau block comment, không có dòng trắng

import (
    "fmt"
    "unsafe"
)

func main() {
    C.hello_from_c()

    result := C.add(3, 4)
    fmt.Println(int(result))  // 7

    // C string → Go string
    cs := C.CString("Hello from Go")
    defer C.free(unsafe.Pointer(cs))  // BẮT BUỘC free!
    C.puts(cs)

    // Go string → C string
    goStr := C.GoString(cs)
    fmt.Println(goStr)

    // C struct
    p := C.Point{x: 10, y: 20}
    fmt.Println(int(p.x), int(p.y))
}

// Cross-compile với CGo — phức tạp hơn nhiều
// CGO_ENABLED=0 go build → static binary (không CGo)

// Nhược điểm CGo:
// - Slow cross-function calls (overhead ~20-100x)
// - Phức tạp khi cross-compile
// - Không dùng được goroutines trong C callbacks (directly)
// - GC không biết về C memory → phải manual free
```

---

## 24. Design Patterns trong Go

### Repository Pattern

```go
// Interface tại nơi DÙNG (service), không phải nơi implement (db)
type UserRepository interface {
    FindByID(ctx context.Context, id int64) (*User, error)
    FindByEmail(ctx context.Context, email string) (*User, error)
    Save(ctx context.Context, user *User) error
    Delete(ctx context.Context, id int64) error
    List(ctx context.Context, filter UserFilter) ([]*User, int64, error)
}

type PostgresUserRepository struct {
    db *sql.DB
}

func (r *PostgresUserRepository) FindByID(ctx context.Context, id int64) (*User, error) {
    var u User
    err := r.db.QueryRowContext(ctx,
        `SELECT id, name, email, created_at FROM users WHERE id = $1 AND deleted_at IS NULL`,
        id,
    ).Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt)
    if errors.Is(err, sql.ErrNoRows) {
        return nil, ErrNotFound
    }
    return &u, err
}
```

### Dependency Injection

```go
// Không dùng global state — inject dependencies
type UserService struct {
    repo   UserRepository    // interface, không phải concrete
    mailer EmailSender       // interface
    cache  Cache             // interface
    logger *slog.Logger
    clock  func() time.Time  // testable!
}

func NewUserService(
    repo UserRepository,
    mailer EmailSender,
    cache Cache,
    logger *slog.Logger,
) *UserService {
    return &UserService{
        repo:   repo,
        mailer: mailer,
        cache:  cache,
        logger: logger,
        clock:  time.Now,
    }
}
```

### Builder Pattern

```go
type QueryBuilder struct {
    table      string
    conditions []string
    args       []any
    order      string
    limit      int
    offset     int
}

func NewQuery(table string) *QueryBuilder {
    return &QueryBuilder{table: table, limit: 100}
}

func (b *QueryBuilder) Where(cond string, args ...any) *QueryBuilder {
    b.conditions = append(b.conditions, cond)
    b.args = append(b.args, args...)
    return b
}

func (b *QueryBuilder) OrderBy(field string) *QueryBuilder {
    b.order = field
    return b
}

func (b *QueryBuilder) Limit(n int) *QueryBuilder {
    b.limit = n
    return b
}

func (b *QueryBuilder) Build() (string, []any) {
    q := fmt.Sprintf("SELECT * FROM %s", b.table)
    if len(b.conditions) > 0 {
        q += " WHERE " + strings.Join(b.conditions, " AND ")
    }
    if b.order != "" {
        q += " ORDER BY " + b.order
    }
    if b.limit > 0 {
        q += fmt.Sprintf(" LIMIT %d", b.limit)
    }
    return q, b.args
}

// Usage
query, args := NewQuery("users").
    Where("age > $1", 18).
    Where("active = $2", true).
    OrderBy("name").
    Limit(50).
    Build()
```

### Event Bus

```go
type EventHandler func(data any)

type EventBus struct {
    handlers map[string][]EventHandler
    mu       sync.RWMutex
}

func NewEventBus() *EventBus {
    return &EventBus{handlers: make(map[string][]EventHandler)}
}

func (eb *EventBus) Subscribe(event string, handler EventHandler) func() {
    eb.mu.Lock()
    eb.handlers[event] = append(eb.handlers[event], handler)
    eb.mu.Unlock()

    // Trả về unsubscribe function
    return func() {
        eb.mu.Lock()
        defer eb.mu.Unlock()
        handlers := eb.handlers[event]
        for i, h := range handlers {
            if reflect.ValueOf(h).Pointer() == reflect.ValueOf(handler).Pointer() {
                eb.handlers[event] = append(handlers[:i], handlers[i+1:]...)
                break
            }
        }
    }
}

func (eb *EventBus) Publish(event string, data any) {
    eb.mu.RLock()
    handlers := make([]EventHandler, len(eb.handlers[event]))
    copy(handlers, eb.handlers[event])
    eb.mu.RUnlock()

    for _, h := range handlers {
        go h(data)  // async dispatch
    }
}
```

---

## 25. Microservices & gRPC

```proto
// user.proto
syntax = "proto3";
option go_package = "github.com/yourname/myapp/proto/userpb";

service UserService {
    rpc GetUser(GetUserRequest) returns (User);
    rpc ListUsers(ListUsersRequest) returns (stream User);
    rpc CreateUser(CreateUserRequest) returns (User);
}

message User {
    int64 id = 1;
    string name = 2;
    string email = 3;
}
```

```go
// Generate: protoc --go_out=. --go-grpc_out=. user.proto

// Server implementation
type userServer struct {
    pb.UnimplementedUserServiceServer  // PHẢI embed này
    svc *UserService
}

func (s *userServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
    user, err := s.svc.GetUser(ctx, req.Id)
    if err != nil {
        if errors.Is(err, ErrNotFound) {
            return nil, status.Errorf(codes.NotFound, "user %d not found", req.Id)
        }
        return nil, status.Error(codes.Internal, "internal error")
    }
    return toProto(user), nil
}

// gRPC server với interceptors
func main() {
    lis, _ := net.Listen("tcp", ":50051")
    s := grpc.NewServer(
        grpc.UnaryInterceptor(
            grpc_middleware.ChainUnaryServer(
                grpc_recovery.UnaryServerInterceptor(),
                grpc_zap.UnaryServerInterceptor(logger),
            ),
        ),
    )
    pb.RegisterUserServiceServer(s, &userServer{})
    reflection.Register(s)  // cho grpcurl, evans
    s.Serve(lis)
}

// gRPC Client
conn, _ := grpc.NewClient(
    "localhost:50051",
    grpc.WithTransportCredentials(insecure.NewCredentials()),
)
client := pb.NewUserServiceClient(conn)

ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
user, err := client.GetUser(ctx, &pb.GetUserRequest{Id: 1})
```

### Health Check & Metrics

```go
import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
    requestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total HTTP requests",
        },
        []string{"method", "path", "status"},
    )
    requestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "HTTP request duration",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "path"},
    )
)

// /metrics endpoint
mux.Handle("GET /metrics", promhttp.Handler())

// Health/Readiness endpoints
mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))
})

mux.HandleFunc("GET /readyz", func(w http.ResponseWriter, r *http.Request) {
    if err := db.PingContext(r.Context()); err != nil {
        http.Error(w, "db not ready", http.StatusServiceUnavailable)
        return
    }
    w.WriteHeader(http.StatusOK)
})
```

---

## 26. Database Patterns

```go
// database/sql — standard library
import "database/sql"

// Connection pool
db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
if err != nil {
    log.Fatal(err)
}
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(5)
db.SetConnMaxLifetime(5 * time.Minute)
db.SetConnMaxIdleTime(1 * time.Minute)

// Verify connection
if err := db.PingContext(ctx); err != nil {
    log.Fatal("cannot ping db:", err)
}

// Query
rows, err := db.QueryContext(ctx,
    "SELECT id, name FROM users WHERE active = $1 ORDER BY name",
    true,
)
if err != nil {
    return err
}
defer rows.Close()

var users []User
for rows.Next() {
    var u User
    if err := rows.Scan(&u.ID, &u.Name); err != nil {
        return err
    }
    users = append(users, u)
}
if err := rows.Err(); err != nil {
    return err
}

// Prepared statement (reuse, prevent injection)
stmt, err := db.PrepareContext(ctx,
    "SELECT * FROM users WHERE id = $1",
)
if err != nil {
    return err
}
defer stmt.Close()
row := stmt.QueryRowContext(ctx, userID)

// Transaction
tx, err := db.BeginTx(ctx, &sql.TxOptions{
    Isolation: sql.LevelReadCommitted,
})
if err != nil {
    return err
}
defer tx.Rollback()  // no-op nếu đã committed

_, err = tx.ExecContext(ctx,
    "UPDATE accounts SET balance = balance - $1 WHERE id = $2",
    amount, fromID,
)
if err != nil {
    return fmt.Errorf("debit: %w", err)
}

_, err = tx.ExecContext(ctx,
    "UPDATE accounts SET balance = balance + $1 WHERE id = $2",
    amount, toID,
)
if err != nil {
    return fmt.Errorf("credit: %w", err)
}

return tx.Commit()

// pgx v5 — recommended PostgreSQL driver
import "github.com/jackc/pgx/v5/pgxpool"

pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
rows, err := pool.Query(ctx,
    "SELECT id, name FROM users WHERE active = $1",
    true,
)
// pgx tự collect results:
users, err := pgx.CollectRows(rows, pgx.RowToStructByName[User])

// sqlc — generate type-safe Go từ SQL
// Viết SQL → chạy sqlc generate → nhận Go code type-safe
// Không cần ORM, không cần reflection
```

---

## 27. Go Runtime Internals — GMP Model

### GMP Scheduler

```
GMP = Goroutine, Machine, Processor

G = Goroutine
  - Unit of work (như thread nhưng siêu nhẹ)
  - Stack: bắt đầu 2KB, tự grow đến ~1GB
  - States: Runnable → Running → Waiting → Dead

M = Machine (OS Thread)
  - 1 M = 1 OS thread
  - Execute code của G
  - Blocked on syscall → M bị detach, P tìm M khác

P = Processor (Logical CPU)
  - Số P = GOMAXPROCS (default = NumCPU)
  - Mỗi P có Local Run Queue (LRQ) chứa runnable Gs
  - P "owns" M trong khi chạy

Scheduling Algorithm:
1. G sẵn sàng chạy → thêm vào P's LRQ (hoặc GRQ nếu LRQ đầy)
2. P lấy G từ LRQ để run trên M
3. Nếu P's LRQ rỗng → "work stealing" từ P khác (lấy 1/2 LRQ)
4. G bị block (channel, syscall, mutex) → G vào wait queue, M tìm G khác
5. G unblocked → trở lại LRQ

Preemption:
- Cooperative: tại function call boundaries
- Asynchronous: SIGURG signal cho goroutines tight-loop (>10ms)
```

```go
// GOMAXPROCS
import "runtime"

runtime.GOMAXPROCS(4)     // set 4 logical processors
n := runtime.GOMAXPROCS(0) // query current value (0 = don't change)

// Go 1.25+: Container-aware GOMAXPROCS
// Tự động respects cgroup CPU limits trên Linux
// GODEBUG=containermaxprocs=0 để disable

// Debug scheduler
// GODEBUG=schedtrace=1000 go run main.go    → print mỗi 1000ms
// GODEBUG=scheddetail=1,schedtrace=1000 go run main.go

// Runtime info
fmt.Println("Goroutines:", runtime.NumGoroutine())
fmt.Println("CPUs:", runtime.NumCPU())
fmt.Println("GOMAXPROCS:", runtime.GOMAXPROCS(0))

// Stack trace tất cả goroutines
buf := make([]byte, 1<<20)
n := runtime.Stack(buf, true)  // true = all goroutines
fmt.Printf("%s", buf[:n])
```

---

## 28. Memory Management & GC

### Garbage Collector — Tri-Color Mark-and-Sweep

```
Phases của GC cycle:

1. Mark Setup (STW — Stop The World, rất ngắn ~100µs)
   - Pause tất cả goroutines
   - Enable write barrier

2. Marking (Concurrent — goroutines vẫn chạy!)
   - Tri-color algorithm:
     * White = unmarked (chưa scan)
     * Gray = reachable nhưng references chưa scan
     * Black = scanned và tất cả references đã scan
   - Write barrier: track pointer writes trong khi scan
   - Bắt đầu từ GC roots (stacks, globals)
   - Traverse toàn bộ object graph

3. Mark Termination (STW)
   - Flush write barrier buffers
   - Disable write barrier

4. Sweeping (Concurrent)
   - White objects = garbage → reclaim memory

GC Trigger:
- GOGC=100 (default): GC trigger khi heap tăng 100% so với sau GC cuối
- GOMEMLIMIT=512MiB (Go 1.19+): soft memory limit
```

### Green Tea GC (Go 1.25+)

```
Experimental GC mới, enable bằng GOEXPERIMENT=greenteagc

Vấn đề của GC cũ:
- Tri-color scan nhảy khắp memory → poor spatial locality
- 85% time GC ở scan loop, >35% stalled on memory

Green Tea solution:
- Scan memory theo SPANS (blocks lớn) thay vì từng object
- Better cache locality → ít cache miss hơn
- Expected 10-40% reduction in GC overhead cho GC-heavy workloads

Enable:
GOEXPERIMENT=greenteagc go build ./...
GOEXPERIMENT=greenteagc go test ./...

Go 1.26: dự kiến là default GC
```

### GC Tuning

```go
// GOGC — set GC percentage
// GOGC=off → disable GC (batch processing)
// GOGC=200 → GC ít thường xuyên hơn, dùng nhiều memory hơn
// GOGC=50  → GC thường xuyên hơn, dùng ít memory hơn
import "runtime/debug"
debug.SetGCPercent(200)

// GOMEMLIMIT — soft memory limit (Go 1.19+)
debug.SetMemoryLimit(512 * 1024 * 1024)  // 512 MB

// Force GC
runtime.GC()

// Disable GC (cho batch/startup)
debug.SetGCPercent(-1)
// ... do work ...
debug.SetGCPercent(100)  // restore

// GC stats
var stats runtime.MemStats
runtime.ReadMemStats(&stats)
fmt.Printf("GC count: %d\n", stats.NumGC)
fmt.Printf("GC pause total: %v\n", time.Duration(stats.PauseTotalNs))
fmt.Printf("Heap in use: %d MB\n", stats.HeapInuse/1024/1024)
```

### Escape Analysis

```go
// Variable escape lên heap khi:
// 1. Pointer đến local variable return ra ngoài function
// 2. Variable quá lớn cho stack
// 3. Variable được gán vào interface
// 4. Variable trong goroutine closure
// 5. Slice/map grow dynamically

// Stack allocation (nhanh hơn)
func stackAlloc() int {
    x := 42      // → stack
    return x
}

// Heap allocation (GC managed)
func heapAlloc() *int {
    x := 42      // → heap (escapes!)
    return &x
}

// Interface boxing → heap allocation
var i interface{} = 42        // 42 escapes to heap

// Kiểm tra escape analysis:
// go build -gcflags="-m" ./...
// go build -gcflags="-m=2" 2>&1 | grep escape
```

---

## 29. Performance Optimization

### DO NOT OPTIMIZE PREMATURELY! Profile → Find Bottleneck → Optimize → Benchmark

```go
// 1. String concatenation
// BAD: O(n²) — tạo new string mỗi iteration
s := ""
for _, word := range words {
    s += word + " "
}

// GOOD: O(n) — single allocation
var b strings.Builder
b.Grow(estimatedSize)  // pre-allocate
for _, word := range words {
    b.WriteString(word)
    b.WriteByte(' ')
}
result := b.String()

// 2. Slice pre-allocation
// BAD: nhiều reallocations
var result []int
for i := 0; i < 1000; i++ {
    result = append(result, i)
}

// GOOD: single allocation
result := make([]int, 0, 1000)
for i := 0; i < 1000; i++ {
    result = append(result, i)
}

// 3. Struct field alignment (giảm padding)
// Dùng: go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest
// fieldalignment -fix ./...

type Optimized struct {
    b int64    // 8 bytes
    c float64  // 8 bytes
    a bool     // 1 byte
    d bool     // 1 byte + 6 padding = 24 bytes total
}
// vs unoptimized: có thể lên đến 32 bytes

// 4. Avoid interface boxing in hot paths
func hotPath(nums []int) int {
    // BAD: boxing vào interface cho mỗi element
    var sum interface{} = 0
    for _, n := range nums {
        sum = sum.(int) + n
    }

    // GOOD: concrete type
    total := 0
    for _, n := range nums {
        total += n
    }
    return total
}

// 5. sync.Pool — reuse objects
var jsonPool = sync.Pool{
    New: func() any { return new(json.Encoder) },
}

// 6. AOS vs SOA — CPU cache optimization
// Array of Structs (AOS) — xấu cho bulk operations trên 1 field
type PointAOS struct{ X, Y, Z float64 }
points := []PointAOS{{1,2,3}, {4,5,6}}

// Struct of Arrays (SOA) — tốt cho SIMD/bulk ops trên 1 field
type PointsSOA struct {
    X []float64
    Y []float64
    Z []float64
}

// 7. Batch operations
// BAD: N round-trips
for _, id := range ids {
    db.QueryRow("SELECT ... WHERE id = ?", id)
}

// GOOD: 1 round-trip
db.Query("SELECT ... WHERE id = ANY($1)", pq.Array(ids))

// 8. Buffer I/O
// BAD: unbuffered
f, _ := os.Open("large.txt")
scanner := bufio.NewScanner(f)  // GOOD: buffered

// 9. Benchmarking đúng cách
func BenchmarkMyFunc(b *testing.B) {
    // Setup — không tính vào benchmark
    data := prepareData()
    b.ResetTimer()
    b.ReportAllocs()

    for b.Loop() {         // Go 1.24+
        myFunc(data)
    }
}
// go test -bench=BenchmarkMyFunc -benchmem -count=5 -benchtime=10s
// benchstat old.txt new.txt  // so sánh kết quả
```

---

## 30. Tools Ecosystem

```bash
# ===== STATIC ANALYSIS =====
go vet ./...                    # built-in (phải chạy!)
golangci-lint run ./...         # meta-linter (cài: https://golangci-lint.run)
staticcheck ./...               # advanced static analysis
go tool vet -shadow ./...       # shadow variables
govulncheck ./...               # security vulnerabilities

# golangci-lint config (.golangci.yml)
linters:
  enable:
    - errcheck      # check ignored errors
    - govet         # go vet
    - staticcheck   # staticcheck
    - gosimple      # simplifications
    - ineffassign   # unused assignments
    - unused        # unused code
    - gosec         # security
    - gofmt         # formatting

# ===== CODE GENERATION =====
go generate ./...               # chạy //go:generate directives
stringer -type=Status           # generate String() cho enums
mockgen -source=repo.go -destination=mock_repo.go  # generate mocks
protoc --go_out=. user.proto    # generate protobuf

# ===== DEPENDENCIES =====
go mod tidy                     # sync go.mod/go.sum
go mod download                 # pre-download dependencies
go mod verify                   # verify checksums
go mod graph | dot -Tsvg > deps.svg  # visualize
govulncheck ./...               # check CVEs

# ===== BUILDING =====
go build ./...
go build -trimpath ./...                      # remove local paths
go build -ldflags="-s -w" ./...              # strip debug (smaller binary)
go build -ldflags="-X main.version=v1.2.3"  # inject version
GOOS=linux GOARCH=amd64 go build ./...       # cross compile
CGO_ENABLED=0 go build ./...                 # static binary
go build -race ./...                          # race detector

# Build constraints
//go:build linux && amd64
//go:build !cgo

# ===== TESTING =====
go test ./...                              # all tests
go test -v ./...                           # verbose
go test -run TestGetUser ./...             # specific test
go test -bench=. -benchmem ./...          # benchmarks
go test -race ./...                        # race detector (BẮT BUỘC!)
go test -cover -coverprofile=cover.out ./...
go tool cover -html=cover.out              # coverage UI
go test -fuzz=FuzzMyFunc -fuzztime=60s     # fuzz testing
go test -count=5 ./...                     # run 5 times (detect flaky)
go test -json ./... | gotestfmt            # pretty output
go test -short ./...                       # skip slow tests

# ===== DOCUMENTATION =====
go doc io.Reader                           # terminal docs
godoc -http=:6060                          # local doc server
# pkg.go.dev — online docs

# ===== DEBUGGING =====
dlv debug ./cmd/server                    # Delve debugger
dlv test ./internal/user -- -test.run TestGetUser
GOTRACEBACK=all go run main.go            # full stack on crash
GOMAXPROCS=1 go test ./...               # single-threaded (find races)

# ===== TOOLS DEPENDENCIES (Go 1.24+) =====
go get -tool golang.org/x/tools/cmd/stringer
go get -tool github.com/golangci/golangci-lint/cmd/golangci-lint
go tool stringer -type=Status ./...
```

---

## 31. Go 1.22 → 1.25 — Tính Năng Mới Quan Trọng

### Go 1.22 (February 2024)

```go
// 1. Routing với method và path parameters
mux.HandleFunc("GET /users/{id}", getUser)
mux.HandleFunc("POST /users", createUser)
// r.PathValue("id") — lấy path parameter

// 2. Loop variable semantics FIX
// Trước 1.22: i được share → closure capture cùng giá trị
// Từ 1.22: mỗi iteration có biến riêng
for i := 0; i < 3; i++ {
    go func() {
        fmt.Println(i)  // OK từ 1.22! In đúng 0, 1, 2
    }()
}

// 3. Range over integers
for i := range 5 {        // i = 0, 1, 2, 3, 4
    fmt.Println(i)
}
for range 3 {              // lặp 3 lần, không cần biến
    doSomething()
}
```

### Go 1.23 (August 2024)

```go
// 1. Range over functions (iterators)
// func(yield func() bool) — range over custom iterator
func Backwards[S ~[]E, E any](s S) iter.Seq2[int, E] {
    return func(yield func(int, E) bool) {
        for i := len(s) - 1; i >= 0; i-- {
            if !yield(i, s[i]) {
                return
            }
        }
    }
}

for i, v := range Backwards([]string{"a", "b", "c"}) {
    fmt.Println(i, v)  // 2 c, 1 b, 0 a
}

// slices và maps packages với iterators
for k := range maps.Keys(myMap) { ... }       // iter.Seq
for v := range slices.Values(mySlice) { ... } // iter.Seq

// 2. Timer cleanup (Go 1.23)
// time.After trong hot paths → dùng time.NewTimer thay
timer := time.NewTimer(5 * time.Second)
defer timer.Stop()  // Từ 1.23: Stop() ngăn GC leak
select {
case <-timer.C:
    fmt.Println("timeout")
case result := <-ch:
    timer.Stop()
    processResult(result)
}

// 3. GODEBUG per module (go.mod level)
```

### Go 1.24 (February 2025)

```go
// 1. Generic type aliases
type Node[T any] struct {
    Value T
    Next  *Node[T]
}
type GenericNode[T any] = Node[T]  // Go 1.24+ generic alias!

// 2. Tool dependencies trong go.mod
// go get -tool golang.org/x/tools/cmd/stringer
// go.mod:
//   tool golang.org/x/tools/cmd/stringer
// go tool stringer -type=Status

// 3. Swiss Tables map implementation
// ~2-60% faster map operations (transparent, không cần code change)

// 4. testing.B.Loop (đơn giản hơn b.N)
func BenchmarkSomething(b *testing.B) {
    for b.Loop() {    // NEW: thay cho "for i := 0; i < b.N; i++"
        expensive()
    }
}

// 5. os.Root — sandbox file operations
root, err := os.OpenRoot("/safe/directory")
if err != nil { ... }
defer root.Close()
// Tất cả operations qua root không thể escape ra ngoài directory
f, err := root.Open("file.txt")   // SAFE: không thể path traverse

// 6. strings/bytes iterators
for line := range strings.Lines(text) {
    process(line)
}
for word := range strings.SplitSeq(text, ",") {
    process(word)
}
```

### Go 1.25 (August 2025)

```go
// 1. sync.WaitGroup.Go — GAME CHANGER!
// Trước 1.25:
var wg sync.WaitGroup
wg.Add(1)
go func() {
    defer wg.Done()
    doWork()
}()

// Từ 1.25:
var wg sync.WaitGroup
wg.Go(func() {    // Add(1) + go + Done() tự động!
    doWork()
})
wg.Wait()

// 2. Green Tea GC (experimental)
// GOEXPERIMENT=greenteagc go build ./...
// 10-40% reduction in GC overhead
// Production-ready, dùng tại Google

// 3. Container-aware GOMAXPROCS
// Tự động respects cgroup CPU limits (Linux)
// Không cần thư viện automaxprocs nữa!
// GODEBUG=containermaxprocs=0 để disable

// 4. FlightRecorder — runtime tracing
import "runtime/trace"

recorder := trace.NewFlightRecorder()
recorder.Start()
// ... program runs ...
// Khi có vấn đề:
var buf bytes.Buffer
recorder.WriteTo(&buf)
os.WriteFile("trace.out", buf.Bytes(), 0644)

// 5. testing.T.Attr — metadata cho tests
func TestSomething(t *testing.T) {
    t.Attr("category", "integration")
    t.Attr("priority", "high")
    // Hiển thị trong -json output
}

// 6. JSON v2 (experimental)
// GOEXPERIMENT=jsonv2 go build ./...
// encoding/json/v2 — major revision
// Stricter by default, more options

// 7. slog.GroupAttrs
attrs := slog.GroupAttrs("request",
    slog.String("method", "GET"),
    slog.String("path", "/users"),
)
logger.Info("incoming", attrs)

// 8. reflect.TypeAssert[T] (Go 1.25)
import "reflect"
// Type-safe reflection assertions
```

---

## 32. Production Checklist

### Graceful Shutdown

```go
func main() {
    srv := &http.Server{Addr: ":8080", Handler: handler}

    // Channel để nhận OS signals
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

    // Start server
    go func() {
        if err := srv.ListenAndServe(); err != http.ErrServerClosed {
            log.Fatal(err)
        }
    }()
    log.Println("Server started on :8080")

    // Block until signal
    <-quit
    log.Println("Shutting down server...")

    // Graceful shutdown với timeout
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        log.Printf("Server forced to shutdown: %v", err)
    }

    // Cleanup resources
    db.Close()
    redisConn.Close()

    log.Println("Server stopped gracefully")
}
```

### Full Production Checklist

```
CORRECTNESS:
✅ Context propagation — mọi I/O operation nhận context
✅ Error handling — không discard errors, wrap với context
✅ Input validation — validate tại entry points
✅ Race condition check — go test -race (CI bắt buộc)
✅ nil pointer protection — check trước khi dereference
✅ Goroutine leak prevention — mọi goroutine có exit condition

OBSERVABILITY:
✅ Structured logging — log/slog với JSON output
✅ Distributed tracing — OpenTelemetry
✅ Metrics — Prometheus counters, histograms, gauges
✅ Health checks — /healthz (alive), /readyz (ready to serve)
✅ Profiling endpoint — pprof (development/internal)

RESILIENCE:
✅ Timeout trên tất cả external calls
✅ Retry với exponential backoff
✅ Circuit breaker — prevent cascade failures
✅ Rate limiting — protect từ abuse/overload
✅ Graceful shutdown — drain requests trước khi exit

PERFORMANCE:
✅ Connection pooling — HTTP client, DB, Redis
✅ Caching strategy — Redis/in-memory với TTL
✅ Batch operations — giảm roundtrips
✅ Pagination — không return toàn bộ dataset

SECURITY:
✅ Authentication & Authorization
✅ HTTPS với valid certificate
✅ Security headers — CORS, CSP, HSTS, X-Frame-Options
✅ Secret management — environment variables, vault
✅ SQL injection prevention — prepared statements
✅ Rate limiting
✅ govulncheck trong CI

DEPLOYMENT:
✅ Config từ environment — 12-factor app
✅ Dependency injection — testable architecture
✅ Dockerfile multi-stage build — minimal image
✅ Liveness & Readiness probes
✅ Resource limits — CPU, memory requests/limits
✅ Horizontal scaling ready — stateless design
```

---

## 33. Những Sai Lầm Phổ Biến Nhất

### 1. Goroutine Leak

```go
// BAD: goroutine blocked mãi mãi
func leaky() {
    ch := make(chan int)
    go func() {
        val := <-ch  // block forever!
    }()
    // return mà không gửi vào ch, không close ch
}

// GOOD: có exit condition
func notLeaky(ctx context.Context) {
    ch := make(chan int)
    go func() {
        select {
        case val := <-ch:
            process(val)
        case <-ctx.Done():
            return  // clean exit
        }
    }()
}
```

### 2. Ignoring Errors

```go
// BAD — NEVER in production
result, _ := doSomething()
os.Remove(filename)  // silently fail

// GOOD
result, err := doSomething()
if err != nil {
    return fmt.Errorf("doSomething: %w", err)
}
if err := os.Remove(filename); err != nil {
    log.Warn("failed to remove file", "error", err, "file", filename)
}
```

### 3. Nil Interface Gotcha

```go
// BAD
func getError() error {
    var err *MyError = nil
    return err   // trả về non-nil interface! (type=*MyError, value=nil)
}

// GOOD
func getError() error {
    return nil   // trả về nil interface
}
```

### 4. String Concatenation trong Loop

```go
// BAD: O(n²)
s := ""
for _, w := range words { s += w }

// GOOD: O(n)
var b strings.Builder
for _, w := range words { b.WriteString(w) }
s := b.String()
```

### 5. Copying Mutex

```go
// BAD: copy mutex sau khi đã dùng
type Counter struct {
    mu    sync.Mutex
    count int
}
c := Counter{}
c2 := c    // BAD! copy mutex
// go vet sẽ warn

// GOOD: truyền pointer
func doSomething(c *Counter) { ... }
```

### 6. Map Concurrent Access → Panic

```go
// BAD: panic tại runtime
m := make(map[string]int)
go func() { m["key"] = 1 }()
go func() { _ = m["key"] }()  // panic!

// GOOD: sync.Map hoặc sync.RWMutex
var m sync.Map
m.Store("key", 1)
v, _ := m.Load("key")
```

### 7. Default HTTP Client Không Có Timeout

```go
// BAD: http.DefaultClient không có timeout → hang forever!
resp, err := http.Get("https://slow-server.com/api")

// GOOD: client với timeout
client := &http.Client{Timeout: 30 * time.Second}
resp, err := client.Get("https://server.com/api")
```

### 8. Không Close Response Body

```go
// BAD: resource leak
resp, err := client.Get(url)
if err != nil { return err }
data, _ := io.ReadAll(resp.Body)
// Body không được close!

// GOOD
resp, err := client.Get(url)
if err != nil { return err }
defer resp.Body.Close()   // LUÔN defer close
data, err := io.ReadAll(resp.Body)
```

### 9. time.After Trong Loop → Memory Leak

```go
// BAD: mỗi iteration tạo timer, GC không collect ngay
for {
    select {
    case <-ch:
        process()
    case <-time.After(5 * time.Second):  // LEAK!
        return
    }
}

// GOOD: tạo timer một lần, reuse
timer := time.NewTimer(5 * time.Second)
defer timer.Stop()
for {
    timer.Reset(5 * time.Second)
    select {
    case <-ch:
        process()
    case <-timer.C:
        return
    }
}
```

### 10. Slice Sharing Bug

```go
a := []int{1, 2, 3, 4, 5}
b := a[1:3]    // b share underlying array với a
b[0] = 99      // a[1] cũng thành 99!

// GOOD: explicit copy
b := make([]int, 2)
copy(b, a[1:3])
```

### 11. Context Không Được Cancel

```go
// BAD: memory leak
ctx, _ = context.WithTimeout(parent, 5*time.Second)
// cancel function bị bỏ qua → resources không được release

// GOOD
ctx, cancel := context.WithTimeout(parent, 5*time.Second)
defer cancel()  // LUÔN defer cancel!
```

### 12. Integer Overflow

```go
// BAD
func addUint8(a, b uint8) uint8 {
    return a + b  // overflow nếu a+b > 255!
}

// GOOD
func addUint8(a, b uint8) (uint8, error) {
    if int(a)+int(b) > math.MaxUint8 {
        return 0, fmt.Errorf("overflow: %d + %d > 255", a, b)
    }
    return a + b, nil
}
```

### 13. Variable Shadowing

```go
x := 1
if condition {
    x := 2         // shadow outer x!
    fmt.Println(x) // 2
}
fmt.Println(x)     // 1 — outer x không thay đổi

// Nguy hiểm với err:
err := doSomething()
if err != nil {
    err := fmt.Errorf("context: %w", err)  // shadow!
    // outer err không thay đổi
    return err  // đây là inner err
}
```

### 14. Ranging Map và Modify

```go
// SAFE: delete trong range (Go cho phép)
for k := range m {
    if shouldDelete(k) {
        delete(m, k)  // OK trong Go
    }
}

// UNSAFE: add trong range (behavior undefined-ish)
for k, v := range m {
    m[newKey] = v  // có thể thấy hoặc không thấy new entry
}
```

---

## 34. Bài Tập Thực Hành

### Phase 1 — Beginner

1. **Fibonacci** — viết cả iterative và recursive, so sánh performance bằng benchmark
2. **Palindrome Unicode** — kiểm tra palindrome cho string có emoji và ký tự đặc biệt
3. **Stack** — implement generic Stack với Push/Pop/Peek/Len
4. **Word Counter** — đếm tần suất từng từ, sort theo frequency

### Phase 2 — Intermediate

5. **LRU Cache** — implement với map + doubly linked list (O(1) get/put)
6. **URL Shortener** — REST API với in-memory storage, CRUD operations
7. **CLI Task Manager** — với file persistence (JSON), support add/done/list/delete
8. **Rate Limiter** — implement token bucket từ scratch

### Phase 3 — Advanced

9. **Worker Pool** — generic worker pool với context cancellation và error handling
10. **Chat Server** — TCP chat server với multiple rooms, WebSockets
11. **Mini Redis** — in-memory key-value store với basic commands (GET/SET/DEL/EXPIRE)
12. **Job Queue** — persistent queue với PostgreSQL, retry logic, dead letter queue

### Phase 4 — Expert

13. **Distributed Key-Value Store** — với Raft consensus (không dùng library)
14. **Build Your Own HTTP Router** — pattern matching, middleware chain
15. **Go Plugin System** — dynamic loading với os/exec và gRPC
16. **Simple Interpreter** — parse và execute simple language trong Go

---

## 35. Tài Nguyên Học Tập

### Sách — Theo Thứ Tự Học

| # | Sách | Level | Ghi Chú |
|---|---|---|---|
| 1 | **The Go Programming Language** — Donovan & Kernighan | Beginner | "Go Bible" |
| 2 | **Learning Go** — Jon Bodner (2nd Ed, 2024) | Intermediate | Idiomatic Go, hiện đại |
| 3 | **Concurrency in Go** — Katherine Cox-Buday | Inter-Advanced | Deep dive concurrency |
| 4 | **100 Go Mistakes** — Teiva Harsanyi | Advanced | PHẢI ĐỌC! |
| 5 | **Let's Go / Let's Go Further** — Alex Edwards | Advanced | REST API production |
| 6 | **Cloud Native Go** — Matthew Titmus | Expert | Distributed systems |
| 7 | **Domain-Driven Design with Golang** — M. Boyle | Expert | Architecture |

### Online Resources

**Official:**
- `go.dev/tour` — Tour of Go (bắt đầu tại đây)
- `go.dev/doc` — Official documentation
- `go.dev/blog` — Go team blog (must read!)
- `pkg.go.dev` — Package documentation
- `go.dev/ref/spec` — Language specification

**Practice:**
- `exercism.org/tracks/go` — Exercises với mentoring
- `gophercises.com` — Real-world exercises (Jon Calhoun)
- `github.com/quii/learn-go-with-tests` — TDD approach (22k+ stars)
- `play.golang.org` — Go Playground (test nhanh không cần setup)

**Advanced:**
- `ardanlabs.com/blog` — Ultimate Go (William Kennedy)
- `dave.cheney.net` — High quality Go articles
- `100go.co` — 100 mistakes online
- `go.dev/blog/greenteagc` — Green Tea GC deep dive

**Community:**
- `gophers.slack.com`
- `reddit.com/r/golang`
- `gophercon.com` — talks archive

---

## 36. Milestones & Checkpoints

### ✅ Checkpoint 1 — Beginner (Tuần 1-6)
- [ ] Viết và chạy chương trình Go cơ bản
- [ ] Hiểu slice internals: len, cap, sharing behavior
- [ ] Viết struct với value và pointer receivers
- [ ] Handle errors idiomatic (không dùng try-catch)
- [ ] Viết unit tests với table-driven approach
- [ ] Dùng `go fmt`, `go vet` thành thói quen

### ✅ Checkpoint 2 — Intermediate (Tuần 6-16)
- [ ] Implement interfaces và dependency injection
- [ ] Viết concurrent program với goroutines + channels
- [ ] Tránh data races (dùng `-race` detector thường xuyên)
- [ ] Viết HTTP server từ standard library (Go 1.22 routing)
- [ ] Dùng context đúng cách trong toàn bộ call chain
- [ ] Viết benchmarks và profile code
- [ ] Hiểu và dùng generics cơ bản

### ✅ Checkpoint 3 — Senior (Tuần 16-30)
- [ ] Thiết kế API với clean architecture
- [ ] Implement worker pool, rate limiter, circuit breaker
- [ ] Integration tests với real database
- [ ] Deploy service lên production với observability
- [ ] Hiểu và tune GC với GOGC/GOMEMLIMIT
- [ ] Profile CPU/memory và optimize real bottlenecks
- [ ] Review code và đưa ra architectural feedback

### ✅ Checkpoint 4 — Expert (Tuần 30+)
- [ ] Giải thích GMP scheduler model sâu sắc
- [ ] Dùng escape analysis để optimize allocations
- [ ] Contribute open source Go project
- [ ] Viết custom sync primitive an toàn
- [ ] Mentor người khác về Go idioms và patterns
- [ ] Build production-grade distributed system

---

## 📌 Quick Reference — Go Cheatsheet

```go
// Khai báo biến
var x int = 42
x := 42            // short (trong function)
const Pi = 3.14    // constant

// Function
func name(param Type) ReturnType { }
func multi() (int, error) { }         // multiple returns

// Interface
type I interface { Method() string }

// Goroutine + Channel
go func() { }()
ch := make(chan int, 10)
ch <- val; val = <-ch

// Error
if err != nil { return fmt.Errorf("context: %w", err) }

// Defer
defer func() { recover() }()    // panic recovery

// Type assertion
if v, ok := i.(string); ok { }

// Go 1.22+ HTTP routing
mux.HandleFunc("GET /path/{id}", handler)
id := r.PathValue("id")

// Go 1.25 WaitGroup
var wg sync.WaitGroup
wg.Go(func() { doWork() })
wg.Wait()

// Green Tea GC
// GOEXPERIMENT=greenteagc go build ./...

// Container-aware GOMAXPROCS (auto từ Go 1.25)
// Không cần automaxprocs library nữa!
```

---

*Tài liệu tổng hợp từ: Go official docs (go.dev), Go 1.22-1.25 Release Notes, "100 Go Mistakes" (Teiva Harsanyi), "Concurrency in Go" (Cox-Buday), "The Go Programming Language" (Donovan & Kernighan), Ardan Labs Ultimate Go, Go blog (go.dev/blog), GopherCon talks.*

*Cập nhật: May 2026 — Go 1.25.x*