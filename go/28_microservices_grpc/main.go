// Bài 28: Microservices & gRPC trong Go
// Giải thích concepts, patterns, và workflow — không cần protoc để chạy demo
// Để chạy full gRPC: xem README trong folder này
// Chạy: go run .
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"
)

func main() {
	fmt.Println("=== Microservices & gRPC ===")

	fmt.Println("\n=== 1. gRPC Concepts ===")
	explainGRPC()

	fmt.Println("\n=== 2. Proto Definition ===")
	showProtoDefinition()

	fmt.Println("\n=== 3. gRPC Server Pattern ===")
	showServerPattern()

	fmt.Println("\n=== 4. gRPC Client Pattern ===")
	showClientPattern()

	fmt.Println("\n=== 5. gRPC Error Handling ===")
	showGRPCErrors()

	fmt.Println("\n=== 6. REST vs gRPC Comparison ===")
	demoRESTvsgRPC()

	fmt.Println("\n=== 7. Service-to-Service Communication ===")
	demoServiceMesh()

	fmt.Println("\n=== 8. gRPC Interceptors (Middleware) ===")
	showInterceptors()
}

func explainGRPC() {
	fmt.Println("  gRPC = Google Remote Procedure Call")
	fmt.Println("  - Protocol: HTTP/2 (multiplexing, bidirectional streaming)")
	fmt.Println("  - Serialization: Protocol Buffers (protobuf) — smaller, faster than JSON")
	fmt.Println("  - Code generation: protoc tự generate client + server stubs")
	fmt.Println()
	fmt.Println("  Streaming types:")
	fmt.Println("  - Unary: client gửi 1 request, server trả 1 response (như HTTP)")
	fmt.Println("  - Server streaming: server trả nhiều responses")
	fmt.Println("  - Client streaming: client gửi nhiều requests")
	fmt.Println("  - Bidirectional: cả hai sides đều stream")
	fmt.Println()
	fmt.Println("  Khi nào dùng gRPC thay REST:")
	fmt.Println("  ✓ Internal service-to-service communication")
	fmt.Println("  ✓ Cần strong typing và contract (proto file = API contract)")
	fmt.Println("  ✓ High throughput, low latency")
	fmt.Println("  ✓ Streaming data")
	fmt.Println("  ✗ Public APIs (browser support hạn chế)")
	fmt.Println("  ✗ Human-readable debugging")
}

func showProtoDefinition() {
	proto := `
  // user.proto
  syntax = "proto3";
  package user;
  option go_package = "./pb";

  service UserService {
    // Unary RPC
    rpc GetUser(GetUserRequest) returns (UserResponse);

    // Server-side streaming: trả về nhiều users
    rpc ListUsers(ListUsersRequest) returns (stream UserResponse);

    // Client-side streaming: upload nhiều users
    rpc BulkCreateUsers(stream CreateUserRequest) returns (BulkCreateResponse);

    // Bidirectional streaming
    rpc Chat(stream ChatMessage) returns (stream ChatMessage);
  }

  message GetUserRequest { string user_id = 1; }
  message UserResponse {
    string id         = 1;
    string name       = 2;
    string email      = 3;
    string created_at = 4;
  }
  message ListUsersRequest {
    int32 limit  = 1;
    int32 offset = 2;
  }
  message CreateUserRequest {
    string name  = 1;
    string email = 2;
  }
  message BulkCreateResponse { int32 created = 1; }
  message ChatMessage {
    string from    = 1;
    string content = 2;
  }`
	fmt.Println(proto)
	fmt.Println()
	fmt.Println("  Để generate Go code từ proto:")
	fmt.Println("  $ brew install protobuf")
	fmt.Println("  $ go install google.golang.org/protobuf/cmd/protoc-gen-go@latest")
	fmt.Println("  $ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest")
	fmt.Println("  $ protoc --go_out=. --go-grpc_out=. proto/user.proto")
}

