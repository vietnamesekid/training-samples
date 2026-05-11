# Go Learning Repository — 34 Comprehensive Modules

**Comprehensive Go learning path from basics to advanced topics**, covering everything from syntax to production patterns. 34 standalone modules, each with runnable examples.

## Quick Start

```bash
cd go/01_hello_world
go run .
```

Each folder is a standalone Go module — just `go run .` any folder to see examples.

## Module Overview

### Core Topics (01–25)

| # | Topic | Key Concepts |
|---|-------|--------------|
| **01** | Hello World | package main, import, fmt, os.Args, go run vs build |
| **02** | Basic Syntax | var/:=, zero values, types, const, iota, type conversion, fmt verbs |
| **03** | Strings, Slices, Maps | UTF-8, arrays vs slices (value/header), slice gotchas, maps, range |
| **04** | Structs & Methods | constructors, value/pointer receivers, embedding, struct tags, fmt.Stringer |
| **05** | Control Flow & Functions | all for forms, switch, defer LIFO, closures, variadic, panic/recover, init() |
| **06** | Interfaces | duck typing, composition, nil interface gotcha, type assert/switch, functional options |
| **07** | Error Handling | errors.New, fmt.Errorf %w, custom errors, errors.Is/As, errors.Join |
| **08** | Packages & Modules | exported/unexported, init() order, internal/ package, go.mod structure |
| **09** | Pointers | &/* semantics, nil panic, new(), value vs pointer receiver, escape analysis |
| **10** | Goroutines & Channels | goroutine creation, WaitGroup, wg.Go (1.25+), directional channels, select, pipeline |
| **11** | Sync Primitives | Mutex, RWMutex, Once, Pool, sync.Map, atomic.Int64, sync.Cond |
| **12** | Context | WithCancel, WithTimeout, WithValue, propagation, goroutine leak prevention |
| **13** | Race Conditions | -race detector, 3 fixes (Mutex/atomic/channel), concurrent map panic |
| **14** | Advanced Concurrency | worker pool, errgroup, semaphore, rate limiter, fan-out/fan-in, pipeline |
| **15** | HTTP Server | Go 1.22+ mux (method+path routing), CRUD handlers, middleware chain, graceful shutdown |
| **16** | Testing | table-driven tests, t.Parallel, B.Loop() (1.24+), fuzz, httptest, TestMain |
| **17** | Logging (slog) | text/JSON handler, structured attrs, custom Handler |
| **18** | Generics | type params, constraints (~int), Stack[T], Map/Filter/Reduce, iter.Seq |
| **19** | Reflection | TypeOf/ValueOf, struct tags, mini validator, performance caveat |
| **20** | Unsafe Package | Sizeof/Alignof/Offsetof, struct layout, zero-copy string↔[]byte (1.20+) |
| **21** | CGo | inline C, C.CString/free, Go↔C string, CGo overhead |
| **22** | Design Patterns | functional options, builder, repository, DI, observer/eventbus, strategy, singleton, middleware |
| **23** | Profiling & Performance | pprof endpoints, escape analysis, strings.Builder, pre-alloc, sync.Pool, struct padding |
| **24** | Runtime Internals | GMP model, GOMAXPROCS, NumGoroutine, GC phases, MemStats, Green Tea GC (1.25+) |
| **25** | Go 1.22–1.25 Features | HTTP mux (1.22), loop var fix, iter.Seq (1.23), generic type alias (1.24), wg.Go (1.25) |

### Bonus Topics (26–34)

| # | Topic | Coverage |
|---|-------|----------|
| **26** | Common Mistakes | 14 BAD/GOOD pairs: goroutine leak, nil interface, mutex copy, concurrent map, loop capture, defer in loop, slice header copy, ignored errors, string concat in loop, empty interface overuse, time format, map order, init() side effects, no synchronization |
| **27** | Production Patterns | Circuit Breaker (Closed/Open/HalfOpen), Retry+Exponential Backoff, Config from Environment (12-Factor), Health Checks (readiness/liveness), Graceful Shutdown |
| **28** | Microservices & gRPC | gRPC concepts, Proto definition, server/client patterns, status codes, interceptors, REST vs gRPC tradeoffs |
| **29** | Database Patterns | database/sql pool, QueryContext/PreparedStatement, Transactions with defer rollback, Repository pattern, mocks |
| **30** | JSON Encoding | Marshal/Unmarshal, struct tags, custom marshaling, json.RawMessage, Encoder/Decoder (streaming), nullable fields |
| **31** | I/O & Files | Reader/Writer composition, bufio (Scanner, Reader, Writer), io.Copy/TeeReader, embed.FS, temp files/dirs |
| **32** | CLI Tools | flag package, subcommands, os/exec, stdin/stdout, shell safety |
| **33** | Advanced Testing | hand-written mocks, golden files, t.TempDir(), t.Cleanup(), table-driven+subtests, t.Parallel, benchmarks |
| **34** | Embed & Build Constraints | //go:embed (files, directories, patterns), //go:generate, //go:build (OS/arch constraints, custom tags), cross-compilation |

## Running Examples

### Run a single module
```bash
cd go/10_goroutines_channels
go run .
```

### Run tests
```bash
cd go/16_testing
go test ./... -v

cd go/33_testing_advanced
go test ./... -v
```

### Run with race detector
```bash
cd go/13_race_conditions
go run -race .
```

### Build for multiple platforms
```bash
cd go/34_embed_build_constraints
GOOS=linux GOARCH=amd64 go build -o myapp-linux-amd64 .
GOOS=darwin GOARCH=arm64 go build -o myapp-darwin-arm64 .
GOOS=windows GOARCH=amd64 go build -o myapp-windows-amd64.exe .
```

## Code Style

- **Comments**: Vietnamese, explaining concepts and gotchas
- **Function/Variable Names**: English (Go convention)
- **Examples**: Runnable code with output shown
- **Patterns**: Best practices, NGUYÊN TẮC (principles), CẢNH BÁO (warnings), GOTCHA notes

## Features Covered

- ✅ Loop variable per-iteration fix (Go 1.22)
- ✅ Range over integers (Go 1.22)
- ✅ Enhanced HTTP mux with method+path routing (Go 1.22)
- ✅ iter.Seq iterators (Go 1.23)
- ✅ strings.Lines/SplitSeq (Go 1.23)
- ✅ Timer fix — no drain needed (Go 1.23)
- ✅ Generic type aliases (Go 1.24)
- ✅ testing.B.Loop() (Go 1.24)
- ✅ Swiss Tables map implementation (Go 1.24)
- ✅ sync.WaitGroup.Go() (Go 1.25)
- ✅ Container-aware GOMAXPROCS (Go 1.25)
- ✅ Green Tea GC experiment (Go 1.25)
- ✅ testing.T.Attr() (Go 1.25)
- ✅ trace.NewFlightRecorder (Go 1.25)

## Prerequisites

- **Go 1.24.1+** (for full feature support)
- **go build**, **go run**, **go test** working

Optional for advanced topics:
- **CGo**: C compiler (gcc/clang)
- **gRPC**: `protoc`, `protoc-gen-go`, `protoc-gen-go-grpc` (for code generation)
- **Database**: PostgreSQL (or just study the patterns)

## Verification

All modules pass:
```bash
for d in */; do cd $d && go build ./... && go vet ./... && cd ..; done
```

Tests:
```bash
go test ./16_testing/...
go test ./33_testing_advanced/...
```

## Learning Path Suggestion

### Beginner (01–09)
Start with syntax, types, and interfaces. Good foundation for any Go developer.

### Intermediate (10–17)
Concurrency, testing, HTTP servers. Essential for production Go.

### Advanced (18–25)
Generics, reflection, runtime. For system programming and performance-critical code.

### Professional (26–34)
Production patterns, databases, testing, deployment. Ready for real-world applications.

## Directory Structure

```
go/
├── 01_hello_world/
├── 02_basic_syntax/
├── ...
├── 33_testing_advanced/
│   ├── go.mod
│   ├── main.go
│   ├── service.go
│   ├── service_test.go
│   └── testdata/
│       └── format_user.golden
├── 34_embed_build_constraints/
│   ├── go.mod
│   ├── main.go
│   ├── platform_unix.go
│   ├── platform_windows.go
│   ├── VERSION
│   ├── static/
│   │   ├── index.html
│   │   └── app.css
│   └── templates/
│       └── welcome.tmpl
└── README.md (this file)
```

Each numbered folder is a **standalone Go module** with its own `go.mod`.

## Contributing / Improvements

Want to add more topics or improve examples? The structure makes it easy:
1. Create a new numbered folder: `NN_topic_name/`
2. Add `go.mod` with `module github.com/tranminhquang/training-samples/go/NN_topic_name`
3. Write `main.go` and supporting files
4. Update this README

## Resources

- **Official**: [golang.org](https://golang.org), [pkg.go.dev](https://pkg.go.dev)
- **Books**: The Go Programming Language (Donovan & Kernighan)
- **References**: [Effective Go](https://go.dev/doc/effective_go)
- **Community**: [Go Forum](https://forum.golang.org), [r/golang](https://reddit.com/r/golang)

---

**Last updated**: May 2026  
**Go version**: 1.24.1+  
**Status**: Complete (34/34 modules)
