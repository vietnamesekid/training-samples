package main

import (
	"fmt"
	"math"
)

type Rect struct {
	width, height float64
}

func (r Rect) Area() float64 {
	return r.width * r.height
}

func (r Rect) Perimeter() float64 {
	return 2 * (r.width + r.height)
}

type Circle struct {
	radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.radius * c.radius
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.radius
}

type Geometry interface {
	Area() float64
	Perimeter() float64
}

type Status int

const (
	Pending Status = iota + 1
	Active
	Inactive
)

func main() {
	rect := Rect{width: 10, height: 5}
	square := Rect{width: 7, height: 7}
	circle := Circle{radius: 5}

	geometries := []Geometry{rect, square, circle}

	for _, g := range geometries {
		fmt.Printf("Area: %.2f, Perimeter: %.2f\n", g.Area(), g.Perimeter())
	}

	// printing status values
	fmt.Println("Status values:")
	fmt.Printf("Pending: %d\n", Pending)
	fmt.Printf("Active: %d\n", Active)
	fmt.Printf("Inactive: %d\n", Inactive)
}
