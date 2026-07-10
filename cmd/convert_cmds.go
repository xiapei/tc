package cmd

import (
	"fmt"
	"strings"

	convertutil "github.com/user/tc/internal/convert"
	"github.com/user/tc/internal/pipeline"
	"github.com/spf13/cobra"
)

var convertCmd = &cobra.Command{
	Use:   "convert <j2y|y2j> [file]",
	Short: "JSON / YAML 互转",
	Long: `JSON 和 YAML 格式互转。

子命令：
  j2y   JSON → YAML
  y2j   YAML → JSON

支持从文件读取或 stdin 管道输入。`,
	Example: `  cat data.json | tc convert j2y
  tc convert j2y data.json
  cat config.yaml | tc convert y2j
  tc convert y2j config.yaml
  cat data.json | tc convert j2y > config.yaml`,
	Args: cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		direction := strings.ToLower(args[0])
		if direction != "j2y" && direction != "y2j" {
			return fmt.Errorf("无效的转换方向: %s（请用 j2y 或 y2j）", args[0])
		}

		fileArgs := args[1:]
		data, err := pipeline.ReadInput(fileArgs)
		if err != nil {
			return err
		}

		if strings.TrimSpace(string(data)) == "" {
			return fmt.Errorf("输入为空")
		}

		var result string
		switch direction {
		case "j2y":
			result, err = convertutil.JSONToYAML(data)
		case "y2j":
			result, err = convertutil.YAMLToJSON(data)
		}

		if err != nil {
			return err
		}

		fmt.Print(result)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(convertCmd)
}
