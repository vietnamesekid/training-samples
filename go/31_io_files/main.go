// Bài 31: I/O & File Operations
// File I/O, bufio, io.Copy, embed.FS, os.Root (Go 1.24+), io.Reader/Writer composition
// Chạy: go run .
package main

import (
	"bufio"
	"bytes"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

//go:embed assets/*
var embeddedAssets embed.FS

func main() {
	fmt.Println("=== I/O & File Operations ===")

	fmt.Println("\n=== 1. io.Reader & io.Writer Interfaces ===")
	demoReaderWriter()

	fmt.Println("\n=== 2. File Operations ===")
	demoFileOps()

	fmt.Println("\n=== 3. bufio — Buffered I/O ===")
	demoBufio()

	fmt.Println("\n=== 4. io.Copy & io.TeeReader ===")
	demoIOCopy()

	fmt.Println("\n=== 5. embed.FS (Go 1.16+) ===")
	demoEmbedFS()

	fmt.Println("\n=== 6. Temp Files & Directories ===")
	demoTempFiles()
}

// ============================================================
// 1. io.Reader & io.Writer
// ============================================================

// CountingReader đếm số bytes đã đọc
type CountingReader struct {
	r     io.Reader
	count int64
}

func (cr *CountingReader) Read(p []byte) (int, error) {
	n, err := cr.r.Read(p)
	cr.count += int64(n)
	return n, err
}

// UpperCaseWriter chuyển tất cả thành uppercase khi write
type UpperCaseWriter struct {
	w io.Writer
}

func (uw *UpperCaseWriter) Write(p []byte) (int, error) {
	upper := bytes.ToUpper(p)
	return uw.w.Write(upper)
}

func demoReaderWriter() {
	// io.Reader composition
	original := strings.NewReader("Hello, Go I/O!")
	cr := &CountingReader{r: original}

	data, _ := io.ReadAll(cr)
	fmt.Printf("  Read %d bytes: %s\n", cr.count, data)

	// io.Writer composition
	var buf bytes.Buffer
	uw := &UpperCaseWriter{w: &buf}
	io.WriteString(uw, "hello world")
	fmt.Printf("  UpperCaseWriter: %s\n", buf.String())

	// io.MultiWriter — write to multiple destinations
	var buf1, buf2 bytes.Buffer
	mw := io.MultiWriter(&buf1, &buf2)
	fmt.Fprintln(mw, "broadcast to multiple writers")
	fmt.Printf("  MultiWriter buf1: %s", buf1.String())
	fmt.Printf("  MultiWriter buf2: %s", buf2.String())

	// io.LimitReader — đọc tối đa N bytes
	src := strings.NewReader("This is a long string")
	limited := io.LimitReader(src, 7)
	data, _ = io.ReadAll(limited)
	fmt.Printf("  LimitReader(7): %q\n", data)

	// io.Pipe — synchronous in-memory pipe
	pr, pw := io.Pipe()
	go func() {
		fmt.Fprintln(pw, "piped data")
		pw.Close()
	}()
	data, _ = io.ReadAll(pr)
	fmt.Printf("  io.Pipe: %s", data)
}

// ============================================================
// 2. File Operations
// ============================================================

func demoFileOps() {
	tmpDir, err := os.MkdirTemp("", "go-io-demo-*")
	if err != nil {
		fmt.Printf("  MkdirTemp error: %v\n", err)
		return
	}
	defer os.RemoveAll(tmpDir)

	filePath := filepath.Join(tmpDir, "example.txt")

	// Write file
	content := "Line 1\nLine 2\nLine 3\n"
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		fmt.Printf("  WriteFile error: %v\n", err)
		return
	}
	fmt.Printf("  Written: %s (%d bytes)\n", filepath.Base(filePath), len(content))

	// Read entire file
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("  ReadFile error: %v\n", err)
		return
	}
	fmt.Printf("  ReadFile: %d bytes\n", len(data))

	// Append to file
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("  OpenFile error: %v\n", err)
		return
	}
	fmt.Fprintln(f, "Line 4 (appended)")
	f.Close()

	// File info
	info, _ := os.Stat(filePath)
	fmt.Printf("  File: name=%s size=%d mod=%s\n",
		info.Name(), info.Size(), info.ModTime().Format(time.RFC3339))

	// Directory listing
	subDir := filepath.Join(tmpDir, "subdir")
	os.Mkdir(subDir, 0755)
	os.WriteFile(filepath.Join(subDir, "a.txt"), []byte("a"), 0644)
	os.WriteFile(filepath.Join(subDir, "b.txt"), []byte("b"), 0644)

	fmt.Printf("  WalkDir:\n")
	fs.WalkDir(os.DirFS(tmpDir), ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			fmt.Printf("    DIR  %s/\n", path)
		} else {
			fmt.Printf("    FILE %s\n", path)
		}
		return nil
	})
}

// ============================================================
// 3. bufio — Buffered I/O
// ============================================================

