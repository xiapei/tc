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
