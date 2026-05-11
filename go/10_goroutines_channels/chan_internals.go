//go:build ignore

// === Channel Internals — mô phỏng hchan (Go runtime internal) ===
//
// File này là "documentation-as-code" — mô phỏng cách Go runtime
// implement channel bên trong. Chạy: go run chan_internals.go
//
// Cấu trúc thật trong runtime/chan.go:
//   type hchan struct {
//       qcount   uint           // số phần tử trong buffer
//       dataqsiz uint           // capacity của buffer
//       buf      unsafe.Pointer // buffer (circular array)
//       elemsize uint16         // size của mỗi element
//       closed   uint32
//       sendx    uint           // index gửi tiếp theo
//       recvx    uint           // index nhận tiếp theo
//       recvq    waitq          // queue goroutines đang chờ nhận
//       sendq    waitq          // queue goroutines đang chờ gửi
//       lock     mutex
//   }

package main

import (
	"fmt"
	"sync"
)

// sudog đại diện cho một Goroutine đang bị blocked (chờ gửi hoặc nhận)
type sudog struct {
	g    *goroutine
	elem any
	next *sudog
	prev *sudog
}

// goroutine giả lập một Goroutine đơn giản
type goroutine struct {
	id     int
	result any
	done   chan struct{}
}

// waitq là danh sách liên kết 2 chiều của các sudog
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

// hchan mô phỏng cấu trúc nội bộ của channel trong Go runtime
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

// send mô phỏng ch <- val
//
// 3 trường hợp:
//  1. Có Goroutine đang chờ nhận (recvq) → chuyển thẳng, không qua buffer
//  2. Buffer còn chỗ → ghi vào buffer
//  3. Buffer đầy hoặc Unbuffered → block Goroutine hiện tại vào sendq
func (c *hchan) send(val any, g *goroutine) bool {
	c.lock.Lock()

	if c.closed {
		c.lock.Unlock()
		panic("send on closed channel")
	}

	// TH1: Có Goroutine đang chờ nhận → chuyển thẳng
	if recv := c.recvq.dequeue(); recv != nil {
		recv.g.result = val
		close(recv.g.done) // "đánh thức" Goroutine đang chờ nhận
		c.lock.Unlock()
		fmt.Printf("[send] G%d → G%d (direct): %v\n", g.id, recv.g.id, val)
		return true
	}

	// TH2: Buffer còn chỗ
	if c.qcount < c.dataqsiz {
		c.buf[c.sendx] = val
		c.sendx = (c.sendx + 1) % c.dataqsiz
		c.qcount++
		c.lock.Unlock()
		fmt.Printf("[send] G%d → buffer: %v (qcount=%d)\n", g.id, val, c.qcount)
		return true
	}

	// TH3: Block — thêm Goroutine vào sendq
	s := &sudog{g: g, elem: val}
	c.sendq.enqueue(s)
	c.lock.Unlock()
	fmt.Printf("[send] G%d blocked (sendq), waiting to send: %v\n", g.id, val)
	<-g.done // chờ được đánh thức
	return true
}

// recv mô phỏng val := <-ch
func (c *hchan) recv(g *goroutine) (any, bool) {
	c.lock.Lock()

	if c.closed && c.qcount == 0 {
		c.lock.Unlock()
		return nil, false
	}

	// TH1: Có Goroutine đang chờ gửi → nhận thẳng
	if send := c.sendq.dequeue(); send != nil {
		val := send.elem
		close(send.g.done) // "đánh thức" Goroutine đang chờ gửi
		c.lock.Unlock()
		fmt.Printf("[recv] G%d ← G%d (direct): %v\n", g.id, send.g.id, val)
		return val, true
	}

	// TH2: Buffer có dữ liệu
	if c.qcount > 0 {
		val := c.buf[c.recvx]
		c.buf[c.recvx] = nil
		c.recvx = (c.recvx + 1) % c.dataqsiz
		c.qcount--
		c.lock.Unlock()
		fmt.Printf("[recv] G%d ← buffer: %v (qcount=%d)\n", g.id, val, c.qcount)
		return val, true
	}

	// TH3: Block — thêm Goroutine vào recvq
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
