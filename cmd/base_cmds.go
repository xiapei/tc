package cmd

import (
	"fmt"
	"strconv"
	"strings"

	baseutil "github.com/user/tc/internal/base"
	"github.com/spf13/cobra"
)

var baseCmd = &cobra.Command{
	Use:   "base <from-base> <value> [value...]",
	Short: "进制转换",
	Long: `进制转换，支持 2~36 进制互转。

常见进制：
  2   - 二进制
  8   - 八进制
  10  - 十进制
  16  - 十六进制

前缀自动识别：0x → 16进制, 0b → 2进制, 0o → 8进制

使用 --all 显示所有进制表示。`,
	Example: `  tc base 16 FF                 # 16进制 → 10进制: 255
  tc base 10 255                # 10进制 → 16进制: 0xFF
  tc base 2 1010                # 2进制 → 10进制: 10
  tc base 16 "FF" "1A" "2B"     # 批量转换
  tc base 16 ff --all           # 显示所有进制
  tc base 0xFF                  # 自动识别前缀（进制省略时自动识别）`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		showAll, _ := cmd.Flags().GetBool("all")

		var fromBase int
		var values []string

		// 判断第一个参数是进制还是值
		if base, err := strconv.Atoi(args[0]); err == nil && len(args) > 1 {
			fromBase = baseutil.NormalizeBase(base)
			values = args[1:]
		} else {
			// 省略进制，自动识别
			fromBase = 0
			values = args
		}

		for _, v := range values {
			v = strings.TrimSpace(v)
			if v == "" {
				continue
			}

			var result *baseutil.ConvertResult
			var err error

			if fromBase > 0 {
				result, err = baseutil.Convert(v, fromBase)
			} else {
				result, err = baseutil.ParseAuto(v)
			}

			if err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "错误: %s — %v\n", v, err)
				continue
			}

			if showAll {
				fmt.Printf("%s\n%s\n", v, baseutil.ShowAll(result))
			} else {
				base := fromBase
				if base == 0 {
					base = 10
				}
				switch base {
				case 2:
					fmt.Printf("0b%s (2) = %d (10)\n", v, result.Decimal)
				case 8:
					fmt.Printf("0o%s (8) = %d (10)\n", v, result.Decimal)
				case 10:
					fmt.Printf("%s (10) = 0x%s (16)\n", v, result.Hex)
				case 16:
					fmt.Printf("0x%s (16) = %d (10)\n", strings.ToUpper(v), result.Decimal)
				default:
					fmt.Printf("%s (%d) = %d (10)\n", v, base, result.Decimal)
				}
			}
		}

		return nil
	},
}

func init() {
	baseCmd.Flags().BoolP("all", "a", false, "显示所有进制表示")
	rootCmd.AddCommand(baseCmd)
}
