package portutil

import (
	"fmt"
	"net"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

// Result 端口检查结果
type Result struct {
	Port      int
	Available bool
	PID       int
	Process   string
}

// Check 检查端口占用情况
func Check(port int) Result {
	r := Result{Port: port}

	// 1. 尝试监听端口判断是否可用
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err == nil {
		ln.Close()
		r.Available = true
		return r
	}
	r.Available = false

	// 2. Windows 下尝试查 netstat 找进程
	if runtime.GOOS == "windows" {
		pid, proc := findProcessWindows(port)
		r.PID = pid
		r.Process = proc
	} else {
		pid, proc := findProcessUnix(port)
		r.PID = pid
		r.Process = proc
	}

	return r
}

func findProcessWindows(port int) (int, string) {
	out, err := exec.Command("netstat", "-ano").Output()
	if err != nil {
		return 0, ""
	}

	re := regexp.MustCompile(fmt.Sprintf(`:%d\s+.*?LISTENING\s+(\d+)`, port))
	matches := re.FindStringSubmatch(string(out))
	if len(matches) < 2 {
		return 0, ""
	}

	pid, _ := strconv.Atoi(matches[1])
	procName := getProcessName(pid)
	return pid, procName
}

func findProcessUnix(port int) (int, string) {
	cmd := exec.Command("lsof", "-i", fmt.Sprintf(":%d", port), "-sTCP:LISTEN", "-P", "-n")
	out, err := cmd.Output()
	if err != nil {
		return 0, ""
	}

	lines := strings.Split(string(out), "\n")
	if len(lines) < 2 {
		return 0, ""
	}

	fields := strings.Fields(lines[1])
	if len(fields) >= 2 {
		pid, _ := strconv.Atoi(fields[1])
		return pid, fields[0]
	}
	return 0, ""
}

func getProcessName(pid int) string {
	out, err := exec.Command("tasklist", "/FI", fmt.Sprintf("PID eq %d", pid), "/NH", "/FO", "CSV").Output()
	if err != nil {
		return ""
	}
	// 输出格式: "映像名称","PID","会话名","会话#","内存使用"
	parts := strings.Split(string(out), ",")
	if len(parts) >= 1 {
		return strings.Trim(parts[0], `" `)
	}
	return ""
}

// Format 格式化输出
func Format(r Result) string {
	if r.Available {
		return fmt.Sprintf("✓ 端口 %d 可用", r.Port)
	}

	desc := fmt.Sprintf("✗ 端口 %d 被占用", r.Port)
	if r.PID > 0 {
		desc += fmt.Sprintf(" (PID: %d", r.PID)
		if r.Process != "" {
			desc += fmt.Sprintf(", 进程: %s", r.Process)
		}
		desc += ")"
	}
	return desc
}
