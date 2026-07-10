package cmd

import (
	"fmt"

	jwtutil "github.com/user/tc/internal/jwt"
	"github.com/user/tc/internal/pipeline"
	"github.com/spf13/cobra"
)

var jwtDecodeCmd = &cobra.Command{
	Use:   "decode [token]",
	Short: "解码 JWT（不验证签名）",
	Long: `解码 JWT token，打印 header 和 payload 的格式化 JSON。

支持从参数或 stdin 传入 token。
注意：本工具只解码查看内容，不验证签名。`,
	Example: `  tc jwt decode eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U
  echo "eyJhbGciOiJIUzI1NiJ9.eyJkYXRhIjoidGVzdCJ9.ZWZlZmUzZ" | tc jwt decode`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var tokenStr string

		if len(args) > 0 {
			tokenStr = args[0]
		} else {
			data, err := pipeline.ReadInput(nil)
			if err != nil {
				return err
			}
			tokenStr = string(data)
		}

		if !jwtutil.IsJWT(tokenStr) {
			return fmt.Errorf("无效的 JWT 格式：需要包含三段的 token（header.payload.signature）")
		}

		token, err := jwtutil.Decode(tokenStr)
		if err != nil {
			return err
		}

		result, err := jwtutil.Format(token)
		if err != nil {
			return err
		}

		fmt.Print(result)
		return nil
	},
}

var jwtCmd = &cobra.Command{
	Use:   "jwt",
	Short: "JWT 工具（解码）",
	Long:  `JWT 工具：解码 JWT Token。`,
	Example: `  tc jwt decode <token>`,
}

func init() {
	jwtCmd.AddCommand(jwtDecodeCmd)
	rootCmd.AddCommand(jwtCmd)
}
