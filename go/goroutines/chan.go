//go:build ignore

package main

import (
	"fmt"
	"sync"
)

// sudog đại diện cho một Goroutine đang bị blocked (chờ gửi hoặc nhận)
type sudog struct {
	g    *goroutine // Goroutine đang bị blocked
	elem any        // Dữ liệu cần gửi hoặc nơi nhận dữ liệu
	next *sudog     // Con trỏ tới sudog tiếp theo trong hàng đợi
	prev *sudog     // Con trỏ tới sudog trước đó trong hàng đợi
}

// goroutine giả lập một Goroutine đơn giản
type goroutine struct {
	id     int
	result any
	done   chan struct{}
}

// waitq là danh sách liên kết 2 chiều của các sudog (Goroutine đang chờ)
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

// hchan là cấu trúc nội bộ của một channel trong Go runtime
type hchan struct {
	// --- Buffer (vòng tròn) ---
	buf      []any // Mảng vòng chứa dữ liệu (chỉ dùng khi có buffer)
	dataqsiz uint  // Sức chứa tối đa của buffer (= 0 nếu Unbuffered)
	qcount   uint  // Số phần tử hiện có trong buffer
	sendx    uint  // Index tiếp theo để ghi vào buffer
	recvx    uint  // Index tiếp theo để đọc từ buffer

	// --- Hàng đợi Goroutine ---
	recvq waitq // Các Goroutine đang bị blocked vì chờ NHẬN (buffer rỗng)
	sendq waitq // Các Goroutine đang bị blocked vì chờ GỬI (buffer đầy)

	// --- Trạng thái ---
	closed bool

	// --- Đồng bộ hóa ---
	lock sync.Mutex
}

func makeHchan(size uint) *hchan {
	return &hchan{
		buf:      make([]any, size),
		dataqsiz: size,
	}
}

// send mô phỏng ch <- val
//
// Có 3 trường hợp:
//  1. Có Goroutine đang chờ nhận (recvq) → chuyển thẳng, không qua buffer
//  2. Buffer còn chỗ → ghi vào buffer
//  3. Buffer đầy (hoặc Unbuffered) → block Goroutine hiện tại vào sendq
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
		fmt.Printf("[send] G%d → buffer[%d]: %v (qcount=%d)\n", g.id, c.sendx-1, val, c.qcount)
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
//
// Có 3 trường hợp:
//  1. Có Goroutine đang chờ gửi (sendq) → nhận thẳng từ Goroutine đó
//  2. Buffer có dữ liệu → đọc từ buffer
//  3. Buffer rỗng (hoặc Unbuffered) → block Goroutine hiện tại vào recvq
func (c *hchan) recv(g *goroutine) (any, bool) {
	c.lock.Lock()

	// Channel đóng và buffer rỗng → trả về zero value
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
		fmt.Printf("[recv] G%d ← buffer[%d]: %v (qcount=%d)\n", g.id, c.recvx-1, val, c.qcount)
		return val, true
	}

	// TH3: Block — thêm Goroutine vào recvq
	s := &sudog{g: g}
	c.recvq.enqueue(s)
	c.lock.Unlock()
	fmt.Printf("[recv] G%d blocked (recvq), waiting to recv...\n", g.id)
	<-g.done // chờ được đánh thức
	return g.result, true
}

func (c *hchan) close() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.closed = true

	// Đánh thức tất cả Goroutine đang chờ nhận
	for {
		s := c.recvq.dequeue()
		if s == nil {
			break
		}
		close(s.g.done)
	}
}

// --- Demo ---

func newG(id int) *goroutine {
	return &goroutine{id: id, done: make(chan struct{})}
}

func main() {
	fmt.Println("=== Buffered channel (size=2) ===")
	ch := makeHchan(2)

	g1 := newG(1)
	g2 := newG(2)
	g3 := newG(3)

	// Gửi 2 giá trị vào buffer
	ch.send(10, g1)
	ch.send(20, g2)

	// Nhận từ buffer
	val, _ := ch.recv(g3)
	fmt.Println("Received:", val)

	fmt.Println("\n=== Unbuffered channel ===")
	uch := makeHchan(0)

	sender := newG(10)
	receiver := newG(11)

	// Sender bị block trước
	go func() {
		uch.send(42, sender)
	}()

	// Receiver nhận trực tiếp từ sender
	val, _ = uch.recv(receiver)
	fmt.Println("Received:", val)
}

// for {
// 	select {
// 	case val, ok := <-one:
// 		if !ok {
// 			one = nil // Đóng channel one để không nhận thêm
// 			continue
// 		}
// 		fmt.Println("Received from one:", val)
// 	case val, ok := <-two:
// 		if !ok {
// 			two = nil // Đóng channel two để không nhận thêm
// 			continue
// 		}
// 		fmt.Println("Received from two:", val)
// 	}

// 	if one == nil && two == nil {
// 		break
// 	}
// }
