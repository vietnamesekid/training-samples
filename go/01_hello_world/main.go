// Bài 1: Hello World — điểm xuất phát của mọi chương trình Go
// Chạy: go run .
// Build: go build -o hello && ./hello
package main

import (
	"fmt"
	"os"
)

func main() {
	// fmt.Println tự thêm newline, dùng khi in nhiều arguments
	fmt.Println("Hello, World!")

	// fmt.Printf dùng format string — không tự thêm newline
	name := "Gopher"
	fmt.Printf("Xin chào, %s!\n", name)

	// fmt.Sprintf trả về string thay vì in ra
	msg := fmt.Sprintf("Go version: %s, Platform: %s/%s",
		"1.26", "darwin", "arm64")
	fmt.Println(msg)

	// os.Args chứa các argument từ command line
	// os.Args[0] là tên chương trình, os.Args[1:] là các argument
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

	// os.Exit — thoát chương trình với exit code
	// os.Exit(0) = thành công, os.Exit(1) = lỗi
	// CẢNH BÁO: defer KHÔNG chạy khi gọi os.Exit
	fmt.Println("\nChương trình kết thúc bình thường (exit code 0)")
}
