// Bài 30: JSON Encoding & Decoding
// json.Marshal/Unmarshal, struct tags, custom marshaling, streaming, RawMessage
// Chạy: go run .
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"
)

func main() {
	fmt.Println("=== JSON Encoding & Decoding ===")

	fmt.Println("\n=== 1. Basic Marshal/Unmarshal ===")
	demoBasicJSON()

	fmt.Println("\n=== 2. Struct Tags ===")
	demoStructTags()

	fmt.Println("\n=== 3. Custom Marshaling ===")
	demoCustomMarshal()

	fmt.Println("\n=== 4. json.RawMessage ===")
	demoRawMessage()

	fmt.Println("\n=== 5. Streaming JSON ===")
	demoStreaming()

	fmt.Println("\n=== 6. Common Patterns ===")
	demoCommonPatterns()
}

// ============================================================
// 1. Basic Marshal / Unmarshal
// ============================================================

type Person struct {
	Name string
	Age  int
}

func demoBasicJSON() {
	// Marshal (struct → JSON)
	p := Person{Name: "Alice", Age: 30}
	data, err := json.Marshal(p)
	if err != nil {
		fmt.Printf("  marshal error: %v\n", err)
		return
	}
	fmt.Printf("  Marshal: %s\n", data)

	// Unmarshal (JSON → struct)
	var p2 Person
	if err := json.Unmarshal([]byte(`{"Name":"Bob","Age":25}`), &p2); err != nil {
		fmt.Printf("  unmarshal error: %v\n", err)
		return
	}
	fmt.Printf("  Unmarshal: %+v\n", p2)

	// MarshalIndent — pretty print
	pretty, _ := json.MarshalIndent(p, "  ", "  ")
	fmt.Printf("  MarshalIndent:\n%s\n", pretty)

	// Marshal slice/map
	scores := map[string]int{"Alice": 95, "Bob": 87}
	data, _ = json.Marshal(scores)
	fmt.Printf("  Map: %s\n", data)

	nums := []int{1, 2, 3, 4, 5}
	data, _ = json.Marshal(nums)
	fmt.Printf("  Slice: %s\n", data)
}

// ============================================================
// 2. Struct Tags
// ============================================================

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`         // luôn bỏ qua khi marshal/unmarshal
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty"` // bỏ qua nếu zero value
	Internal  string    `json:"internal,omitempty"`
}

func demoStructTags() {
	u := User{
		ID:       1,
		Name:     "Alice",
		Email:    "alice@example.com",
		Password: "secret123", // sẽ không xuất hiện trong JSON
	}

	data, _ := json.MarshalIndent(u, "  ", "  ")
	fmt.Printf("  Output:\n%s\n", data)
	fmt.Println("  Note: Password (tag \"-\") bị loại, UpdatedAt/Internal (\"omitempty\") bị loại vì zero")

	// Unmarshal với tags
	jsonStr := `{"id":2,"name":"Bob","email":"bob@example.com","password":"ignored"}`
	var u2 User
	json.Unmarshal([]byte(jsonStr), &u2)
	fmt.Printf("  Unmarshal: id=%d name=%s password=%q (empty — tag \"-\")\n",
		u2.ID, u2.Name, u2.Password)
}

// ============================================================
// 3. Custom Marshaling
// ============================================================

// Duration wrapper — marshal/unmarshal thành "1h30m" thay vì nanoseconds
type Duration struct {
	time.Duration
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Duration.String())
}

func (d *Duration) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	dur, err := time.ParseDuration(s)
	if err != nil {
		return fmt.Errorf("invalid duration %q: %w", s, err)
	}
	d.Duration = dur
	return nil
}

// Money — marshal thành cents integer thay vì float (tránh float precision)
type Money struct {
	Cents int64
}

func (m Money) MarshalJSON() ([]byte, error) {
	// Serialize as {"amount": 1099, "currency": "USD"}
	return json.Marshal(struct {
		Amount   int64  `json:"amount"`
		Currency string `json:"currency"`
	}{Amount: m.Cents, Currency: "USD"})
}

