package main

import "fmt"

// === Nil Interface Gotcha — LỖI PHỔ BIẾN NHẤT VỚI INTERFACE ===
//
// Interface value thực ra là một pair: (type, value)
// nil interface: (nil, nil) — cả type lẫn value đều nil
// non-nil interface chứa nil pointer: (*MyError, nil) — type không nil!

type MyError struct {
	Message string
}

func (e *MyError) Error() string {
	return e.Message
}

// BAD: function này có bug tinh tế
func badGetError(fail bool) error {
	var err *MyError // err = nil pointer
	if fail {
		err = &MyError{Message: "something failed"}
	}
	return err // NGUY HIỂM: trả về (*MyError, nil) — interface KHÔNG nil!
}

// GOOD: trả về nil interface thật sự
func goodGetError(fail bool) error {
	if fail {
		return &MyError{Message: "something failed"}
	}
	return nil // trả về (nil, nil) — interface thật sự nil
}

func demoNilInterfaceGotcha() {
	fmt.Println("--- BAD: (*MyError, nil) != nil interface ---")
	err1 := badGetError(false)
	fmt.Printf("  badGetError(false) == nil: %t\n", err1 == nil)
	// ← IN FALSE! Vì err1 = (*MyError, nil), không phải (nil, nil)
	fmt.Printf("  Actual value: (%T, %v)\n", err1, err1)

	fmt.Println("\n--- GOOD: nil interface thật sự ---")
	err2 := goodGetError(false)
	fmt.Printf("  goodGetError(false) == nil: %t\n", err2 == nil)
	// ← IN TRUE! Vì err2 = (nil, nil)

	err3 := goodGetError(true)
	fmt.Printf("  goodGetError(true) == nil: %t\n", err3 == nil)
	fmt.Printf("  Error message: %v\n", err3)

	fmt.Println("\n--- Interface internals ---")
	fmt.Println("  Interface = (type, value) pair")
	fmt.Println("  nil interface: both type and value are nil → == nil is true")
	fmt.Println("  (*MyError)(nil): type=*MyError, value=nil → == nil is FALSE!")
	fmt.Println()
	fmt.Println("  NGUYÊN TẮC: Không bao giờ return typed nil từ function trả về interface")
}
