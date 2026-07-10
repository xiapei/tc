package genutil

import (
	"fmt"
	"math/rand"
	"strings"
)

// PasswordConfig 密码生成配置
type PasswordConfig struct {
	Length      int
	UseUpper    bool
	UseLower    bool
	UseDigit    bool
	UseSymbol   bool
	MinUpper    int
	MinLower    int
	MinDigit    int
	MinSymbol   int
}

const (
	upperChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lowerChars = "abcdefghijklmnopqrstuvwxyz"
	digitChars = "0123456789"
	symbolChars = "!@#$%^&*_+-=,.<>?"
)

// DefaultConfig 返回默认配置
func DefaultConfig() PasswordConfig {
	return PasswordConfig{
		Length:    16,
		UseUpper:  true,
		UseLower:  true,
		UseDigit:  true,
		UseSymbol: true,
	}
}

// GeneratePassword 生成密码
func GeneratePassword(cfg PasswordConfig) (string, error) {
	if cfg.Length < 1 {
		return "", fmt.Errorf("密码长度必须大于 0")
	}

	// 构建可用字符集
	var charset strings.Builder
	charCategories := make(map[rune]string) // 记录每个字符的来源类别

	if cfg.UseUpper {
		charset.WriteString(upperChars)
		for _, c := range upperChars {
			charCategories[c] = "upper"
		}
	}
	if cfg.UseLower {
		charset.WriteString(lowerChars)
		for _, c := range lowerChars {
			charCategories[c] = "lower"
		}
	}
	if cfg.UseDigit {
		charset.WriteString(digitChars)
		for _, c := range digitChars {
			charCategories[c] = "digit"
		}
	}
	if cfg.UseSymbol {
		charset.WriteString(symbolChars)
		for _, c := range symbolChars {
			charCategories[c] = "symbol"
		}
	}

	if charset.Len() == 0 {
		return "", fmt.Errorf("必须至少启用一种字符类型")
	}

	chars := charset.String()

	// 如果启用了多种类型，用洗牌确保混合
	if cfg.UseUpper && cfg.UseLower || cfg.UseDigit || cfg.UseSymbol {
		return generateMixed(chars, charCategories, cfg)
	}

	// 单一类型直接生成
	buf := make([]byte, cfg.Length)
	for i := 0; i < cfg.Length; i++ {
		buf[i] = chars[rand.Intn(len(chars))]
	}
	return string(buf), nil
}

func generateMixed(chars string, categories map[rune]string, cfg PasswordConfig) (string, error) {
	n := cfg.Length
	buf := make([]byte, n)

	// 收集各类别字符
	upperList := []byte(upperChars)
	lowerList := []byte(lowerChars)
	digitList := []byte(digitChars)
	symbolList := []byte(symbolChars)

	// 随机打乱位置
	positions := rand.Perm(n)

	idx := 0

	// 先放必要的小写字母
	if cfg.UseLower {
		for i := 0; i < min(cfg.MinLower, len(positions)-idx); i++ {
			buf[positions[idx]] = lowerList[rand.Intn(len(lowerList))]
			idx++
		}
	}
	// 必要的大写字母
	if cfg.UseUpper {
		for i := 0; i < min(cfg.MinUpper, len(positions)-idx); i++ {
			buf[positions[idx]] = upperList[rand.Intn(len(upperList))]
			idx++
		}
	}
	// 必要的数字
	if cfg.UseDigit {
		for i := 0; i < min(cfg.MinDigit, len(positions)-idx); i++ {
			buf[positions[idx]] = digitList[rand.Intn(len(digitList))]
			idx++
		}
	}
	// 必要的符号
	if cfg.UseSymbol {
		for i := 0; i < min(cfg.MinSymbol, len(positions)-idx); i++ {
			buf[positions[idx]] = symbolList[rand.Intn(len(symbolList))]
			idx++
		}
	}

	// 剩余位置随机填充
	allChars := []byte(chars)
	for i := idx; i < len(positions); i++ {
		buf[positions[i]] = allChars[rand.Intn(len(allChars))]
	}

	// 二次洗牌
	rand.Shuffle(n, func(i, j int) { buf[i], buf[j] = buf[j], buf[i] })

	return string(buf), nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Entropy 计算密码熵值（bit）
func Entropy(cfg PasswordConfig) float64 {
	var poolSize int
	if cfg.UseUpper {
		poolSize += 26
	}
	if cfg.UseLower {
		poolSize += 26
	}
	if cfg.UseDigit {
		poolSize += 10
	}
	if cfg.UseSymbol {
		poolSize += len(symbolChars)
	}
	if poolSize == 0 {
		return 0
	}
	return float64(cfg.Length) * (float64(poolSize) / float64(cfg.Length))
}

// PasswordStrength 返回密码强度评级
func PasswordStrength(entropy float64) string {
	switch {
	case entropy < 30:
		return "弱"
	case entropy < 50:
		return "中"
	case entropy < 70:
		return "强"
	default:
		return "极强"
	}
}

// FormatConfigSummary 返回配置摘要字符串
func (c PasswordConfig) FormatConfigSummary() string {
	var parts []string
	if c.UseUpper {
		parts = append(parts, "大写")
	}
	if c.UseLower {
		parts = append(parts, "小写")
	}
	if c.UseDigit {
		parts = append(parts, "数字")
	}
	if c.UseSymbol {
		parts = append(parts, "符号")
	}
	return fmt.Sprintf("%d位 (%s)", c.Length, strings.Join(parts, "+"))
}
