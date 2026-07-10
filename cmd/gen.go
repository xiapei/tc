package cmd

import (
	"github.com/spf13/cobra"
)

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "生成随机数据（密码、Token 等）",
	Long:  `生成随机数据：密码、Token 等。`,
	Example: `  tc gen password           # 16位随机密码
  tc gen password 20        # 20位密码
  tc gen password --no-sym  # 不含符号`,
}

func init() {
	rootCmd.AddCommand(genCmd)
}
