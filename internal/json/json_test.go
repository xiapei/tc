package jsonutil

import (
	"strings"
	"testing"
)

func TestFormat(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:  "simple object",
			input: `{"name":"John","age":30}`,
			want: `{
  "name": "John",
  "age": 30
}`,
		},
		{
			name:  "nested object",
			input: `{"user":{"name":"John","address":{"city":"Beijing"}}}`,
			want: `{
  "user": {
    "name": "John",
    "address": {
      "city": "Beijing"
    }
  }
}`,
		},
		{
			name:    "invalid json",
			input:   `{invalid}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Format([]byte(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("Format() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Format() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestMinify(t *testing.T) {
	input := `{
  "name": "John",
  "age": 30
}`
	got, err := Minify([]byte(input))
	if err != nil {
		t.Fatalf("Minify() error = %v", err)
	}
	want := `{"name":"John","age":30}`
	if got != want {
		t.Errorf("Minify() = %q, want %q", got, want)
	}
}

func TestGet(t *testing.T) {
	data := []byte(`{"user":{"name":"John","emails":["a@b.com","c@d.com"]}}`)

	tests := []struct {
		path string
		want string
	}{
		{"user.name", "John"},
		{"user.emails.0", "a@b.com"},
		{"user.emails.1", "c@d.com"},
		{"user.age", "null"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got, err := Get(data, tt.path)
			if err != nil {
				t.Fatalf("Get() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("Get(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}

func TestFilter(t *testing.T) {
	data := []byte(`[{"name":"Alice","age":25},{"name":"Bob","age":17},{"name":"Charlie","age":30}]`)

	tests := []struct {
		expr    string
		wantLen int
	}{
		{"age > 18", 2},
		{"age < 20", 1},
		{"name == \"Alice\"", 1},
		{"name ~ \"li\"", 2},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			got, err := Filter(data, tt.expr)
			if err != nil {
				t.Fatalf("Filter() error = %v", err)
			}
			// Count objects in result
			count := 0
			for _, c := range got {
				if c == '{' {
					count++
				}
			}
			if count != tt.wantLen {
				t.Errorf("Filter(%q) returned %d items, want %d", tt.expr, count, tt.wantLen)
			}
		})
	}
}

func TestKeys(t *testing.T) {
	data := []byte(`{"name":"John","age":30,"city":"Beijing"}`)
	got, err := Keys(data)
	if err != nil {
		t.Fatalf("Keys() error = %v", err)
	}
	// Should contain all three keys
	for _, key := range []string{"name", "age", "city"} {
		if !contains(got, key) {
			t.Errorf("Keys() = %q, missing %q", got, key)
		}
	}
}

func TestPaths(t *testing.T) {
	data := []byte(`{"a":{"b":1},"c":[2,3]}`)
	paths, err := Paths(data)
	if err != nil {
		t.Fatalf("Paths() error = %v", err)
	}
	expected := []string{"a", "a.b", "c", "c.0", "c.1"}
	if len(paths) != len(expected) {
		t.Errorf("Paths() returned %d paths, want %d", len(paths), len(expected))
	}
}

func TestDiff(t *testing.T) {
	old := []byte(`{"name":"John","age":30}`)
	new := []byte(`{"name":"John","age":31,"city":"Beijing"}`)
	result, err := Diff(old, new)
	if err != nil {
		t.Fatalf("Diff() error = %v", err)
	}
	if result == "无差异\n" {
		t.Error("Diff() should detect age change")
	}
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func TestFilter_Compound(t *testing.T) {
	data := []byte(`[{"name":"Alice","age":25,"status":"active"},{"name":"Bob","age":17,"status":"active"},{"name":"Charlie","age":30,"status":"inactive"}]`)

	tests := []struct {
		name    string
		expr    string
		wantLen int
	}{
		{"AND simple", `age > 18 && status == "active"`, 1},
		{"OR simple", `age > 30 || age < 18`, 1},
		{"OR match all", `age > 18 || status == "inactive"`, 2},
		{"AND no match", `age > 100 && status == "active"`, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Filter(data, tt.expr)
			if err != nil {
				t.Fatalf("Filter() error = %v", err)
			}
			count := 0
			for _, c := range got {
				if c == '{' {
					count++
				}
			}
			if count != tt.wantLen {
				t.Errorf("Filter(%q) returned %d items, want %d", tt.expr, count, tt.wantLen)
			}
		})
	}
}

func TestSet(t *testing.T) {
	data := []byte(`{"name":"John","age":30}`)

	result, err := Set(data, "name", "Jane")
	if err != nil {
		t.Fatalf("Set() error = %v", err)
	}
	if !strings.Contains(result, "Jane") {
		t.Errorf("Set() = %q, expected to contain Jane", result)
	}

	result, err = Set(data, "age", "25")
	if err != nil {
		t.Fatalf("Set() error = %v", err)
	}
	if !strings.Contains(result, "25") {
		t.Errorf("Set() = %q, expected to contain 25", result)
	}

	data = []byte(`{"user":{"name":"John"}}`)
	result, err = Set(data, "user.name", "Jane")
	if err != nil {
		t.Fatalf("Set() error = %v", err)
	}
	if !strings.Contains(result, "Jane") {
		t.Errorf("Set() = %q, expected to contain Jane", result)
	}
}

func TestDelete(t *testing.T) {
	data := []byte(`{"name":"John","age":30,"city":"Beijing"}`)
	result, err := Delete(data, "age")
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	if strings.Contains(result, "age") {
		t.Errorf("Delete() = %q, should not contain age", result)
	}
	if !strings.Contains(result, "name") || !strings.Contains(result, "city") {
		t.Errorf("Delete() = %q, should contain name and city", result)
	}
}

func TestMerge(t *testing.T) {
	base := []byte(`{"name":"John","age":30}`)
	patch := []byte(`{"age":31,"city":"Beijing"}`)
	result, err := Merge(base, patch)
	if err != nil {
		t.Fatalf("Merge() error = %v", err)
	}
	if !strings.Contains(result, "31") {
		t.Errorf("Merge() = %q, expected age=31", result)
	}
	if !strings.Contains(result, "Beijing") {
		t.Errorf("Merge() = %q, expected city=Beijing", result)
	}
}

func TestMerge_Nested(t *testing.T) {
	base := []byte(`{"user":{"name":"John","age":30},"version":1}`)
	patch := []byte(`{"user":{"city":"Beijing"},"version":2}`)
	result, err := Merge(base, patch)
	if err != nil {
		t.Fatalf("Merge() error = %v", err)
	}
	if !strings.Contains(result, "John") {
		t.Errorf("Merge() lost existing nested field 'name'")
	}
	if !strings.Contains(result, "Beijing") {
		t.Errorf("Merge() missing new nested field 'city'")
	}
	if !strings.Contains(result, "2") {
		t.Errorf("Merge() expected version=2")
	}
}
