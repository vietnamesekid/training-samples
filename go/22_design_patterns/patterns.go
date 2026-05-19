package main

import (
	"fmt"
	"sync"
	"time"
)

// === 1. Functional Options Pattern ===

type Server struct {
	host     string
	port     int
	timeout  time.Duration
	maxConns int
}

type Option func(*Server)

func WithPort(port int) Option       { return func(s *Server) { s.port = port } }
func WithTimeout(d time.Duration) Option { return func(s *Server) { s.timeout = d } }
func WithMaxConns(n int) Option      { return func(s *Server) { s.maxConns = n } }

func NewServer(host string, opts ...Option) *Server {
	s := &Server{host: host, port: 8080, timeout: 30 * time.Second, maxConns: 1000}
	for _, o := range opts {
		o(s)
	}
	return s
}

func demoFunctionalOptions() {
	s1 := NewServer("localhost")
	fmt.Printf("  Default: %s:%d (timeout=%s)\n", s1.host, s1.port, s1.timeout)

	s2 := NewServer("api.example.com",
		WithPort(443),
		WithTimeout(60*time.Second),
		WithMaxConns(5000),
	)
	fmt.Printf("  Custom: %s:%d (timeout=%s, maxConns=%d)\n",
		s2.host, s2.port, s2.timeout, s2.maxConns)
}

// === 2. Builder Pattern — build complex objects step by step ===

type Query struct {
	table      string
	conditions []string
	orderBy    string
	limit      int
	offset     int
}

type QueryBuilder struct {
	query Query
}

func NewQuery(table string) *QueryBuilder {
	return &QueryBuilder{query: Query{table: table, limit: -1}}
}

func (b *QueryBuilder) Where(condition string) *QueryBuilder {
	b.query.conditions = append(b.query.conditions, condition)
	return b // fluent API: return self for chaining
}

func (b *QueryBuilder) OrderBy(field string) *QueryBuilder {
	b.query.orderBy = field
	return b
}

func (b *QueryBuilder) Limit(n int) *QueryBuilder {
	b.query.limit = n
	return b
}

func (b *QueryBuilder) Offset(n int) *QueryBuilder {
	b.query.offset = n
	return b
}

func (b *QueryBuilder) Build() string {
	sql := fmt.Sprintf("SELECT * FROM %s", b.query.table)
	for i, c := range b.query.conditions {
		if i == 0 {
			sql += " WHERE " + c
		} else {
			sql += " AND " + c
		}
	}
	if b.query.orderBy != "" {
		sql += " ORDER BY " + b.query.orderBy
	}
	if b.query.limit > 0 {
		sql += fmt.Sprintf(" LIMIT %d", b.query.limit)
	}
	if b.query.offset > 0 {
		sql += fmt.Sprintf(" OFFSET %d", b.query.offset)
	}
	return sql
}

func demoBuilder() {
	sql := NewQuery("users").
		Where("age > 18").
		Where("active = true").
		OrderBy("created_at DESC").
		Limit(10).
		Offset(20).
		Build()
	fmt.Printf("  SQL: %s\n", sql)

	// Simple query
	simple := NewQuery("products").Where("price < 100").Build()
	fmt.Printf("  SQL: %s\n", simple)
}

// === 3. Repository Pattern + Dependency Injection ===

type User struct {
	ID   int
	Name string
}

// Interface defined at the use site (service package), not the implementation site
type UserRepository interface {
	FindByID(id int) (*User, error)
	Save(u *User) error
}

type UserService struct {
	repo   UserRepository
	logger func(string) // injected logger
	clock  func() time.Time // injected clock (easy to test)
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{
		repo:   repo,
		logger: func(s string) { fmt.Println(" ", s) },
		clock:  time.Now,
	}
}

func (s *UserService) GetUser(id int) (*User, error) {
	s.logger(fmt.Sprintf("getting user %d", id))
	return s.repo.FindByID(id)
}

// In-memory implementation (production will use PostgresRepo, MongoRepo, etc.)
type InMemoryUserRepo struct {
	users map[int]*User
}

func NewInMemoryUserRepo() *InMemoryUserRepo {
	return &InMemoryUserRepo{users: map[int]*User{
		1: {ID: 1, Name: "Alice"},
		2: {ID: 2, Name: "Bob"},
	}}
}

func (r *InMemoryUserRepo) FindByID(id int) (*User, error) {
	u, ok := r.users[id]
	if !ok {
		return nil, fmt.Errorf("user %d not found", id)
	}
	return u, nil
}

func (r *InMemoryUserRepo) Save(u *User) error {
	r.users[u.ID] = u
	return nil
}

func demoRepository() {
	repo := NewInMemoryUserRepo()
	svc := NewUserService(repo)

	u, err := svc.GetUser(1)
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	} else {
		fmt.Printf("  User: %+v\n", u)
	}

	_, err = svc.GetUser(99)
	fmt.Printf("  GetUser(99): %v\n", err)
}

// === 4. Observer / Event Bus ===

type Event struct {
	Type string
	Data any
}

type Handler func(Event)

type EventBus struct {
	mu       sync.RWMutex
	handlers map[string][]Handler
}

func NewEventBus() *EventBus {
	return &EventBus{handlers: make(map[string][]Handler)}
}

func (b *EventBus) Subscribe(eventType string, handler Handler) func() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[eventType] = append(b.handlers[eventType], handler)

	// Returns unsubscribe function
	idx := len(b.handlers[eventType]) - 1
	return func() {
		b.mu.Lock()
		defer b.mu.Unlock()
		handlers := b.handlers[eventType]
		b.handlers[eventType] = append(handlers[:idx], handlers[idx+1:]...)
	}
}