func showServerPattern() {
	code := `
  // server/main.go
  package main

  import (
      "context"
      "net"
      "google.golang.org/grpc"
      "google.golang.org/grpc/codes"
      "google.golang.org/grpc/status"
      pb "github.com/you/app/pb"
  )

  type userServer struct {
      pb.UnimplementedUserServiceServer // embed để forward-compatible
      store UserStore
  }

  // Unary RPC
  func (s *userServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
      user, err := s.store.FindByID(ctx, req.UserId)
      if err != nil {
          // Map lỗi domain sang gRPC status codes
          return nil, status.Errorf(codes.NotFound, "user %s not found", req.UserId)
      }
      return &pb.UserResponse{Id: user.ID, Name: user.Name, Email: user.Email}, nil
  }

  // Server streaming RPC
  func (s *userServer) ListUsers(req *pb.ListUsersRequest, stream pb.UserService_ListUsersServer) error {
      users, _ := s.store.List(stream.Context(), int(req.Limit), int(req.Offset))
      for _, u := range users {
          if err := stream.Send(&pb.UserResponse{Id: u.ID, Name: u.Name}); err != nil {
              return err // client disconnected
          }
      }
      return nil
  }

  func main() {
      lis, _ := net.Listen("tcp", ":50051")
      srv := grpc.NewServer(
          grpc.ChainUnaryInterceptor(loggingInterceptor, authInterceptor),
      )
      pb.RegisterUserServiceServer(srv, &userServer{})
      srv.Serve(lis)
  }`
	fmt.Println(code)
}

func showClientPattern() {
	code := `
  // client/main.go
  package main

  import (
      "context"
      "io"
      "time"
      "google.golang.org/grpc"
      "google.golang.org/grpc/credentials/insecure"
      pb "github.com/you/app/pb"
  )

  func main() {
      // Tạo connection — không block, lazy connect
      conn, err := grpc.NewClient("localhost:50051",
          grpc.WithTransportCredentials(insecure.NewCredentials()),
      )
      if err != nil { panic(err) }
      defer conn.Close()

      client := pb.NewUserServiceClient(conn)

      // Unary call với timeout
      ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
      defer cancel()

      resp, err := client.GetUser(ctx, &pb.GetUserRequest{UserId: "123"})
      if err != nil {
          // status.Code(err) → codes.NotFound
          handleGRPCError(err)
          return
      }
      fmt.Printf("User: %+v\n", resp)

      // Server streaming
      stream, _ := client.ListUsers(ctx, &pb.ListUsersRequest{Limit: 10})
      for {
          user, err := stream.Recv()
          if err == io.EOF { break }
          if err != nil { break }
          fmt.Printf("  User: %s\n", user.Name)
      }
  }`
	fmt.Println(code)
}

func showGRPCErrors() {
	fmt.Println("  gRPC Status Codes (ánh xạ với HTTP):")

	codes := []struct {
		code    string
		http    int
		meaning string
	}{
		{"codes.OK", 200, "Success"},
		{"codes.InvalidArgument", 400, "Bad request — client error"},
		{"codes.NotFound", 404, "Resource not found"},
		{"codes.AlreadyExists", 409, "Duplicate resource"},
		{"codes.PermissionDenied", 403, "Không có quyền"},
		{"codes.Unauthenticated", 401, "Chưa xác thực"},
		{"codes.ResourceExhausted", 429, "Rate limit exceeded"},
		{"codes.Internal", 500, "Server error"},
		{"codes.Unavailable", 503, "Service down"},
		{"codes.DeadlineExceeded", 504, "Timeout"},
	}

	for _, c := range codes {
		fmt.Printf("  %-30s HTTP %-3d — %s\n", c.code, c.http, c.meaning)
	}

	fmt.Println()
	fmt.Println("  Cách dùng:")
	fmt.Printf("  return nil, status.Errorf(codes.NotFound, \"user %%s not found\", id)\n")
	fmt.Println()
	fmt.Println("  Client check:")
	fmt.Println(`  if status.Code(err) == codes.NotFound { ... }`)
}

// ============================================================
// REST vs gRPC demo dùng net/http (không cần gRPC binary)
// ============================================================

type UserDTO struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func demoRESTvsgRPC() {
	// REST API simulation
	mux := http.NewServeMux()
	mux.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		user := UserDTO{ID: r.PathValue("id"), Name: "Alice", Email: "alice@example.com"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	})

	req := httptest.NewRequest("GET", "/users/123", nil)
	rw := httptest.NewRecorder()
	mux.ServeHTTP(rw, req)

	fmt.Println("  REST call: GET /users/123")
	fmt.Printf("  Response (JSON): %s", rw.Body.String())

	// gRPC là binary — không thể demo trực tiếp ở đây
	// nhưng đây là protobuf binary equivalent (ước tính):
	user := UserDTO{ID: "123", Name: "Alice", Email: "alice@example.com"}
	jsonBytes, _ := json.Marshal(user)

	// Protobuf sẽ nhỏ hơn ~3-4x:
	// field 1 (id): 0x0a 0x03 "123" = 5 bytes
	// field 2 (name): 0x12 0x05 "Alice" = 7 bytes
	// field 3 (email): 0x1a ... = ...
	protoApproxBytes := len(jsonBytes) / 3 // ước tính

	fmt.Printf("  JSON size: %d bytes\n", len(jsonBytes))
	fmt.Printf("  Proto approx: ~%d bytes (3-4x smaller)\n", protoApproxBytes)

	fmt.Println()
	fmt.Println("  REST vs gRPC tradeoffs:")
	fmt.Printf("  %-20s %-15s %-15s\n", "Feature", "REST/JSON", "gRPC/Protobuf")
	fmt.Printf("  %-20s %-15s %-15s\n", "Human readable", "✓", "✗ (binary)")
	fmt.Printf("  %-20s %-15s %-15s\n", "Browser support", "✓", "✗ (grpc-web needed)")
	fmt.Printf("  %-20s %-15s %-15s\n", "Performance", "Medium", "Fast")
	fmt.Printf("  %-20s %-15s %-15s\n", "Type safety", "Manual", "✓ (proto contract)")
	fmt.Printf("  %-20s %-15s %-15s\n", "Streaming", "Limited", "✓ Native")
	fmt.Printf("  %-20s %-15s %-15s\n", "Code generation", "Optional", "✓ Required")
}

