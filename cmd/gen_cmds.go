package cmd

import (
	"fmt"
	"strconv"

	genutil "github.com/user/tc/internal/gen"
	"github.com/spf13/cobra"
)

var genPasswordCmd = &cobra.Command{
	Use:   "password [length]",
	Short: "生成随机密码",
	Long: `生成随机密码，支持自定义字符组合和长度。

默认包含大小写字母、数字和符号，长度 16 位。
可用 --no-xxx 排除某类字符。`,
	Example: `  tc gen password              # 16位，大小写+数字+符号
  tc gen password 20           # 20位
  tc gen password --no-sym     # 不含符号
  tc gen password --no-upper --no-sym  # 仅小写+数字
  tc gen password 32 --no-sym --digit-only  # 仅数字 32 位`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := genutil.DefaultConfig()

		// 可选长度
		if len(args) > 0 {
			n, err := strconv.Atoi(args[0])
			if err != nil || n <= 0 {
				return fmt.Errorf("无效的长度: %s", args[0])
			}
			cfg.Length = n
		}

		// 读取 flags
		if cmd.Flags().Changed("no-upper") {
			cfg.UseUpper = false
		}
		if cmd.Flags().Changed("no-lower") {
			cfg.UseLower = false
		}
		if cmd.Flags().Changed("no-digit") {
			cfg.UseDigit = false
		}
		if cmd.Flags().Changed("no-sym") {
			cfg.UseSymbol = false
		}

		// 兼容旧式 flag
		upperOnly, _ := cmd.Flags().GetBool("upper-only")
		digitOnly, _ := cmd.Flags().GetBool("digit-only")
		if upperOnly {
			cfg.UseLower = false
			cfg.UseDigit = false
			cfg.UseSymbol = false
		}
		if digitOnly {
			cfg.UseUpper = false
			cfg.UseLower = false
			cfg.UseSymbol = false
		}

		password, err := genutil.GeneratePassword(cfg)
		if err != nil {
			return err
		}

		entropy := genutil.Entropy(cfg)
		strength := genutil.PasswordStrength(entropy)

		fmt.Println(password)
		fmt.Fprintf(cmd.ErrOrStderr(), "# %s | 熵值: %.1f bit | 强度: %s\n",
			cfg.FormatConfigSummary(), entropy, strength)

		return nil
	},
}

func init() {
	genPasswordCmd.Flags().Bool("no-upper", false, "不含大写字母")
	genPasswordCmd.Flags().Bool("no-lower", false, "不含小写字母")
	genPasswordCmd.Flags().Bool("no-digit", false, "不含数字")
	genPasswordCmd.Flags().Bool("no-sym", false, "不含符号")
	genPasswordCmd.Flags().Bool("upper-only", false, "仅大写字母")
	genPasswordCmd.Flags().Bool("digit-only", false, "仅数字")

	genCmd.AddCommand(genPasswordCmd)
}
