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
