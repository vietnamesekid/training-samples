// Lesson 8: Packages & Modules — organizing code in Go
// Run: go run .
package main

import (
	"fmt"

	"github.com/tranminhquang/training-samples/go/08_packages_modules/internal/config"
	"github.com/tranminhquang/training-samples/go/08_packages_modules/pkg/validator"
)

// Package-level variable — initialized before init()
var greeting = "Hello"

// init() runs after package-level vars, before main()
// A package can have multiple init() functions
// Order: imported packages init() → current package vars → current package init() → main()
func init() {
	greeting = greeting + ", World"
	fmt.Println("[init] greeting initialized:", greeting)
}

func init() {
	// Multiple init() functions are allowed — they run in the order they appear in the file
	fmt.Println("[init] second init() runs after first")
}

func main() {
	fmt.Println("\n=== Package Structure ===")
	fmt.Println("greeting:", greeting)

	fmt.Println("\n=== Internal Package ===")
	// internal/ package is only accessible within this module
	cfg := config.Load()
	fmt.Printf("Config: %+v\n", cfg)
	fmt.Printf("DB URL: %s\n", cfg.DatabaseURL())

	fmt.Println("\n=== Public Package (pkg/) ===")
	v := validator.New()
	results := v.ValidateUser("Alice", "alice@example.com", 25)
	for _, r := range results {
		fmt.Printf("  %s\n", r)
	}

	results2 := v.ValidateUser("", "invalid-email", -1)
	for _, r := range results2 {
		fmt.Printf("  %s\n", r)
	}

	fmt.Println("\n=== Package Concepts ===")
	fmt.Println("Exported (PascalCase): Config, Load, ValidateUser")
	fmt.Println("Unexported (camelCase): chỉ dùng trong package")
	fmt.Println()
	fmt.Println("go.mod structure:")
	fmt.Println("  module github.com/yourname/project  ← module path")
	fmt.Println("  go 1.26.1                           ← minimum Go version")
	fmt.Println()
	fmt.Println("Các lệnh module:")
	fmt.Println("  go mod init <module>   → tạo go.mod")
	fmt.Println("  go get <pkg>@<version> → thêm dependency")
	fmt.Println("  go mod tidy            → dọn dẹp go.mod và go.sum")
	fmt.Println("  go mod download        → download dependencies")
	fmt.Println("  go list -m all         → liệt kê tất cả dependencies")
	fmt.Println()
	fmt.Println("Blank import (side effect):")
	fmt.Println("  import _ \"net/http/pprof\"  → chỉ chạy init(), register handlers")
	fmt.Println("  import _ \"image/png\"       → register PNG decoder")
}
