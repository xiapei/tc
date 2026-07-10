package convertutil

import (
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// JSONToYAML 将 JSON 转换为 YAML
func JSONToYAML(data []byte) (string, error) {
	// 解析 JSON
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return "", fmt.Errorf("JSON 解析失败: %w", err)
	}

	// 标准化 JSON（展开 \t 内的字符串等）
	v = normalizeJSON(v)

	// 转为 YAML
	out, err := yaml.Marshal(v)
	if err != nil {
		return "", fmt.Errorf("YAML 序列化失败: %w", err)
	}

	return string(out), nil
}

// YAMLToJSON 将 YAML 转换为 JSON（格式化）
func YAMLToJSON(data []byte) (string, error) {
	// 解析 YAML
	var v interface{}
	if err := yaml.Unmarshal(data, &v); err != nil {
		return "", fmt.Errorf("YAML 解析失败: %w", err)
	}

	// 标准化 YAML 结果为 JSON 兼容类型
	v = yamlToJSON(v)

	// 转为 JSON
	out, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", fmt.Errorf("JSON 序列化失败: %w", err)
	}

	return string(out), nil
}

// YAMLToJSONCompact 将 YAML 转换为 JSON（压缩）
func YAMLToJSONCompact(data []byte) (string, error) {
	var v interface{}
	if err := yaml.Unmarshal(data, &v); err != nil {
		return "", fmt.Errorf("YAML 解析失败: %w", err)
	}
	v = yamlToJSON(v)

	out, err := json.Marshal(v)
	if err != nil {
		return "", fmt.Errorf("JSON 序列化失败: %w", err)
	}

	return string(out), nil
}

// normalizeJSON 递归处理 JSON 值，确保一致性
func normalizeJSON(v interface{}) interface{} {
	switch val := v.(type) {
	case map[string]interface{}:
		m := make(map[string]interface{}, len(val))
		for k, vv := range val {
			m[k] = normalizeJSON(vv)
		}
		return m
	case []interface{}:
		for i, vv := range val {
			val[i] = normalizeJSON(vv)
		}
		return val
	case string:
		// 尝试解析 JSON 字符串中的嵌套 JSON
		return val
	default:
		return v
	}
}

// yamlToJSON 递归转换 YAML 结果为 JSON 兼容类型
// YAML 数值默认是 int，JSON 需要 float64
// YAML 可能生成 map[interface{}]interface{}，JSON 需要 map[string]interface{}
func yamlToJSON(v interface{}) interface{} {
	switch val := v.(type) {
	case map[string]interface{}:
		m := make(map[string]interface{}, len(val))
		for k, vv := range val {
			m[k] = yamlToJSON(vv)
		}
		return m
	case map[interface{}]interface{}:
		m := make(map[string]interface{}, len(val))
		for k, vv := range val {
			key := fmt.Sprintf("%v", k)
			m[key] = yamlToJSON(vv)
		}
		return m
	case []interface{}:
		for i, vv := range val {
			val[i] = yamlToJSON(vv)
		}
		return val
	case int:
		return float64(val)
	case int64:
		return float64(val)
	case float64:
		return val
	case bool:
		return val
	case string:
		return val
	case nil:
		return nil
	default:
		return fmt.Sprintf("%v", val)
	}
}

// DetectFormat 检测输入是否 JSON 或 YAML
func DetectFormat(data []byte) string {
	s := strings.TrimSpace(string(data))
	if len(s) == 0 {
		return "unknown"
	}

	// JSON 通常以 { 或 [ 开头
	if s[0] == '{' || s[0] == '[' {
		var js json.RawMessage
		if json.Unmarshal([]byte(s), &js) == nil {
			return "json"
		}
	}

	// YAML 检测：包含 yaml 特征
	if strings.Contains(s, ":") || strings.Contains(s, "---") {
		var v interface{}
		if yaml.Unmarshal([]byte(s), &v) == nil {
			if _, ok := v.(map[string]interface{}); ok {
				return "yaml"
			}
		}
	}

	return "unknown"
}
