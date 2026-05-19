// Lesson 10: Goroutines & Channels — concurrency primitives in Go
// Run: go run .
package main

import "fmt"

func main() {
	fmt.Println("=== GOROUTINES ===")
	demoGoroutines()

	fmt.Println("\n=== CHANNELS ===")
	demoChannels()

	fmt.Println("\n=== PIPELINE PATTERN ===")
	demoPipeline()
}
