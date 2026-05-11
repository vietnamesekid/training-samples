package main

import "iter"

type Student struct {
	Name string
	GPA  float64
}

func ScholarshipStudents(students []Student) iter.Seq[Student] {
	return func(yield func(Student) bool) {
		for _, student := range students {
			if student.GPA > 8.0 {
				if !yield(student) {
					return
				}
			}
		}
	}
}

type SliceSeq[T any] []T

func (s SliceSeq[T]) Map(f func(T) T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, item := range s {
			if !yield(f(item)) {
				return
			}
		}
	}
}

func main() {
	students := []Student{
		{Name: "Alice", GPA: 9.0},
		{Name: "Bob", GPA: 7.5},
		{Name: "Charlie", GPA: 8.5},
	}

	for s := range ScholarshipStudents(students) {
		println(s.Name)
	}

	slice := SliceSeq[int]{1, 2, 3, 4, 5}
	doubled := slice.Map(func(x int) int { return x * 2 })

	for x := range doubled {
		println(x)
	}
}
