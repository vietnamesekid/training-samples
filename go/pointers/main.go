package main

import "fmt"

func zeroval(ival int) {
	ival = 0
}

func zeroptr(iptr *int) {
	*iptr = 0
}

type Counter struct {
	value int
}

func (c *Counter) Increment() {
	c.value++
}

func (c *Counter) Value() int {
	return c.value
}

func main() {
	i := 1

	fmt.Println("Initial value of i:", i)

	zeroval(i)
	fmt.Println("After zeroval(i), i is still:", i)

	zeroptr(&i)
	fmt.Println("After zeroptr(&i), i is now:", i)

	fmt.Println("Pointer to i:", &i)

	// & get the pointer to i
	// * dereference the pointer to get the value of i

	counter := &Counter{}

	fmt.Println("Initial counter value:", counter.Value())

	counter.Increment()
	fmt.Println("Counter value after incrementing:", counter.Value())

	counter.Increment()
	fmt.Println("Counter value after incrementing again:", counter.Value())
}
