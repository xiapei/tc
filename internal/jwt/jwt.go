package jwtutil

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

// Token JWT 解码结果
type Token struct {
	Raw       string
	Header    map[string]interface{}
	Payload   map[string]interface{}
	Signature string
}

// Decode 解码 JWT（不验证签名）
func Decode(tokenStr string) (*Token, error) {
	tokenStr = strings.TrimSpace(tokenStr)
	parts := strings.Split(tokenStr, ".")

	if len(parts) < 2 {
		return nil, fmt.Errorf("无效的 JWT 格式：至少需要包含两个点分隔的部分（header.payload）")
	}
	if len(parts) < 3 {
		return nil, fmt.Errorf("无效的 JWT 格式：缺少签名部分（header.payload.signature）")
	}

	t := &Token{Raw: tokenStr}

	// 解码 Header
	headerJSON, err := base64Decode(parts[0])
	if err != nil {
		return nil, fmt.Errorf("header 解码失败: %w", err)
	}
	if err := json.Unmarshal(headerJSON, &t.Header); err != nil {
		return nil, fmt.Errorf("header JSON 解析失败: %w", err)
	}

	// 解码 Payload
	payloadJSON, err := base64Decode(parts[1])
	if err != nil {
		return nil, fmt.Errorf("payload 解码失败: %w", err)
	}
	if err := json.Unmarshal(payloadJSON, &t.Payload); err != nil {
		return nil, fmt.Errorf("payload JSON 解析失败: %w", err)
	}

	// 签名原文
	t.Signature = parts[2]

	return t, nil
}

// base64Decode JWT 使用 URL-safe Base64，需要补充填充
func base64Decode(s string) ([]byte, error) {
	// JWT 使用 URL-safe Base64（- 代替 +，_ 代替 /）
	// 且无填充，需要补全
	s = strings.ReplaceAll(s, "-", "+")
	s = strings.ReplaceAll(s, "_", "/")

	// 补全填充
	switch len(s) % 4 {
	case 2:
		s += "=="
	case 3:
		s += "="
	}

	return base64.StdEncoding.DecodeString(s)
}

// Format 格式化输出 JWT 解码结果
func Format(t *Token) (string, error) {
	var sb strings.Builder

	// Header
	sb.WriteString("=== Header ===\n")
	h, err := json.MarshalIndent(t.Header, "", "  ")
	if err != nil {
		return "", err
	}
	sb.WriteString(string(h))
	sb.WriteString("\n\n")

	// Payload
	sb.WriteString("=== Payload ===\n")
	p, err := json.MarshalIndent(t.Payload, "", "  ")
	if err != nil {
		return "", err
	}
	sb.WriteString(string(p))
	sb.WriteString("\n\n")

	// Signature
	sb.WriteString("=== Signature (raw) ===\n")
	sig := t.Signature
	if len(sig) > 40 {
		sig = sig[:40] + "..."
	}
	sb.WriteString(sig + "\n")
	sb.WriteString("\n")

	// 警告
	sb.WriteString("⚠  JWT 解码未验证签名，请勿在不可信来源中使用本工具处理敏感信息\n")

	return sb.String(), nil
}

// IsJWT 快速判断是否是 JWT 格式
func IsJWT(s string) bool {
	parts := strings.Split(strings.TrimSpace(s), ".")
	return len(parts) == 3 && len(parts[0]) > 0 && len(parts[1]) > 0 && len(parts[2]) > 0
}
