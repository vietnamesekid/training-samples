// Package validator — public package, có thể import từ bên ngoài
// Đây là ví dụ về pkg/ — reusable, public-facing package
package validator

import (
	"fmt"
	"strings"
)

// Validator thực hiện validation
type Validator struct {
	rules []rule
}

type rule struct {
	name  string
	check func(string) bool
}

// New tạo Validator mới
func New() *Validator {
	return &Validator{}
}

// ValidateUser validate thông tin user, trả về danh sách lỗi
func (v *Validator) ValidateUser(name, email string, age int) []string {
	var errs []string

	if strings.TrimSpace(name) == "" {
		errs = append(errs, "name: required")
	} else if len(name) < 2 {
		errs = append(errs, "name: minimum 2 characters")
	}

	if email == "" {
		errs = append(errs, "email: required")
	} else if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		errs = append(errs, fmt.Sprintf("email: %q is invalid format", email))
	}

	if age < 0 {
		errs = append(errs, fmt.Sprintf("age: %d must be non-negative", age))
	} else if age > 150 {
		errs = append(errs, fmt.Sprintf("age: %d exceeds maximum 150", age))
	}

	if len(errs) == 0 {
		return []string{fmt.Sprintf("✓ Valid: name=%s, email=%s, age=%d", name, email, age)}
	}
	return errs
}
