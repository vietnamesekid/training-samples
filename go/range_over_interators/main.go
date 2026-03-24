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

func main() {
	students := []Student{
		{Name: "Alice", GPA: 9.0},
		{Name: "Bob", GPA: 7.5},
		{Name: "Charlie", GPA: 8.5},
	}

	for s := range ScholarshipStudents(students) {
		println(s.Name)
	}
}
