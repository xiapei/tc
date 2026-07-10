package rxutil

import (
	"testing"
)

func TestMatch(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		pattern string
		want    int // number of matched lines
	}{
		{"digits", "abc123\ndef456\nghi", `\d+`, 2},
		{"email", "test@gmail.com\nhello world\nfoo@bar.com", `@\w+\.com`, 2},
		{"no match", "hello\nworld", `\d+`, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Match([]byte(tt.input), tt.pattern)
			if err != nil {
				t.Fatalf("Match() error = %v", err)
			}
			lines := 0
			if len(result) > 0 {
				for _, c := range result {
					if c == '\n' {
						lines++
					}
				}
			}
			if lines != tt.want {
				t.Errorf("Match() got %d lines, want %d", lines, tt.want)
			}
		})
	}
}

func TestExtract(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		pattern string
		want    string
	}{
		{
			name:    "date extraction",
			input:   "2024-01-15 error occurred",
			pattern: `(\d{4}-\d{2}-\d{2})`,
			want:    "2024-01-15\n",
		},
		{
			name:    "key-value",
			input:   "name=John age=30",
			pattern: `(\w+)=(\w+)`,
			want:    "name\tJohn\nage\t30\n",
		},
		{
			name:    "no match",
			input:   "hello world",
			pattern: `(\d+)`,
			want:    "\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Extract([]byte(tt.input), tt.pattern)
			if err != nil {
				t.Fatalf("Extract() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("Extract() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestReplace(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		pattern     string
		replacement string
		want        string
	}{
		{"simple", "hello world", "world", "Go", "hello Go"},
		{"group", "foo123bar", `(\d+)`, "NUM:$1", "fooNUM:123bar"},
		{"global", "a1b2c3", `\d`, "X", "aXbXcX"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Replace([]byte(tt.input), tt.pattern, tt.replacement)
			if err != nil {
				t.Fatalf("Replace() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("Replace() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGrep(t *testing.T) {
	input := "line1 ERROR\nline2 INFO\nline3 ERROR\nline4 DEBUG"

	result, err := Grep([]byte(input), "ERROR", false)
	if err != nil {
		t.Fatalf("Grep() error = %v", err)
	}

	lines := 0
	for _, c := range result {
		if c == '\n' {
			lines++
		}
	}
	if lines != 2 {
		t.Errorf("Grep() got %d lines, want 2", lines)
	}

	// Test invert
	result, err = Grep([]byte(input), "ERROR", true)
	if err != nil {
		t.Fatalf("Grep(invert) error = %v", err)
	}
	lines = 0
	for _, c := range result {
		if c == '\n' {
			lines++
		}
	}
	if lines != 2 {
		t.Errorf("Grep(invert) got %d lines, want 2", lines)
	}
}

func TestCount(t *testing.T) {
	input := "ERROR: something\nINFO: ok\nERROR: again"
	count, err := Count([]byte(input), "ERROR")
	if err != nil {
		t.Fatalf("Count() error = %v", err)
	}
	if count != 2 {
		t.Errorf("Count() = %d, want 2", count)
	}
}

func TestFindAll(t *testing.T) {
	input := "a1b2c3"
	matches, err := FindAll([]byte(input), `\d`)
	if err != nil {
		t.Fatalf("FindAll() error = %v", err)
	}
	if len(matches) != 3 {
		t.Errorf("FindAll() got %d matches, want 3", len(matches))
	}
}