func (b *EventBus) Publish(e Event) {
	b.mu.RLock()
	handlers := b.handlers[e.Type]
	b.mu.RUnlock()

	for _, h := range handlers {
		h(e)
	}
}

func demoObserver() {
	bus := NewEventBus()

	// Subscribe to user events
	unsubEmail := bus.Subscribe("user.created", func(e Event) {
		fmt.Printf("  [Email] New user: %v\n", e.Data)
	})

	bus.Subscribe("user.created", func(e Event) {
		fmt.Printf("  [Analytics] Track user signup: %v\n", e.Data)
	})

	bus.Subscribe("user.deleted", func(e Event) {
		fmt.Printf("  [Cleanup] Remove user data: %v\n", e.Data)
	})

	bus.Publish(Event{Type: "user.created", Data: map[string]any{"id": 1, "name": "Alice"}})
	bus.Publish(Event{Type: "user.created", Data: map[string]any{"id": 2, "name": "Bob"}})
	bus.Publish(Event{Type: "user.deleted", Data: map[string]any{"id": 1}})

	// Unsubscribe email handler
	unsubEmail()
	fmt.Println("  [After unsubscribe email]")
	bus.Publish(Event{Type: "user.created", Data: map[string]any{"id": 3, "name": "Carol"}})
}

// === 5. Strategy Pattern ===

type SortStrategy interface {
	Sort(data []int) []int
	Name() string
}

type BubbleSort struct{}

func (b BubbleSort) Name() string { return "BubbleSort" }
func (b BubbleSort) Sort(data []int) []int {
	result := make([]int, len(data))
	copy(result, data)
	n := len(result)
	for i := range n {
		for j := 0; j < n-i-1; j++ {
			if result[j] > result[j+1] {
				result[j], result[j+1] = result[j+1], result[j]
			}
		}
	}
	return result
}

type QuickSort struct{}

func (q QuickSort) Name() string { return "QuickSort" }
func (q QuickSort) Sort(data []int) []int {
	if len(data) <= 1 {
		return data
	}
	result := make([]int, len(data))
	copy(result, data)
	quickSort(result, 0, len(result)-1)
	return result
}

func quickSort(arr []int, lo, hi int) {
	if lo >= hi {
		return
	}
	pivot := arr[hi]
	i := lo - 1
	for j := lo; j < hi; j++ {
		if arr[j] <= pivot {
			i++
			arr[i], arr[j] = arr[j], arr[i]
		}
	}
	arr[i+1], arr[hi] = arr[hi], arr[i+1]
	quickSort(arr, lo, i)
	quickSort(arr, i+2, hi)
}

type Sorter struct {
	strategy SortStrategy
}

func (s *Sorter) SetStrategy(strategy SortStrategy) { s.strategy = strategy }
func (s *Sorter) Sort(data []int) []int             { return s.strategy.Sort(data) }

func demoStrategy() {
	data := []int{64, 34, 25, 12, 22, 11, 90}
	sorter := &Sorter{}

	for _, strategy := range []SortStrategy{BubbleSort{}, QuickSort{}} {
		sorter.SetStrategy(strategy)
		result := sorter.Sort(data)
		fmt.Printf("  %s: %v\n", strategy.Name(), result)
	}
}

// === 6. Singleton with sync.Once ===

type DBConnection struct {
	DSN string
}

var (
	dbConn *DBConnection
	dbOnce sync.Once
)

func GetDBConnection() *DBConnection {
	dbOnce.Do(func() {
		fmt.Println("  [Singleton] Creating DB connection...")
		dbConn = &DBConnection{DSN: "postgres://localhost:5432/myapp"}
	})
	return dbConn
}

func demoSingleton() {
	// Called multiple times — only initialized once
	for range 3 {
		db := GetDBConnection()
		fmt.Printf("  DB: %s\n", db.DSN)
	}
}

// === 7. Middleware Chain ===

type HandlerFunc func(request string) string
type MiddlewareFunc func(HandlerFunc) HandlerFunc

func Chain(h HandlerFunc, middlewares ...MiddlewareFunc) HandlerFunc {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

func LoggingMW(next HandlerFunc) HandlerFunc {
	return func(req string) string {
		fmt.Printf("  [LOG] Before: %q\n", req)
		resp := next(req)
		fmt.Printf("  [LOG] After: %q\n", resp)
		return resp
	}
}

func AuthMW(next HandlerFunc) HandlerFunc {
	return func(req string) string {
		if req == "" {
			return "401 Unauthorized"
		}
		return next(req)
	}
}

func CacheMW(next HandlerFunc) HandlerFunc {
	cache := make(map[string]string)
	return func(req string) string {
		if v, ok := cache[req]; ok {
			fmt.Printf("  [CACHE] Hit for %q\n", req)
			return v
		}
		resp := next(req)
		cache[req] = resp
		return resp
	}
}

func demoMiddlewareChain() {
	handler := func(req string) string {
		return "Hello, " + req + "!"
	}

	// Chain: LoggingMW → AuthMW → CacheMW → handler
	chained := Chain(handler, LoggingMW, AuthMW, CacheMW)

	fmt.Println("  Request 1:")
	resp := chained("Gopher")
	fmt.Printf("  Response: %s\n", resp)

	fmt.Println("  Request 2 (cached):")
	resp = chained("Gopher")
	fmt.Printf("  Response: %s\n", resp)

	fmt.Println("  Request 3 (unauthorized):")
	resp = chained("")
	fmt.Printf("  Response: %s\n", resp)
}
