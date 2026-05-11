// Bài 4: Structs & Methods — kiểu dữ liệu tự định nghĩa
// Chạy: go run .
package main

import (
	"fmt"
	"math"
	"reflect"
	"unsafe"
)

// === Struct cơ bản ===

// Quy ước: PascalCase cho exported, camelCase cho unexported
type Person struct {
	Name    string  // exported — accessible từ bên ngoài package
	Age     int
	Email   string
	address Address // unexported — chỉ dùng trong package này
	score   int     // unexported
}

type Address struct {
	Street string
	City   string
	Zip    string
}

// Constructor pattern — Go không có constructor, dùng hàm NewXxx
func NewPerson(name string, age int, email string) (*Person, error) {
	if name == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}
	if age < 0 || age > 150 {
		return nil, fmt.Errorf("invalid age: %d", age)
	}
	return &Person{
		Name:  name,
		Age:   age,
		Email: email,
	}, nil
}

// Value receiver — nhận COPY của struct, KHÔNG mutate được
// Dùng khi: struct nhỏ, không cần mutate, type là primitive/immutable
func (p Person) String() string {
	return fmt.Sprintf("%s (age=%d, email=%s)", p.Name, p.Age, p.Email)
}

func (p Person) IsAdult() bool {
	return p.Age >= 18
}

// Pointer receiver — nhận POINTER, CÓ THỂ mutate
// NGUYÊN TẮC: Nếu bất kỳ method nào cần pointer receiver,
// dùng pointer receiver cho TẤT CẢ methods của type đó.
func (p *Person) Birthday() {
	p.Age++
}

func (p *Person) SetAddress(street, city, zip string) {
	p.address = Address{Street: street, City: city, Zip: zip}
}

func (p *Person) GetAddress() Address {
	return p.address
}

// === Embedding — composition, không phải inheritance ===

type Employee struct {
	Person            // embedded — promotes all methods and fields
	Company   string
	Salary    float64
	JobTitle  string
}

// Employee có thể override method của Person
func (e Employee) String() string {
	return fmt.Sprintf("%s @ %s (salary=%.0f)", e.Name, e.Company, e.Salary)
}

// === Struct Tags — metadata cho reflection ===
type User struct {
	ID        int     `json:"id"`
	Username  string  `json:"username"`
	Password  string  `json:"-"`                          // không bao giờ marshal
	Email     string  `json:"email,omitempty"`            // bỏ qua nếu empty
	Age       int     `json:"age" validate:"min=0,max=150"`
	Score     float64 `json:"score" db:"user_score"`
}

// === Anonymous struct — dùng cho one-off data shapes ===

// === Struct comparison ===

type Point struct{ X, Y float64 }

func (p Point) Distance(other Point) float64 {
	dx := p.X - other.X
	dy := p.Y - other.Y
	return math.Sqrt(dx*dx + dy*dy)
}

// === Struct padding và memory layout ===

// BadLayout: 24 bytes (padding waste)
type BadLayout struct {
	A bool    // 1 byte
	// 7 bytes padding
	B float64 // 8 bytes
	C bool    // 1 byte
	// 7 bytes padding
}

// GoodLayout: 16 bytes (fields sorted by size, largest first)
type GoodLayout struct {
	B float64 // 8 bytes
	C bool    // 1 byte
	A bool    // 1 byte
	// 6 bytes padding
}

func main() {
	fmt.Println("=== 1. Tạo Struct ===")

	// Tạo với struct literal
	p1 := Person{Name: "Alice", Age: 30, Email: "alice@example.com"}
	fmt.Println("Literal:", p1) // gọi p1.String() qua fmt.Stringer

	// Tạo với constructor
	p2, err := NewPerson("Bob", 25, "bob@example.com")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Constructor:", p2)
	}

	// Zero value struct — tất cả fields là zero value
	var p3 Person
	fmt.Printf("Zero value: %+v\n", p3) // %+v in kèm field names

	// Anonymous struct — dùng cho one-off, config, test data
	config := struct {
		Host string
		Port int
		TLS  bool
	}{
		Host: "localhost",
		Port: 8080,
		TLS:  false,
	}
	fmt.Printf("Anonymous struct: %+v\n", config)

	fmt.Println("\n=== 2. Methods ===")
	p2.Birthday()
	fmt.Printf("Sau Birthday: %s\n", p2)
	fmt.Printf("IsAdult: %t\n", p2.IsAdult())

	p2.SetAddress("123 Main St", "Hanoi", "10000")
	addr := p2.GetAddress()
	fmt.Printf("Address: %+v\n", addr)

	fmt.Println("\n=== 3. Embedding & Composition ===")
	emp := Employee{
		Person:   Person{Name: "Carol", Age: 28, Email: "carol@corp.com"},
		Company:  "Gophers Inc",
		Salary:   50000,
		JobTitle: "Software Engineer",
	}

	// Promoted fields và methods
	fmt.Printf("emp.Name = %s (promoted từ Person)\n", emp.Name)    // emp.Person.Name
	fmt.Printf("emp.IsAdult() = %t (promoted)\n", emp.IsAdult())     // emp.Person.IsAdult()
	emp.Birthday()  // emp.Person.Birthday()
	fmt.Printf("emp.String() = %s (override)\n", emp)               // Employee.String()

	fmt.Println("\n=== 4. Struct Tags & Reflection ===")
	u := User{ID: 1, Username: "gopher", Password: "secret", Email: "gopher@go.dev", Age: 5}
	t := reflect.TypeOf(u)
	for i := range t.NumField() {
		field := t.Field(i)
		fmt.Printf("  Field: %-10s json=%q validate=%q\n",
			field.Name,
			field.Tag.Get("json"),
			field.Tag.Get("validate"),
		)
	}

	fmt.Println("\n=== 5. Struct So Sánh ===")
	pt1 := Point{1, 2}
	pt2 := Point{1, 2}
	pt3 := Point{3, 4}
	fmt.Printf("pt1 == pt2: %t\n", pt1 == pt2)
	fmt.Printf("pt1 == pt3: %t\n", pt1 == pt3)
	fmt.Printf("Distance(pt1, pt3): %.3f\n", pt1.Distance(pt3))

	// reflect.DeepEqual — dùng khi struct chứa slice/map (không comparable)
	type Config struct {
		Tags []string
	}
	c1 := Config{Tags: []string{"a", "b"}}
	c2 := Config{Tags: []string{"a", "b"}}
	// c1 == c2 ← compile error: invalid operation
	fmt.Printf("reflect.DeepEqual: %t\n", reflect.DeepEqual(c1, c2))

	fmt.Println("\n=== 6. Struct Memory Layout ===")
	fmt.Printf("BadLayout:  %d bytes\n", unsafe.Sizeof(BadLayout{}))
	fmt.Printf("GoodLayout: %d bytes\n", unsafe.Sizeof(GoodLayout{}))
	fmt.Println("  → Sắp xếp fields từ lớn đến nhỏ tiết kiệm memory")
	fmt.Printf("  GoodLayout.B offset: %d\n", unsafe.Offsetof(GoodLayout{}.B))
	fmt.Printf("  GoodLayout.C offset: %d\n", unsafe.Offsetof(GoodLayout{}.C))
	fmt.Printf("  GoodLayout.A offset: %d\n", unsafe.Offsetof(GoodLayout{}.A))
}
