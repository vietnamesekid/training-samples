// Lesson 22: Design Patterns in Go
// Run: go run .
package main

import "fmt"

func main() {
	fmt.Println("=== FUNCTIONAL OPTIONS ===")
	demoFunctionalOptions()

	fmt.Println("\n=== BUILDER PATTERN ===")
	demoBuilder()

	fmt.Println("\n=== REPOSITORY PATTERN & DI ===")
	demoRepository()

	fmt.Println("\n=== OBSERVER / EVENT BUS ===")
	demoObserver()

	fmt.Println("\n=== STRATEGY PATTERN ===")
	demoStrategy()

	fmt.Println("\n=== SINGLETON (sync.Once) ===")
	demoSingleton()

	fmt.Println("\n=== MIDDLEWARE CHAIN ===")
	demoMiddlewareChain()
}
