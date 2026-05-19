// Lesson 18: Generics — type parameters in Go 1.18+
// Run: go run .
package main

import "fmt"

func main() {
	fmt.Println("=== GENERIC FUNCTIONS ===")
	demoGenericFunctions()

	fmt.Println("\n=== GENERIC TYPES ===")
	demoGenericTypes()

	fmt.Println("\n=== STDLIB GENERICS (slices, maps, iter) ===")
	demoStdlibGenerics()

	fmt.Println("\n=== KHI NÀO KHÔNG DÙNG GENERICS ===")
	demoWhenNotToUse()
}
