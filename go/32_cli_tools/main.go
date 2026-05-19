// Lesson 32: CLI Tools with Go
// flag package, subcommands, os/exec, stdin/stdout pipelines
// Run: go run . [command] [flags]
// Examples:
//   go run . help
//   go run . greet -name Alice
//   go run . calc -op add -a 10 -b 20
package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func main() {
	// Subcommand pattern — like git, go, docker
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(0)
	}

	switch os.Args[1] {
	case "greet":
		runGreet(os.Args[2:])
	case "calc":
		runCalc(os.Args[2:])
	case "env":
		runEnv(os.Args[2:])
	case "run":
		runExec(os.Args[2:])
	case "help", "-h", "--help":
		printHelp()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		printHelp()
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println("Usage: cli [command] [flags]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  greet   Print greeting message")
	fmt.Println("  calc    Simple calculator")
	fmt.Println("  env     Show environment info")
	fmt.Println("  run     Execute a shell command")
	fmt.Println("  help    Show this help")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  go run . greet -name Alice -times 3")
	fmt.Println("  go run . calc -op add -a 10 -b 5")
	fmt.Println("  go run . env")
	fmt.Println("  go run . run -cmd 'echo hello'")
	fmt.Println()
	fmt.Println("=== flag package demo (running all subcommands) ===")
	demoAllSubcommands()
}

// ============================================================
// Subcommand: greet
// ============================================================

func runGreet(args []string) {
	fs := flag.NewFlagSet("greet", flag.ExitOnError)
	name := fs.String("name", "World", "Name to greet")
	times := fs.Int("times", 1, "Number of times to greet")
	upper := fs.Bool("upper", false, "Uppercase the greeting")

	fs.Parse(args)

	for range *times {
		msg := fmt.Sprintf("Hello, %s!", *name)
		if *upper {
			msg = strings.ToUpper(msg)
		}
		fmt.Println(msg)
	}
}

// ============================================================
// Subcommand: calc
// ============================================================

func runCalc(args []string) {
	fs := flag.NewFlagSet("calc", flag.ExitOnError)
	op := fs.String("op", "add", "Operation: add, sub, mul, div, sqrt")
	a := fs.Float64("a", 0, "First operand")
	b := fs.Float64("b", 0, "Second operand")

	fs.Parse(args)

	var result float64
	var err error

	switch *op {
	case "add":
		result = *a + *b
	case "sub":
		result = *a - *b
	case "mul":
		result = *a * *b
	case "div":
		if *b == 0 {
			err = fmt.Errorf("division by zero")
		} else {
			result = *a / *b
		}
	case "sqrt":
		if *a < 0 {
			err = fmt.Errorf("cannot sqrt negative number")
		} else {
			result = math.Sqrt(*a)
		}
	default:
		fmt.Fprintf(os.Stderr, "unknown operation: %s\n", *op)
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Result: %g\n", result)
}

// ============================================================
// Subcommand: env
// ============================================================

func runEnv(args []string) {
	fs := flag.NewFlagSet("env", flag.ExitOnError)
	verbose := fs.Bool("v", false, "Show all env vars")
	fs.Parse(args)

	fmt.Printf("OS: %s\n", runtime.GOOS)
	fmt.Printf("Arch: %s\n", runtime.GOARCH)
	fmt.Printf("CPUs: %d\n", runtime.NumCPU())
	fmt.Printf("Go version: %s\n", runtime.Version())
	fmt.Printf("Goroutines: %d\n", runtime.NumGoroutine())

	if *verbose {
		fmt.Println("\nEnvironment variables:")
		for _, env := range os.Environ() {
			fmt.Printf("  %s\n", env)
		}
	}
}

// ============================================================
// Subcommand: run (execute shell command)
// ============================================================

func runExec(args []string) {
	fs := flag.NewFlagSet("run", flag.ExitOnError)
	cmdFlag := fs.String("cmd", "", "Command to execute")
	timeout := fs.Int("timeout", 30, "Timeout in seconds")
	fs.Parse(args)

	if *cmdFlag == "" {
		fmt.Fprintln(os.Stderr, "error: -cmd is required")
		fs.Usage()
		os.Exit(1)
	}

	_ = timeout // not implemented in demo

	// Execute command
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", *cmdFlag)
	} else {
		cmd = exec.Command("sh", "-c", *cmdFlag)
	}

	// Capture output
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

// ============================================================
// Demo all subcommands in sequence (for go run . help)
// ============================================================

func demoAllSubcommands() {
	fmt.Println("--- greet ---")
	runGreet([]string{"-name", "Gopher", "-times", "2"})

	fmt.Println()
	fmt.Println("--- calc ---")
	for _, op := range []struct{ op, a, b string }{
		{"add", "10", "5"},
		{"mul", "7", "8"},
		{"div", "22", "7"},
		{"sqrt", "144", "0"},
	} {
		fmt.Printf("  %s(%s, %s) = ", op.op, op.a, op.b)
		a, _ := strconv.ParseFloat(op.a, 64)
		b, _ := strconv.ParseFloat(op.b, 64)
		switch op.op {
		case "add":
			fmt.Printf("%g\n", a+b)
		case "mul":
			fmt.Printf("%g\n", a*b)
		case "div":
			fmt.Printf("%g\n", a/b)
		case "sqrt":
			fmt.Printf("%g\n", math.Sqrt(a))
		}
	}

	fmt.Println()
	fmt.Println("--- env ---")
	runEnv([]string{})

	fmt.Println()
	fmt.Println("--- os/exec ---")
	demoExec()

	fmt.Println()
	fmt.Println("--- flag best practices ---")
	showFlagBestPractices()
}

func demoExec() {
	// exec.LookPath — find binary in PATH
	path, err := exec.LookPath("go")
	if err != nil {
		fmt.Printf("  go binary not found: %v\n", err)
	} else {
		fmt.Printf("  go binary: %s\n", path)
	}

	// exec.Command with captured output
	cmd := exec.Command("go", "version")
	out, err := cmd.Output() // captures stdout
	if err != nil {
		fmt.Printf("  exec error: %v\n", err)
	} else {
		fmt.Printf("  go version: %s", out)
	}

	// CombinedOutput — stdout + stderr
	cmd2 := exec.Command("go", "env", "GOPATH")
	combined, _ := cmd2.CombinedOutput()
	fmt.Printf("  GOPATH: %s", combined)

	// os/exec PRINCIPLE:
	fmt.Println("  CẢNH BÁO: không dùng shell injection!")
	fmt.Println("  BAD:  exec.Command(\"sh\", \"-c\", \"ls \" + userInput) // command injection!")
	fmt.Println("  GOOD: exec.Command(\"ls\", userInput) // args separated, safe")
}

func showFlagBestPractices() {
	fmt.Println("  flag package tips:")
	fmt.Println()
	fmt.Println("  1. Use flag.NewFlagSet for subcommands (don't share global flags)")
	fmt.Println("  2. flag.ExitOnError → program exits if flag parse fails")
	fmt.Println("  3. Set meaningful defaults, clear descriptions")
	fmt.Println("  4. Check flag.Args() for positional arguments")
	fmt.Println("  5. Use os.Stderr for errors, os.Stdout for output")
	fmt.Println()
	fmt.Println("  Alternatives:")
	fmt.Println("  - github.com/spf13/cobra: feature-rich CLI framework (kubectl, hugo)")
	fmt.Println("  - github.com/urfave/cli: simpler alternative")
	fmt.Println("  - github.com/alecthomas/kong: struct-based CLI parsing")
}
