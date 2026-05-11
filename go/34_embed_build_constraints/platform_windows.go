//go:build windows

package main

import "fmt"

func platformInfo() {
	fmt.Println("  Platform: Windows")
	fmt.Println("  Path separator: \\")
	fmt.Println("  Line ending: CRLF (\\r\\n)")
}
