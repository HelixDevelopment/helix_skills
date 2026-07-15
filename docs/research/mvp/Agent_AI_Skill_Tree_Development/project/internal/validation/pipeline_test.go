package validation

import "testing"

// ---------------------------------------------------------------------------
// extractCodeBlocks: pure fenced-code-block extraction from markdown, no I/O.
// ---------------------------------------------------------------------------

func TestExtractCodeBlocks(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    []codeSnippet
	}{
		{
			name:    "no fenced blocks yields nothing",
			content: "just plain prose with no code",
			want:    nil,
		},
		{
			name: "single fenced block with language tag",
			content: "Some text.\n" +
				"```go\n" +
				"package main\n" +
				"func main() {}\n" +
				"```\n" +
				"More text.",
			want: []codeSnippet{
				{Code: "package main\nfunc main() {}", Language: "go"},
			},
		},
		{
			name: "fenced block with no language tag",
			content: "```\n" +
				"echo hello\n" +
				"```",
			want: []codeSnippet{
				{Code: "echo hello", Language: ""},
			},
		},
		{
			name: "multiple fenced blocks preserve order and languages",
			content: "```python\n" +
				"print('a')\n" +
				"```\n" +
				"text between\n" +
				"```javascript\n" +
				"console.log('b')\n" +
				"```",
			want: []codeSnippet{
				{Code: "print('a')", Language: "python"},
				{Code: "console.log('b')", Language: "javascript"},
			},
		},
		{
			name: "unterminated fenced block is dropped (no closing fence)",
			content: "```go\n" +
				"package main\n",
			want: nil,
		},
		{
			name: "leading/trailing whitespace inside block is trimmed",
			content: "```go\n" +
				"\n\n  package main  \n\n" +
				"```",
			want: []codeSnippet{
				{Code: "package main", Language: "go"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractCodeBlocks(tt.content)
			if !codeSnippetsEqual(got, tt.want) {
				t.Errorf("extractCodeBlocks(%q) = %+v, want %+v", tt.content, got, tt.want)
			}
		})
	}
}

func codeSnippetsEqual(a, b []codeSnippet) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
