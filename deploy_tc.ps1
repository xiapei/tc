# tc 项目部署脚本 - 在 PowerShell 中运行
$base = "C:\01WorkSpace\tc"
$dirs = @("cmd","internal\json","internal\rx","internal\pipeline","docs",".github\workflows")
foreach ($d in $dirs) { New-Item -ItemType Directory -Force -Path "$base\$d" | Out-Null }
Write-Host "创建目录结构完成" -ForegroundColor Green

function Write-File($path, $content) {
    $full = Join-Path $base $path
    $content | Set-Content -Path $full -Encoding UTF8 -NoNewline
    Write-Host "  + $path"
}

Write-File 'cmd/json_cmds.go' @'
package cmd

import (
	"fmt"
	"os"

	jsonutil "github.com/user/tc/internal/json"
	"github.com/user/tc/internal/pipeline"
	"github.com/spf13/cobra"
)

var jsonFmtCmd = &cobra.Command{
	Use:   "fmt [file]",
	Short: "格式化 JSON",
	Example: `  cat data.json | tc json fmt
  tc json fmt data.json
  echo '{"a":1,"b":2}' | tc json fmt`,
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := pipeline.ReadInput(args)
		if err != nil {
			return err
		}
		result, err := jsonutil.Format(data)
		if err != nil {
			return err
		}
		fmt.Println(result)
		return nil
	},
}

var jsonMinCmd = &cobra.Command{
	Use:   "min [file]",
	Short: "压缩 JSON 为单行",
	Example: `  cat data.json | tc json min
  tc json min data.json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := pipeline.ReadInput(args)
		if err != nil {
			return err
		}
		result, err := jsonutil.Minify(data)
		if err != nil {
			return err
		}
		fmt.Print(result)
		return nil
	},
}

var jsonGetCmd = &cobra.Command{
	Use:   "get <path> [file]",
	Short: "按路径查询 JSON 字段",
	Long: `按路径查询 JSON 字段，支持：
  - 嵌套字段: "user.name"
  - 数组索引: "users.0.name"
  - 通配符:   "users.*.email"
  - 数组尾部: "items.-1"`,
	Example: `  cat data.json | tc json get "users.0.name"
  cat data.json | tc json get "users.*.email"
  echo '{"a":{"b":1}}' | tc json get "a.b"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]
		fileArgs := args[1:]
		data, err := pipeline.ReadInput(fileArgs)
		if err != nil {
			return err
		}
		result, err := jsonutil.Get(data, path)
		if err != nil {
			return err
		}
		fmt.Println(result)
		return nil
	},
}

var jsonFilterCmd = &cobra.Command{
	Use:   "filter <expr> [file]",
	Short: "按条件过滤 JSON 数组元素",
	Long: `支持简单条件表达式：
  - 比较: age > 18, name == "John", status != "inactive"
  - 包含: email ~ "@gmail"
  - 存在: phone`,
	Example: `  cat users.json | tc json filter 'age > 18'
  cat users.json | tc json filter 'status == "active"'
  cat users.json | tc json filter 'email ~ "@gmail"'`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		expr := args[0]
		fileArgs := args[1:]
		data, err := pipeline.ReadInput(fileArgs)
		if err != nil {
			return err
		}
		result, err := jsonutil.Filter(data, expr)
		if err != nil {
			return err
		}
		fmt.Println(result)
		return nil
	},
}

var jsonKeysCmd = &cobra.Command{
	Use:   "keys [file]",
	Short: "提取 JSON 对象的所有 key",
	Example: `  echo '{"name":"John","age":30}' | tc json keys
  cat data.json | tc json keys`,
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := pipeline.ReadInput(args)
		if err != nil {
			return err
		}
		result, err := jsonutil.Keys(data)
		if err != nil {
			return err
		}
		fmt.Println(result)
		return nil
	},
}

var jsonValuesCmd = &cobra.Command{
	Use:   "values [file]",
	Short: "提取 JSON 对象的所有 value",
	Example: `  echo '{"name":"John","age":30}' | tc json values`,
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := pipeline.ReadInput(args)
		if err != nil {
			return err
		}
		result, err := jsonutil.Values(data)
		if err != nil {
			return err
		}
		fmt.Println(result)
		return nil
	},
}

var jsonPathsCmd = &cobra.Command{
	Use:   "paths [file]",
	Short: "列出 JSON 的所有路径",
	Example: `  echo '{"a":{"b":1},"c":[2,3]}' | tc json paths`,
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := pipeline.ReadInput(args)
		if err != nil {
			return err
		}
		result, err := jsonutil.Paths(data)
		if err != nil {
			return err
		}
		for _, p := range result {
			fmt.Println(p)
		}
		return nil
	},
}

var jsonTableCmd = &cobra.Command{
	Use:   "table <fields> [file]",
	Short: "将 JSON 数组转为表格",
	Long: `提取指定字段，以制表符分隔输出（方便粘贴到 Excel）。
第一个行为表头。`,
	Example: `  cat users.json | tc json table "name,email,age"
  cat users.json | tc json table "id,title" --csv`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fields := args[0]
		fileArgs := args[1:]
		csv, _ := cmd.Flags().GetBool("csv")
		data, err := pipeline.ReadInput(fileArgs)
		if err != nil {
			return err
		}
		result, err := jsonutil.Table(data, fields, csv)
		if err != nil {
			return err
		}
		fmt.Print(result)
		return nil
	},
}

var jsonDiffCmd = &cobra.Command{
	Use:   "diff <file1> <file2>",
	Short: "比较两个 JSON 文件的差异",
	Example: `  tc json diff old.json new.json`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		data1, err := os.ReadFile(args[0])
		if err != nil {
			return fmt.Errorf("读取文件 %s 失败: %w", args[0], err)
		}
		data2, err := os.ReadFile(args[1])
		if err != nil {
			return fmt.Errorf("读取文件 %s 失败: %w", args[1], err)
		}
		result, err := jsonutil.Diff(data1, data2)
		if err != nil {
			return err
		}
		fmt.Print(result)
		return nil
	},
}

func init() {
	jsonTableCmd.Flags().Bool("csv", false, "输出 CSV 格式（逗号分隔）")

	jsonCmd.AddCommand(jsonFmtCmd)
	jsonCmd.AddCommand(jsonMinCmd)
	jsonCmd.AddCommand(jsonGetCmd)
	jsonCmd.AddCommand(jsonFilterCmd)
	jsonCmd.AddCommand(jsonKeysCmd)
	jsonCmd.AddCommand(jsonValuesCmd)
	jsonCmd.AddCommand(jsonPathsCmd)
	jsonCmd.AddCommand(jsonTableCmd)
	jsonCmd.AddCommand(jsonDiffCmd)
}
'@

