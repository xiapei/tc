package pipeline

import (
	"fmt"
	"io"
	"os"
)

// ReadInput 从文件参数或 stdin 读取数据
// 如果 args 非空，读取第一个文件；否则读取 stdin
func ReadInput(args []string) ([]byte, error) {
	if len(args) > 0 {
		data, err := os.ReadFile(args[0])
		if err != nil {
			return nil, fmt.Errorf("读取文件 %s 失败: %w", args[0], err)
		}
		return data, nil
	}
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return nil, fmt.Errorf("读取 stdin 失败: %w", err)
	}
	return data, nil
}
