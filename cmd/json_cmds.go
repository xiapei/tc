package cmd

import (
	"fmt"
	"os"
	"strings"

	jsonutil "github.com/user/tc/internal/json"
	"github.com/user/tc/internal/pipeline"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var jsonFmtCmd = &cobra.Command{
	Use:   "fmt [file]",
	Short: "格式化 JSON",
	Example: `  cat data.json | tc json fmt
  tc json fmt data.json
  echo '{"a":1,"b":2}' | tc json fmt`,
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := pipeline.ReadInput(args)
		if err != nil {
			return err
		}
		result, err := jsonutil.Format(data)
		if err != nil {
			return err
		}
		fmt.Println(result)
		return nil
	},
}

var jsonMinCmd = &cobra.Command{
	Use:   "min [file]",
	Short: "压缩 JSON 为单行",
	Example: `  cat data.json | tc json min
  tc json min data.json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := pipeline.ReadInput(args)
		if err != nil {
			return err
		}
		result, err := jsonutil.Minify(data)
		if err != nil {
			return err
		}
		fmt.Print(result)
		return nil
	},
}

var jsonGetCmd = &cobra.Command{
	Use:   "get <path> [file]",
	Short: "按路径查询 JSON 字段",
	Long: `按路径查询 JSON 字段，支持：
  - 嵌套字段: "user.name"
  - 数组索引: "users.0.name"
  - 通配符:   "users.*.email"
  - 数组尾部: "items.-1"`,
	Example: `  cat data.json | tc json get "users.0.name"
  cat data.json | tc json get "users.*.email"
  echo '{"a":{"b":1}}' | tc json get "a.b"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]
		fileArgs := args[1:]
		data, err := pipeline.ReadInput(fileArgs)
		if err != nil {
			return err
		}
		result, err := jsonutil.Get(data, path)
		if err != nil {
			return err
		}
		fmt.Println(result)
		return nil
	},
}

var jsonFilterCmd = &cobra.Command{
	Use:   "filter <expr> [file]",
	Short: "按条件过滤 JSON 数组元素",
	Long: `支持简单条件表达式和复合条件：
  - 比较: age > 18, name == "John", status != "inactive"
  - 包含: email ~ "@gmail"
  - 存在: phone
  - 复合: age > 18 && status == "active"
  - 复合: age < 18 || age > 60`,
	Example: `  cat users.json | tc json filter 'age > 18'
  cat users.json | tc json filter 'status == "active"'
  cat users.json | tc json filter 'email ~ "@gmail"'
  cat users.json | tc json filter 'age > 18 && status == "active"'
  cat users.json | tc json filter 'age < 18 || age > 60'`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		expr := args[0]
		fileArgs := args[1:]
		data, err := pipeline.ReadInput(fileArgs)
		if err != nil {
			return err
		}
		result, err := jsonutil.Filter(data, expr)
		if err != nil {
			return err
		}
		fmt.Println(result)
		return nil
	},
}

var jsonKeysCmd = &cobra.Command{
	Use:   "keys [file]",
	Short: "提取 JSON 对象的所有 key",
	Example: `  echo '{"name":"John","age":30}' | tc json keys
  cat data.json | tc json keys`,
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := pipeline.ReadInput(args)
		if err != nil {
			return err
		}
		result, err := jsonutil.Keys(data)
		if err != nil {
			return err
		}
		fmt.Println(result)
		return nil
	},
}

var jsonValuesCmd = &cobra.Command{
	Use:   "values [file]",
	Short: "提取 JSON 对象的所有 value",
	Example: `  echo '{"name":"John","age":30}' | tc json values`,
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := pipeline.ReadInput(args)
		if err != nil {
			return err
		}
		result, err := jsonutil.Values(data)
		if err != nil {
			return err
		}
		fmt.Println(result)
		return nil
	},
}

var jsonPathsCmd = &cobra.Command{
	Use:   "paths [file]",
	Short: "列出 JSON 的所有路径",
	Example: `  echo '{"a":{"b":1},"c":[2,3]}' | tc json paths`,
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := pipeline.ReadInput(args)
		if err != nil {
			return err
		}
		result, err := jsonutil.Paths(data)
		if err != nil {
			return err
		}
		for _, p := range result {
			fmt.Println(p)
		}
		return nil
	},
}

var jsonTableCmd = &cobra.Command{
	Use:   "table <fields> [file]",
	Short: "将 JSON 数组转为表格",
	Long: `提取指定字段，以制表符分隔输出（方便粘贴到 Excel）。
第一个行为表头。`,
	Example: `  cat users.json | tc json table "name,email,age"
  cat users.json | tc json table "id,title" --csv`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fields := args[0]
		fileArgs := args[1:]
		csv, _ := cmd.Flags().GetBool("csv")
		data, err := pipeline.ReadInput(fileArgs)
		if err != nil {
			return err
		}
		result, err := jsonutil.Table(data, fields, csv)
		if err != nil {
			return err
		}
		fmt.Print(result)
		return nil
	},
}