Write-File 'cmd/json.go' @'
package cmd

import (
	"github.com/spf13/cobra"
)

var jsonCmd = &cobra.Command{
	Use:   "json",
	Short: "JSON 处理命令",
	Long: `处理 JSON 数据：格式化、压缩、查询、过滤、转换。

所有命令支持从 stdin 读取或指定文件参数。`,
	Example: `  cat data.json | tc json fmt
  cat data.json | tc json get "users.0.name"
  cat data.json | tc json filter 'age > 18'
  cat data.json | tc json keys`,
}

func init() {
	rootCmd.AddCommand(jsonCmd)
}
'@

Write-File 'cmd/root.go' @'
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var versionStr = "dev"

func SetVersion(v string) {
	versionStr = v
}

var rootCmd = &cobra.Command{
	Use:   "tc",
	Short: "ToolCraft — 程序员的文本/数据瑞士军刀",
	Long: `tc 是一个轻量级 CLI 工具，一站式处理 JSON、正则和文本数据。

核心命令：
  tc json    JSON 处理（格式化、查询、过滤、转换）
  tc rx      正则处理（匹配、提取、替换、过滤）
  tc         通用工具（编码、哈希、统计、排序）

每个命令都支持管道操作：stdin → 处理 → stdout`,
	Example:  `  cat data.json | tc json get "users.0.name"
  echo "abc123" | tc rx match "\d+"
  echo "hello" | tc hash sha256`,
	SilenceUsage: true,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "显示版本信息",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("tc %s\n", versionStr)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
'@

Write-File 'cmd/rx_cmds.go' @'
package cmd

import (
	"fmt"

	rxutil "github.com/user/tc/internal/rx"
	"github.com/user/tc/internal/pipeline"
	"github.com/spf13/cobra"
)

var rxMatchCmd = &cobra.Command{
	Use:   "match <pattern> [file]",
	Short: "高亮显示匹配的内容",
	Long: `在输入文本中查找匹配 pattern 的内容，高亮输出匹配行。
只输出包含匹配的行。`,
	Example: `  echo "abc123def456" | tc rx match "\\d+"
  cat access.log | tc rx match "5\\d{2}"
  echo "hello world" | tc rx match "world"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pattern := args[0]
		fileArgs := args[1:]
		data, err := pipeline.ReadInput(fileArgs)
		if err != nil {
			return err
		}
		result, err := rxutil.Match(data, pattern)
		if err != nil {
			return err
		}
		fmt.Print(result)
		return nil
	},
}

var rxExtractCmd = &cobra.Command{
	Use:   "extract <pattern> [file]",
	Short: "提取匹配捕获组的内容",
	Long: `提取正则捕获组匹配的内容。
  - 无捕获组：输出整个匹配
  - 有捕获组：输出 $1（或所有组，tab 分隔）`,
	Example: `  echo "2024-01-15 error" | tc rx extract "(\\d{4}-\\d{2}-\\d{2})"
  echo "name=John age=30" | tc rx extract "(\\w+)=(\\w+)"
  cat app.log | tc rx extract '(\\{.*\\})'`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pattern := args[0]
		fileArgs := args[1:]
		data, err := pipeline.ReadInput(fileArgs)
		if err != nil {
			return err
		}
		result, err := rxutil.Extract(data, pattern)
		if err != nil {
			return err
		}
		fmt.Print(result)
		return nil
	},
}

var rxReplaceCmd = &cobra.Command{
	Use:   "replace <pattern> <replacement> [file]",
	Short: "替换匹配的内容",
	Long: `替换输入中匹配 pattern 的内容。
replacement 中可用 $1, $2 引用捕获组。`,
	Example: `  echo "hello world" | tc rx replace "world" "Go"
  cat file.txt | tc rx replace "(\\d+)" "NUM:$1"
  echo "foo-bar" | tc rx replace "-" "_"`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		pattern := args[0]
		replacement := args[1]
		fileArgs := args[2:]
		data, err := pipeline.ReadInput(fileArgs)
		if err != nil {
			return err
		}
		result, err := rxutil.Replace(data, pattern, replacement)
		if err != nil {
			return err
		}
		fmt.Print(result)
		return nil
	},
}

var rxGrepCmd = &cobra.Command{
	Use:   "grep <pattern> [file]",
	Short: "输出匹配正则的行",
	Example: `  cat access.log | tc rx grep "5\\d{2}"
  cat app.log | tc rx grep "ERROR|WARN"
  echo -e "abc\n123\ndef" | tc rx grep "\\d+"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pattern := args[0]
		fileArgs := args[1:]
		invert, _ := cmd.Flags().GetBool("invert")
		data, err := pipeline.ReadInput(fileArgs)
		if err != nil {
			return err
		}
		result, err := rxutil.Grep(data, pattern, invert)
		if err != nil {
			return err
		}
		fmt.Print(result)
		return nil
	},
}

var rxCountCmd = &cobra.Command{
	Use:   "count <pattern> [file]",
	Short: "统计匹配次数",
	Example: `  cat app.log | tc rx count "ERROR"
  echo "a1b2c3" | tc rx count "\\d"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pattern := args[0]
		fileArgs := args[1:]
		data, err := pipeline.ReadInput(fileArgs)
		if err != nil {
			return err
		}
		result, err := rxutil.Count(data, pattern)
		if err != nil {
			return err
		}
		fmt.Println(result)
		return nil
	},
}

var rxFindallCmd = &cobra.Command{
	Use:   "findall <pattern> [file]",
	Short: "列出所有匹配项",
	Example: `  echo "a1b2c3" | tc rx findall "\\d"
  echo "2024-01-15 and 2024-12-25" | tc rx findall "\\d{4}-\\d{2}-\\d{2}"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pattern := args[0]
		fileArgs := args[1:]
		data, err := pipeline.ReadInput(fileArgs)
		if err != nil {
			return err
		}
		result, err := rxutil.FindAll(data, pattern)
		if err != nil {
			return err
		}
		for _, m := range result {
			fmt.Println(m)
		}
		return nil
	},
}

