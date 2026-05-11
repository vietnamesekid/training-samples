// Bài 13: Race Conditions — phát hiện và sửa data race
// Chạy bình thường: go run .
// Chạy với race detector: go run -race .
// Build với race detector: go build -race .
//
// Race detector của Go phát hiện:
//   - Concurrent reads và writes không synchronized
//   - Goroutine leak patterns
//   WARNING: -race làm chậm chương trình ~10x và dùng nhiều memory hơn
//   Chỉ dùng trong development/testing, không phải production
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// === Ví dụ 1: Data Race — đọc/ghi không synchronized ===

func demoDataRace() {
	fmt.Println("\n--- Data Race Demo ---")
	fmt.Println("  CẢNH BÁO: Code này có intentional data race")
	fmt.Println("  Chạy: go run -race . để thấy race detector output")

	// TRẠNG THÁI: counter được access từ nhiều goroutines không có sync
	counter := 0
	var wg sync.WaitGroup

	for range 1000 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter++ // DATA RACE: đọc + ghi không atomic!
		}()
	}
	wg.Wait()
	// Kết quả không phải 1000 vì race condition
	fmt.Printf("  counter (racey): %d (expected 1000, may differ)\n", counter)
}

// === Fix 1: Dùng sync.Mutex ===

func fixWithMutex() {
	fmt.Println("\n--- Fix 1: sync.Mutex ---")
	counter := 0
	var mu sync.Mutex
	var wg sync.WaitGroup

	for range 1000 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			counter++
			mu.Unlock()
		}()
	}
	wg.Wait()
	fmt.Printf("  counter (mutex): %d\n", counter)
}

// === Fix 2: Dùng atomic ===

func fixWithAtomic() {
	fmt.Println("\n--- Fix 2: sync/atomic ---")
	var counter atomic.Int64
	var wg sync.WaitGroup

	for range 1000 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Add(1) // atomic: không cần lock
		}()
	}
	wg.Wait()
	fmt.Printf("  counter (atomic): %d\n", counter.Load())
}

// === Fix 3: Dùng channel (message passing) ===

func fixWithChannel() {
	fmt.Println("\n--- Fix 3: Channel (message passing) ---")
	counterCh := make(chan int, 1000)
	var wg sync.WaitGroup

	for range 1000 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counterCh <- 1
		}()
	}

	wg.Wait()
	close(counterCh)

	total := 0
	for v := range counterCh {
		total += v
	}
	fmt.Printf("  counter (channel): %d\n", total)
}

// === Race Condition vs Data Race ===
// Data Race: concurrent reads/writes không synchronized → undefined behavior
// Race Condition: timing-dependent incorrect behavior, có thể không có data race

// Ví dụ race condition: check-then-act không atomic
type BankAccount struct {
	mu      sync.Mutex
	balance float64
}

func (a *BankAccount) Deposit(amount float64) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.balance += amount
}

// BAD: check-then-act không atomic
func (a *BankAccount) WithdrawBad(amount float64) bool {
	if a.balance >= amount { // check
		// RACE CONDITION: giữa check và act, goroutine khác có thể withdraw!
		a.balance -= amount // act
		return true
	}
	return false
}

// GOOD: check-then-act trong 1 critical section
func (a *BankAccount) WithdrawGood(amount float64) bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.balance >= amount {
		a.balance -= amount
		return true
	}
	return false
}

// === Concurrent Map Panic (khác data race) ===

func demoConcurrentMapPanic() {
	fmt.Println("\n--- Concurrent Map Panic ---")
	fmt.Println("  Đọc/ghi map đồng thời = panic (không chỉ data race!)")
	fmt.Println("  Fix 1: sync.Mutex bao quanh map operations")
	fmt.Println("  Fix 2: sync.Map (built-in concurrent-safe map)")
	fmt.Println("  Fix 3: Chỉ write trong 1 goroutine, đọc sau khi done")

	// GOOD: dùng sync.Mutex với regular map
	type SafeMap struct {
		mu   sync.Mutex
		data map[string]int
	}
	sm := &SafeMap{data: make(map[string]int)}

	var wg sync.WaitGroup
	for i := range 10 {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", n%5)
			sm.mu.Lock()
			sm.data[key]++
			sm.mu.Unlock()
		}(i)
	}
	wg.Wait()
	fmt.Printf("  SafeMap: %v\n", sm.data)
}

func main() {
	fmt.Println("=== RACE CONDITIONS ===")
	fmt.Println("Chạy với: go run -race . để detect races")

	demoDataRace()
	fixWithMutex()
	fixWithAtomic()
	fixWithChannel()
	demoConcurrentMapPanic()

	fmt.Println("\n=== Bank Account: Race Condition ===")
	account := &BankAccount{balance: 100}
	var wg sync.WaitGroup

	// Concurrent withdrawals với good implementation
	success := 0
	var mu sync.Mutex
	for range 5 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if account.WithdrawGood(30) {
				mu.Lock()
				success++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	fmt.Printf("  Successful withdrawals (of 5 attempts of $30 from $100): %d\n", success)

	fmt.Println("\n=== Tóm Tắt ===")
	fmt.Println("  go run -race . : detect data races tại runtime")
	fmt.Println("  go test -race ./... : run tests với race detector")
	fmt.Println("  Fixes: Mutex, atomic, channel (chọn tùy context)")
	fmt.Println("  Mutex: simple, general-purpose")
	fmt.Println("  Atomic: performance-critical, single variable")
	fmt.Println("  Channel: pipeline, fan-out patterns")
}