func demoBufio() {
	// bufio.Scanner — đọc line-by-line
	text := "first line\nsecond line\nthird line\n"
	scanner := bufio.NewScanner(strings.NewReader(text))
	fmt.Println("  Scanner line-by-line:")
	for scanner.Scan() {
		fmt.Printf("    %q\n", scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("  scanner error: %v\n", err)
	}

	// bufio.Reader — đọc với buffer (tránh syscall mỗi byte)
	reader := bufio.NewReader(strings.NewReader("Hello\nWorld"))
	line, _ := reader.ReadString('\n')
	fmt.Printf("  ReadString: %q\n", line)
	rest, _ := io.ReadAll(reader)
	fmt.Printf("  Rest: %q\n", rest)

	// bufio.Writer — batch writes, giảm syscalls
	var buf bytes.Buffer
	bw := bufio.NewWriter(&buf)
	for i := range 5 {
		fmt.Fprintf(bw, "item %d\n", i)
	}
	bw.Flush() // QUAN TRỌNG: flush buffer vào underlying writer
	fmt.Printf("  BufferedWriter: %d bytes flushed\n", buf.Len())

	// Word scanner
	wordText := "the quick brown fox"
	ws := bufio.NewScanner(strings.NewReader(wordText))
	ws.Split(bufio.ScanWords) // scan by word instead of line
	var words []string
	for ws.Scan() {
		words = append(words, ws.Text())
	}
	fmt.Printf("  ScanWords: %v\n", words)
}

// ============================================================
// 4. io.Copy & io.TeeReader
// ============================================================

func demoIOCopy() {
	// io.Copy — efficient copy, uses internal 32KB buffer
	src := strings.NewReader("Source data for copying")
	var dst bytes.Buffer

	n, err := io.Copy(&dst, src)
	if err != nil {
		fmt.Printf("  Copy error: %v\n", err)
		return
	}
	fmt.Printf("  io.Copy: copied %d bytes: %s\n", n, dst.String())

	// io.TeeReader — đọc VÀ ghi sang writer khác đồng thời
	// Dùng để: hash + save, log + process, etc.
	var logger bytes.Buffer
	tee := io.TeeReader(strings.NewReader("Data being processed"), &logger)

	processed, _ := io.ReadAll(tee)
	fmt.Printf("  TeeReader processed: %s\n", processed)
	fmt.Printf("  TeeReader logged: %s\n", logger.String())

	// io.CopyN — copy tối đa N bytes
	src2 := strings.NewReader("Only first 5 chars")
	var dst2 bytes.Buffer
	io.CopyN(&dst2, src2, 5)
	fmt.Printf("  CopyN(5): %q\n", dst2.String())
}

// ============================================================
// 5. embed.FS
// ============================================================

func demoEmbedFS() {
	// Đọc file từ embedded FS
	data, err := embeddedAssets.ReadFile("assets/hello.txt")
	if err != nil {
		fmt.Printf("  ReadFile error: %v\n", err)
	} else {
		fmt.Printf("  Embedded file: %s", data)
	}

	// List embedded directory
	entries, _ := embeddedAssets.ReadDir("assets")
	fmt.Printf("  Embedded assets (%d files):\n", len(entries))
	for _, e := range entries {
		fmt.Printf("    %s\n", e.Name())
	}

	// Dùng embed.FS với fs.WalkDir
	fs.WalkDir(embeddedAssets, "assets", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		data, _ := embeddedAssets.ReadFile(path)
		fmt.Printf("  %s: %d bytes\n", path, len(data))
		return nil
	})

	fmt.Println()
	fmt.Println("  embed.FS use cases:")
	fmt.Println("  - Serve static files (HTML, CSS, JS) trong single binary")
	fmt.Println("  - Embed SQL migration files")
	fmt.Println("  - Embed config templates")
	fmt.Println("  - Embed certificates/keys")
}

// ============================================================
// 6. Temp Files & Directories
// ============================================================

func demoTempFiles() {
	// os.CreateTemp — temp file (auto unique name)
	tmpFile, err := os.CreateTemp("", "go-demo-*.txt")
	if err != nil {
		fmt.Printf("  CreateTemp error: %v\n", err)
		return
	}
	defer os.Remove(tmpFile.Name()) // cleanup
	defer tmpFile.Close()

	fmt.Fprintln(tmpFile, "Temporary content")
	tmpFile.Close()
	fmt.Printf("  Temp file: %s\n", filepath.Base(tmpFile.Name()))

	// Đọc lại
	data, _ := os.ReadFile(tmpFile.Name())
	fmt.Printf("  Content: %s", data)

	// os.MkdirTemp — temp directory
	tmpDir, err := os.MkdirTemp("", "go-workdir-*")
	if err != nil {
		fmt.Printf("  MkdirTemp error: %v\n", err)
		return
	}
	defer os.RemoveAll(tmpDir)
	fmt.Printf("  Temp dir: %s\n", filepath.Base(tmpDir))

	// Tạo files trong temp dir
	os.WriteFile(filepath.Join(tmpDir, "work.txt"), []byte("working"), 0644)
	entries, _ := os.ReadDir(tmpDir)
	fmt.Printf("  Files in temp dir: %d\n", len(entries))

	// NGUYÊN TẮC: luôn defer os.Remove/RemoveAll ngay sau tạo temp file/dir
	// Dùng t.TempDir() trong tests — tự cleanup sau test
}