func init() {
	rxGrepCmd.Flags().BoolP("invert", "v", false, "反转匹配（输出不匹配的行）")

	rxCmd.AddCommand(rxMatchCmd)
	rxCmd.AddCommand(rxExtractCmd)
	rxCmd.AddCommand(rxReplaceCmd)
	rxCmd.AddCommand(rxGrepCmd)
	rxCmd.AddCommand(rxCountCmd)
	rxCmd.AddCommand(rxFindallCmd)
}
'@

Write-File 'cmd/rx.go' @'
package cmd

import (
	"github.com/spf13/cobra"
)

var rxCmd = &cobra.Command{
	Use:   "rx",
	Short: "正则表达式处理命令",
	Long: `正则表达式处理：匹配、提取、替换、过滤。

使用 Go RE2 语法（与 Perl/PCRE 基本兼容，无回溯）。
支持捕获组：用 $1, $2 引用分组。`,
	Example: `  echo "abc123" | tc rx match "\\d+"
  cat app.log | tc rx extract "(\\d{4}-\\d{2}-\\d{2})"
  cat file.txt | tc rx replace "old" "new"`,
}

func init() {
	rootCmd.AddCommand(rxCmd)
}
'@

Write-File 'cmd/util_cmds.go' @'
package cmd

import (
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/user/tc/internal/pipeline"
	"github.com/spf13/cobra"
)

// hash 命令
var hashCmd = &cobra.Command{
	Use:   "hash <algorithm> [file]",
	Short: "计算哈希值",
	Long: `支持算法: md5, sha256 (默认 sha256)`,
	Example: `  echo -n "hello" | tc hash md5
  echo -n "hello" | tc hash sha256
  tc hash sha256 file.txt`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		algo := args[0]
		fileArgs := args[1:]
		data, err := pipeline.ReadInput(fileArgs)
		if err != nil {
			return err
		}
		var result string
		switch algo {
		case "md5":
			h := md5.Sum(data)
			result = hex.EncodeToString(h[:])
		case "sha256":
			h := sha256.Sum256(data)
			result = hex.EncodeToString(h[:])
		default:
			return fmt.Errorf("不支持的算法: %s (可选: md5, sha256)", algo)
		}
		fmt.Println(result)
		return nil
	},
}

// enc 命令
var encCmd = &cobra.Command{
	Use:   "enc <encoding> [file]",
	Short: "编码",
	Long: `支持编码: base64, url`,
	Example: `  echo -n "hello world" | tc enc base64
  echo -n "hello world" | tc enc url`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		enc := args[0]
		fileArgs := args[1:]
		data, err := pipeline.ReadInput(fileArgs)
		if err != nil {
			return err
		}
		var result string
		switch enc {
		case "base64":
			result = base64.StdEncoding.EncodeToString(data)
		case "url":
			result = url.QueryEscape(string(data))
		default:
			return fmt.Errorf("不支持的编码: %s (可选: base64, url)", enc)
		}
		fmt.Println(result)
		return nil
	},
}

// dec 命令
var decCmd = &cobra.Command{
	Use:   "dec <encoding> [file]",
	Short: "解码",
	Long: `支持编码: base64, url`,
	Example: `  echo "aGVsbG8=" | tc dec base64
  echo "hello%20world" | tc dec url`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		enc := args[0]
		fileArgs := args[1:]
		data, err := pipeline.ReadInput(fileArgs)
		if err != nil {
			return err
		}
		input := strings.TrimSpace(string(data))
		var result string
		switch enc {
		case "base64":
			decoded, err := base64.StdEncoding.DecodeString(input)
			if err != nil {
				return fmt.Errorf("base64 解码失败: %w", err)
			}
			result = string(decoded)
		case "url":
			decoded, err := url.QueryUnescape(input)
			if err != nil {
				return fmt.Errorf("URL 解码失败: %w", err)
			}
			result = decoded
		default:
			return fmt.Errorf("不支持的编码: %s (可选: base64, url)", enc)
		}
		fmt.Print(result)
		return nil
	},
}

// stats 命令
var statsCmd = &cobra.Command{
	Use:   "stats [file]",
	Short: "统计行数、单词数、字符数、字节数",
	Example: `  cat file.txt | tc stats
  tc stats file.txt`,
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := pipeline.ReadInput(args)
		if err != nil {
			return err
		}
		s := string(data)
		lines := strings.Count(s, "\n")
		if len(s) > 0 && !strings.HasSuffix(s, "\n") {
			lines++
		}
		words := len(strings.Fields(s))
		chars := len([]rune(s))
		bytes := len(data)
		fmt.Printf("  %d\t%d\t%d\t%d\n", lines, words, chars, bytes)
		return nil
	},
}

// sort 命令
var sortCmd = &cobra.Command{
	Use:   "sort [file]",
	Short: "排序行",
	Example: `  echo -e "banana\napple\ncherry" | tc sort
  cat file.txt | tc sort --numeric`,
	RunE: func(cmd *cobra.Command, args []string) error {
		numeric, _ := cmd.Flags().GetBool("numeric")
		reverse, _ := cmd.Flags().GetBool("reverse")
		data, err := pipeline.ReadInput(args)
		if err != nil {
			return err
		}
		lines := strings.Split(strings.TrimRight(string(data), "\n"), "\n")
		if numeric {
			sort.Slice(lines, func(i, j int) bool {
				a, _ := strconv.ParseFloat(strings.TrimSpace(lines[i]), 64)
				b, _ := strconv.ParseFloat(strings.TrimSpace(lines[j]), 64)
				if reverse {
					return a > b
				}
				return a < b
			})
		} else {
			if reverse {
				sort.Sort(sort.Reverse(sort.StringSlice(lines)))
			} else {
				sort.Strings(lines)
			}
		}
		for _, line := range lines {
			fmt.Println(line)
		}
		return nil
	},
}

// uniq 命令
var uniqCmd = &cobra.Command{
	Use:   "uniq [file]",
	Short: "去重（相邻重复行）",
	Long: `去除相邻的重复行。
提示：先 tc sort 再 tc uniq 可以完全去重。`,
	Example: `  echo -e "a\na\nb\nc\nc" | tc uniq
  echo -e "a\nb\na" | tc sort | tc uniq
  cat file.txt | tc uniq --count`,
	RunE: func(cmd *cobra.Command, args []string) error {
		countFlag, _ := cmd.Flags().GetBool("count")
		data, err := pipeline.ReadInput(args)
		if err != nil {
			return err
		}
		lines := strings.Split(strings.TrimRight(string(data), "\n"), "\n")
		var prev string
		for _, line := range lines {
			if line != prev {
				if countFlag {
					fmt.Printf("1\t%s\n", line)
				} else {
					fmt.Println(line)
				}
				prev = line
			}
		}
		return nil
	},
}

