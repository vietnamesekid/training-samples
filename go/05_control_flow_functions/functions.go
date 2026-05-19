package main

import (
	"errors"
	"fmt"
	"math"
)

// === Basic Function ===

// Multiple return values — very common in Go
func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

// Named return values — use sparingly, avoid overuse
// Useful when: returning multiple values of the same type, defer modifies return values
func minMax(nums []int) (min, max int) {
	if len(nums) == 0 {
		return // naked return: returns zero values for min and max
	}
	min, max = nums[0], nums[0]
	for _, n := range nums[1:] {
		if n < min {
			min = n
		}
		if n > max {
			max = n
		}
	}
	return // naked return: returns current min and max
}

// === Variadic Functions ===

// ...T: accepts 0 or more arguments — internally treated as []T
func sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

// Pass a slice to a variadic function with ...
func joinStrings(sep string, strs ...string) string {
	result := ""
	for i, s := range strs {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}

// === Functions as First-Class Values ===

// Function type
type MathFunc func(float64, float64) float64

// Higher-order function: accepts a function as an argument
func apply(f MathFunc, a, b float64) float64 {
	return f(a, b)
}

// Higher-order function: returns a function (function factory)
func multiplier(factor float64) func(float64) float64 {
	return func(x float64) float64 {
		return x * factor
	}
}

// === Closures ===

// Closure captures variables from the outer scope
func makeCounter() func() int {
	count := 0
	return func() int {
		count++ // capture and modify count
		return count
	}
}

// Closure with memoization
func makeMemoFib() func(int) int {
	cache := map[int]int{0: 0, 1: 1}
	var fib func(int) int
	fib = func(n int) int {
		if v, ok := cache[n]; ok {
			return v
		}
		cache[n] = fib(n-1) + fib(n-2)
		return cache[n]
	}
	return fib
}

// === panic & recover ===

// safeDivide wraps a function that may panic and converts it into an error
func safeDivide(a, b int) (result int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered from panic: %v", r)
		}
	}()
	result = a / b // will panic if b == 0
	return
}

// === init() function ===

// init() runs automatically before main(), after package-level vars
// A package can have multiple init() functions
// Order: package vars → init() → main()
var initMessage string

func init() {
	initMessage = "init() đã chạy trước main()"
}

// === Recursion ===

func factorial(n int) int {
	if n <= 1 {
		return 1
	}
	return n * factorial(n-1)
}

// Fibonacci in tail recursion style (Go does not optimize tail calls)
func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

func demoFunctions() {
	fmt.Println("init() result:", initMessage)

	fmt.Println("\n--- Basic & Multiple Returns ---")
	result, err := divide(10, 3)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("10/3 = %.4f\n", result)
	}

	_, err = divide(10, 0)
	fmt.Printf("10/0 error: %v\n", err)

	min, max := minMax([]int{3, 1, 4, 1, 5, 9, 2, 6})
	fmt.Printf("minMax: min=%d, max=%d\n", min, max)

	fmt.Println("\n--- Variadic ---")
	fmt.Printf("sum(): %d\n", sum())
	fmt.Printf("sum(1,2,3): %d\n", sum(1, 2, 3))
	nums := []int{1, 2, 3, 4, 5}
	fmt.Printf("sum(nums...): %d\n", sum(nums...)) // unpack slice

	fmt.Printf("joinStrings: %s\n", joinStrings(", ", "Go", "Python", "Rust"))

	fmt.Println("\n--- First-Class Functions ---")
	add := func(a, b float64) float64 { return a + b }
	mul := func(a, b float64) float64 { return a * b }
	fmt.Printf("apply(add, 3, 4) = %.0f\n", apply(add, 3, 4))
	fmt.Printf("apply(mul, 3, 4) = %.0f\n", apply(mul, 3, 4))
	fmt.Printf("apply(math.Pow, 2, 10) = %.0f\n", apply(math.Pow, 2, 10))

	double := multiplier(2)
	triple := multiplier(3)
	fmt.Printf("double(5) = %.0f, triple(5) = %.0f\n", double(5), triple(5))

	fmt.Println("\n--- Closures ---")
	counter1 := makeCounter()
	counter2 := makeCounter() // separate counter
	fmt.Printf("counter1: %d, %d, %d\n", counter1(), counter1(), counter1())
	fmt.Printf("counter2: %d (riêng biệt)\n", counter2())

	fib := makeMemoFib()
	fmt.Printf("fibonacci(10) = %d (memoized)\n", fib(10))
	fmt.Printf("fibonacci(20) = %d (memoized)\n", fib(20))

	fmt.Println("\n--- Immediately Invoked Function (IIFE) ---")
	result2 := func(x, y int) int {
		return x * y
	}(6, 7)
	fmt.Printf("IIFE result: %d\n", result2)

	fmt.Println("\n--- panic & recover ---")
	r1, err1 := safeDivide(10, 2)
	fmt.Printf("safeDivide(10,2): result=%d, err=%v\n", r1, err1)
	r2, err2 := safeDivide(10, 0)
	fmt.Printf("safeDivide(10,0): result=%d, err=%v\n", r2, err2)

	fmt.Println("\n--- Recursion ---")
	for i := range 8 {
		fmt.Printf("  factorial(%d) = %d\n", i, factorial(i))
	}
	for i := range 10 {
		fmt.Printf("  fibonacci(%d) = %d\n", i, fibonacci(i))
	}
}
