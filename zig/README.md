# Zig Training Samples

Practice projects for learning [Zig](https://ziglang.org/).

## Projects

### [hello_zig](hello_zig/)

A starter project covering Zig fundamentals:

- Basic project structure with `build.zig` and `build.zig.zon`
- Module system: separating business logic (`root.zig`) from CLI entry point (`main.zig`)
- Buffered stdout writer
- Unit testing with `std.testing`
- Fuzz testing
- Memory leak detection with `std.testing.allocator`

**Requirements:** Zig `>= 0.15.2`

**Usage:**

```bash
cd hello_zig

# Build and run
zig build run

# Run tests
zig build test

# Run fuzz tests
zig build test --fuzz
```