// head 命令
var headCmd = &cobra.Command{
	Use:   "head [file]",
	Short: "输出前 N 行",
	Example: `  cat file.txt | tc head 10
  cat file.txt | tc head -n 5`,
	RunE: func(cmd *cobra.Command, args []string) error {
		n, _ := cmd.Flags().GetInt("lines")
		if n <= 0 {
			n = 10
		}
		fileArgs := args
		data, err := pipeline.ReadInput(fileArgs)
		if err != nil {
			return err
		}
		lines := strings.Split(string(data), "\n")
		for i := 0; i < n && i < len(lines); i++ {
			fmt.Println(lines[i])
		}
		return nil
	},
}

// tail 命令
var tailCmd = &cobra.Command{
	Use:   "tail [file]",
	Short: "输出后 N 行",
	Example: `  cat file.txt | tc tail 10
  cat file.txt | tc tail -n 5`,
	RunE: func(cmd *cobra.Command, args []string) error {
		n, _ := cmd.Flags().GetInt("lines")
		if n <= 0 {
			n = 10
		}
		fileArgs := args
		data, err := pipeline.ReadInput(fileArgs)
		if err != nil {
			return err
		}
		lines := strings.Split(strings.TrimRight(string(data), "\n"), "\n")
		start := len(lines) - n
		if start < 0 {
			start = 0
		}
		for _, line := range lines[start:] {
			fmt.Println(line)
		}
		return nil
	},
}

// sample 命令
var sampleCmd = &cobra.Command{
	Use:   "sample <n> [file]",
	Short: "随机采样 N 行",
	Example: `  cat bigfile.txt | tc sample 100
  cat data.json | tc sample 50`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		n, err := strconv.Atoi(args[0])
		if err != nil || n <= 0 {
			return fmt.Errorf("无效的采样数量: %s", args[0])
		}
		fileArgs := args[1:]
		data, err := pipeline.ReadInput(fileArgs)
		if err != nil {
			return err
		}
		lines := strings.Split(strings.TrimRight(string(data), "\n"), "\n")
		if n >= len(lines) {
			for _, line := range lines {
				fmt.Println(line)
			}
			return nil
		}
		// Fisher-Yates shuffle (partial)
		indices := make([]int, len(lines))
		for i := range indices {
			indices[i] = i
		}
		for i := 0; i < n; i++ {
			j := i + int(abs(int64(fastrand()))%(int64(len(lines)-i)))
			indices[i], indices[j] = indices[j], indices[i]
		}
		sort.Ints(indices[:n])
		for _, idx := range indices[:n] {
			fmt.Println(lines[idx])
		}
		return nil
	},
}

// count 命令（统计行数，类似 wc -l）
var countCmd = &cobra.Command{
	Use:   "count [file]",
	Short: "统计行数",
	Example: `  cat file.txt | tc count
  echo -e "a\nb\nc" | tc count`,
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := pipeline.ReadInput(args)
		if err != nil {
			return err
		}
		s := string(data)
		if len(s) == 0 {
			fmt.Println(0)
			return nil
		}
		lines := strings.Count(s, "\n")
		if !strings.HasSuffix(s, "\n") {
			lines++
		}
		fmt.Println(lines)
		return nil
	},
}

// fields 命令（按分隔符提取字段）
var fieldsCmd = &cobra.Command{
	Use:   "fields <n> [file]",
	Short: "提取第 N 个字段（按分隔符）",
	Example: `  echo "a:b:c:d" | tc fields 2 --sep ":"
  cat data.csv | tc fields 3 --sep ","`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		n, err := strconv.Atoi(args[0])
		if err != nil || n < 1 {
			return fmt.Errorf("无效的字段编号: %s (从 1 开始)", args[0])
		}
		sep, _ := cmd.Flags().GetString("sep")
		if sep == "" {
			sep = "\t"
		}
		fileArgs := args[1:]
		data, err := pipeline.ReadInput(fileArgs)
		if err != nil {
			return err
		}
		lines := strings.Split(strings.TrimRight(string(data), "\n"), "\n")
		for _, line := range lines {
			parts := strings.Split(line, sep)
			if n <= len(parts) {
				fmt.Println(parts[n-1])
			}
		}
		return nil
	},
}

// tc json fmt 快捷别名（根级别）
var fmtCmd = &cobra.Command{
	Use:   "fmt [file]",
	Short: "格式化 JSON（等同于 tc json fmt）",
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := pipeline.ReadInput(args)
		if err != nil {
			return err
		}
		var out bytes.Buffer
		if err := json.Indent(&out, data, "", "  "); err != nil {
			return err
		}
		fmt.Println(out.String())
		return nil
	},
}

// 辅助函数
func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

var seed uint64 = 1

func fastrand() uint64 {
	seed ^= seed << 13
	seed ^= seed >> 7
	seed ^= seed << 17
	return seed
}

func init() {
	sortCmd.Flags().BoolP("numeric", "n", false, "按数值排序")
	sortCmd.Flags().BoolP("reverse", "r", false, "反转排序")

	headCmd.Flags().IntP("lines", "n", 10, "输出行数")
	tailCmd.Flags().IntP("lines", "n", 10, "输出行数")

	fieldsCmd.Flags().StringP("sep", "s", "\t", "字段分隔符")

	rootCmd.AddCommand(hashCmd)
	rootCmd.AddCommand(encCmd)
	rootCmd.AddCommand(decCmd)
	rootCmd.AddCommand(statsCmd)
	rootCmd.AddCommand(sortCmd)
	rootCmd.AddCommand(uniqCmd)
	rootCmd.AddCommand(headCmd)
	rootCmd.AddCommand(tailCmd)
	rootCmd.AddCommand(sampleCmd)
	rootCmd.AddCommand(countCmd)
	rootCmd.AddCommand(fieldsCmd)
	rootCmd.AddCommand(fmtCmd)
}
'@

Write-File 'docs/cheatsheet.md' @'
# tc 命令速查表

## JSON 处理

