package main

import (
	"fmt"
	"io"
	"strings"
)

// === Interface Best Practices ===

// NGUYÊN TẮC 1: "Accept interfaces, return concrete types"
// Nhận interface → flexible, testable
// Trả về concrete type → caller có đủ thông tin

// BAD: trả về interface khi concrete type đủ rồi
// func NewBuffer() io.Writer { ... }

// GOOD: trả về concrete, nhận interface
type Buffer struct {
	builder strings.Builder
}

func NewBuffer() *Buffer {
	return &Buffer{}
}

func (b *Buffer) Write(p []byte) (n int, err error) {
	return b.builder.Write(p)
}

func (b *Buffer) String() string {
	return b.builder.String()
}

// Hàm nhận io.Writer → hoạt động với bất kỳ implementation nào
func writeGreeting(w io.Writer, name string) {
	fmt.Fprintf(w, "Hello, %s!\n", name)
}

// NGUYÊN TẮC 2: Interface nhỏ là interface tốt (Interface Segregation)
// io.Reader chỉ có 1 method → cực kỳ mạnh vì mọi thứ đều implement được

// BAD: interface quá lớn → khó mock, khó implement
type BigRepository interface {
	FindByID(id int) (*User2, error)
	FindAll() ([]*User2, error)
	Create(u *User2) error
	Update(u *User2) error
	Delete(id int) error
	FindByEmail(email string) (*User2, error)
	Count() (int, error)
}

type User2 struct {
	ID    int
	Name  string
	Email string
}

// GOOD: chia thành interface nhỏ — mỗi use case cần interface riêng
type UserFinder interface {
	FindByID(id int) (*User2, error)
}

type UserCreator interface {
	Create(u *User2) error
}

type UserUpdater interface {
	Update(u *User2) error
}

// NGUYÊN TẮC 3: Định nghĩa interface ở nơi SỬ DỤNG, không phải nơi implement
// → "Define interfaces at use site, not at implement site"
//
// Bad: package "store" define interface, package "service" import "store"
// Good: package "service" define interface nó cần, "store" implement ngầm

// UserService định nghĩa interface nó cần (ở đây, trong service)
type UserRepository interface {
	FindByID(id int) (*User2, error)
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUserName(id int) (string, error) {
	u, err := s.repo.FindByID(id)
	if err != nil {
		return "", err
	}
	return u.Name, nil
}

// In-memory implementation (testing/demo)
type InMemoryUserRepo struct {
	users map[int]*User2
}

func NewInMemoryUserRepo() *InMemoryUserRepo {
	return &InMemoryUserRepo{
		users: map[int]*User2{
			1: {ID: 1, Name: "Alice", Email: "alice@example.com"},
			2: {ID: 2, Name: "Bob", Email: "bob@example.com"},
		},
	}
}

func (r *InMemoryUserRepo) FindByID(id int) (*User2, error) {
	u, ok := r.users[id]
	if !ok {
		return nil, fmt.Errorf("user %d not found", id)
	}
	return u, nil
}

func demoInterfaceBestPractices() {
	fmt.Println("--- Accept interfaces, return concrete types ---")

	buf := NewBuffer()
	writeGreeting(buf, "Gopher") // truyền *Buffer qua io.Writer
	fmt.Printf("  buffer contains: %q\n", buf.String())

	// Cũng hoạt động với strings.Builder
	var sb strings.Builder
	writeGreeting(&sb, "World")
	fmt.Printf("  strings.Builder: %q\n", sb.String())

	fmt.Println("\n--- Interface at use site ---")
	repo := NewInMemoryUserRepo()
	svc := NewUserService(repo)

	name, err := svc.GetUserName(1)
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	} else {
		fmt.Printf("  User 1: %s\n", name)
	}

	_, err = svc.GetUserName(99)
	fmt.Printf("  User 99 error: %v\n", err)

	fmt.Println("\n--- Interface không phải silver bullet ---")
	fmt.Println("  Khi nào KHÔNG cần interface:")
	fmt.Println("  - Chỉ có 1 implementation → dùng concrete type")
	fmt.Println("  - Internal code → không cần abstraction sớm")
	fmt.Println("  - Performance critical → interface có overhead (1 indirect call)")
	fmt.Println("  Khi nào CẦN interface:")
	fmt.Println("  - Cần swap implementation (testing, DI)")
	fmt.Println("  - Nhiều types cần xử lý giống nhau (polymorphism)")
	fmt.Println("  - Tách biệt boundary giữa packages")
}
