package cmd

import (
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
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
		if len(lines) == 0 {
			return nil
		}
		prev := lines[0]
		count := 1
		for _, line := range lines[1:] {
			if line == prev {
				count++
			} else {
				if countFlag {
					fmt.Printf("%d\t%s\n", count, prev)
				} else {
					fmt.Println(prev)
				}
				prev = line
				count = 1
			}
		}
		if countFlag {
			fmt.Printf("%d\t%s\n", count, prev)
		} else {
			fmt.Println(prev)
		}
		return nil
	},
}

// head 命令
var headCmd = &cobra.Command{
	Use:   "head [n] [file]",
	Short: "输出前 N 行（默认 10）",
	Example: `  cat file.txt | tc head 10
  cat file.txt | tc head -n 5`,
	RunE: func(cmd *cobra.Command, args []string) error {
		n, _ := cmd.Flags().GetInt("lines")
		if n <= 0 {
			n = 10
		}
		fileArgs := args
		if len(args) > 0 {
			if parsed, err := strconv.Atoi(args[0]); err == nil && parsed > 0 {
				n = parsed
				fileArgs = args[1:]
			}
		}
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
	Use:   "tail [n] [file]",
	Short: "输出后 N 行（默认 10）",
	Example: `  cat file.txt | tc tail 10
  cat file.txt | tc tail -n 5`,
	RunE: func(cmd *cobra.Command, args []string) error {
		n, _ := cmd.Flags().GetInt("lines")
		if n <= 0 {
			n = 10
		}
		fileArgs := args
		if len(args) > 0 {
			if parsed, err := strconv.Atoi(args[0]); err == nil && parsed > 0 {
				n = parsed
				fileArgs = args[1:]
			}
		}
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
		indices := rand.Perm(len(lines))
		for _, idx := range indices[:n] {
			fmt.Println(lines[idx])
		}
		return nil
	},
}

// b64encode 命令
var b64encodeCmd = &cobra.Command{
	Use:   "b64encode [data...]",
	Short: "Base64 编码",
	Long: `将字符串或 stdin 编码为 Base64。

支持标准 Base64 和 URL-safe Base64（使用 --url 标志）。`,
	Example: `  tc b64encode "hello world"
  tc b64encode "hello" "world"
  echo -n "hello world" | tc b64encode
  tc b64encode --url "hello world"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		urlMode, _ := cmd.Flags().GetBool("url")

		var data []byte
		if len(args) > 0 {
			data = []byte(strings.Join(args, " "))
		} else {
			piped, err := pipeline.ReadInput(nil)
			if err != nil {
				return err
			}
			data = piped
		}

		var result string
		if urlMode {
			result = base64.URLEncoding.EncodeToString(data)
		} else {
			result = base64.StdEncoding.EncodeToString(data)
		}
		fmt.Print(result)
		return nil
	},
}

// b64decode 命令
var b64decodeCmd = &cobra.Command{
	Use:   "b64decode [data...]",
	Short: "Base64 解码",
	Long: `将 Base64 字符串或 stdin 解码为原文。

支持标准 Base64 和 URL-safe Base64（使用 --url 标志）。`,
	Example: `  tc b64decode "aGVsbG8gd29ybGQ="
  tc b64decode --url "aGVsbG8gd29ybGQ="
  echo "aGVsbG8gd29ybGQ=" | tc b64decode`,
	RunE: func(cmd *cobra.Command, args []string) error {
		urlMode, _ := cmd.Flags().GetBool("url")

		var input string
		if len(args) > 0 {
			input = strings.Join(args, " ")
		} else {
			piped, err := pipeline.ReadInput(nil)
			if err != nil {
				return err
			}
			input = strings.TrimSpace(string(piped))
		}

		var decoded []byte
		var err error
		if urlMode {
			decoded, err = base64.URLEncoding.DecodeString(input)
		} else {
			decoded, err = base64.StdEncoding.DecodeString(input)
		}
		if err != nil {
			return fmt.Errorf("base64 解码失败: %w", err)
		}
		fmt.Print(string(decoded))
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

func init() {
	sortCmd.Flags().BoolP("numeric", "n", false, "按数值排序")
	sortCmd.Flags().BoolP("reverse", "r", false, "反转排序")

	headCmd.Flags().IntP("lines", "n", 10, "输出行数")
	tailCmd.Flags().IntP("lines", "n", 10, "输出行数")

	fieldsCmd.Flags().StringP("sep", "s", "\t", "字段分隔符")
	uniqCmd.Flags().BoolP("count", "c", false, "输出每行出现次数")

	b64encodeCmd.Flags().Bool("url", false, "URL-safe Base64 编码")
	b64decodeCmd.Flags().Bool("url", false, "URL-safe Base64 解码")

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
	rootCmd.AddCommand(b64encodeCmd)
	rootCmd.AddCommand(b64decodeCmd)
}
