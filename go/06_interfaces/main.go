// Lesson 6: Interfaces — the most important concept in Go
// Run: go run .
package main

import "fmt"

func main() {
	fmt.Println("=== INTERFACES ===")
	demoInterfaces()

	fmt.Println("\n=== NIL INTERFACE GOTCHA ===")
	demoNilInterfaceGotcha()

	fmt.Println("\n=== INTERFACE BEST PRACTICES ===")
	demoInterfaceBestPractices()

	fmt.Println("\n=== FUNCTIONAL OPTIONS PATTERN ===")
	demoFunctionalOptions()
}
