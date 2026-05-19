package main

import (
	"fmt"
	"io"
	"strings"
)

// === Interface Best Practices ===

// PRINCIPLE 1: "Accept interfaces, return concrete types"
// Accept interfaces → flexible and testable
// Return concrete types → caller has full information

// BAD: returning an interface when a concrete type is sufficient
// func NewBuffer() io.Writer { ... }

// GOOD: return concrete, accept interface
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

// Function accepts io.Writer → works with any implementation
func writeGreeting(w io.Writer, name string) {
	fmt.Fprintf(w, "Hello, %s!\n", name)
}

// PRINCIPLE 2: Small interfaces are good interfaces (Interface Segregation)
// io.Reader has only 1 method → extremely powerful because everything can implement it

// BAD: interface too large → hard to mock, hard to implement
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

// GOOD: split into small interfaces — each use case gets its own interface
type UserFinder interface {
	FindByID(id int) (*User2, error)
}

type UserCreator interface {
	Create(u *User2) error
}

type UserUpdater interface {
	Update(u *User2) error
}

// PRINCIPLE 3: Define interfaces at the USE SITE, not the implementation site
// → "Define interfaces at use site, not at implement site"
//
// Bad: package "store" defines the interface, package "service" imports "store"
// Good: package "service" defines the interface it needs, "store" implements it implicitly

// UserService defines the interface it needs (here, within the service)
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
	writeGreeting(buf, "Gopher") // pass *Buffer via io.Writer
	fmt.Printf("  buffer contains: %q\n", buf.String())

	// Also works with strings.Builder
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
