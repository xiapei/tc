package jsonutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// Format 格式化 JSON
func Format(data []byte) (string, error) {
	var out bytes.Buffer
	if err := json.Indent(&out, bytes.TrimSpace(data), "", "  "); err != nil {
		return "", fmt.Errorf("无效的 JSON: %w", err)
	}
	return out.String(), nil
}

// Minify 压缩 JSON
func Minify(data []byte) (string, error) {
	var out bytes.Buffer
	if err := json.Compact(&out, bytes.TrimSpace(data)); err != nil {
		return "", fmt.Errorf("无效的 JSON: %w", err)
	}
	return out.String(), nil
}

// Get 按路径查询 JSON
func Get(data []byte, path string) (string, error) {
	result := gjson.GetBytes(data, path)
	if !result.Exists() {
		return "null", nil
	}
	if result.Type == gjson.String {
		return result.String(), nil
	}
	return result.Raw, nil
}

// condition 表示一个简单条件: field op value
type condition struct {
	field, op, value string
}

// splitRespectingQuotes 按分隔符拆分字符串，尊重引号内的内容
func splitRespectingQuotes(s, sep string) []string {
	var parts []string
	inQuote := false
	var quoteChar byte
	start := 0

	for i := 0; i < len(s); i++ {
		if inQuote {
			if s[i] == quoteChar {
				inQuote = false
			}
			continue
		}
		if s[i] == '"' || s[i] == '\'' {
			inQuote = true
			quoteChar = s[i]
			continue
		}
		if i+len(sep) <= len(s) && s[i:i+len(sep)] == sep {
			parts = append(parts, s[start:i])
			i += len(sep) - 1
			start = i + 1
		}
	}
	parts = append(parts, s[start:])
	return parts
}

// parseFilterExpr 解析表达式为 OR 组 + AND 组
// 返回 [][]condition，外层 OR，内层 AND
// 例: "a>1 && b==2 || c<3" → [[{a,>,1},{b,==,2}], [{c,<,3}]]
func parseFilterExpr(expr string) ([][]condition, error) {
	orParts := splitRespectingQuotes(expr, "||")
	var groups [][]condition
	for _, orPart := range orParts {
		andParts := splitRespectingQuotes(strings.TrimSpace(orPart), "&&")
		var group []condition
		for _, andPart := range andParts {
			field, op, value, err := parseExpr(strings.TrimSpace(andPart))
			if err != nil {
				return nil, err
			}
			group = append(group, condition{field: field, op: op, value: value})
		}
		groups = append(groups, group)
	}
	return groups, nil
}

// matchFilter 根据条件组判断 item 是否匹配
func matchFilter(item gjson.Result, groups [][]condition) bool {
	for _, group := range groups {
		allMatch := true
		for _, c := range group {
			fieldVal := item.Get(c.field)
			if !fieldVal.Exists() || !matchCondition(fieldVal, c.op, c.value) {
				allMatch = false
				break
			}
		}
		if allMatch {
			return true
		}
	}
	return false
}

// Filter 按条件过滤 JSON 数组，支持 && 和 || 复合条件
func Filter(data []byte, expr string) (string, error) {
	result := gjson.ParseBytes(data)
	if !result.IsArray() {
		return "", fmt.Errorf("输入不是 JSON 数组")
	}

	groups, err := parseFilterExpr(expr)
	if err != nil {
		return "", err
	}

	var filtered []string
	result.ForEach(func(_, item gjson.Result) bool {
		if matchFilter(item, groups) {
			filtered = append(filtered, item.Raw)
		}
		return true
	})

	return "[" + strings.Join(filtered, ",") + "]", nil
}

// Set 设置 JSON 路径的值
func Set(data []byte, path string, rawValue string) (string, error) {
	var v interface{}
	if err := json.Unmarshal([]byte(rawValue), &v); err != nil {
		v = rawValue
	}
	result, err := sjson.SetBytes(data, path, v)
	if err != nil {
		return "", fmt.Errorf("设置 JSON 路径失败: %w", err)
	}
	return Format(result)
}

