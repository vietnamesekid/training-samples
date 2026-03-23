package main

import "fmt"

func plus(a, b int) int {
	return a + b
}

func minus(a, b int) int {
	return a - b
}

func sum(nums ...int) (total int) {
	for _, num := range nums {
		total += num
	}

	return total
}

func intSequence() func() int {
	i := 0

	return func() int {
		i += 1
		return i
	}
}

func factorial(n int) int {
	if n == 0 {
		return 1
	}

	return n * factorial(n-1)
}

func fibonacci(n int) int {
	if n <= 1 {
		return n
	}

	return fibonacci(n-1) + fibonacci(n-2)
}

func main() {
	fmt.Println("3 + 5 =", plus(3, 5))
	fmt.Println("10 - 4 =", minus(10, 4))
	fmt.Println("Sum of 1, 2, 3, 4, 5 =", sum(1, 2, 3, 4, 5))

	nextInt := intSequence()
	fmt.Println("Next integer:", nextInt())
	fmt.Println("Next integer:", nextInt())
	fmt.Println("Next integer:", nextInt())

	fmt.Println("Factorial of 5 =", factorial(5))
	fmt.Println("Fibonacci of 10 =", fibonacci(10))

	kvs := map[string]string{"a": "apple", "b": "banana", "c": "cherry"}
	for k, v := range kvs {
		fmt.Printf("%s - %s\n", k, v)
	}
}
