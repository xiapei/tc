package cmd

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/spf13/pflag"
)

// resetCommandsFlags 重置所有子命令的 flag 值到默认值，防止跨测试干扰
func resetCommandsFlags() {
	for _, cmd := range rootCmd.Commands() {
		cmd.Flags().VisitAll(func(f *pflag.Flag) {
			f.Value.Set(f.DefValue)
			f.Changed = false
		})
	}
}

// runCmd 执行命令并捕获 stdout，同时通过 stdin 提供输入
func runCmd(args []string, stdin string) (string, error) {
	resetCommandsFlags()

	oldStdin := os.Stdin
	oldStdout := os.Stdout

	// stdin pipe
	stdinR, stdinW, err := os.Pipe()
	if err != nil {
		return "", err
	}
	stdinW.Write([]byte(stdin))
	stdinW.Close()
	os.Stdin = stdinR

	// stdout pipe
	stdoutR, stdoutW, err := os.Pipe()
	if err != nil {
		return "", err
	}
	os.Stdout = stdoutW

	rootCmd.SetArgs(args)
	execErr := rootCmd.Execute()

	stdoutW.Close()
	var buf bytes.Buffer
	io.Copy(&buf, stdoutR)

	os.Stdin = oldStdin
	os.Stdout = oldStdout
	stdinR.Close()
	stdoutR.Close()

	return strings.TrimRight(buf.String(), "\n"), execErr
}