// Delete 删除 JSON 路径
func Delete(data []byte, path string) (string, error) {
	result, err := sjson.DeleteBytes(data, path)
	if err != nil {
		return "", fmt.Errorf("删除 JSON 路径失败: %w", err)
	}
	return Format(result)
}

// Merge 合并两个 JSON 对象（深合并）
func Merge(base, patch []byte) (string, error) {
	var baseMap, patchMap map[string]interface{}
	if err := json.Unmarshal(base, &baseMap); err != nil {
		return "", fmt.Errorf("第一个 JSON 无效: %w", err)
	}
	if err := json.Unmarshal(patch, &patchMap); err != nil {
		return "", fmt.Errorf("第二个 JSON 无效: %w", err)
	}
	deepMerge(baseMap, patchMap)
	result, err := json.MarshalIndent(baseMap, "", "  ")
	if err != nil {
		return "", err
	}
	return string(result), nil
}

func deepMerge(base, patch map[string]interface{}) {
	for k, v := range patch {
		if baseVal, ok := base[k]; ok {
			baseMap, baseOk := baseVal.(map[string]interface{})
			patchMap, patchOk := v.(map[string]interface{})
			if baseOk && patchOk {
				deepMerge(baseMap, patchMap)
				continue
			}
		}
		base[k] = v
	}
}

// Keys 提取 JSON 对象的所有 key
func Keys(data []byte) (string, error) {
	result := gjson.ParseBytes(data)
	if !result.IsObject() {
		return "", fmt.Errorf("输入不是 JSON 对象")
	}
	var keys []string
	result.ForEach(func(key, _ gjson.Result) bool {
		keys = append(keys, key.String())
		return true
	})
	out, _ := json.Marshal(keys)
	return string(out), nil
}

// Values 提取 JSON 对象的所有 value
func Values(data []byte) (string, error) {
	result := gjson.ParseBytes(data)
	if !result.IsObject() {
		return "", fmt.Errorf("输入不是 JSON 对象")
	}
	var vals []json.RawMessage
	result.ForEach(func(_, val gjson.Result) bool {
		vals = append(vals, json.RawMessage(val.Raw))
		return true
	})
	out, _ := json.Marshal(vals)
	return string(out), nil
}

// Paths 列出 JSON 的所有路径
func Paths(data []byte) ([]string, error) {
	result := gjson.ParseBytes(data)
	var paths []string
	walkPaths(result, "", &paths)
	sort.Strings(paths)
	return paths, nil
}

func walkPaths(node gjson.Result, prefix string, paths *[]string) {
	if prefix != "" {
		*paths = append(*paths, prefix)
	}
	if node.IsObject() {
		node.ForEach(func(key, val gjson.Result) bool {
			childPath := prefix + "." + key.String()
			if prefix == "" {
				childPath = key.String()
			}
			walkPaths(val, childPath, paths)
			return true
		})
	} else if node.IsArray() {
		for i := 0; i < len(node.Array()); i++ {
			childPath := prefix + "." + strconv.Itoa(i)
			if prefix == "" {
				childPath = strconv.Itoa(i)
			}
			walkPaths(node.Array()[i], childPath, paths)
		}
	}
}

// Table 将 JSON 数组转为表格
func Table(data []byte, fieldsStr string, csv bool) (string, error) {
	result := gjson.ParseBytes(data)
	if !result.IsArray() {
		return "", fmt.Errorf("输入不是 JSON 数组")
	}

	fields := strings.Split(fieldsStr, ",")
	for i := range fields {
		fields[i] = strings.TrimSpace(fields[i])
	}

	sep := "\t"
	if csv {
		sep = ","
	}

	var buf bytes.Buffer
	// 表头
	buf.WriteString(strings.Join(fields, sep) + "\n")

	// 数据行
	result.ForEach(func(_, item gjson.Result) bool {
		var row []string
		for _, f := range fields {
			val := item.Get(f)
			if val.Exists() {
				row = append(row, val.String())
			} else {
				row = append(row, "")
			}
		}
		buf.WriteString(strings.Join(row, sep) + "\n")
		return true
	})

	return buf.String(), nil
}

