package main

import "fmt"

// === Nil Interface Gotcha — THE MOST COMMON MISTAKE WITH INTERFACES ===
//
// An interface value is actually a pair: (type, value)
// nil interface: (nil, nil) — both type and value are nil
// non-nil interface holding a nil pointer: (*MyError, nil) — type is NOT nil!

type MyError struct {
	Message string
}

func (e *MyError) Error() string {
	return e.Message
}

// BAD: this function has a subtle bug
func badGetError(fail bool) error {
	var err *MyError // err = nil pointer
	if fail {
		err = &MyError{Message: "something failed"}
	}
	return err // DANGER: returns (*MyError, nil) — interface is NOT nil!
}

// GOOD: returns a truly nil interface
func goodGetError(fail bool) error {
	if fail {
		return &MyError{Message: "something failed"}
	}
	return nil // returns (nil, nil) — truly nil interface
}

func demoNilInterfaceGotcha() {
	fmt.Println("--- BAD: (*MyError, nil) != nil interface ---")
	err1 := badGetError(false)
	fmt.Printf("  badGetError(false) == nil: %t\n", err1 == nil)
	// ← PRINTS FALSE! Because err1 = (*MyError, nil), not (nil, nil)
	fmt.Printf("  Actual value: (%T, %v)\n", err1, err1)

	fmt.Println("\n--- GOOD: nil interface thật sự ---")
	err2 := goodGetError(false)
	fmt.Printf("  goodGetError(false) == nil: %t\n", err2 == nil)
	// ← PRINTS TRUE! Because err2 = (nil, nil)

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
