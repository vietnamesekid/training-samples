// Lesson 19: Reflection — inspect and manipulate types at runtime
// Run: go run .
// WARNING: Reflection is ~10-100x slower than regular code
// Only use when truly necessary: frameworks, serialization, generic code
package main

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// === Struct used for demo ===

type Address struct {
	Street string `json:"street" validate:"required"`
	City   string `json:"city" validate:"required"`
	Zip    string `json:"zip" validate:"len=5"`
}

type User struct {
	ID      int     `json:"id"`
	Name    string  `json:"name" validate:"required,min=2"`
	Email   string  `json:"email" validate:"required,email"`
	Age     int     `json:"age" validate:"min=0,max=150"`
	Score   float64 `json:"score,omitempty"`
	Address Address `json:"address"`
}

// === Mini Struct Validator using Reflection ===

type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("field %s: %s", e.Field, e.Message)
}

// Validate checks struct fields according to `validate` tags
func Validate(v any) []error {
	val := reflect.ValueOf(v)
	typ := reflect.TypeOf(v)

	// Dereference pointer if needed
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}

	if val.Kind() != reflect.Struct {
		return []error{fmt.Errorf("validate: expected struct, got %s", val.Kind())}
	}

	var errs []error
	for i := range typ.NumField() {
		field := typ.Field(i)
		value := val.Field(i)

		validateTag := field.Tag.Get("validate")
		if validateTag == "" {
			continue
		}

		rules := strings.Split(validateTag, ",")
		for _, rule := range rules {
			if err := applyRule(field.Name, value, rule); err != nil {
				errs = append(errs, err)
			}
		}
	}

	return errs
}

func applyRule(fieldName string, v reflect.Value, rule string) error {
	parts := strings.SplitN(rule, "=", 2)
	ruleName := parts[0]

	switch ruleName {
	case "required":
		if v.IsZero() {
			return ValidationError{Field: fieldName, Message: "is required"}
		}
	case "min":
		if len(parts) < 2 {
			return nil
		}
		minVal, _ := strconv.ParseInt(parts[1], 10, 64)
		switch v.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if v.Int() < minVal {
				return ValidationError{Field: fieldName, Message: fmt.Sprintf("must be >= %d", minVal)}
			}
		case reflect.String:
			if int64(len(v.String())) < minVal {
				return ValidationError{Field: fieldName, Message: fmt.Sprintf("length must be >= %d", minVal)}
			}
		}
	case "max":
		if len(parts) < 2 {
			return nil
		}
		maxVal, _ := strconv.ParseInt(parts[1], 10, 64)
		if v.Int() > maxVal {
			return ValidationError{Field: fieldName, Message: fmt.Sprintf("must be <= %d", maxVal)}
		}
	case "email":
		s := v.String()
		if !strings.Contains(s, "@") {
			return ValidationError{Field: fieldName, Message: "must be valid email"}
		}
	}
	return nil
}

// === Struct to Map (using reflection) ===

func StructToMap(v any) map[string]any {
	result := make(map[string]any)
	val := reflect.ValueOf(v)
	typ := reflect.TypeOf(v)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}

	for i := range typ.NumField() {
		field := typ.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" {
			jsonTag = field.Name
		}
		// Strip options like omitempty
		name := strings.Split(jsonTag, ",")[0]
		if name == "-" {
			continue
		}
		result[name] = val.Field(i).Interface()
	}

	return result
}

func main() {
	fmt.Println("=== 1. reflect.TypeOf & reflect.ValueOf ===")

	values := []any{42, "hello", 3.14, true, []int{1, 2, 3}, User{ID: 1, Name: "Alice"}}
	for _, v := range values {
		t := reflect.TypeOf(v)
		rv := reflect.ValueOf(v)
		fmt.Printf("  TypeOf: %-20s Kind: %-10s Value: %v\n", t, t.Kind(), rv)
	}

	fmt.Println("\n=== 2. Modify value through pointer ===")
	x := 42
	rv := reflect.ValueOf(&x).Elem() // Elem() to deref pointer
	fmt.Printf("  Before: %d\n", x)
	rv.SetInt(100) // can set because we have a pointer
	fmt.Printf("  After SetInt(100): %d\n", x)

	fmt.Println("\n=== 3. Inspect Struct Fields & Tags ===")
	u := User{ID: 1, Name: "Alice", Email: "alice@example.com", Age: 30, Score: 9.5}
	t := reflect.TypeOf(u)
	v := reflect.ValueOf(u)

	fmt.Printf("  Struct: %s (%d fields)\n", t.Name(), t.NumField())
	for i := range t.NumField() {
		field := t.Field(i)
		value := v.Field(i)
		fmt.Printf("  [%d] %-10s %-15s json=%-20q validate=%q\n",
			i, field.Name, field.Type,
			field.Tag.Get("json"),
			field.Tag.Get("validate"),
		)
		_ = value
	}

	fmt.Println("\n=== 4. Mini Struct Validator ===")
	validUser := User{
		ID:    1,
		Name:  "Alice",
		Email: "alice@example.com",
		Age:   30,
	}
	invalidUser := User{
		ID:    0,
		Name:  "A", // too short
		Email: "not-an-email",
		Age:   -5,
	}

	fmt.Println("  Valid user:")
	if errs := Validate(validUser); len(errs) == 0 {
		fmt.Println("    ✓ No validation errors")
	} else {
		for _, e := range errs {
			fmt.Printf("    ✗ %v\n", e)
		}
	}

	fmt.Println("  Invalid user:")
	for _, e := range Validate(invalidUser) {
		fmt.Printf("    ✗ %v\n", e)
	}

	fmt.Println("\n=== 5. StructToMap ===")
	m := StructToMap(u)
	for k, v := range m {
		fmt.Printf("  %s: %v\n", k, v)
	}

	fmt.Println("\n=== 6. reflect.DeepEqual ===")
	a1 := []int{1, 2, 3}
	a2 := []int{1, 2, 3}
	a3 := []int{1, 2, 4}
	fmt.Printf("  DeepEqual([1,2,3], [1,2,3]): %t\n", reflect.DeepEqual(a1, a2))
	fmt.Printf("  DeepEqual([1,2,3], [1,2,4]): %t\n", reflect.DeepEqual(a1, a3))

	fmt.Println("\n=== 7. Performance Caveat ===")
	fmt.Println("  Reflection overhead ~ 10-100x so với direct access")
	fmt.Println("  Dùng reflection khi: frameworks, serialization, testing")
	fmt.Println("  Không dùng reflection khi: hot paths, performance critical")
	fmt.Println("  Alternative: code generation (go generate + text/template)")
}
