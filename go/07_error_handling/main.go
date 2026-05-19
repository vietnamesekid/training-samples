// Lesson 7: Error Handling — handling errors idiomatically in Go
// Run: go run .
package main

import (
	"errors"
	"fmt"
	"strconv"
)

// === Sentinel Errors — error values that can be compared ===
// Defined at package level, names start with Err

var (
	ErrNotFound     = errors.New("not found")
	ErrUnauthorized = errors.New("unauthorized")
	ErrInvalidInput = errors.New("invalid input")
)

// === Custom Error Types ===

// ValidationError implements the error interface
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: field=%s, message=%s", e.Field, e.Message)
}

// DBError with Unwrap — allows errors.Is/As to traverse the error chain
type DBError struct {
	Code    int
	Message string
	Err     error // wrapped error
}

func (e *DBError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("db error %d: %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("db error %d: %s", e.Code, e.Message)
}

func (e *DBError) Unwrap() error {
	return e.Err // allows errors.Is/As to traverse the chain
}

// HTTPError with multiple wrapped errors (Go 1.20+ errors.Join)
type MultiError struct {
	Errors []error
}

func (m *MultiError) Error() string {
	if len(m.Errors) == 0 {
		return "no errors"
	}
	result := fmt.Sprintf("%d errors occurred:", len(m.Errors))
	for i, e := range m.Errors {
		result += fmt.Sprintf("\n  [%d] %v", i+1, e)
	}
	return result
}

func (m *MultiError) Unwrap() []error {
	return m.Errors
}

// === Functions demonstrating error patterns ===

func findUser(id int) (*struct{ Name string }, error) {
	if id <= 0 {
		return nil, &ValidationError{Field: "id", Message: "must be positive"}
	}
	if id > 100 {
		return nil, fmt.Errorf("findUser: %w", ErrNotFound)
	}
	return &struct{ Name string }{Name: fmt.Sprintf("User%d", id)}, nil
}

func getFromDB(userID int) error {
	_, err := findUser(userID)
	if err != nil {
		// Wrap the error with context — add information to the error chain
		return &DBError{
			Code:    500,
			Message: "query failed",
			Err:     err, // err is wrapped into DBError
		}
	}
	return nil
}

func processRequest(userID int) error {
	err := getFromDB(userID)
	if err != nil {
		// Continue wrapping with context
		return fmt.Errorf("processRequest(userID=%d): %w", userID, err)
	}
	return nil
}

func parseAge(s string) (int, error) {
	age, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("parseAge: invalid age %q: %w", s, err)
	}
	if age < 0 || age > 150 {
		return 0, fmt.Errorf("parseAge: %w: age=%d out of range [0,150]", ErrInvalidInput, age)
	}
	return age, nil
}

func main() {
	fmt.Println("=== 1. errors.New & fmt.Errorf ===")
	err1 := errors.New("simple error")
	err2 := fmt.Errorf("wrapped: %w", err1) // %w creates a wrapped error
	fmt.Printf("err1: %v\n", err1)
	fmt.Printf("err2: %v\n", err2)
	fmt.Printf("errors.Is(err2, err1): %t\n", errors.Is(err2, err1))

	fmt.Println("\n=== 2. Sentinel Errors ===")
	_, err := findUser(999)
	fmt.Printf("findUser(999): %v\n", err)
	fmt.Printf("errors.Is(err, ErrNotFound): %t\n", errors.Is(err, ErrNotFound))

	_, err = findUser(-1)
	fmt.Printf("findUser(-1): %v\n", err)
	var valErr *ValidationError
	fmt.Printf("errors.As(err, *ValidationError): %t\n", errors.As(err, &valErr))
	if valErr != nil {
		fmt.Printf("  Field=%s, Message=%s\n", valErr.Field, valErr.Message)
	}

	fmt.Println("\n=== 3. Error Chain & Unwrap ===")
	err3 := processRequest(999)
	fmt.Printf("processRequest(999):\n  %v\n", err3)

	// errors.Is traverses the entire error chain
	fmt.Printf("errors.Is(err3, ErrNotFound): %t\n", errors.Is(err3, ErrNotFound))

	// errors.As finds a specific type in the chain
	var dbErr *DBError
	if errors.As(err3, &dbErr) {
		fmt.Printf("errors.As DBError: code=%d, msg=%s\n", dbErr.Code, dbErr.Message)
	}

	// ValidationError can also be found through the chain
	var valErr2 *ValidationError
	if errors.As(err3, &valErr2) {
		fmt.Printf("errors.As ValidationError: field=%s\n", valErr2.Field)
	}

	fmt.Println("\n=== 4. Early Return Pattern ===")
	processWithEarlyReturn(1, "25")   // success
	processWithEarlyReturn(-1, "25")  // validation error
	processWithEarlyReturn(1, "abc")  // parse error
	processWithEarlyReturn(1, "999")  // range error

	fmt.Println("\n=== 5. errors.Join (Go 1.20+) ===")
	errs := []error{
		fmt.Errorf("error 1: %w", ErrNotFound),
		fmt.Errorf("error 2: %w", ErrInvalidInput),
		nil, // nil is ignored
	}
	joinErr := errors.Join(errs...)
	fmt.Printf("errors.Join: %v\n", joinErr)
	fmt.Printf("errors.Is(joined, ErrNotFound): %t\n", errors.Is(joinErr, ErrNotFound))
	fmt.Printf("errors.Is(joined, ErrInvalidInput): %t\n", errors.Is(joinErr, ErrInvalidInput))

	fmt.Println("\n=== 6. MultiError (custom) ===")
	multiErr := &MultiError{
		Errors: []error{
			&ValidationError{Field: "name", Message: "required"},
			&ValidationError{Field: "email", Message: "invalid format"},
		},
	}
	fmt.Printf("MultiError: %v\n", multiErr)
	var ve *ValidationError
	fmt.Printf("errors.As first ValidationError: %t\n", errors.As(multiErr, &ve))
	if ve != nil {
		fmt.Printf("  Field=%s\n", ve.Field)
	}
}

func processWithEarlyReturn(userID int, ageStr string) {
	// Early return pattern: handle errors immediately, avoid nesting
	if userID <= 0 {
		fmt.Printf("  processWithEarlyReturn(%d, %q): %v\n",
			userID, ageStr,
			&ValidationError{Field: "userID", Message: "must be positive"})
		return
	}

	age, err := parseAge(ageStr)
	if err != nil {
		fmt.Printf("  processWithEarlyReturn(%d, %q): %v\n", userID, ageStr, err)
		return
	}

	fmt.Printf("  processWithEarlyReturn(%d, %q): OK, age=%d\n", userID, ageStr, age)
}
