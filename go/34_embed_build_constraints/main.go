// Bài 34: Embed & Build Constraints
// //go:embed, //go:generate, //go:build constraints, GOOS/GOARCH, cross-compile
// Chạy: go run .
package main

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"runtime"
	"text/template"
)

// embed directives phải ngay trước var declaration (không có blank line)

//go:embed templates/*.tmpl
var templateFiles embed.FS

//go:embed static
var staticFiles embed.FS

//go:embed VERSION
var version string

func main() {
	fmt.Println("=== Embed & Build Constraints ===")

	fmt.Println("\n=== 1. //go:embed Directives ===")
	demoEmbed()

	fmt.Println("\n=== 2. Template Rendering với embed.FS ===")
	demoTemplates()

	fmt.Println("\n=== 3. Build Constraints ===")
	demoBuildConstraints()

	fmt.Println("\n=== 4. Cross Compilation ===")
	demoCrossCompile()

	fmt.Println("\n=== 5. //go:generate ===")
	demoGenerate()
}

func demoEmbed() {
	// embed single string
	fmt.Printf("  Version (//go:embed VERSION): %q\n", version)

	// embed.FS — đọc files
	data, err := staticFiles.ReadFile("static/index.html")
	if err != nil {
		fmt.Printf("  ReadFile error: %v\n", err)
	} else {
		fmt.Printf("  static/index.html: %d bytes\n", len(data))
	}

	// Walk embedded FS
	fmt.Println("  Embedded files:")
	fs.WalkDir(staticFiles, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		data, _ := staticFiles.ReadFile(path)
		fmt.Printf("    %s (%d bytes)\n", path, len(data))
		return nil
	})

	fmt.Println()
	fmt.Println("  embed.FS vs os.Open tradeoffs:")
	fmt.Println("  + Single binary deployment — không cần ship separate assets")
	fmt.Println("  + Files always present — không có missing file at runtime")
	fmt.Println("  - Binary size tăng")
	fmt.Println("  - Không thể update assets mà không rebuild")
	fmt.Println("  - Không phù hợp cho user-generated content")

	fmt.Println()
	fmt.Println("  Embed patterns:")
	fmt.Println(`  //go:embed file.txt               → single file to string/[]byte/embed.FS`)
	fmt.Println(`  //go:embed dir/                   → entire directory`)
	fmt.Println(`  //go:embed dir/*.sql               → glob pattern`)
	fmt.Println(`  //go:embed a.txt b.txt dir/        → multiple patterns`)
	fmt.Println(`  // NOTE: files starting with . or _ are excluded by default`)
	fmt.Println(`  //go:embed all:dir/                → include hidden files too`)
}

func demoTemplates() {
	// Parse templates từ embedded FS
	tmpl, err := template.ParseFS(templateFiles, "templates/*.tmpl")
	if err != nil {
		fmt.Printf("  ParseFS error: %v\n", err)
		return
	}

	data := struct {
		Name    string
		Version string
		Items   []string
	}{
		Name:    "Go App",
		Version: version,
		Items:   []string{"Feature A", "Feature B", "Feature C"},
	}

	fmt.Println("  Rendering welcome.tmpl:")
	if err := tmpl.ExecuteTemplate(os.Stdout, "welcome.tmpl", data); err != nil {
		fmt.Printf("  template error: %v\n", err)
	}
}

func demoBuildConstraints() {
	fmt.Println("  Build constraints syntax (Go 1.17+):")
	fmt.Println()
	fmt.Println("  //go:build linux")
	fmt.Println("  //go:build linux && amd64")
	fmt.Println("  //go:build linux || darwin")
	fmt.Println("  //go:build !windows")
	fmt.Println("  //go:build go1.21        // Go version constraint")
	fmt.Println("  //go:build integration   // custom tag: go test -tags integration")
	fmt.Println()
	fmt.Println("  File naming convention (alternative to build tags):")
	fmt.Println("  - file_linux.go         → only builds on Linux")
	fmt.Println("  - file_windows.go       → only builds on Windows")
	fmt.Println("  - file_linux_amd64.go   → only on Linux/amd64")
	fmt.Println("  - file_test.go          → only in tests")
	fmt.Println()
	fmt.Printf("  Current: GOOS=%s GOARCH=%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Println()

	// platform-specific code đang chạy
	platformInfo()
}

func demoCrossCompile() {
	fmt.Println("  Cross compilation — không cần toolchain trên target platform:")
	fmt.Println()

	targets := []struct{ os, arch string }{
		{"linux", "amd64"},
		{"linux", "arm64"},
		{"darwin", "amd64"},
		{"darwin", "arm64"},
		{"windows", "amd64"},
		{"windows", "arm64"},
	}

	for _, t := range targets {
		fmt.Printf("  GOOS=%-10s GOARCH=%-8s → go build ./...\n", t.os, t.arch)
	}

	fmt.Println()
	fmt.Println("  Build for all platforms:")
	fmt.Println("  for os in linux darwin windows; do")
	fmt.Println("    for arch in amd64 arm64; do")
	fmt.Println("      GOOS=$os GOARCH=$arch go build -o dist/app-$os-$arch .")
	fmt.Println("    done")
	fmt.Println("  done")
	fmt.Println()
	fmt.Println("  Useful build flags:")
	fmt.Println("  -ldflags \"-s -w\"                  → strip debug info, reduce binary size")
	fmt.Println("  -ldflags \"-X main.Version=1.2.3\"  → inject version at build time")
	fmt.Println("  -trimpath                          → remove local paths from binary")
	fmt.Println("  CGO_ENABLED=0                      → static binary (no libc dependency)")
}

func demoGenerate() {
	fmt.Println("  //go:generate directive:")
	fmt.Println()
	fmt.Println("  //go:generate stringer -type=Status")
	fmt.Println("  //go:generate mockgen -source=interface.go -destination=mock.go")
	fmt.Println("  //go:generate protoc --go_out=. proto/user.proto")
	fmt.Println("  //go:generate go run ./cmd/gen/main.go")
	fmt.Println()
	fmt.Println("  Run: go generate ./...")
	fmt.Println("  go generate tìm tất cả //go:generate directives và chạy commands")
	fmt.Println()
	fmt.Println("  NGUYÊN TẮC:")
	fmt.Println("  - Commit generated files vào repo (không generate trong CI)")
	fmt.Println("  - //go:generate nằm trong file .go, thường gần code nó generate cho")
	fmt.Println("  - Dùng go:generate + stringer cho enum String() method")
	fmt.Println()
	fmt.Println("  stringer example:")
	fmt.Println("  type Weekday int")
	fmt.Println("  const (")
	fmt.Println("      Monday Weekday = iota")
	fmt.Println("      Tuesday")
	fmt.Println("      // ...")
	fmt.Println("  )")
	fmt.Println("  // Generated: func (w Weekday) String() string { ... }")
}
