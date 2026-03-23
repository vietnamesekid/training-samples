package main

import "fmt"

func main() {
	for i := 0; i < 5; i++ {
		fmt.Println(i)
	}

	fmt.Println("=====================")

	for j := 0; j < 5; {
		fmt.Println(j)
		j++
	}

	fmt.Println("=====================")

	whoAmI := func(t interface{}) {
		switch t.(type) {
		case int:
			fmt.Println("I am an int")
		case string:
			fmt.Println("I am a string")
		case bool:
			fmt.Println("I am a bool")
		default:
			fmt.Println("I am something else")
		}
	}

	whoAmI(42)
	whoAmI("hello")
	whoAmI(true)
	whoAmI(3.14)

}
