// Bài 11: Sync Primitives — công cụ đồng bộ hóa
// Chạy: go run .
package main

import "fmt"

func main() {
	fmt.Println("=== MUTEX & RWMUTEX ===")
	demoMutex()

	fmt.Println("\n=== ONCE & POOL ===")
	demoOncePool()

	fmt.Println("\n=== SYNC.MAP & ATOMIC ===")
	demoMapAtomic()

	fmt.Println("\n=== SYNC.COND ===")
	demoCond()
}
