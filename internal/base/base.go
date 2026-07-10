package baseutil

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Convert 进制转换，将 value 从 fromBase 转换到 toBase，返回各进制表示
type ConvertResult struct {
	Decimal int64
	Hex     string
	Octal   string
	Binary  string
	Base36  string
}

// Convert 将 value 从指定进制转换为十进制，并显示各进制结果
func Convert(value string, fromBase int) (*ConvertResult, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, fmt.Errorf("值不能为空")
	}

	// 忽略常见前缀
	value = strings.TrimPrefix(value, "0x")
	value = strings.TrimPrefix(value, "0X")
	value = strings.TrimPrefix(value, "0b")
	value = strings.TrimPrefix(value, "0B")
	value = strings.TrimPrefix(value, "0o")
	value = strings.TrimPrefix(value, "0O")

	// 解析为十进制
	n, err := strconv.ParseInt(value, fromBase, 64)
	if err != nil {
		// 尝试无符号
		u, err2 := strconv.ParseUint(value, fromBase, 64)
		if err2 != nil {
			return nil, fmt.Errorf("无法将 %q 从 %s 进制转换为数值", value, baseName(fromBase))
		}
		n = int64(u)
		_ = err // 忽略有符号的错误
	}

	return &ConvertResult{
		Decimal: n,
		Hex:     strings.ToUpper(strconv.FormatInt(n, 16)),
		Octal:   strconv.FormatInt(n, 8),
		Binary:  strconv.FormatInt(n, 2),
		Base36:  strings.ToUpper(strconv.FormatInt(n, 36)),
	}, nil
}

// FormatRange 输出范围展示（显示从-from到-to的转换）
func FormatRange(results []*ConvertResult, toBase int) string {
	var sb strings.Builder
	for _, r := range results {
		sb.WriteString(formatSingle(r, toBase) + "\n")
	}
	return sb.String()
}

func formatSingle(r *ConvertResult, toBase int) string {
	switch toBase {
	case 2:
		return fmt.Sprintf("%d (10) = %s (2)", r.Decimal, r.Binary)
	case 8:
		return fmt.Sprintf("%d (10) = %s (8)", r.Decimal, r.Octal)
	case 10:
		return fmt.Sprintf("%s (16) = %d (10)", "0x"+r.Hex, r.Decimal)
	case 16:
		return fmt.Sprintf("%d (10) = %s (16)", r.Decimal, "0x"+r.Hex)
	default:
		return fmt.Sprintf("%d (10) = %s (%d)", r.Decimal, strconv.FormatInt(r.Decimal, toBase), toBase)
	}
}

// ShowAll 显示所有进制表示
func ShowAll(r *ConvertResult) string {
	return fmt.Sprintf(`  十进制: %d
  十六进制: 0x%s
  八进制: 0o%s
  二进制: 0b%s
  Base36: %s`, r.Decimal, r.Hex, r.Octal, r.Binary, r.Base36)
}

func baseName(base int) string {
	switch base {
	case 2:
		return "二进制"
	case 8:
		return "八进制"
	case 10:
		return "十进制"
	case 16:
		return "十六进制"
	default:
		return strconv.Itoa(base)
	}
}

// ParseAuto 自动推断进制并转换（0x→16, 0b→2, 纯数字→10）
func ParseAuto(value string) (*ConvertResult, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, fmt.Errorf("值不能为空")
	}

	if strings.HasPrefix(value, "0x") || strings.HasPrefix(value, "0X") {
		return Convert(value, 16)
	}
	if strings.HasPrefix(value, "0b") || strings.HasPrefix(value, "0B") {
		return Convert(value, 2)
	}
	if strings.HasPrefix(value, "0o") || strings.HasPrefix(value, "0O") {
		return Convert(value, 8)
	}
	// 尝试十进制
	return Convert(value, 10)
}

// NormalizeBase 规范化进制：输入 16 或 "hex" 都返回 16
func NormalizeBase(base int) int {
	if base <= 0 || base > 36 {
		if math.Abs(float64(base)) == 16 || base == -16 {
			return 16
		}
		return 10 // 默认十进制
	}
	return base
}
