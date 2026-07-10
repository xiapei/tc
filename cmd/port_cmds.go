package cmd

import (
	"fmt"
	"strconv"

	portutil "github.com/user/tc/internal/port"
	"github.com/spf13/cobra"
)

var portCmd = &cobra.Command{
	Use:   "port <port>",
	Short: "检查端口占用",
	Long: `检查指定端口是否被占用，以及被哪个进程占用。

Windows 下可通过 netstat 查询占用进程的 PID 和进程名。`,
	Example: `  tc port 8080
  tc port 3000
  tc port 80`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		port, err := strconv.Atoi(args[0])
		if err != nil || port < 1 || port > 65535 {
			return fmt.Errorf("无效端口号: %s（请输入 1-65535）", args[0])
		}

		result := portutil.Check(port)
		fmt.Println(portutil.Format(result))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(portCmd)
}
