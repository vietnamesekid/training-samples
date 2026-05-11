// Bài 29: Database Patterns trong Go
// database/sql pool, QueryContext, prepared statements, transactions, repository pattern
// Chạy: go run .
// Note: dùng in-memory SQLite simulation — không cần database thật
package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

func main() {
	fmt.Println("=== Database Patterns ===")

	fmt.Println("\n=== 1. database/sql Concepts ===")
	explainDatabaseSQL()

	fmt.Println("\n=== 2. Connection Pool Configuration ===")
	showConnectionPool()

	fmt.Println("\n=== 3. Query Patterns ===")
	demoQueryPatterns()

	fmt.Println("\n=== 4. Transaction Pattern ===")
	demoTransactions()

	fmt.Println("\n=== 5. Repository Pattern ===")
	demoRepository()

	fmt.Println("\n=== 6. Common Mistakes ===")
	showDatabaseMistakes()
}

func explainDatabaseSQL() {
	fmt.Println("  database/sql package:")
	fmt.Println("  - Abstract interface — driver-agnostic")
	fmt.Println("  - Built-in connection pooling")
	fmt.Println("  - Context-aware queries")
	fmt.Println()
	fmt.Println("  Popular drivers:")
	fmt.Println("  - PostgreSQL: github.com/jackc/pgx/v5 (recommended)")
	fmt.Println("                github.com/lib/pq (older, still used)")
	fmt.Println("  - MySQL:      github.com/go-sql-driver/mysql")
	fmt.Println("  - SQLite:     github.com/mattn/go-sqlite3 (CGo)")
	fmt.Println("                modernc.org/sqlite (pure Go)")
	fmt.Println()
	fmt.Println("  Driver registration:")
	fmt.Println(`  import _ "github.com/jackc/pgx/v5/stdlib" // side-effect import`)
	fmt.Println(`  db, err := sql.Open("pgx", "postgres://user:pass@host/dbname")`)
}

func showConnectionPool() {
	// sql.Open() không tạo connection ngay — lazy
	// db là connection pool, không phải single connection
	fmt.Println("  Connection Pool setup:")
	fmt.Println()

	code := `
  db, err := sql.Open("pgx", dsn)
  if err != nil { return err }

  // Cấu hình pool — QUAN TRỌNG cho production
  db.SetMaxOpenConns(25)          // max connections đang dùng
  db.SetMaxIdleConns(10)          // max connections idle (cached)
  db.SetConnMaxLifetime(5*time.Minute) // max tuổi thọ connection
  db.SetConnMaxIdleTime(1*time.Minute) // max idle time trước khi đóng

  // Verify connection
  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()
  if err := db.PingContext(ctx); err != nil {
      return fmt.Errorf("cannot connect to database: %w", err)
  }`

	fmt.Println(code)
	fmt.Println()
	fmt.Println("  Rule of thumb pool sizing:")
	fmt.Println("  MaxOpenConns = min(numCPU * 4, DB_max_connections / num_replicas)")
	fmt.Println("  MaxIdleConns = MaxOpenConns / 2")
}

// ============================================================
// Simulated DB cho demo (không cần actual database driver)
// ============================================================

// MockDB simulates database/sql interface
type MockDB struct {
	users map[int]*DBUser
	seq   int
}

type DBUser struct {
	ID        int
	Name      string
	Email     string
	CreatedAt time.Time
}

func newMockDB() *MockDB {
	return &MockDB{users: make(map[int]*DBUser)}
}

// Simulate sql.ErrNoRows
var errNoRows = errors.New("sql: no rows in result set")

func (db *MockDB) insertUser(ctx context.Context, name, email string) (int, error) {
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
	}
	db.seq++
	db.users[db.seq] = &DBUser{
		ID: db.seq, Name: name, Email: email,
		CreatedAt: time.Now(),
	}
	return db.seq, nil
}

func (db *MockDB) getUserByID(ctx context.Context, id int) (*DBUser, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	u, ok := db.users[id]
	if !ok {
		return nil, errNoRows
	}
	return u, nil
}

