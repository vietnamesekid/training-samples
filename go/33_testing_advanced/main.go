// Bài 33: Advanced Testing Patterns
// Hand-written mocks, golden files, t.TempDir, t.Cleanup, integration build tags
// Chạy tests: go test ./... -v
// Chạy với race detector: go test -race ./...
package main

import "fmt"

func main() {
	fmt.Println("=== Advanced Testing Patterns ===")
	fmt.Println("Chạy: go test ./... -v để xem tất cả tests")
	fmt.Println()
	fmt.Println("Patterns covered trong test files:")
	fmt.Println("  1. Hand-written mocks (không dùng mockery/gomock)")
	fmt.Println("  2. Golden files (.golden) cho snapshot testing")
	fmt.Println("  3. t.TempDir() — auto cleanup temp directories")
	fmt.Println("  4. t.Cleanup() — register cleanup functions")
	fmt.Println("  5. Table-driven tests với subtests")
	fmt.Println("  6. t.Parallel() — chạy tests song song")
	fmt.Println("  7. Build tags: //go:build integration")
	fmt.Println("  8. Test helpers với t.Helper()")
}
