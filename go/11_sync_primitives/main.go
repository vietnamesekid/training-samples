// Lesson 11: Sync Primitives — synchronization tools
// Run: go run .
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
