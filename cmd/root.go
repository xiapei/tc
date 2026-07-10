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
