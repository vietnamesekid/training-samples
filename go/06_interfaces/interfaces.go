package main

import (
	"fmt"
	"io"
	"math"
	"strings"
)

// === Duck Typing — implicit interface implementation ===

// Interface is implemented implicitly — no need to declare "implements"
type Shape interface {
	Area() float64
	Perimeter() float64
}

type Rectangle struct {
	Width, Height float64
}

func (r Rectangle) Area() float64      { return r.Width * r.Height }
func (r Rectangle) Perimeter() float64 { return 2 * (r.Width + r.Height) }

type Circle struct {
	Radius float64
}

func (c Circle) Area() float64      { return math.Pi * c.Radius * c.Radius }
func (c Circle) Perimeter() float64 { return 2 * math.Pi * c.Radius }

// Polymorphism via interface
func printShape(s Shape) {
	fmt.Printf("  %T: area=%.2f, perimeter=%.2f\n", s, s.Area(), s.Perimeter())
}

func totalArea(shapes []Shape) float64 {
	total := 0.0
	for _, s := range shapes {
		total += s.Area()
	}
	return total
}

// === Interface Composition ===

// io.Reader, io.Writer, io.Closer from stdlib
type ReadWriter interface {
	io.Reader
	io.Writer
}

// Custom interface composition
type Storage interface {
	io.Reader
	io.Writer
	io.Closer
	Name() string
}

// === fmt.Stringer — the most commonly used interface ===
type Color struct {
	R, G, B uint8
}

// Implement fmt.Stringer: fmt calls this automatically with %s, %v, Println
func (c Color) String() string {
	return fmt.Sprintf("#%02X%02X%02X", c.R, c.G, c.B)
}

// === empty interface / any ===

func printAny(v any) {
	fmt.Printf("  type=%T, value=%v\n", v, v)
}

// === Type Assertion ===

func demoTypeAssertion() {
	var s Shape = Rectangle{Width: 3, Height: 4}

	// Safe type assertion — use the 2-value form
	if r, ok := s.(Rectangle); ok {
		fmt.Printf("  Rectangle: %+v\n", r)
	}

	// Unsafe assertion — panics if the type is wrong
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("  Unsafe assertion panic: %v\n", r)
		}
	}()
	c := s.(Circle) // ← PANIC because s is Rectangle, not Circle
	fmt.Println("  Circle:", c)
}

func demoInterfaces() {
	fmt.Println("\n--- Polymorphism ---")
	shapes := []Shape{
		Rectangle{Width: 3, Height: 4},
		Circle{Radius: 5},
		Rectangle{Width: 2, Height: 6},
	}
	for _, s := range shapes {
		printShape(s)
	}
	fmt.Printf("  Total area: %.2f\n", totalArea(shapes))

	fmt.Println("\n--- fmt.Stringer ---")
	colors := []Color{{255, 0, 0}, {0, 255, 0}, {0, 0, 255}}
	for _, c := range colors {
		fmt.Println(" ", c) // calls c.String() automatically
	}

	fmt.Println("\n--- Interface với io package ---")
	// strings.NewReader implements io.Reader
	r := strings.NewReader("Hello, Go!")
	buf := make([]byte, 4)
	for {
		n, err := r.Read(buf)
		if err == io.EOF {
			break
		}
		fmt.Printf("  Read %d bytes: %q\n", n, buf[:n])
	}

	fmt.Println("\n--- empty interface (any) ---")
	values := []any{42, "hello", 3.14, true, []int{1, 2, 3}}
	for _, v := range values {
		printAny(v)
	}

	fmt.Println("\n--- Type Switch ---")
	for _, v := range values {
		switch t := v.(type) {
		case int:
			fmt.Printf("  int: %d (double=%d)\n", t, t*2)
		case string:
			fmt.Printf("  string: %q (len=%d)\n", t, len(t))
		case float64:
			fmt.Printf("  float64: %.2f\n", t)
		case bool:
			fmt.Printf("  bool: %t\n", t)
		default:
			fmt.Printf("  other: %T = %v\n", t, t)
		}
	}

	fmt.Println("\n--- Type Assertion ---")
	demoTypeAssertion()
}