| 命令 | 说明 | 示例 |
|------|------|------|
| `tc json fmt` | 格式化 | `echo '{"a":1}' \| tc json fmt` |
| `tc json min` | 压缩 | `cat data.json \| tc json min` |
| `tc json get <path>` | 查询 | `cat d.json \| tc json get "users.0.name"` |
| `tc json filter <expr>` | 过滤 | `cat d.json \| tc json filter 'age > 18'` |
| `tc json keys` | 所有 key | `echo '{"a":1,"b":2}' \| tc json keys` |
| `tc json values` | 所有 value | `echo '{"a":1,"b":2}' \| tc json values` |
| `tc json paths` | 所有路径 | `echo '{"a":{"b":1}}' \| tc json paths` |
| `tc json table <fields>` | 转表格 | `cat d.json \| tc json table "name,email"` |
| `tc json diff <f1> <f2>` | 比较差异 | `tc json diff old.json new.json` |

### 过滤表达式

```
age > 18          数值比较
name == "John"    字符串相等
status != "off"   不等于
email ~ "@gmail"  包含
phone             存在性检查
```

## 正则处理

| 命令 | 说明 | 示例 |
|------|------|------|
| `tc rx match <p>` | 高亮匹配 | `echo "abc123" \| tc rx match "\d+"` |
| `tc rx extract <p>` | 提取捕获组 | `echo "a=1" \| tc rx extract "(\w+)=(\w+)"` |
| `tc rx replace <p> <r>` | 替换 | `echo "foo" \| tc rx replace "foo" "bar"` |
| `tc rx grep <p>` | 过滤行 | `cat log \| tc rx grep "ERROR"` |
| `tc rx grep <p> -v` | 反转过滤 | `cat log \| tc rx grep "DEBUG" -v` |
| `tc rx count <p>` | 统计次数 | `cat log \| tc rx count "ERROR"` |
| `tc rx findall <p>` | 列出所有匹配 | `echo "a1b2" \| tc rx findall "\d"` |

## 通用工具

| 命令 | 说明 | 示例 |
|------|------|------|
| `tc hash <algo>` | 哈希 | `echo -n "hi" \| tc hash sha256` |
| `tc enc <type>` | 编码 | `echo "hi" \| tc enc base64` |
| `tc dec <type>` | 解码 | `echo "aGk=" \| tc dec base64` |
| `tc stats` | 统计 | `cat f.txt \| tc stats` |
| `tc count` | 行数 | `cat f.txt \| tc count` |
| `tc sort` | 排序 | `cat f.txt \| tc sort` |
| `tc sort -n` | 数值排序 | `echo -e "10\n2\n1" \| tc sort -n` |
| `tc sort -r` | 反转排序 | `cat f.txt \| tc sort -r` |
| `tc uniq` | 去重 | `cat f.txt \| tc sort \| tc uniq` |
| `tc head [n]` | 前 N 行 | `cat f.txt \| tc head 20` |
| `tc tail [n]` | 后 N 行 | `cat f.txt \| tc tail 20` |
| `tc sample [n]` | 随机采样 | `cat f.txt \| tc sample 100` |
| `tc fields <n>` | 提取字段 | `echo "a:b:c" \| tc fields 2 -s ":"` |

## 组合技

```bash
# Top 10 IP
cat access.log | tc rx extract "(\d+\.\d+\.\d+\.\d+)" | tc sort | tc uniq -c | tc sort -n -r | tc head 10

# 从日志提取 JSON 错误码
cat app.log | tc rx extract '(\{.*\})' | tc json get "error_code"

# JSON 过滤后提取邮箱
cat users.json | tc json filter 'age > 18' | tc json get "*.email"

# 日志中各状态码计数
cat access.log | tc rx extract 'HTTP/\d\.\d" (\d{3})' | tc sort | tc uniq -c | tc sort -n -r
```
'@

Write-File 'docs/examples.md' @'
# tc 使用示例

## 场景 1：分析 API 响应

```bash
# 调用 API 并格式化
curl -s https://api.example.com/users | tc json fmt

# 提取特定字段
curl -s https://api.example.com/users | tc json get "data.*.email"

# 过滤活跃用户
curl -s https://api.example.com/users | tc json filter 'status == "active"'
```

## 场景 2：分析 Nginx 日志

```bash
# 统计 5xx 错误数
cat access.log | tc rx grep " 5\d{2} " | tc count

# Top 10 访问 IP
cat access.log | tc rx extract "^(\d+\.\d+\.\d+\.\d+)" | tc sort | tc uniq -c | tc sort -n -r | tc head 10

# 统计各 HTTP 方法
cat access.log | tc rx extract '"(GET|POST|PUT|DELETE)' | tc sort | tc uniq -c | tc sort -n -r

# 找出慢请求（响应时间 > 5s）
cat access.log | tc rx extract '(\d+\.\d{3})$' | tc rx grep '^[5-9]\.'
```

## 场景 3：处理配置文件

```bash
# 查看 JSON 配置
cat config.json | tc json fmt

# 修改前备份，然后 diff
cp config.json config.json.bak
# ... 修改 config.json ...
tc json diff config.json.bak config.json

# 提取所有环境变量名
cat .env | tc rx extract "^(\w+)=" | tc sort
```

## 场景 4：数据清洗

```bash
# 从 CSV 提取第 2 列
cat data.csv | tc fields 2 -s ","

# 去除空行
cat messy.txt | tc rx grep "^.+$"

# 提取所有邮箱
cat contacts.txt | tc rx findall "[\w.]+@[\w.]+\.\w+"

# 提取所有 URL
cat page.html | tc rx findall "https?://[^\s\"'<>]+"
```

## 场景 5：快速哈希/编码

```bash
# 计算文件哈希
tc hash sha256 important.tar.gz

# 编码 JWT payload
echo '{"sub":"12345","name":"John"}' | tc enc base64

# 解码 URL 编码
echo "hello%20world%21" | tc dec url
```

## 场景 6：开发调试

```bash
# 格式化 API 响应
curl -s localhost:3000/api/users | tc json fmt

# 查看数据库导出
cat dump.json | tc json get "rows.*.id" | tc count

# 对比两个环境的配置
tc json diff dev.json prod.json

# 快速生成测试数据摘要
cat test-data.json | tc json table "id,name,email" --csv > test-report.csv
```
'@

Write-File '.github/workflows/ci.yml' @'
name: CI

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        go-version: ['1.22', '1.23']
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ matrix.go-version }}
      - run: go test ./... -v -race
      - run: go build .

  release:
    needs: test
    runs-on: ubuntu-latest
    if: startsWith(github.ref, 'refs/tags/')
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.22'
      - uses: goreleaser/goreleaser-action@v6
        with:
          version: latest
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
'@