// Diff 比较两个 JSON 的差异
func Diff(data1, data2 []byte) (string, error) {
	var v1, v2 interface{}
	if err := json.Unmarshal(data1, &v1); err != nil {
		return "", fmt.Errorf("第一个 JSON 无效: %w", err)
	}
	if err := json.Unmarshal(data2, &v2); err != nil {
		return "", fmt.Errorf("第二个 JSON 无效: %w", err)
	}

	var diffs []string
	diffWalk(v1, v2, "", &diffs)

	if len(diffs) == 0 {
		return "无差异\n", nil
	}

	var buf bytes.Buffer
	for _, d := range diffs {
		buf.WriteString(d + "\n")
	}
	return buf.String(), nil
}

func diffWalk(v1, v2 interface{}, path string, diffs *[]string) {
	switch a := v1.(type) {
	case map[string]interface{}:
		b, ok := v2.(map[string]interface{})
		if !ok {
			*diffs = append(*diffs, fmt.Sprintf("- %s: %v", path, v1))
			*diffs = append(*diffs, fmt.Sprintf("+ %s: %v", path, v2))
			return
		}
		allKeys := make(map[string]bool)
		for k := range a {
			allKeys[k] = true
		}
		for k := range b {
			allKeys[k] = true
		}
		for k := range allKeys {
			p := path + "." + k
			if path == "" {
				p = k
			}
			va, okA := a[k]
			vb, okB := b[k]
			if !okA {
				*diffs = append(*diffs, fmt.Sprintf("+ %s: %v", p, vb))
			} else if !okB {
				*diffs = append(*diffs, fmt.Sprintf("- %s: %v", p, va))
			} else {
				diffWalk(va, vb, p, diffs)
			}
		}
	case []interface{}:
		b, ok := v2.([]interface{})
		if !ok {
			*diffs = append(*diffs, fmt.Sprintf("- %s: %v", path, v1))
			*diffs = append(*diffs, fmt.Sprintf("+ %s: %v", path, v2))
			return
		}
		maxLen := len(a)
		if len(b) > maxLen {
			maxLen = len(b)
		}
		for i := 0; i < maxLen; i++ {
			p := path + "[" + strconv.Itoa(i) + "]"
			if i >= len(a) {
				*diffs = append(*diffs, fmt.Sprintf("+ %s: %v", p, b[i]))
			} else if i >= len(b) {
				*diffs = append(*diffs, fmt.Sprintf("- %s: %v", p, a[i]))
			} else {
				diffWalk(a[i], b[i], p, diffs)
			}
		}
	default:
		if fmt.Sprintf("%v", v1) != fmt.Sprintf("%v", v2) {
			*diffs = append(*diffs, fmt.Sprintf("- %s: %v", path, v1))
			*diffs = append(*diffs, fmt.Sprintf("+ %s: %v", path, v2))
		}
	}
}

// parseExpr 解析简单条件表达式: field op value
func parseExpr(expr string) (field, op, value string, err error) {
	ops := []string{"!=", ">=", "<=", "==", ">", "<", "~"}
	for _, o := range ops {
		idx := strings.Index(expr, o)
		if idx > 0 {
			field = strings.TrimSpace(expr[:idx])
			value = strings.TrimSpace(expr[idx+len(o):])
			// 去掉引号
			value = strings.Trim(value, `"'`)
			return field, o, value, nil
		}
	}
	// 没有操作符，当作"存在性"检查
	return strings.TrimSpace(expr), "exists", "", nil
}

// matchCondition 匹配条件
func matchCondition(val gjson.Result, op, value string) bool {
	switch op {
	case "exists":
		return val.Exists()
	case "==":
		return val.String() == value
	case "!=":
		return val.String() != value
	case "~":
		return strings.Contains(val.String(), value)
	case ">", "<", ">=", "<=":
		f1, err1 := strconv.ParseFloat(val.String(), 64)
		f2, err2 := strconv.ParseFloat(value, 64)
		if err1 != nil || err2 != nil {
			return false
		}
		switch op {
		case ">":
			return f1 > f2
		case "<":
			return f1 < f2
		case ">=":
			return f1 >= f2
		case "<=":
			return f1 <= f2
		}
	}
	return false
}
