package validation

import (
	"context"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/helixdevelopment/skill-system/internal/models"
)

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
// G02/G16: the static code check NEVER executes untrusted code.
//
// Paired-mutation-real: this test fails the moment any host-execution path is
// reintroduced into staticCheckCode — the sentinel files it "would" create when
// run must NOT exist after the (non-executing) check.
// ---------------------------------------------------------------------------

func TestStaticCheckCode_NeverExecutesUntrustedCode(t *testing.T) {
	dir := t.TempDir()
	goSentinel := filepath.Join(dir, "go_pwned")
	shSentinel := filepath.Join(dir, "sh_pwned")

	// A valid Go program that WOULD create a file if executed.
	goCode := "package main\n\nimport \"os\"\n\nfunc main() {\n\t_ = os.WriteFile(\"" + goSentinel + "\", []byte(\"x\"), 0644)\n}\n"
	// A shell command that WOULD create a file if executed.
	shCode := "echo pwned > " + shSentinel

	snippets := []codeSnippet{
		{Code: goCode, Language: "go"},
		{Code: shCode, Language: "bash"},
	}

	rep := staticCheckCode(snippets)

	if _, err := os.Stat(goSentinel); !os.IsNotExist(err) {
		t.Fatalf("go snippet was EXECUTED: sentinel %s exists (err=%v)", goSentinel, err)
	}
	if _, err := os.Stat(shSentinel); !os.IsNotExist(err) {
		t.Fatalf("shell snippet was EXECUTED: sentinel %s exists (err=%v)", shSentinel, err)
	}

	if rep.Total != 2 {
		t.Errorf("Total = %d, want 2", rep.Total)
	}
	if rep.Checked != 1 { // go parsed in-process
		t.Errorf("Checked = %d, want 1", rep.Checked)
	}
	if rep.Unchecked != 1 { // bash has no in-process front-end
		t.Errorf("Unchecked = %d, want 1", rep.Unchecked)
	}
}

// ---------------------------------------------------------------------------
// Tier B: the default IsolatedExecutor always SKIPs (fail-closed), never runs.
// ---------------------------------------------------------------------------

func TestSkipIsolatedExecutor_AlwaysSkips(t *testing.T) {
	res, err := NewSkipIsolatedExecutor().Execute(context.Background(), "rm -rf /", "bash", time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Status != StageSkip {
		t.Errorf("Status = %q, want %q", res.Status, StageSkip)
	}
	if res.Reason != "isolation_runtime_absent" {
		t.Errorf("Reason = %q, want %q", res.Reason, "isolation_runtime_absent")
	}
	if res.Stdout != "" || res.Stderr != "" {
		t.Errorf("expected no output from a skip, got stdout=%q stderr=%q", res.Stdout, res.Stderr)
	}
}

// ---------------------------------------------------------------------------
// G21: SSRF egress guard blocks internal / metadata targets.
// ---------------------------------------------------------------------------

func TestIsBlockedIP(t *testing.T) {
	tests := []struct {
		ip      string
		blocked bool
	}{
		{"169.254.169.254", true}, // cloud metadata (link-local)
		{"169.254.10.1", true},    // link-local
		{"127.0.0.1", true},       // loopback
		{"::1", true},             // loopback v6
		{"10.1.2.3", true},        // RFC1918
		{"192.168.1.1", true},     // RFC1918
		{"172.16.5.5", true},      // RFC1918
		{"0.0.0.0", true},         // unspecified
		{"fe80::1", true},         // link-local v6
		{"fc00::1", true},         // unique-local v6
		{"224.0.0.1", true},       // multicast
		{"100.100.100.200", true}, // CGNAT — Alibaba Cloud metadata endpoint
		{"100.64.0.1", true},      // CGNAT (RFC 6598)
		{"240.0.0.1", true},       // reserved (Class E)
		{"192.0.0.1", true},       // IETF protocol assignments
		{"198.18.0.1", true},      // benchmarking (RFC 2544)
		{"255.255.255.255", true}, // limited broadcast
		{"64:ff9b::1", true},      // NAT64 well-known prefix
		{"8.8.8.8", false},        // public
		{"1.1.1.1", false},        // public
		{"93.184.216.34", false},  // public (example.com)
	}
	for _, tt := range tests {
		t.Run(tt.ip, func(t *testing.T) {
			got, reason := isBlockedIP(net.ParseIP(tt.ip))
			if got != tt.blocked {
				t.Errorf("isBlockedIP(%s) = %v (%q), want %v", tt.ip, got, reason, tt.blocked)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// G03/G05: fail-closed overall-verdict aggregation.
//
// Paired-mutation-real: a SKIP/FAIL/BLOCKED in any stage must never yield PASS.
// ---------------------------------------------------------------------------

func TestComputeOverallVerdict(t *testing.T) {
	tests := []struct {
		name   string
		stages map[string]StageStatus
		want   bool
	}{
		{"empty is fail-closed", map[string]StageStatus{}, false},
		{"all pass", map[string]StageStatus{"a": StagePass, "b": StagePass}, true},
		{"pass plus na", map[string]StageStatus{"a": StagePass, "b": StageNA}, true},
		{"a SKIP never passes", map[string]StageStatus{"a": StagePass, "b": StageSkip}, false},
		{"a FAIL never passes", map[string]StageStatus{"a": StagePass, "b": StageFail}, false},
		{"a BLOCKED never passes", map[string]StageStatus{"a": StagePass, "b": StageBlocked}, false},
		{"all N/A has no positive", map[string]StageStatus{"a": StageNA, "b": StageNA}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := computeOverallVerdict(tt.stages); got != tt.want {
				t.Errorf("computeOverallVerdict(%v) = %v, want %v", tt.stages, got, tt.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// G03 request-path: a client can never self-promote without a passing verdict.
// ---------------------------------------------------------------------------

func TestDecideCreateStatus(t *testing.T) {
	pass := &ValidationResult{Passed: true}
	fail := &ValidationResult{Passed: false}

	tests := []struct {
		name      string
		enabled   bool
		requested models.SkillStatus
		res       *ValidationResult
		want      models.SkillStatus
	}{
		{"disabled -> draft", false, models.SkillStatusActive, nil, models.SkillStatusDraft},
		{"nil result -> draft", true, models.SkillStatusActive, nil, models.SkillStatusDraft},
		{"failed -> draft even if active requested", true, models.SkillStatusActive, fail, models.SkillStatusDraft},
		{"passed + draft requested -> validated", true, models.SkillStatusDraft, pass, models.SkillStatusValidated},
		{"passed + active requested -> active", true, models.SkillStatusActive, pass, models.SkillStatusActive},
		{"passed + no request -> validated", true, "", pass, models.SkillStatusValidated},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DecideCreateStatus(tt.enabled, tt.requested, tt.res); got != tt.want {
				t.Errorf("DecideCreateStatus(%v, %q, %v) = %q, want %q",
					tt.enabled, tt.requested, tt.res, got, tt.want)
			}
		})
	}
}