func (m *Money) UnmarshalJSON(data []byte) error {
	var v struct {
		Amount int64 `json:"amount"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	m.Cents = v.Amount
	return nil
}

type Config struct {
	Name    string   `json:"name"`
	Timeout Duration `json:"timeout"`
	Price   Money    `json:"price"`
}

func demoCustomMarshal() {
	cfg := Config{
		Name:    "myapp",
		Timeout: Duration{90 * time.Second},
		Price:   Money{Cents: 1099},
	}

	data, _ := json.MarshalIndent(cfg, "  ", "  ")
	fmt.Printf("  Custom marshal:\n%s\n", data)

	// Unmarshal back
	jsonStr := `{"name":"app2","timeout":"2m30s","price":{"amount":2500,"currency":"USD"}}`
	var cfg2 Config
	if err := json.Unmarshal([]byte(jsonStr), &cfg2); err != nil {
		fmt.Printf("  unmarshal error: %v\n", err)
		return
	}
	fmt.Printf("  Custom unmarshal: name=%s timeout=%v price=%d cents\n",
		cfg2.Name, cfg2.Timeout.Duration, cfg2.Price.Cents)
}

// ============================================================
// 4. json.RawMessage — defer decode
// ============================================================

type APIResponse struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"` // decode sau dựa vào Type
}

type OrderPayload struct {
	OrderID string  `json:"order_id"`
	Amount  float64 `json:"amount"`
}

type UserPayload struct {
	UserID string `json:"user_id"`
	Action string `json:"action"`
}

func demoRawMessage() {
	messages := []string{
		`{"type":"order","payload":{"order_id":"ORD-123","amount":99.99}}`,
		`{"type":"user","payload":{"user_id":"USR-456","action":"login"}}`,
	}

	for _, msg := range messages {
		var resp APIResponse
		if err := json.Unmarshal([]byte(msg), &resp); err != nil {
			fmt.Printf("  error: %v\n", err)
			continue
		}

		fmt.Printf("  Type: %s\n", resp.Type)
		switch resp.Type {
		case "order":
			var order OrderPayload
			json.Unmarshal(resp.Payload, &order)
			fmt.Printf("  Order: id=%s amount=%.2f\n", order.OrderID, order.Amount)
		case "user":
			var user UserPayload
			json.Unmarshal(resp.Payload, &user)
			fmt.Printf("  User: id=%s action=%s\n", user.UserID, user.Action)
		}
	}

	// RawMessage cũng hữu ích khi forward JSON mà không decode
	raw := json.RawMessage(`{"key":"value","nested":{"a":1}}`)
	combined := struct {
		Meta    string          `json:"meta"`
		Dynamic json.RawMessage `json:"dynamic"`
	}{Meta: "forwarded", Dynamic: raw}

	data, _ := json.Marshal(combined)
	fmt.Printf("  Forwarded: %s\n", data)
}

// ============================================================
// 5. Streaming JSON
// ============================================================

func demoStreaming() {
	// json.Encoder — stream to writer (HTTP response, file)
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")

	users := []Person{
		{"Alice", 30},
		{"Bob", 25},
		{"Charlie", 35},
	}

	fmt.Println("  Encoder (streaming write):")
	for _, u := range users {
		enc.Encode(u) // writes one JSON object + newline per call
	}
	fmt.Print(buf.String())

	// json.Decoder — stream from reader (large files, HTTP request body)
	fmt.Println("  Decoder (streaming read):")
	jsonLines := `{"Name":"Dave","Age":28}
{"Name":"Eve","Age":32}
{"Name":"Frank","Age":27}`

	dec := json.NewDecoder(strings.NewReader(jsonLines))
	for {
		var p Person
		if err := dec.Decode(&p); err == io.EOF {
			break
		} else if err != nil {
			fmt.Printf("  decode error: %v\n", err)
			break
		}
		fmt.Printf("  Decoded: %+v\n", p)
	}

	// NGUYÊN TẮC:
	// - json.Marshal/Unmarshal: in-memory, small payloads
	// - json.Encoder/Decoder: streaming, large payloads, avoid loading all into memory
}

// ============================================================
// 6. Common Patterns
// ============================================================

// Nullable fields với pointer
type Profile struct {
	Name     string  `json:"name"`
	Bio      *string `json:"bio"`       // null nếu nil, omit nếu dùng omitempty
	Age      *int    `json:"age"`
	Verified bool    `json:"verified"`
}

func demoCommonPatterns() {
	// Nullable fields
	bio := "Go developer"
	age := 28
	p := Profile{Name: "Alice", Bio: &bio, Age: &age, Verified: true}

	data, _ := json.Marshal(p)
	fmt.Printf("  With values: %s\n", data)

	p2 := Profile{Name: "Bob"} // Bio và Age là nil → "null"
	data, _ = json.Marshal(p2)
	fmt.Printf("  With nil: %s\n", data)

	// Unknown fields — json.Decoder strict mode
	type Strict struct {
		Name string `json:"name"`
	}
	dec := json.NewDecoder(strings.NewReader(`{"name":"Alice","unknown_field":"value"}`))
	dec.DisallowUnknownFields() // return error nếu có field lạ
	var s Strict
	if err := dec.Decode(&s); err != nil {
		fmt.Printf("  DisallowUnknownFields error: %v\n", err)
	}

	// number decoder — tránh float64 cho large integers
	dec2 := json.NewDecoder(strings.NewReader(`{"id":9007199254740993}`))
	dec2.UseNumber() // parse number thành json.Number thay vì float64
	var m map[string]any
	dec2.Decode(&m)
	num := m["id"].(json.Number)
	fmt.Printf("  UseNumber: %s (giữ nguyên precision)\n", num.String())

	fmt.Println()
	fmt.Println("  Tips:")
	fmt.Println("  - Dùng omitempty cho optional fields")
	fmt.Println("  - Dùng pointer types cho nullable fields")
	fmt.Println("  - Dùng json.Number hoặc string cho large integers (>2^53)")
	fmt.Println("  - Dùng DisallowUnknownFields cho strict API validation")
	fmt.Println("  - Tránh marshal/unmarshal trong hot path — cache hoặc dùng sonic/jsoniter")
}
