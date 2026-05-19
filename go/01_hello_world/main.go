// Lesson 1: Hello World — the starting point of every Go program
// Run: go run .
// Build: go build -o hello && ./hello
package main

import (
	"fmt"
	"os"
)

func main() {
	// fmt.Println automatically adds a newline, use it when printing multiple arguments
	fmt.Println("Hello, World!")

	// fmt.Printf uses a format string — does not add a newline automatically
	name := "Gopher"
	fmt.Printf("Xin chào, %s!\n", name)

	// fmt.Sprintf returns a string instead of printing it
	msg := fmt.Sprintf("Go version: %s, Platform: %s/%s",
		"1.26", "darwin", "arm64")
	fmt.Println(msg)

	// os.Args holds arguments from the command line
	// os.Args[0] is the program name, os.Args[1:] are the arguments
	fmt.Println("\n=== Command-line arguments ===")
	fmt.Printf("Program: %s\n", os.Args[0])
	if len(os.Args) > 1 {
		fmt.Printf("Arguments: %v\n", os.Args[1:])
	} else {
		fmt.Println("Không có argument. Thử: go run . Alice Bob")
	}

	// os.Stdout, os.Stderr — standard output streams
	fmt.Fprintln(os.Stdout, "\nĐây là stdout")
	fmt.Fprintln(os.Stderr, "Đây là stderr (thường dùng cho lỗi)")

	// os.Exit — exits the program with an exit code
	// os.Exit(0) = success, os.Exit(1) = error
	// WARNING: defer does NOT run when os.Exit is called
	fmt.Println("\nChương trình kết thúc bình thường (exit code 0)")
}
