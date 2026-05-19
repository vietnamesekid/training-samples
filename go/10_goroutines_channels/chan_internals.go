//go:build ignore

// === Channel Internals — simulating hchan (Go runtime internal) ===
//
// This file is "documentation-as-code" — simulating how the Go runtime
// implements channels internally. Run: go run chan_internals.go
//
// Actual structure in runtime/chan.go:
//   type hchan struct {
//       qcount   uint           // number of elements in the buffer
//       dataqsiz uint           // buffer capacity
//       buf      unsafe.Pointer // buffer (circular array)
//       elemsize uint16         // size of each element
//       closed   uint32
//       sendx    uint           // next send index
//       recvx    uint           // next receive index
//       recvq    waitq          // queue of goroutines waiting to receive
//       sendq    waitq          // queue of goroutines waiting to send
//       lock     mutex
//   }

package main

import (
	"fmt"
	"sync"
)

// sudog represents a Goroutine that is blocked (waiting to send or receive)
type sudog struct {
	g    *goroutine
	elem any
	next *sudog
	prev *sudog
}

// goroutine simulates a simple Goroutine
type goroutine struct {
	id     int
	result any
	done   chan struct{}
}

// waitq is a doubly linked list of sudogs
type waitq struct {
	first *sudog
	last  *sudog
}

func (q *waitq) enqueue(s *sudog) {
	s.next = nil
	if q.last == nil {
		q.first = s
		q.last = s
		return
	}
	s.prev = q.last
	q.last.next = s
	q.last = s
}

func (q *waitq) dequeue() *sudog {
	s := q.first
	if s == nil {
		return nil
	}
	q.first = s.next
	if q.first == nil {
		q.last = nil
	} else {
		q.first.prev = nil
	}
	return s
}

// hchan simulates the internal channel structure in the Go runtime
type hchan struct {
	buf      []any
	dataqsiz uint
	qcount   uint
	sendx    uint
	recvx    uint
	recvq    waitq
	sendq    waitq
	closed   bool
	lock     sync.Mutex
}

func makeHchan(size uint) *hchan {
	return &hchan{
		buf:      make([]any, size),
		dataqsiz: size,
	}
}

// send simulates ch <- val
//
// 3 cases:
//  1. A Goroutine is waiting to receive (recvq) → direct transfer, bypasses buffer
//  2. Buffer has space → write to buffer
//  3. Buffer is full or Unbuffered → block current Goroutine into sendq
func (c *hchan) send(val any, g *goroutine) bool {
	c.lock.Lock()

	if c.closed {
		c.lock.Unlock()
		panic("send on closed channel")
	}

	// Case 1: A Goroutine is waiting to receive → direct transfer
	if recv := c.recvq.dequeue(); recv != nil {
		recv.g.result = val
		close(recv.g.done) // "wake up" the waiting receiver Goroutine
		c.lock.Unlock()
		fmt.Printf("[send] G%d → G%d (direct): %v\n", g.id, recv.g.id, val)
		return true
	}

	// Case 2: Buffer has space
	if c.qcount < c.dataqsiz {
		c.buf[c.sendx] = val
		c.sendx = (c.sendx + 1) % c.dataqsiz
		c.qcount++
		c.lock.Unlock()
		fmt.Printf("[send] G%d → buffer: %v (qcount=%d)\n", g.id, val, c.qcount)
		return true
	}

	// Case 3: Block — add Goroutine to sendq
	s := &sudog{g: g, elem: val}
	c.sendq.enqueue(s)
	c.lock.Unlock()
	fmt.Printf("[send] G%d blocked (sendq), waiting to send: %v\n", g.id, val)
	<-g.done // wait to be woken up
	return true
}

// recv simulates val := <-ch
func (c *hchan) recv(g *goroutine) (any, bool) {
	c.lock.Lock()

	if c.closed && c.qcount == 0 {
		c.lock.Unlock()
		return nil, false
	}

	// Case 1: A Goroutine is waiting to send → receive directly
	if send := c.sendq.dequeue(); send != nil {
		val := send.elem
		close(send.g.done) // "wake up" the waiting sender Goroutine
		c.lock.Unlock()
		fmt.Printf("[recv] G%d ← G%d (direct): %v\n", g.id, send.g.id, val)
		return val, true
	}

	// Case 2: Buffer has data
	if c.qcount > 0 {
		val := c.buf[c.recvx]
		c.buf[c.recvx] = nil
		c.recvx = (c.recvx + 1) % c.dataqsiz
		c.qcount--
		c.lock.Unlock()
		fmt.Printf("[recv] G%d ← buffer: %v (qcount=%d)\n", g.id, val, c.qcount)
		return val, true
	}

	// Case 3: Block — add Goroutine to recvq
	s := &sudog{g: g}
	c.recvq.enqueue(s)
	c.lock.Unlock()
	fmt.Printf("[recv] G%d blocked (recvq), waiting to recv...\n", g.id)
	<-g.done
	return g.result, true
}

func (c *hchan) close() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.closed = true
	for {
		s := c.recvq.dequeue()
		if s == nil {
			break
		}
		close(s.g.done)
	}
}

func newG(id int) *goroutine {
	return &goroutine{id: id, done: make(chan struct{})}
}

func main() {
	fmt.Println("=== Buffered channel (size=2) ===")
	ch := makeHchan(2)
	g1, g2, g3 := newG(1), newG(2), newG(3)

	ch.send(10, g1)
	ch.send(20, g2)
	val, _ := ch.recv(g3)
	fmt.Println("Received:", val)

	fmt.Println("\n=== Unbuffered channel ===")
	uch := makeHchan(0)
	sender := newG(10)
	receiver := newG(11)

	go func() {
		uch.send(42, sender)
	}()

	val, _ = uch.recv(receiver)
	fmt.Println("Received:", val)

	fmt.Println("\n=== Tóm tắt Channel Internals ===")
	fmt.Println("Channel = hchan struct với:")
	fmt.Println("  - Circular buffer (cho buffered channel)")
	fmt.Println("  - recvq: queue goroutines đang chờ nhận")
	fmt.Println("  - sendq: queue goroutines đang chờ gửi")
	fmt.Println("  - lock: mutex bảo vệ cấu trúc")
	fmt.Println()
	fmt.Println("3 trường hợp khi SEND:")
	fmt.Println("  1. recvq có goroutine → direct transfer (không qua buffer)")
	fmt.Println("  2. Buffer còn chỗ → ghi vào buffer")
	fmt.Println("  3. Buffer đầy → goroutine block vào sendq")
}
