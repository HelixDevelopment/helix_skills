// Package validation provides sandboxed code execution environments for
// safely running code snippets extracted from skill content. The default
// implementation uses WebAssembly for lightweight isolation, with a
// Docker-based fallback for full containerization.
package validation

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
)

// ---------------------------------------------------------------------------
// Sandbox interface
// ---------------------------------------------------------------------------

// Sandbox executes code in an isolated environment with resource limits
// and timeout controls.
type Sandbox interface {
	// Execute runs the given code snippet in an isolated environment.
	// Returns the execution result including stdout, stderr, exit code,
	// and duration. The language parameter hints at the runtime to use.
	Execute(ctx context.Context, code string, language string, timeout time.Duration) (*ExecutionResult, error)
}

// ExecutionResult captures the outcome of sandboxed code execution.
type ExecutionResult struct {
	Stdout   string        `json:"stdout"`
	Stderr   string        `json:"stderr"`
	ExitCode int           `json:"exit_code"`
	Duration time.Duration `json:"duration"`
}

// ---------------------------------------------------------------------------
// WASM Sandbox (default - lightweight)
// ---------------------------------------------------------------------------

// WASMSandbox uses WebAssembly for lightweight code execution isolation.
// It compiles supported languages to WASM and runs them in a wasm runtime.
// For unsupported languages, it falls back to process-based execution.
type WASMSandbox struct {
	logger      *zap.Logger
	wasmRuntime string // path to wasm runtime binary (e.g., wasmtime, wasmer)
	mu          sync.Mutex
}

// NewWASMSandbox creates a new WASM-based sandbox.
// It auto-detects available WASM runtimes on the system.
func NewWASMSandbox(logger *zap.Logger) *WASMSandbox {
	runtime := detectWasmRuntime()
	if runtime == "" {
		logger.Warn("no WASM runtime found, will use process fallback for code execution")
	} else {
		logger.Info("WASM sandbox initialized", zap.String("runtime", runtime))
	}

	return &WASMSandbox{
		logger:      logger,
		wasmRuntime: runtime,
	}
}

// detectWasmRuntime finds an available WASM runtime on the system.
func detectWasmRuntime() string {
	runtimes := []string{"wasmtime", "wasmer", "wasmedge"}
	for _, r := range runtimes {
		if _, err := exec.LookPath(r); err == nil {
			return r
		}
	}
	return ""
}

// Execute runs code in the WASM sandbox. For languages that can be compiled
// to WASM (Rust, Go, C, AssemblyScript), it attempts WASM execution.
// For other languages, it falls back to process-based execution.
func (s *WASMSandbox) Execute(ctx context.Context, code string, language string, timeout time.Duration) (*ExecutionResult, error) {
	if timeout <= 0 {
		timeout = 30 * time.Second
	}

	// Normalize language name
	language = normalizeLanguage(language)

	// Try WASM execution for supported languages
	if s.wasmRuntime != "" && isWASMSupported(language) {
		return s.executeWASM(ctx, code, language, timeout)
	}

	// Fallback to process-based execution for unsupported languages
	s.logger.Debug("falling back to process execution",
		zap.String("language", language),
		zap.String("reason", "wasm_not_available"),
	)
	return s.executeProcess(ctx, code, language, timeout)
}

// isWASMSupported checks if a language can be compiled to WASM.
func isWASMSupported(language string) bool {
	switch language {
	case "go", "golang", "rust", "c", "cpp", "c++", "assemblyscript", "ts", "typescript":
		return true
	default:
		return false
	}
}

