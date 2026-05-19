package main

import (
	"fmt"
	"time"
)

func demoControlFlow() {
	fmt.Println("\n--- For Loops ---")

	// Method 1: C-style for loop
	for i := 0; i < 3; i++ {
		fmt.Printf("  C-style: %d\n", i)
	}

	// Method 2: while-style (condition only)
	n := 0
	for n < 3 {
		fmt.Printf("  while-style: %d\n", n)
		n++
	}

	// Method 3: range over integer (Go 1.22+)
	for i := range 3 {
		fmt.Printf("  range int: %d\n", i)
	}

	// Method 4: range over slice
	nums := []int{10, 20, 30}
	for i, v := range nums {
		fmt.Printf("  range slice: index=%d, value=%d\n", i, v)
	}

	// Method 5: range with index only
	for i := range nums {
		fmt.Printf("  range index only: %d\n", i)
	}

	// Method 6: range over map (random order)
	m := map[string]int{"a": 1, "b": 2}
	for k, v := range m {
		fmt.Printf("  range map: %s=%d\n", k, v)
	}

	// Method 7: range over string (returns runes)
	for i, r := range "Go🎯" {
		fmt.Printf("  range string: index=%d, rune=%c\n", i, r)
	}

	// Method 8: infinite loop with break
	count := 0
	for {
		if count >= 3 {
			break
		}
		fmt.Printf("  infinite+break: %d\n", count)
		count++
	}

	// continue
	fmt.Println("  continue (skip chẵn):")
	for i := range 6 {
		if i%2 == 0 {
			continue
		}
		fmt.Printf("    %d\n", i)
	}

	// Labeled break/continue — break/continue an outer loop from an inner loop
	fmt.Println("  Labeled break:")
outer:
	for i := range 3 {
		for j := range 3 {
			if i+j >= 3 {
				break outer // break both loops
			}
			fmt.Printf("    i=%d, j=%d\n", i, j)
		}
	}

	fmt.Println("\n--- Switch ---")

	// Expression switch
	day := time.Thursday
	switch day {
	case time.Saturday, time.Sunday:
		fmt.Println("  Weekend!")
	case time.Monday:
		fmt.Println("  Thứ Hai blues...")
	default:
		fmt.Printf("  Ngày thường: %v\n", day)
	}

	// Switch without expression (like an if-else chain)
	score := 85
	switch {
	case score >= 90:
		fmt.Println("  Grade: A")
	case score >= 80:
		fmt.Println("  Grade: B")
	case score >= 70:
		fmt.Println("  Grade: C")
	default:
		fmt.Println("  Grade: F")
	}

	// Switch with initialization statement
	switch x := 42; {
	case x > 100:
		fmt.Println("  > 100")
	case x > 50:
		fmt.Println("  > 50")
	default:
		fmt.Printf("  x = %d\n", x)
	}

	// Type switch
	fmt.Println("  Type switch:")
	values := []any{42, "hello", 3.14, true, []int{1, 2, 3}}
	for _, v := range values {
		switch t := v.(type) {
		case int:
			fmt.Printf("    int: %d\n", t)
		case string:
			fmt.Printf("    string: %q\n", t)
		case float64:
			fmt.Printf("    float64: %.2f\n", t)
		case bool:
			fmt.Printf("    bool: %t\n", t)
		default:
			fmt.Printf("    unknown type: %T\n", t)
		}
	}

	// fallthrough — rarely used, executes the next case
	fmt.Println("  fallthrough:")
	switch 1 {
	case 1:
		fmt.Println("    case 1")
		fallthrough
	case 2:
		fmt.Println("    case 2 (do fallthrough)")
	case 3:
		fmt.Println("    case 3 (không chạy)")
	}

	fmt.Println("\n--- if với initialization ---")
	// if statement can have an initialization before the condition
	if err := riskyOp(); err != nil {
		fmt.Printf("  Error: %v\n", err)
	} else {
		fmt.Println("  OK")
	}

	fmt.Println("\n--- defer ---")
	// defer: delays execution until the surrounding function returns
	// LIFO order: the last defer runs first
	demonstrateDefer()
}

func riskyOp() error {
	return nil // success
}

func demonstrateDefer() {
	fmt.Println("  defer LIFO order:")
	defer fmt.Println("    defer 1 (chạy cuối cùng)")
	defer fmt.Println("    defer 2")
	defer fmt.Println("    defer 3 (chạy đầu tiên)")
	fmt.Println("    main function body")

	// GOTCHA: Arguments to defer are evaluated immediately when defer is called
	x := 10
	defer fmt.Printf("    defer value of x: %d (evaluated tại thời điểm defer)\n", x)
	x = 20
	fmt.Printf("    current x: %d\n", x)

	// Common use case: cleanup resources
	// f, _ := os.Open("file.txt")
	// defer f.Close()  // ← ensures Close is always called regardless of where the function returns

	// GOTCHA: defer in a loop accumulates — avoid using defer inside large loops
	// for rows.Next() {
	//     row := rows.Scan(...)
	//     defer row.Close()  // ← ALL defers run when the function returns, not at the end of the loop!
	// }
}