func TestHashCmd(t *testing.T) {
	tests := []struct {
		name string
		args []string
		stdin string
		want string
	}{
		{"sha256", []string{"hash", "sha256"}, "hello", "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"},
		{"md5", []string{"hash", "md5"}, "hello", "5d41402abc4b2a76b9719d911017c592"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := runCmd(tt.args, tt.stdin)
			if err != nil {
				t.Fatalf("runCmd() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("hash got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestHashCmd_InvalidAlgo(t *testing.T) {
	_, err := runCmd([]string{"hash", "sha1"}, "hello")
	if err == nil {
		t.Error("expected error for unsupported algorithm")
	}
}

func TestEncCmd(t *testing.T) {
	got, err := runCmd([]string{"enc", "base64"}, "hello world")
	if err != nil {
		t.Fatalf("runCmd() error = %v", err)
	}
	want := "aGVsbG8gd29ybGQ="
	if got != want {
		t.Errorf("enc base64 got %q, want %q", got, want)
	}
}

func TestEncCmd_URL(t *testing.T) {
	got, err := runCmd([]string{"enc", "url"}, "hello world")
	if err != nil {
		t.Fatalf("runCmd() error = %v", err)
	}
	want := "hello+world"
	if got != want {
		t.Errorf("enc url got %q, want %q", got, want)
	}
}

func TestDecCmd(t *testing.T) {
	got, err := runCmd([]string{"dec", "base64"}, "aGVsbG8gd29ybGQ=")
	if err != nil {
		t.Fatalf("runCmd() error = %v", err)
	}
	want := "hello world"
	if got != want {
		t.Errorf("dec base64 got %q, want %q", got, want)
	}
}

func TestDecCmd_URL(t *testing.T) {
	got, err := runCmd([]string{"dec", "url"}, "hello+world")
	if err != nil {
		t.Fatalf("runCmd() error = %v", err)
	}
	if got != "hello world" {
		t.Errorf("dec url got %q, want %q", got, "hello world")
	}
}

func TestStatsCmd(t *testing.T) {
	got, err := runCmd([]string{"stats"}, "hello world\nfoo bar\nbaz")
	if err != nil {
		t.Fatalf("runCmd() error = %v", err)
	}
	// Format: "  lines\twords\tchars\tbytes"
	wantPrefix := "  3"
	if !strings.HasPrefix(got, wantPrefix) {
		t.Errorf("stats got %q, expected prefix %q", got, wantPrefix)
	}
	if !strings.Contains(got, "\t") {
		t.Errorf("stats output should be tab-separated, got %q", got)
	}
}

func TestStatsCmd_Empty(t *testing.T) {
	got, err := runCmd([]string{"stats"}, "")
	if err != nil {
		t.Fatalf("runCmd() error = %v", err)
	}
	want := "  0\t0\t0\t0"
	if got != want {
		t.Errorf("stats empty got %q, want %q", got, want)
	}
}

func TestSortCmd(t *testing.T) {
	got, err := runCmd([]string{"sort"}, "c\nb\na")
	if err != nil {
		t.Fatalf("runCmd() error = %v", err)
	}
	want := "a\nb\nc"
	if got != want {
		t.Errorf("sort got %q, want %q", got, want)
	}
}

func TestSortCmd_Numeric(t *testing.T) {
	got, err := runCmd([]string{"sort", "-n"}, "10\n2\n30\n1")
	if err != nil {
		t.Fatalf("runCmd() error = %v", err)
	}
	want := "1\n2\n10\n30"
	if got != want {
		t.Errorf("sort -n got %q, want %q", got, want)
	}
}

func TestSortCmd_Reverse(t *testing.T) {
	got, err := runCmd([]string{"sort", "--reverse"}, "a\nb\nc")
	if err != nil {
		t.Fatalf("runCmd() error = %v", err)
	}
	want := "c\nb\na"
	if got != want {
		t.Errorf("sort --reverse got %q, want %q", got, want)
	}
}

func TestUniqCmd(t *testing.T) {
	got, err := runCmd([]string{"uniq"}, "a\na\nb\nc\nc")
	if err != nil {
		t.Fatalf("runCmd() error = %v", err)
	}
	want := "a\nb\nc"
	if got != want {
		t.Errorf("uniq got %q, want %q", got, want)
	}
}

func TestUniqCmd_Count(t *testing.T) {
	got, err := runCmd([]string{"uniq", "--count"}, "a\na\nb\nc\nc\nc")
	if err != nil {
		t.Fatalf("runCmd() error = %v", err)
	}
	if !strings.Contains(got, "2\ta") {
		t.Errorf("uniq -c expected '2\\ta', got %q", got)
	}
	if !strings.Contains(got, "1\tb") {
		t.Errorf("uniq -c expected '1\\tb', got %q", got)
	}
	if !strings.Contains(got, "3\tc") {
		t.Errorf("uniq -c expected '3\\tc', got %q", got)
	}
}

func TestHeadCmd(t *testing.T) {
	got, err := runCmd([]string{"head", "2"}, "a\nb\nc\nd\ne")
	if err != nil {
		t.Fatalf("runCmd() error = %v", err)
	}
	want := "a\nb"
	if got != want {
		t.Errorf("head 2 got %q, want %q", got, want)
	}
}

func TestHeadCmd_Default(t *testing.T) {
	lines := ""
	for i := 1; i <= 20; i++ {
		lines += string(rune('0'+i)) + "\n"
	}
	lines = "a\nb\nc\nd\ne\nf\ng\nh\ni\nj\nk\nl\nm\nn"
	got, err := runCmd([]string{"head"}, lines)
	if err != nil {
		t.Fatalf("runCmd() error = %v", err)
	}
	count := strings.Count(got, "\n") + 1
	if count != 10 {
		t.Errorf("head default should output 10 lines, got %d", count)
	}
}

func TestTailCmd(t *testing.T) {
	got, err := runCmd([]string{"tail", "2"}, "a\nb\nc\nd\ne")
	if err != nil {
		t.Fatalf("runCmd() error = %v", err)
	}
	want := "d\ne"
	if got != want {
		t.Errorf("tail 2 got %q, want %q", got, want)
	}
}

func TestSampleCmd(t *testing.T) {
	got, err := runCmd([]string{"sample", "3"}, "a\nb\nc\nd\ne")
	if err != nil {
		t.Fatalf("runCmd() error = %v", err)
	}
	count := strings.Count(got, "\n") + 1
	if count != 3 {
		t.Errorf("sample 3 should output 3 lines, got %d", count)
	}
}

func TestSampleCmd_Invalid(t *testing.T) {
	_, err := runCmd([]string{"sample", "0"}, "a\nb\nc")
	if err == nil {
		t.Error("expected error for invalid sample count")
	}
}

func TestCountCmd(t *testing.T) {
	got, err := runCmd([]string{"count"}, "a\nb\nc")
	if err != nil {
		t.Fatalf("runCmd() error = %v", err)
	}
	if got != "3" {
		t.Errorf("count got %q, want %q", got, "3")
	}
}

func TestCountCmd_Empty(t *testing.T) {
	got, err := runCmd([]string{"count"}, "")
	if err != nil {
		t.Fatalf("runCmd() error = %v", err)
	}
	if got != "0" {
		t.Errorf("count empty got %q, want %q", got, "0")
	}
}

func TestFieldsCmd(t *testing.T) {
	got, err := runCmd([]string{"fields", "2", "--sep", ":"}, "a:b:c:d")
	if err != nil {
		t.Fatalf("runCmd() error = %v", err)
	}
	want := "b"
	if got != want {
		t.Errorf("fields got %q, want %q", got, want)
	}
}

func TestFmtCmd(t *testing.T) {
	got, err := runCmd([]string{"fmt"}, `{"name":"John","age":30}`)
	if err != nil {
		t.Fatalf("runCmd() error = %v", err)
	}
	if !strings.Contains(got, `"name"`) {
		t.Errorf("fmt output should contain fields, got %q", got)
	}
	if !strings.Contains(got, "\n") {
		t.Errorf("fmt should format with newlines, got %q", got)
	}
}
