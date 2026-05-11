package main

import (
	"fmt"
	"iter"
)

// === Generic Types ===

// Stack[T] — type-safe stack
type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(v T) {
	s.items = append(s.items, v)
}

func (s *Stack[T]) Pop() (T, bool) {
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}
	top := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return top, true
}

func (s *Stack[T]) Peek() (T, bool) {
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}
	return s.items[len(s.items)-1], true
}

func (s *Stack[T]) Len() int { return len(s.items) }

func (s *Stack[T]) IsEmpty() bool { return len(s.items) == 0 }

// All — iterate qua stack (dùng iter.Seq, Go 1.23+)
func (s *Stack[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, v := range s.items {
			if !yield(v) {
				return
			}
		}
	}
}

// Optional[T] — safe nullable value (thay thế cho *T trong một số cases)
type Optional[T any] struct {
	value *T
}

func Some[T any](v T) Optional[T] {
	return Optional[T]{value: &v}
}

func None[T any]() Optional[T] {
	return Optional[T]{}
}

func (o Optional[T]) IsPresent() bool {
	return o.value != nil
}

func (o Optional[T]) Get() (T, bool) {
	if o.value == nil {
		var zero T
		return zero, false
	}
	return *o.value, true
}

func (o Optional[T]) OrElse(defaultVal T) T {
	if o.value == nil {
		return defaultVal
	}
	return *o.value
}

// Pair[A, B] — generic tuple
type Pair[A, B any] struct {
	First  A
	Second B
}

func MakePair[A, B any](a A, b B) Pair[A, B] {
	return Pair[A, B]{First: a, Second: b}
}

// Zip: kết hợp 2 slices thành slice of pairs
func Zip[A, B any](as []A, bs []B) []Pair[A, B] {
	n := len(as)
	if len(bs) < n {
		n = len(bs)
	}
	result := make([]Pair[A, B], n)
	for i := range n {
		result[i] = MakePair(as[i], bs[i])
	}
	return result
}

// Result[T] — error handling type (như Rust's Result<T, E>)
type Result[T any] struct {
	value T
	err   error
}

func Ok[T any](v T) Result[T]        { return Result[T]{value: v} }
func Err[T any](err error) Result[T] { return Result[T]{err: err} }

func (r Result[T]) IsOk() bool     { return r.err == nil }
func (r Result[T]) IsErr() bool    { return r.err != nil }
func (r Result[T]) Unwrap() T      { return r.value }
func (r Result[T]) Error() error   { return r.err }

func demoGenericTypes() {
	fmt.Println("\n--- Stack[T] ---")
	s := &Stack[int]{}
	s.Push(1)
	s.Push(2)
	s.Push(3)
	fmt.Printf("  Len: %d\n", s.Len())
	for v := range s.All() {
		fmt.Printf("  item: %d\n", v)
	}
	if v, ok := s.Pop(); ok {
		fmt.Printf("  Pop: %d\n", v)
	}
	if v, ok := s.Peek(); ok {
		fmt.Printf("  Peek: %d\n", v)
	}

	// String stack
	ss := &Stack[string]{}
	ss.Push("go")
	ss.Push("is")
	ss.Push("awesome")
	fmt.Printf("  String stack len: %d\n", ss.Len())

	fmt.Println("\n--- Optional[T] ---")
	name := Some("Alice")
	empty := None[string]()

	if v, ok := name.Get(); ok {
		fmt.Printf("  Some: %s\n", v)
	}
	fmt.Printf("  None.IsPresent: %t\n", empty.IsPresent())
	fmt.Printf("  None.OrElse: %s\n", empty.OrElse("default"))

	fmt.Println("\n--- Pair[A, B] & Zip ---")
	names := []string{"Alice", "Bob", "Carol"}
	ages := []int{30, 25, 35}
	pairs := Zip(names, ages)
	for _, p := range pairs {
		fmt.Printf("  %s: %d\n", p.First, p.Second)
	}

	fmt.Println("\n--- Result[T] ---")
	r1 := Ok(42)
	r2 := Err[int](fmt.Errorf("something failed"))
	fmt.Printf("  Ok: isOk=%t, value=%d\n", r1.IsOk(), r1.Unwrap())
	fmt.Printf("  Err: isOk=%t, err=%v\n", r2.IsOk(), r2.Error())
}