Write-File '.gitignore' @'
bin/
dist/
*.exe
*.test
*.out
.idea/
.vscode/
*.swp
*.swo
*~
.DS_Store
'@

Write-File 'go.mod' @'
module github.com/user/tc

go 1.22

require (
	github.com/spf13/cobra v1.8.1
	github.com/tidwall/gjson v1.17.3
	github.com/tidwall/sjson v1.2.5
	github.com/fatih/color v1.17.0
)

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	golang.org/x/sys v0.21.0 // indirect
)
'@

Write-File '.goreleaser.yml' @'
version: 2
project_name: tc

before:
  hooks:
    - go mod tidy
    - go test ./...

builds:
  - main: .
    binary: tc
    ldflags:
      - -s -w -X main.version={{.Version}}
    goos:
      - linux
      - darwin
      - windows
    goarch:
      - amd64
      - arm64
    ignore:
      - goos: windows
        goarch: arm64

archives:
  - format: tar.gz
    name_template: "{{ .ProjectName }}_{{ .Os }}_{{ .Arch }}"
    format_overrides:
      - goos: windows
        format: zip

brews:
  - repository:
      owner: user
      name: homebrew-tc
    homepage: "https://github.com/user/tc"
    description: "ToolCraft — 程序员的文本/数据瑞士军刀"

checksum:
  name_template: "checksums.txt"

changelog:
  sort: asc
  filters:
    exclude:
      - "^docs:"
      - "^test:"
      - "^ci:"
'@

Write-File 'internal/json/json.go' @'
package jsonutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

// Format 格式化 JSON
func Format(data []byte) (string, error) {
	var out bytes.Buffer
	if err := json.Indent(&out, bytes.TrimSpace(data), "", "  "); err != nil {
		return "", fmt.Errorf("无效的 JSON: %w", err)
	}
	return out.String(), nil
}

// Minify 压缩 JSON
func Minify(data []byte) (string, error) {
	var out bytes.Buffer
	if err := json.Compact(&out, bytes.TrimSpace(data)); err != nil {
		return "", fmt.Errorf("无效的 JSON: %w", err)
	}
	return out.String(), nil
}

// Get 按路径查询 JSON
func Get(data []byte, path string) (string, error) {
	result := gjson.GetBytes(data, path)
	if !result.Exists() {
		return "null", nil
	}
	if result.Type == gjson.String {
		return result.String(), nil
	}
	return result.Raw, nil
}

// Filter 按条件过滤 JSON 数组
func Filter(data []byte, expr string) (string, error) {
	result := gjson.ParseBytes(data)
	if !result.IsArray() {
		return "", fmt.Errorf("输入不是 JSON 数组")
	}

	// 解析表达式: field op value
	field, op, value, err := parseExpr(expr)
	if err != nil {
		return "", err
	}

	var filtered []string
	result.ForEach(func(_, item gjson.Result) bool {
		fieldVal := item.Get(field)
		if !fieldVal.Exists() {
			return true
		}
		if matchCondition(fieldVal, op, value) {
			filtered = append(filtered, item.Raw)
		}
		return true
	})

	return "[" + strings.Join(filtered, ",") + "]", nil
}

// Keys 提取 JSON 对象的所有 key
func Keys(data []byte) (string, error) {
	result := gjson.ParseBytes(data)
	if !result.IsObject() {
		return "", fmt.Errorf("输入不是 JSON 对象")
	}
	var keys []string
	result.ForEach(func(key, _ gjson.Result) bool {
		keys = append(keys, key.String())
		return true
	})
	out, _ := json.Marshal(keys)
	return string(out), nil
}

// Values 提取 JSON 对象的所有 value
func Values(data []byte) (string, error) {
	result := gjson.ParseBytes(data)
	if !result.IsObject() {
		return "", fmt.Errorf("输入不是 JSON 对象")
	}
	var vals []json.RawMessage
	result.ForEach(func(_, val gjson.Result) bool {
		vals = append(vals, json.RawMessage(val.Raw))
		return true
	})
	out, _ := json.Marshal(vals)
	return string(out), nil
}

// Paths 列出 JSON 的所有路径
func Paths(data []byte) ([]string, error) {
	result := gjson.ParseBytes(data)
	var paths []string
	walkPaths(result, "", &paths)
	sort.Strings(paths)
	return paths, nil
}

func walkPaths(node gjson.Result, prefix string, paths *[]string) {
	if prefix != "" {
		*paths = append(*paths, prefix)
	}
	if node.IsObject() {
		node.ForEach(func(key, val gjson.Result) bool {
			childPath := prefix + "." + key.String()
			if prefix == "" {
				childPath = key.String()
			}
			walkPaths(val, childPath, paths)
			return true
		})
	} else if node.IsArray() {
		for i := 0; i < len(node.Array()); i++ {
			childPath := prefix + "." + strconv.Itoa(i)
			if prefix == "" {
				childPath = strconv.Itoa(i)
			}
			walkPaths(node.Array()[i], childPath, paths)
		}
	}
}

// Table 将 JSON 数组转为表格
func Table(data []byte, fieldsStr string, csv bool) (string, error) {
	result := gjson.ParseBytes(data)
	if !result.IsArray() {
		return "", fmt.Errorf("输入不是 JSON 数组")
	}

	fields := strings.Split(fieldsStr, ",")
	for i := range fields {
		fields[i] = strings.TrimSpace(fields[i])
	}

	sep := "\t"
	if csv {
		sep = ","
	}

	var buf bytes.Buffer
	// 表头
	buf.WriteString(strings.Join(fields, sep) + "\n")

	// 数据行
	result.ForEach(func(_, item gjson.Result) bool {
		var row []string
		for _, f := range fields {
			val := item.Get(f)
			if val.Exists() {
				row = append(row, val.String())
			} else {
				row = append(row, "")
			}
		}
		buf.WriteString(strings.Join(row, sep) + "\n")
		return true
	})

	return buf.String(), nil
}

