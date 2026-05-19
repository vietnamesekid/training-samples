// Lesson 16: Testing — writing tests correctly in Go
// Run tests: go test ./...
// Run with verbose: go test -v ./...
// Run benchmark: go test -bench=. -benchmem ./...
// Run fuzz: go test -fuzz=FuzzReverse -fuzztime=10s ./...
// View coverage: go test -cover ./...
// Coverage HTML: go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out
package main

import (
	"errors"
	"fmt"
)

// Calculator — this is the code we want to test
type Calculator struct {
	history []string
}

func NewCalculator() *Calculator {
	return &Calculator{}
}

var ErrDivisionByZero = errors.New("division by zero")

func (c *Calculator) Add(a, b float64) float64 {
	result := a + b
	c.history = append(c.history, fmt.Sprintf("%v+%v=%v", a, b, result))
	return result
}

func (c *Calculator) Sub(a, b float64) float64 {
	result := a - b
	c.history = append(c.history, fmt.Sprintf("%v-%v=%v", a, b, result))
	return result
}

func (c *Calculator) Mul(a, b float64) float64 {
	result := a * b
	c.history = append(c.history, fmt.Sprintf("%v*%v=%v", a, b, result))
	return result
}

func (c *Calculator) Div(a, b float64) (float64, error) {
	if b == 0 {
		return 0, ErrDivisionByZero
	}
	result := a / b
	c.history = append(c.history, fmt.Sprintf("%v/%v=%v", a, b, result))
	return result, nil
}

func (c *Calculator) History() []string {
	return c.history
}

func (c *Calculator) Clear() {
	c.history = nil
}

// Reverse reverses a string — used to demo fuzz testing
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// main: entry point, but the interesting code is in *_test.go files
func main() {
	c := NewCalculator()
	fmt.Printf("3 + 4 = %.1f\n", c.Add(3, 4))
	fmt.Printf("10 / 2 = ")
	if r, err := c.Div(10, 2); err == nil {
		fmt.Printf("%.1f\n", r)
	}
	fmt.Printf("Reverse(%q) = %q\n", "Hello, 🌍!", Reverse("Hello, 🌍!"))
	fmt.Println("\nChạy tests: go test -v ./...")
	fmt.Println("Chạy benchmark: go test -bench=. -benchmem ./...")
}
