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
