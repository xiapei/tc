package rxutil

import (
	"fmt"
	"regexp"
	"strings"
)

// Match 高亮匹配内容
func Match(data []byte, pattern string) (string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", fmt.Errorf("无效的正则表达式: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	var matched []string
	for _, line := range lines {
		if re.MatchString(line) {
			// 高亮匹配部分
			highlighted := re.ReplaceAllString(line, "\033[1;33m$0\033[0m")
			matched = append(matched, highlighted)
		}
	}
	if len(matched) == 0 {
		return "", nil
	}
	return strings.Join(matched, "\n") + "\n", nil
}

// Extract 提取捕获组
func Extract(data []byte, pattern string) (string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", fmt.Errorf("无效的正则表达式: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	var results []string
	for _, line := range lines {
		matches := re.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			if len(match) > 1 {
				// 有捕获组，输出所有组（tab 分隔）
				results = append(results, strings.Join(match[1:], "\t"))
			} else {
				// 无捕获组，输出整个匹配
				results = append(results, match[0])
			}
		}
	}
	if len(results) == 0 {
		return "\n", nil
	}
	return strings.Join(results, "\n") + "\n", nil
}

// Replace 替换匹配内容
func Replace(data []byte, pattern, replacement string) (string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", fmt.Errorf("无效的正则表达式: %w", err)
	}

	result := re.ReplaceAllString(string(data), replacement)
	return result, nil
}

// Grep 过滤匹配行
func Grep(data []byte, pattern string, invert bool) (string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", fmt.Errorf("无效的正则表达式: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	var filtered []string
	for _, line := range lines {
		match := re.MatchString(line)
		if invert {
			match = !match
		}
		if match {
			filtered = append(filtered, line)
		}
	}
	return strings.Join(filtered, "\n") + "\n", nil
}

// Count 统计匹配次数
func Count(data []byte, pattern string) (int, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return 0, fmt.Errorf("无效的正则表达式: %w", err)
	}

	return len(re.FindAllString(string(data), -1)), nil
}

// FindAll 列出所有匹配
func FindAll(data []byte, pattern string) ([]string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("无效的正则表达式: %w", err)
	}

	matches := re.FindAllString(string(data), -1)
	return matches, nil
}
