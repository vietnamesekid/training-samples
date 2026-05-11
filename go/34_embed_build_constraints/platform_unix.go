//go:build !windows

package main

import "fmt"

func platformInfo() {
	fmt.Println("  Platform: Unix-like (Linux/macOS/etc.)")
	fmt.Println("  Path separator: /")
	fmt.Println("  Line ending: LF (\\n)")
}