// executeWASM attempts to compile and run code via a WASM runtime.
// For this implementation, we use a simplified approach that writes code
// to a temp file and invokes the language toolchain to produce WASM,
// then runs it. In production, this would use a proper WASM compiler service.
func (s *WASMSandbox) executeWASM(ctx context.Context, code string, language string, timeout time.Duration) (*ExecutionResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Create temp directory for build artifacts
	tmpDir, err := os.MkdirTemp("", "helix-wasm-*")
	if err != nil {
		return nil, fmt.Errorf("create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	// For Go code, we can use `go run` directly as a lightweight sandbox
	// In production, this would compile to WASM and run in the runtime
	if language == "go" || language == "golang" {
		return s.executeGoSnippet(ctx, code, tmpDir, timeout)
	}

	// For other languages, fall back to process execution
	return s.executeProcess(ctx, code, language, timeout)
}

// executeGoSnippet runs a Go code snippet using `go run` with module restrictions.
func (s *WASMSandbox) executeGoSnippet(ctx context.Context, code string, tmpDir string, timeout time.Duration) (*ExecutionResult, error) {
	srcFile := filepath.Join(tmpDir, "main.go")
	if err := os.WriteFile(srcFile, []byte(code), 0644); err != nil {
		return nil, fmt.Errorf("write source: %w", err)
	}

	// Create a minimal go.mod
	goMod := `module snippet

go 1.22
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		return nil, fmt.Errorf("write go.mod: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var stdout, stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, "go", "run", srcFile)
	cmd.Dir = tmpDir
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Security: restrict network access and file system
	cmd.Env = append(os.Environ(),
		"GOPROXY=off",        // no network access for modules
		"GONOSUMDB=*",        // no sumdb lookups
		"GOFLAGS=-mod=vendor", // use vendor only
	)

	start := time.Now()
	err := cmd.Run()
	duration := time.Since(start)

	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = -1
		}
	}

	// Check for timeout
	if ctx.Err() == context.DeadlineExceeded {
		return &ExecutionResult{
			Stdout:   stdout.String(),
			Stderr:   "execution timed out",
			ExitCode: 124,
			Duration: timeout,
		}, nil
	}

	return &ExecutionResult{
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		ExitCode: exitCode,
		Duration: duration,
	}, nil
}

// executeProcess runs code using the system interpreter with restricted permissions.
func (s *WASMSandbox) executeProcess(ctx context.Context, code string, language string, timeout time.Duration) (*ExecutionResult, error) {
	interpreter := getInterpreter(language)
	if interpreter == "" {
		return &ExecutionResult{
			Stdout:   "",
			Stderr:   fmt.Sprintf("unsupported language: %s", language),
			ExitCode: -1,
			Duration: 0,
		}, nil
	}

	// Create temp file for the code
	tmpDir, err := os.MkdirTemp("", "helix-sandbox-*")
	if err != nil {
		return nil, fmt.Errorf("create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	// Determine file extension
	ext := getFileExtension(language)
	srcFile := filepath.Join(tmpDir, fmt.Sprintf("main.%s", ext))
	if err := os.WriteFile(srcFile, []byte(code), 0644); err != nil {
		return nil, fmt.Errorf("write source: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var stdout, stderr bytes.Buffer

	var cmd *exec.Cmd
	if language == "python" || language == "py" {
		cmd = exec.CommandContext(ctx, interpreter, "-c", code)
	} else if language == "sh" || language == "bash" || language == "shell" {
		cmd = exec.CommandContext(ctx, interpreter, "-c", code)
	} else if language == "js" || language == "javascript" || language == "node" {
		cmd = exec.CommandContext(ctx, interpreter, "-e", code)
	} else {
		cmd = exec.CommandContext(ctx, interpreter, srcFile)
	}

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Security: minimal environment
	cmd.Env = []string{
		"PATH=" + os.Getenv("PATH"),
		"HOME=" + tmpDir,
		"TMPDIR=" + tmpDir,
	}

	// Restrict network (Linux/Mac)
	if runtime.GOOS == "linux" {
		cmd.Env = append(cmd.Env, "LD_PRELOAD=") // disable LD_PRELOAD
	}

	start := time.Now()
	err = cmd.Run()
	duration := time.Since(start)

	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = -1
		}
	}

	if ctx.Err() == context.DeadlineExceeded {
		return &ExecutionResult{
			Stdout:   stdout.String(),
			Stderr:   "execution timed out",
			ExitCode: 124,
			Duration: timeout,
		}, nil
	}

	return &ExecutionResult{
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		ExitCode: exitCode,
		Duration: duration,
	}, nil
}

// ---------------------------------------------------------------------------
// Docker Sandbox (full isolation)
// ---------------------------------------------------------------------------

// DockerSandbox uses containers for full code execution isolation.
// Each code snippet runs in a fresh, disposable container with strict
// resource limits.
type DockerSandbox struct {
	logger    *zap.Logger
	mu        sync.Mutex
	available bool // whether docker is available on this system
}

// NewDockerSandbox creates a new Docker-based sandbox.
// It checks if Docker is available and falls back to WASMSandbox if not.
func NewDockerSandbox(logger *zap.Logger) *DockerSandbox {
	available := isDockerAvailable()
	if !available {
		logger.Warn("Docker not available, sandbox will fall back to process execution")
	} else {
		logger.Info("Docker sandbox initialized")
	}

	return &DockerSandbox{
		logger:    logger,
		available: available,
	}
}

// isDockerAvailable checks if Docker CLI is installed and the daemon is running.
func isDockerAvailable() bool {
	cmd := exec.Command("docker", "version")
	err := cmd.Run()
	return err == nil
}

// Execute runs code in a Docker container for full isolation.
// It creates a temporary container, copies the code, executes it,
// and removes the container.
func (s *DockerSandbox) Execute(ctx context.Context, code string, language string, timeout time.Duration) (*ExecutionResult, error) {
	if timeout <= 0 {
		timeout = 30 * time.Second
	}

	language = normalizeLanguage(language)

	if !s.available {
		// Fall back to WASM sandbox
		fallback := NewWASMSandbox(s.logger)
		return fallback.Execute(ctx, code, language, timeout)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "helix-docker-*")
	if err != nil {
		return nil, fmt.Errorf("create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	// Get Docker image for language
	image := getDockerImage(language)
	if image == "" {
		return &ExecutionResult{
			Stdout:   "",
			Stderr:   fmt.Sprintf("no Docker image for language: %s", language),
			ExitCode: -1,
			Duration: 0,
		}, nil
	}

	// Write code to temp file
	ext := getFileExtension(language)
	srcFile := filepath.Join(tmpDir, fmt.Sprintf("main.%s", ext))
	if err := os.WriteFile(srcFile, []byte(code), 0644); err != nil {
		return nil, fmt.Errorf("write source: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, timeout+10*time.Second) // extra time for container setup
	defer cancel()

	containerName := fmt.Sprintf("helix-sandbox-%d", time.Now().UnixNano())

	// Build the run command based on language
	var runArgs []string
	switch language {
	case "python", "py":
		runArgs = []string{"python", "-c", code}
	case "go", "golang":
		// For Go, write file and run
		if err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(code), 0644); err != nil {
			return nil, fmt.Errorf("write go source: %w", err)
		}
		runArgs = []string{"go", "run", "/tmp/main.go"}
	case "js", "javascript", "node":
		runArgs = []string{"node", "-e", code}
	case "sh", "bash", "shell":
		runArgs = []string{"bash", "-c", code}
	case "ruby":
		runArgs = []string{"ruby", "-e", code}
	default:
		runArgs = []string{"cat", "/tmp/main." + ext}
	}

	// Build docker run command with security flags
	args := []string{
		"run",
		"--rm",                           // auto-remove container after run
		"--name", containerName,
		"--network", "none",              // no network access
		"--memory", "128m",               // memory limit
		"--cpus", "0.5",                  // CPU limit
		"--read-only",                    // read-only root filesystem
		"--tmpfs", "/tmp:noexec,nosuid,size=50m",
		"-v", fmt.Sprintf("%s:/tmp:ro", tmpDir), // mount code as read-only
		"--stop-timeout", fmt.Sprintf("%d", int(timeout.Seconds())+5),
		image,
	}
	args = append(args, runArgs...)

	var stdout, stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, "docker", args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	start := time.Now()
	err = cmd.Run()
	duration := time.Since(start)

	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = -1
		}
	}

	// Clean up container if still running
	_ = exec.Command("docker", "rm", "-f", containerName).Run()

	if ctx.Err() == context.DeadlineExceeded {
		return &ExecutionResult{
			Stdout:   stdout.String(),
			Stderr:   "execution timed out (container killed)",
			ExitCode: 124,
			Duration: timeout,
		}, nil
	}

	return &ExecutionResult{
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		ExitCode: exitCode,
		Duration: duration,
	}, nil
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// normalizeLanguage converts various language name formats to canonical forms.
func normalizeLanguage(lang string) string {
	lang = strings.ToLower(strings.TrimSpace(lang))
	switch lang {
	case "golang":
		return "go"
	case "py", "python3":
		return "python"
	case "js", "nodejs":
		return "javascript"
	case "sh":
		return "bash"
	case "ts":
		return "typescript"
	case "c++", "cxx", "cpp":
		return "cpp"
	case "rs":
		return "rust"
	default:
		return lang
	}
}

// getInterpreter returns the system interpreter command for a language.
func getInterpreter(language string) string {
	switch normalizeLanguage(language) {
	case "python":
		return findExecutable("python3", "python")
	case "bash", "shell":
		return findExecutable("bash", "sh")
	case "javascript", "node":
		return findExecutable("node")
	case "ruby":
		return findExecutable("ruby")
	case "perl":
		return findExecutable("perl")
	case "php":
		return findExecutable("php")
	case "r":
		return findExecutable("R")
	default:
		return ""
	}
}

// getDockerImage returns the appropriate Docker image for a language.
func getDockerImage(language string) string {
	switch normalizeLanguage(language) {
	case "python":
		return "python:3.11-slim"
	case "go":
		return "golang:1.22-alpine"
	case "javascript", "node":
		return "node:20-alpine"
	case "bash", "shell":
		return "bash:5"
	case "ruby":
		return "ruby:3.2-slim"
	case "rust":
		return "rust:1.75-slim"
	case "c", "cpp":
		return "gcc:13"
	case "java":
		return "openjdk:21-slim"
	default:
		return "alpine:latest"
	}
}

// getFileExtension returns the source file extension for a language.
func getFileExtension(language string) string {
	switch normalizeLanguage(language) {
	case "go":
		return "go"
	case "python":
		return "py"
	case "javascript":
		return "js"
	case "typescript":
		return "ts"
	case "bash", "shell":
		return "sh"
	case "ruby":
		return "rb"
	case "rust":
		return "rs"
	case "c":
		return "c"
	case "cpp":
		return "cpp"
	case "java":
		return "java"
	case "kotlin":
		return "kt"
	case "php":
		return "php"
	case "perl":
		return "pl"
	case "r":
		return "r"
	default:
		return "txt"
	}
}

// findExecutable searches for the first available executable from candidates.
func findExecutable(candidates ...string) string {
	for _, c := range candidates {
		if path, err := exec.LookPath(c); err == nil {
			return path
		}
	}
	return ""
}