func (db *MockDB) listUsers(ctx context.Context, limit, offset int) ([]*DBUser, error) {
	var result []*DBUser
	count := 0
	for i := 1; i <= db.seq; i++ {
		if count >= offset {
			if len(result) >= limit {
				break
			}
			if u, ok := db.users[i]; ok {
				result = append(result, u)
			}
		}
		count++
	}
	return result, nil
}

func demoQueryPatterns() {
	db := newMockDB()
	ctx := context.Background()

	// Insert some data
	for _, u := range []struct{ name, email string }{
		{"Alice", "alice@example.com"},
		{"Bob", "bob@example.com"},
		{"Charlie", "charlie@example.com"},
	} {
		id, _ := db.insertUser(ctx, u.name, u.email)
		fmt.Printf("  Inserted user id=%d\n", id)
	}

	// Query single row
	fmt.Println()
	fmt.Println("  QueryRowContext pattern:")
	user, err := db.getUserByID(ctx, 2)
	if err != nil {
		if errors.Is(err, errNoRows) {
			fmt.Println("  user not found (sql.ErrNoRows)")
		} else {
			fmt.Printf("  error: %v\n", err)
		}
	} else {
		fmt.Printf("  Found: id=%d name=%s email=%s\n", user.ID, user.Name, user.Email)
	}

	// Query không tồn tại
	_, err = db.getUserByID(ctx, 999)
	if errors.Is(err, errNoRows) {
		fmt.Println("  id=999: not found (sql.ErrNoRows)")
	}

	// Query multiple rows
	fmt.Println()
	fmt.Println("  QueryContext pattern (multiple rows):")
	users, _ := db.listUsers(ctx, 10, 0)
	for _, u := range users {
		fmt.Printf("  - [%d] %s <%s>\n", u.ID, u.Name, u.Email)
	}

	// Prepared statement pattern
	fmt.Println()
	fmt.Println("  Prepared Statement pattern:")
	fmt.Println(`  stmt, err := db.PrepareContext(ctx, "SELECT * FROM users WHERE id = $1")`)
	fmt.Println(`  defer stmt.Close()`)
	fmt.Println(`  // Reuse stmt nhiều lần — tránh SQL injection, tăng performance`)
	fmt.Println(`  rows, err := stmt.QueryContext(ctx, userID)`)
}

// ============================================================
// Transaction Pattern
// ============================================================

type MockTx struct {
	db       *MockDB
	ops      []string
	rolledBack bool
	committed  bool
}

func (db *MockDB) beginTx(ctx context.Context) (*MockTx, error) {
	return &MockTx{db: db}, nil
}

func (tx *MockTx) insertUser(ctx context.Context, name, email string) (int, error) {
	if tx.rolledBack {
		return 0, fmt.Errorf("transaction already rolled back")
	}
	tx.ops = append(tx.ops, fmt.Sprintf("INSERT user(%s, %s)", name, email))
	return tx.db.insertUser(ctx, name, email)
}

func (tx *MockTx) deductCredits(ctx context.Context, userID, amount int) error {
	if amount > 100 {
		return fmt.Errorf("insufficient credits")
	}
	tx.ops = append(tx.ops, fmt.Sprintf("UPDATE credits user=%d amount=-%d", userID, amount))
	return nil
}

func (tx *MockTx) Commit() error {
	tx.committed = true
	return nil
}

func (tx *MockTx) Rollback() error {
	tx.rolledBack = true
	return nil
}