var jsonDiffCmd = &cobra.Command{
	Use:   "diff <file1> <file2>",
	Short: "比较两个 JSON 文件的差异",
	Example: `  tc json diff old.json new.json`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		data1, err := os.ReadFile(args[0])
		if err != nil {
			return fmt.Errorf("读取文件 %s 失败: %w", args[0], err)
		}
		data2, err := os.ReadFile(args[1])
		if err != nil {
			return fmt.Errorf("读取文件 %s 失败: %w", args[1], err)
		}
		result, err := jsonutil.Diff(data1, data2)
		if err != nil {
			return err
		}
		for _, line := range strings.Split(strings.TrimRight(result, "\n"), "\n") {
			if strings.HasPrefix(line, "- ") {
				fmt.Println(color.RedString(line))
			} else if strings.HasPrefix(line, "+ ") {
				fmt.Println(color.GreenString(line))
			} else {
				fmt.Println(line)
			}
		}
		return nil
	},
}

var jsonSetCmd = &cobra.Command{
	Use:   "set <path> <value> [file]",
	Short: "设置 JSON 路径的值",
	Example: `  echo '{"name":"John"}' | tc json set "name" "Jane"
  echo '{"user":{"name":"John"}}' | tc json set "user.name" "Jane"
  echo '{"count":0}' | tc json set "count" 42
  tc json set "name" "Jane" data.json -o data.json`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]
		value := args[1]
		fileArgs := args[2:]
		data, err := pipeline.ReadInput(fileArgs)
		if err != nil {
			return err
		}
		result, err := jsonutil.Set(data, path, value)
		if err != nil {
			return err
		}
		return writeOutput(cmd, result+"\n")
	},
}

var jsonDelCmd = &cobra.Command{
	Use:   "del <path> [file]",
	Short: "删除 JSON 路径",
	Example: `  echo '{"name":"John","age":30}' | tc json del "age"
  echo '{"user":{"name":"John","age":30}}' | tc json del "user.age"
  tc json del "age" data.json -o data.json`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]
		fileArgs := args[1:]
		data, err := pipeline.ReadInput(fileArgs)
		if err != nil {
			return err
		}
		result, err := jsonutil.Delete(data, path)
		if err != nil {
			return err
		}
		return writeOutput(cmd, result+"\n")
	},
}

var jsonMergeCmd = &cobra.Command{
	Use:   "merge <file> [file2|-]",
	Short: "合并两个 JSON 对象（深合并）",
	Long: `将第二个 JSON 合并到第一个，返回合并结果。
嵌套对象深合并，同名字段以第二个为准。

如果第二个文件传 "-"，从 stdin 读取。
使用 -o 参数将结果写入文件。`,
	Example: `  tc json merge base.json patch.json
  echo '{"b":2}' | tc json merge base.json -
  tc json merge base.json patch.json -o result.json`,
	Args: cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		arg1 := args[0]
		var data1, data2 []byte

		data1, err := os.ReadFile(arg1)
		if err != nil {
			return fmt.Errorf("读取文件 %s 失败: %w", arg1, err)
		}

		if len(args) > 1 && args[1] != "-" {
			data2, err = os.ReadFile(args[1])
			if err != nil {
				return fmt.Errorf("读取文件 %s 失败: %w", args[1], err)
			}
		} else {
			piped, err := pipeline.ReadInput(nil)
			if err != nil {
				return err
			}
			data2 = piped
		}

		result, err := jsonutil.Merge(data1, data2)
		if err != nil {
			return err
		}
		return writeOutput(cmd, result+"\n")
	},
}

// writeOutput 输出结果：指定 -o 则写文件，否则输出到 stdout
func writeOutput(cmd *cobra.Command, data string) error {
	outputFile, _ := cmd.Flags().GetString("output")
	if outputFile != "" {
		if err := os.WriteFile(outputFile, []byte(data), 0644); err != nil {
			return fmt.Errorf("写入文件 %s 失败: %w", outputFile, err)
		}
		return nil
	}
	fmt.Print(data)
	return nil
}

func init() {
	jsonTableCmd.Flags().Bool("csv", false, "输出 CSV 格式（逗号分隔）")

	jsonSetCmd.Flags().StringP("output", "o", "", "输出到文件（默认输出到 stdout）")
	jsonDelCmd.Flags().StringP("output", "o", "", "输出到文件（默认输出到 stdout）")
	jsonMergeCmd.Flags().StringP("output", "o", "", "输出到文件（默认输出到 stdout）")

	jsonCmd.AddCommand(jsonFmtCmd)
	jsonCmd.AddCommand(jsonMinCmd)
	jsonCmd.AddCommand(jsonGetCmd)
	jsonCmd.AddCommand(jsonFilterCmd)
	jsonCmd.AddCommand(jsonKeysCmd)
	jsonCmd.AddCommand(jsonValuesCmd)
	jsonCmd.AddCommand(jsonPathsCmd)
	jsonCmd.AddCommand(jsonTableCmd)
	jsonCmd.AddCommand(jsonDiffCmd)
	jsonCmd.AddCommand(jsonSetCmd)
	jsonCmd.AddCommand(jsonDelCmd)
	jsonCmd.AddCommand(jsonMergeCmd)
}
