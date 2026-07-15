package validation

import "testing"

// ---------------------------------------------------------------------------
// normalizeLanguage: pure string canonicalization, no I/O.
// ---------------------------------------------------------------------------

func TestNormalizeLanguage(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"golang", "go"},
		{"GOLANG", "go"},
		{"  golang  ", "go"},
		{"py", "python"},
		{"python3", "python"},
		{"js", "javascript"},
		{"nodejs", "javascript"},
		{"sh", "bash"},
		{"ts", "typescript"},
		{"c++", "cpp"},
		{"cxx", "cpp"},
		{"cpp", "cpp"},
		{"rs", "rust"},
		{"go", "go"},
		{"ruby", "ruby"},
		{"unknown-lang", "unknown-lang"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			if got := normalizeLanguage(tt.in); got != tt.want {
				t.Errorf("normalizeLanguage(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// isWASMSupported: pure allow-list check.
// ---------------------------------------------------------------------------

func TestIsWASMSupported(t *testing.T) {
	tests := []struct {
		lang string
		want bool
	}{
		{"go", true},
		{"golang", true},
		{"rust", true},
		{"c", true},
		{"cpp", true},
		{"c++", true},
		{"assemblyscript", true},
		{"ts", true},
		{"typescript", true},
		{"python", false},
		{"javascript", false},
		{"ruby", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.lang, func(t *testing.T) {
			if got := isWASMSupported(tt.lang); got != tt.want {
				t.Errorf("isWASMSupported(%q) = %v, want %v", tt.lang, got, tt.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// getDockerImage: pure language -> Docker image mapping (normalizes first).
// ---------------------------------------------------------------------------

func TestGetDockerImage(t *testing.T) {
	tests := []struct {
		lang string
		want string
	}{
		{"python", "python:3.11-slim"},
		{"py", "python:3.11-slim"}, // normalized first
		{"go", "golang:1.22-alpine"},
		{"golang", "golang:1.22-alpine"},
		{"javascript", "node:20-alpine"},
		{"node", "node:20-alpine"},
		{"nodejs", "node:20-alpine"},
		{"bash", "bash:5"},
		{"sh", "bash:5"},
		{"ruby", "ruby:3.2-slim"},
		{"rust", "rust:1.75-slim"},
		{"rs", "rust:1.75-slim"},
		{"c", "gcc:13"},
		{"cpp", "gcc:13"},
		{"c++", "gcc:13"},
		{"java", "openjdk:21-slim"},
		{"totally-unknown-language", "alpine:latest"},
	}

	for _, tt := range tests {
		t.Run(tt.lang, func(t *testing.T) {
			if got := getDockerImage(tt.lang); got != tt.want {
				t.Errorf("getDockerImage(%q) = %q, want %q", tt.lang, got, tt.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// getFileExtension: pure language -> extension mapping (normalizes first).
// ---------------------------------------------------------------------------

func TestGetFileExtension(t *testing.T) {
	tests := []struct {
		lang string
		want string
	}{
		{"go", "go"},
		{"golang", "go"},
		{"python", "py"},
		{"py", "py"},
		{"javascript", "js"},
		{"js", "js"},
		{"typescript", "ts"},
		{"ts", "ts"},
		{"bash", "sh"},
		{"sh", "sh"},
		{"ruby", "rb"},
		{"rust", "rs"},
		{"rs", "rs"},
		{"c", "c"},
		{"cpp", "cpp"},
		{"c++", "cpp"},
		{"java", "java"},
		{"kotlin", "kt"},
		{"php", "php"},
		{"perl", "pl"},
		{"r", "r"},
		{"totally-unknown-language", "txt"},
	}

	for _, tt := range tests {
		t.Run(tt.lang, func(t *testing.T) {
			if got := getFileExtension(tt.lang); got != tt.want {
				t.Errorf("getFileExtension(%q) = %q, want %q", tt.lang, got, tt.want)
			}
		})
	}
}