// transferWithTx minh họa transaction pattern với defer rollback
func transferWithTx(ctx context.Context, db *MockDB, fromUser, toUser string, amount int) error {
	tx, err := db.beginTx(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	// NGUYÊN TẮC: luôn defer rollback — nếu commit đã chạy, rollback sẽ là no-op
	defer tx.Rollback()

	fromID, err := tx.insertUser(ctx, fromUser, fromUser+"@example.com")
	if err != nil {
		return fmt.Errorf("insert from: %w", err)
	}

	if err := tx.deductCredits(ctx, fromID, amount); err != nil {
		return fmt.Errorf("deduct credits: %w", err) // trigger defer Rollback
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	fmt.Printf("  TX ops: %v\n", tx.ops)
	fmt.Printf("  TX committed: %v\n", tx.committed)
	return nil
}

func demoTransactions() {
	db := newMockDB()
	ctx := context.Background()

	// Thành công
	fmt.Println("  Successful transaction:")
	if err := transferWithTx(ctx, db, "Dave", "Eve", 50); err != nil {
		fmt.Printf("  ERROR: %v\n", err)
	} else {
		fmt.Println("  Transfer: OK")
	}

	// Thất bại (amount > 100)
	fmt.Println()
	fmt.Println("  Failed transaction (rollback):")
	if err := transferWithTx(ctx, db, "Frank", "Grace", 200); err != nil {
		fmt.Printf("  Expected error: %v\n", err)
		fmt.Println("  Transaction rolled back automatically via defer")
	}

	// Pattern code:
	fmt.Println()
	fmt.Println("  Transaction pattern:")
	fmt.Println(`  tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})`)
	fmt.Println(`  if err != nil { return err }`)
	fmt.Println(`  defer tx.Rollback() // no-op if committed`)
	fmt.Println()
	fmt.Println(`  // ... operations using tx instead of db ...`)
	fmt.Println()
	fmt.Println(`  return tx.Commit()`)
}

// ============================================================
// Repository Pattern
// ============================================================

// User domain model
type User struct {
	ID        int
	Name      string
	Email     string
	CreatedAt time.Time
}

// UserRepository interface — dependency inversion
// Service layer phụ thuộc vào interface, không phải implementation cụ thể
type UserRepository interface {
	Create(ctx context.Context, name, email string) (*User, error)
	FindByID(ctx context.Context, id int) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	List(ctx context.Context, limit, offset int) ([]*User, error)
	Delete(ctx context.Context, id int) error
}

// InMemoryUserRepo — dùng cho testing
type InMemoryUserRepo struct {
	db  *MockDB
	idx map[string]int // email → id index
}

func NewInMemoryUserRepo() *InMemoryUserRepo {
	return &InMemoryUserRepo{
		db:  newMockDB(),
		idx: make(map[string]int),
	}
}

func (r *InMemoryUserRepo) Create(ctx context.Context, name, email string) (*User, error) {
	if _, exists := r.idx[email]; exists {
		return nil, fmt.Errorf("email %s already exists", email)
	}
	id, err := r.db.insertUser(ctx, name, email)
	if err != nil {
		return nil, err
	}
	r.idx[email] = id
	return &User{ID: id, Name: name, Email: email, CreatedAt: time.Now()}, nil
}

func (r *InMemoryUserRepo) FindByID(ctx context.Context, id int) (*User, error) {
	u, err := r.db.getUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, errNoRows) {
			return nil, fmt.Errorf("user %d: %w", id, sql.ErrNoRows)
		}
		return nil, err
	}
	return &User{ID: u.ID, Name: u.Name, Email: u.Email, CreatedAt: u.CreatedAt}, nil
}

func (r *InMemoryUserRepo) FindByEmail(ctx context.Context, email string) (*User, error) {
	id, ok := r.idx[email]
	if !ok {
		return nil, sql.ErrNoRows
	}
	return r.FindByID(ctx, id)
}

func (r *InMemoryUserRepo) List(ctx context.Context, limit, offset int) ([]*User, error) {
	dbUsers, err := r.db.listUsers(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	users := make([]*User, len(dbUsers))
	for i, u := range dbUsers {
		users[i] = &User{ID: u.ID, Name: u.Name, Email: u.Email}
	}
	return users, nil
}

func (r *InMemoryUserRepo) Delete(ctx context.Context, id int) error {
	if u, ok := r.db.users[id]; ok {
		delete(r.idx, u.Email)
		delete(r.db.users, id)
		return nil
	}
	return sql.ErrNoRows
}

// UserService — business logic layer
type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(ctx context.Context, name, email string) (*User, error) {
	// Business validation
	if name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if !isValidEmail(email) {
		return nil, fmt.Errorf("invalid email: %s", email)
	}

	user, err := s.repo.Create(ctx, name, email)
	if err != nil {
		return nil, fmt.Errorf("register user: %w", err)
	}
	return user, nil
}

func (s *UserService) GetProfile(ctx context.Context, id int) (*User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user %d not found", id)
		}
		return nil, fmt.Errorf("get profile: %w", err)
	}
	return user, nil
}

