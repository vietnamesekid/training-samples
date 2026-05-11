// Bài 10: Goroutines & Channels — concurrency primitives của Go
// Chạy: go run .
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