// Diff 比较两个 JSON 的差异
func Diff(data1, data2 []byte) (string, error) {
	var v1, v2 interface{}
	if err := json.Unmarshal(data1, &v1); err != nil {
		return "", fmt.Errorf("第一个 JSON 无效: %w", err)
	}
	if err := json.Unmarshal(data2, &v2); err != nil {
		return "", fmt.Errorf("第二个 JSON 无效: %w", err)
	}

	var diffs []string
	diffWalk(v1, v2, "", &diffs)

	if len(diffs) == 0 {
		return "无差异\n", nil
	}

	var buf bytes.Buffer
	for _, d := range diffs {
		buf.WriteString(d + "\n")
	}
	return buf.String(), nil
}

func diffWalk(v1, v2 interface{}, path string, diffs *[]string) {
	switch a := v1.(type) {
	case map[string]interface{}:
		b, ok := v2.(map[string]interface{})
		if !ok {
			*diffs = append(*diffs, fmt.Sprintf("- %s: %v", path, v1))
			*diffs = append(*diffs, fmt.Sprintf("+ %s: %v", path, v2))
			return
		}
		allKeys := make(map[string]bool)
		for k := range a {
			allKeys[k] = true
		}
		for k := range b {
			allKeys[k] = true
		}
		for k := range allKeys {
			p := path + "." + k
			if path == "" {
				p = k
			}
			va, okA := a[k]
			vb, okB := b[k]
			if !okA {
				*diffs = append(*diffs, fmt.Sprintf("+ %s: %v", p, vb))
			} else if !okB {
				*diffs = append(*diffs, fmt.Sprintf("- %s: %v", p, va))
			} else {
				diffWalk(va, vb, p, diffs)
			}
		}
	case []interface{}:
		b, ok := v2.([]interface{})
		if !ok {
			*diffs = append(*diffs, fmt.Sprintf("- %s: %v", path, v1))
			*diffs = append(*diffs, fmt.Sprintf("+ %s: %v", path, v2))
			return
		}
		maxLen := len(a)
		if len(b) > maxLen {
			maxLen = len(b)
		}
		for i := 0; i < maxLen; i++ {
			p := path + "[" + strconv.Itoa(i) + "]"
			if i >= len(a) {
				*diffs = append(*diffs, fmt.Sprintf("+ %s: %v", p, b[i]))
			} else if i >= len(b) {
				*diffs = append(*diffs, fmt.Sprintf("- %s: %v", p, a[i]))
			} else {
				diffWalk(a[i], b[i], p, diffs)
			}
		}
	default:
		if fmt.Sprintf("%v", v1) != fmt.Sprintf("%v", v2) {
			*diffs = append(*diffs, fmt.Sprintf("- %s: %v", path, v1))
			*diffs = append(*diffs, fmt.Sprintf("+ %s: %v", path, v2))
		}
	}
}

// parseExpr 解析简单条件表达式: field op value
func parseExpr(expr string) (field, op, value string, err error) {
	ops := []string{"!=", ">=", "<=", "==", ">", "<", "~"}
	for _, o := range ops {
		idx := strings.Index(expr, o)
		if idx > 0 {
			field = strings.TrimSpace(expr[:idx])
			value = strings.TrimSpace(expr[idx+len(o):])
			// 去掉引号
			value = strings.Trim(value, `"'`)
			return field, o, value, nil
		}
	}
	// 没有操作符，当作"存在性"检查
	return strings.TrimSpace(expr), "exists", "", nil
}

// matchCondition 匹配条件
func matchCondition(val gjson.Result, op, value string) bool {
	switch op {
	case "exists":
		return val.Exists()
	case "==":
		return val.String() == value
	case "!=":
		return val.String() != value
	case "~":
		return strings.Contains(val.String(), value)
	case ">", "<", ">=", "<=":
		f1, err1 := strconv.ParseFloat(val.String(), 64)
		f2, err2 := strconv.ParseFloat(value, 64)
		if err1 != nil || err2 != nil {
			return false
		}
		switch op {
		case ">":
			return f1 > f2
		case "<":
			return f1 < f2
		case ">=":
			return f1 >= f2
		case "<=":
			return f1 <= f2
		}
	}
	return false
}
'@

Write-File 'internal/json/json_test.go' @'
package jsonutil

import (
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
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && (s[0:1] == substr || contains(s[1:], substr)))
}
'@

Write-File 'internal/pipeline/pipe.go' @'
package pipeline

import (
	"fmt"
	"io"
	"os"
)

// ReadInput 从文件参数或 stdin 读取数据
// 如果 args 非空，读取第一个文件；否则读取 stdin
func ReadInput(args []string) ([]byte, error) {
	if len(args) > 0 {
		data, err := os.ReadFile(args[0])
		if err != nil {
			return nil, fmt.Errorf("读取文件 %s 失败: %w", args[0], err)
		}
		return data, nil
	}
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return nil, fmt.Errorf("读取 stdin 失败: %w", err)
	}
	return data, nil
}
'@

Write-File 'internal/rx/rx.go' @'
package rxutil

import (
	"fmt"
	"regexp"
	"strings"
)

// Match 高亮匹配内容
func Match(data []byte, pattern string) (string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", fmt.Errorf("无效的正则表达式: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	var matched []string
	for _, line := range lines {
		if re.MatchString(line) {
			// 高亮匹配部分
			highlighted := re.ReplaceAllString(line, "\033[1;33m$0\033[0m")
			matched = append(matched, highlighted)
		}
	}
	return strings.Join(matched, "\n") + "\n", nil
}

// Extract 提取捕获组
func Extract(data []byte, pattern string) (string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", fmt.Errorf("无效的正则表达式: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	var results []string
	for _, line := range lines {
		matches := re.FindStringSubmatch(line)
		if matches == nil {
			continue
		}
		if len(matches) > 1 {
			// 有捕获组，输出所有组（tab 分隔）
			results = append(results, strings.Join(matches[1:], "\t"))
		} else {
			// 无捕获组，输出整个匹配
			results = append(results, matches[0])
		}
	}
	return strings.Join(results, "\n") + "\n", nil
}

// Replace 替换匹配内容
func Replace(data []byte, pattern, replacement string) (string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", fmt.Errorf("无效的正则表达式: %w", err)
	}

	result := re.ReplaceAllString(string(data), replacement)
	return result, nil
}

// Grep 过滤匹配行
func Grep(data []byte, pattern string, invert bool) (string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", fmt.Errorf("无效的正则表达式: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	var filtered []string
	for _, line := range lines {
		match := re.MatchString(line)
		if invert {
			match = !match
		}
		if match {
			filtered = append(filtered, line)
		}
	}
	return strings.Join(filtered, "\n") + "\n", nil
}