func isValidEmail(email string) bool {
	return len(email) > 3 && contains(email, "@")
}

func contains(s, sub string) bool {
	for i := range len(s) - len(sub) + 1 {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

func demoRepository() {
	ctx := context.Background()

	// Dùng InMemoryUserRepo cho demo (production sẽ dùng PostgresUserRepo)
	repo := NewInMemoryUserRepo()
	svc := NewUserService(repo)

	// Create users
	users := []struct{ name, email string }{
		{"Alice", "alice@example.com"},
		{"Bob", "bob@example.com"},
	}
	for _, u := range users {
		created, err := svc.Register(ctx, u.name, u.email)
		if err != nil {
			fmt.Printf("  Error registering %s: %v\n", u.name, err)
			continue
		}
		fmt.Printf("  Registered: id=%d name=%s\n", created.ID, created.Name)
	}

	// Duplicate email
	_, err := svc.Register(ctx, "Alice2", "alice@example.com")
	fmt.Printf("  Duplicate email: %v\n", err)

	// Validation error
	_, err = svc.Register(ctx, "", "invalid")
	fmt.Printf("  Validation error: %v\n", err)

	// FindByID
	user, err := svc.GetProfile(ctx, 1)
	if err != nil {
		fmt.Printf("  GetProfile error: %v\n", err)
	} else {
		fmt.Printf("  GetProfile: %s <%s>\n", user.Name, user.Email)
	}

	// Not found
	_, err = svc.GetProfile(ctx, 999)
	fmt.Printf("  Not found: %v\n", err)

	// List
	all, _ := repo.List(ctx, 10, 0)
	fmt.Printf("  All users (%d):\n", len(all))
	for _, u := range all {
		fmt.Printf("    [%d] %s\n", u.ID, u.Name)
	}

	fmt.Println()
	fmt.Println("  Repository pattern benefits:")
	fmt.Println("  - Service layer không biết gì về SQL/DB implementation")
	fmt.Println("  - Test dùng InMemoryRepo, production dùng PostgresRepo")
	fmt.Println("  - Dễ swap database (MongoDB, DynamoDB) mà không đổi business logic")
}

func showDatabaseMistakes() {
	fmt.Println("  Common database mistakes:")
	fmt.Println()

	mistakes := []struct {
		bad  string
		good string
	}{
		{
			"rows, _ := db.Query(...) // lỗi bị ignore",
			"rows, err := db.Query(...); if err != nil { ... }",
		},
		{
			"defer rows.Close() // sau khi check err",
			"rows, err := db.Query(...)\nif err != nil { return err }\ndefer rows.Close()",
		},
		{
			`db.Exec("SELECT * FROM users WHERE id=" + id) // SQL INJECTION!`,
			`db.ExecContext(ctx, "SELECT * FROM users WHERE id=$1", id) // parameterized`,
		},
		{
			"rows.Next() { rows.Scan(&u) } // không check rows.Err()",
			"rows.Next() { rows.Scan(&u) }\nif err := rows.Err(); err != nil { ... }",
		},
		{
			"db.SetMaxOpenConns(0) // unlimited connections",
			"db.SetMaxOpenConns(25) // set reasonable limit",
		},
		{
			"db.QueryRow(...) // không defer ctx cancel",
			"ctx, cancel := context.WithTimeout(ctx, 5*time.Second)\ndefer cancel()\ndb.QueryRowContext(ctx, ...)",
		},
	}

	for i, m := range mistakes {
		fmt.Printf("  %d. BAD:  %s\n", i+1, m.bad)
		fmt.Printf("     GOOD: %s\n\n", m.good)
	}
}