// ============================================================
// Service mesh patterns
// ============================================================

type ServiceClient struct {
	name    string
	baseURL string
}

func (c *ServiceClient) Call(ctx context.Context, endpoint string) (string, error) {
	// Simulate network call với timeout
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case <-time.After(10 * time.Millisecond):
		return fmt.Sprintf("%s:%s response", c.name, endpoint), nil
	}
}

func demoServiceMesh() {
	fmt.Println("  Microservice communication patterns:")
	fmt.Println()
	fmt.Println("  1. Synchronous (gRPC/REST):")
	fmt.Println("     OrderService → UserService.GetUser()")
	fmt.Println("     OrderService → PaymentService.Charge()")
	fmt.Println()
	fmt.Println("  2. Async (message queue):")
	fmt.Println("     OrderService → [Kafka/NATS] → NotificationService")
	fmt.Println("     OrderService → [Kafka/NATS] → InventoryService")
	fmt.Println()
	fmt.Println("  3. Service discovery:")
	fmt.Println("     DNS-based: user-service.namespace.svc.cluster.local")
	fmt.Println("     Consul/etcd: dynamic registration")
	fmt.Println()

	// Simulate parallel service calls
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	userSvc := &ServiceClient{name: "user-service", baseURL: "localhost:50051"}
	productSvc := &ServiceClient{name: "product-service", baseURL: "localhost:50052"}

	type result struct {
		data string
		err  error
	}

	userCh := make(chan result, 1)
	productCh := make(chan result, 1)

	// Parallel calls
	go func() {
		data, err := userSvc.Call(ctx, "/GetUser")
		userCh <- result{data, err}
	}()
	go func() {
		data, err := productSvc.Call(ctx, "/GetProduct")
		productCh <- result{data, err}
	}()

	userRes := <-userCh
	productRes := <-productCh

	fmt.Printf("  Parallel calls:\n")
	fmt.Printf("    UserService: %s\n", userRes.data)
	fmt.Printf("    ProductService: %s\n", productRes.data)
}

func showInterceptors() {
	code := `
  // Unary interceptor (như middleware trong HTTP)
  func loggingInterceptor(
      ctx context.Context,
      req any,
      info *grpc.UnaryServerInfo,
      handler grpc.UnaryHandler,
  ) (any, error) {
      start := time.Now()
      resp, err := handler(ctx, req)  // gọi handler thực sự
      fmt.Printf("RPC %s took %v, err=%v\n", info.FullMethod, time.Since(start), err)
      return resp, err
  }

  // Auth interceptor
  func authInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
      md, ok := metadata.FromIncomingContext(ctx)
      if !ok { return nil, status.Error(codes.Unauthenticated, "missing metadata") }
      tokens := md.Get("authorization")
      if len(tokens) == 0 { return nil, status.Error(codes.Unauthenticated, "missing token") }
      // validate token...
      return handler(ctx, req)
  }

  // Chain interceptors:
  grpc.NewServer(
      grpc.ChainUnaryInterceptor(loggingInterceptor, authInterceptor, recoveryInterceptor),
  )`

	fmt.Println(code)
	fmt.Println()
	fmt.Println("  Các interceptor phổ biến:")
	interceptors := []string{
		"logging — request/response logging với duration",
		"auth — JWT/API key validation",
		"recovery — catch panics, return Internal error",
		"rate-limiting — giới hạn request rate",
		"tracing — OpenTelemetry span injection",
		"validation — validate request fields",
	}
	for _, i := range interceptors {
		fmt.Printf("  - %s\n", i)
	}

	_ = strings.Builder{}
}