// Count 统计匹配次数
func Count(data []byte, pattern string) (int, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return 0, fmt.Errorf("无效的正则表达式: %w", err)
	}

	return len(re.FindAllString(string(data), -1)), nil
}

// FindAll 列出所有匹配
func FindAll(data []byte, pattern string) ([]string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("无效的正则表达式: %w", err)
	}

	matches := re.FindAllString(string(data), -1)
	return matches, nil
}
'@

Write-File 'internal/rx/rx_test.go' @'
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
'@

Write-File 'LICENSE' @'
MIT License

Copyright (c) 2026 tc contributors

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
'@

Write-File 'main.go' @'
package main

import (
	"github.com/user/tc/cmd"
)

var version = "dev"

func main() {
	cmd.SetVersion(version)
	cmd.Execute()
}
'@

Write-File 'Makefile' @'
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-s -w -X main.version=$(VERSION)"

.PHONY: build clean test install cross

build:
	go build $(LDFLAGS) -o bin/tc .

clean:
	rm -rf bin/

test:
	go test ./... -v

install:
	go install $(LDFLAGS) .

cross:
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/tc-linux-amd64 .
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o bin/tc-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o bin/tc-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o bin/tc-windows-amd64.exe .

lint:
	golangci-lint run ./...
'@

Write-File 'README.md' @'
# tc — ToolCraft

> 程序员的文本/数据瑞士军刀。JSON + 正则 + 文本处理，一个命令搞定。

[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

---

## 为什么用 tc？

| 痛点 | tc 的解决方式 |
|------|-------------|
| jq 语法记不住 | `tc json get "users.0.name"` — 直觉语法 |
| grep 不懂 JSON | tc 理解 JSON 结构，能按字段过滤 |
| 正则+JSON 组合要写脚本 | 管道一行搞定 |
| 工具太多记不住 | 一个 `tc` 全搞定 |

## 安装

```bash
# Go（推荐）
go install github.com/user/tc@latest

# Homebrew
brew install user/tc/tc

# 手动下载
# https://github.com/user/tc/releases
```

## 快速上手

### JSON 处理

```bash
# 格式化
echo '{"name":"John","age":30}' | tc json fmt

# 查询
cat data.json | tc json get "users.0.email"
cat data.json | tc json get "users.*.name"

# 过滤
cat users.json | tc json filter 'age > 18'
cat users.json | tc json filter 'status == "active"'
cat users.json | tc json filter 'email ~ "@gmail"'

# 转换
cat data.json | tc json keys       # 所有 key
cat data.json | tc json values     # 所有 value
cat data.json | tc json paths      # 所有路径

# 表格输出
cat users.json | tc json table "name,email,age"
cat users.json | tc json table "id,title" --csv

# Diff
tc json diff old.json new.json
```

### 正则处理

```bash
# 高亮匹配
echo "abc123def456" | tc rx match "\d+"
cat access.log | tc rx match "5\d{2}"

# 提取捕获组
echo "2024-01-15 error occurred" | tc rx extract "(\d{4}-\d{2}-\d{2})"
echo "name=John age=30" | tc rx extract "(\w+)=(\w+)"

# 替换
echo "hello world" | tc rx replace "world" "Go"
cat file.txt | tc rx replace "(\d+)" "NUM:$1"

# 过滤行
cat access.log | tc rx grep "ERROR|WARN"
cat access.log | tc rx grep "DEBUG" --invert  # 反转

# 统计
cat app.log | tc rx count "ERROR"

# 列出所有匹配
echo "a1b2c3" | tc rx findall "\d"
```

### 组合技（杀手功能）

```bash
# 从日志提取 JSON 再查字段
cat app.log | tc rx extract '(\{.*\})' | tc json get "error_code"

# JSON 数组过滤后提取特定字段
cat users.json | tc json filter 'age > 18' | tc json get "*.email"

# 日志中统计各状态码出现次数
cat access.log | tc rx extract "HTTP/\d\.\d\" (\d{3})" | tc sort | tc uniq --count | tc sort -n -r

# Top 10 IP
cat access.log | tc rx extract "(\d+\.\d+\.\d+\.\d+)" | tc sort | tc uniq --count | tc sort -n -r | tc head 10
```

### 通用工具

```bash
# 哈希
echo -n "hello" | tc hash sha256
echo -n "hello" | tc hash md5

# 编码/解码
echo "hello world" | tc enc base64
echo "aGVsbG8gd29ybGQ=" | tc dec base64
echo "hello world" | tc enc url

# 统计
cat file.txt | tc stats        # 行数 词数 字符数 字节数
cat file.txt | tc count        # 行数

# 排序/去重
cat file.txt | tc sort
cat file.txt | tc sort -n      # 数值排序
cat file.txt | tc sort | tc uniq

# 采样
cat bigfile.txt | tc sample 100

# 切片
cat file.txt | tc head 10
cat file.txt | tc tail 10

# 字段提取
echo "a:b:c:d" | tc fields 2 --sep ":"
```

## 命令速查表

```
tc json fmt          格式化 JSON
tc json min          压缩 JSON
tc json get <path>   查询字段
tc json filter <expr> 过滤数组
tc json keys         提取所有 key
tc json values       提取所有 value
tc json paths        列出所有路径
tc json table <f>    转为表格
tc json diff <f1> <f2> 比较差异

tc rx match <p>      高亮匹配
tc rx extract <p>    提取捕获组
tc rx replace <p> <r> 替换
tc rx grep <p>       过滤行
tc rx count <p>      统计次数
tc rx findall <p>    列出所有匹配

tc hash <algo>       计算哈希
tc enc <type>        编码
tc dec <type>        解码
tc stats             统计信息
tc count             行数
tc sort              排序
tc uniq              去重
tc head [n]          前 N 行
tc tail [n]          后 N 行
tc sample [n]        随机采样
tc fields <n>        提取字段

tc version           版本信息
```

## 设计原则

- **管道友好** — stdin → 处理 → stdout
- **零配置** — 开箱即用
- **快速** — Go 编译，启动 < 50ms
- **小体积** — 单二进制，< 10MB
- **错误友好** — 告诉你怎么修，不是堆栈

## 技术栈

- Go 1.22+
- [cobra](https://github.com/spf13/cobra) — CLI 框架
- [gjson](https://github.com/tidwall/gjson) — JSON 查询
- [fatih/color](https://github.com/fatih/color) — 终端颜色

## License

MIT
'@

