# tc 测试用例

## 1. JSON 处理

### 1.1 格式化 (json fmt)

| 用例 | 命令 | 预期结果 |
|------|------|----------|
| 简单对象格式化 | `echo '{"name":"John","age":30}' | tc json fmt` | 缩进格式化输出 |
| 嵌套对象格式化 | `echo '{"user":{"name":"John","address":{"city":"Beijing"}}}' | tc json fmt` | 多级缩进输出 |
| 无效 JSON | `echo '{invalid}' | tc json fmt` | 报错: "无效的 JSON" |

```bash
$ echo '{"name":"John","age":30}' | tc json fmt
{
  "name": "John",
  "age": 30
}
```

### 1.2 压缩 (json min)

| 用例 | 命令 | 预期结果 |
|------|------|----------|
| 多行压缩为单行 | `echo -e '{\n "name": "John"\n}' | tc json min` | `{"name":"John"}` |

```bash
$ echo '{
  "name": "John"
}' | tc json min
{"name":"John"}
```

### 1.3 查询 (json get)

| 用例 | 命令 | 预期结果 |
|------|------|----------|
| 嵌套字段 | `echo '{"user":{"name":"John"}}' | tc json get "user.name"` | `John` |
| 数组索引 | `echo '{"items":["a","b"]}' | tc json get "items.0"` | `"a"` |
| 通配符 | `echo '{"users":[{"n":"A"},{"n":"B"}]}' | tc json get "users.*.n"` | `["A","B"]` |
| 不存在的路径 | `echo '{"a":1}' | tc json get "b"` | `null` |

```bash
$ echo '{"items":["a","b"]}' | tc json get "items.0"
"a"
```

### 1.4 过滤 (json filter)

| 用例 | 命令 | 预期结果 |
|------|------|----------|
| 数值大于 | `echo '[{"n":"A","age":25},{"n":"B","age":17}]' | tc json filter 'age > 18'` | `[{"n":"A","age":25}]` |
| 字符串等于 | `echo '[{"s":"active"},{"s":"inactive"}]' | tc json filter 's == "active"'` | `[{"s":"active"}]` |
| 模糊包含 | `echo '[{"e":"a@gmail.com"},{"e":"b@yahoo.com"}]' | tc json filter 'e ~ "@gmail"'` | `[{"e":"a@gmail.com"}]` |
| 非数组输入 | `echo '{"a":1}' | tc json filter 'a > 0'` | 报错: "输入不是 JSON 数组" |

```bash
$ echo '[{"n":"A","age":25},{"n":"B","age":17}]' | tc json filter 'age > 18'
[{"n":"A","age":25}]
```

### 1.5 提取 Key/Value (json keys / values / paths)

| 用例 | 命令 | 预期结果 |
|------|------|----------|
| 提取 keys | `echo '{"name":"John","age":30}' | tc json keys` | `["name","age"]` |
| 提取 values | `echo '{"name":"John","age":30}' | tc json values` | `["John",30]` |
| 列出 paths | `echo '{"a":{"b":1},"c":[2,3]}' | tc json paths` | 逐行输出路径 |
| 非对象输入 keys | `echo '["a","b"]' | tc json keys` | 报错: "输入不是 JSON 对象" |

```bash
$ echo '{"a":{"b":1},"c":[2,3]}' | tc json paths
a
a.b
c
c.0
c.1
```

### 1.6 表格输出 (json table)

| 用例 | 命令 | 预期结果 |
|------|------|----------|
| TSV 输出 | `echo '[{"n":"A","age":25},{"n":"B","age":30}]' | tc json table "n,age"` | 制表符分隔表格 |
| CSV 输出 | `echo '[{"n":"A","age":25}]' | tc json table "n,age" --csv` | 逗号分隔表格 |
| 不存在的字段 | `echo '[{"n":"A"}]' | tc json table "n,x"` | 空单元格 |

```bash
$ echo '[{"n":"A","age":25},{"n":"B","age":30}]' | tc json table "n,age"
n	age
A	25
B	30
```

### 1.7 Diff (json diff)

| 用例 | 命令 | 预期结果 |
|------|------|----------|
| 有差异 | `tc json diff old.json new.json` | 输出差异行(-/+) |
| 无差异 | 两个相同文件 | `无差异` |

---

## 2. 正则处理

### 2.1 高亮匹配 (rx match)

| 用例 | 命令 | 预期结果 |
|------|------|----------|
| 匹配数字 | `echo "abc123def456" | tc rx match "\d+"` | 高亮显示 `123` 和 `456` |
| 无匹配 | `echo "hello" | tc rx match "\d+"` | 空输出 |
| 无效正则 | `echo "hello" | tc rx match "[invalid"` | 报错: "无效的正则表达式" |

### 2.2 提取捕获组 (rx extract)

| 用例 | 命令 | 预期结果 |
|------|------|----------|
| 单个捕获组 | `echo "2024-01-15 error" | tc rx extract "(\d{4}-\d{2}-\d{2})"` | `2024-01-15` |
| 多个捕获组 | `echo "name=John age=30" | tc rx extract "(\w+)=(\w+)"` | `name\tJohn\nage\t30` |
| 无捕获组 | `echo "abc123" | tc rx extract "\d+"` | `123` |

### 2.3 替换 (rx replace)

| 用例 | 命令 | 预期结果 |
|------|------|----------|
| 简单替换 | `echo "hello world" | tc rx replace "world" "Go"` | `hello Go` |
| 引用捕获组 | `echo "foo123bar" | tc rx replace "(\d+)" "NUM:$1"` | `fooNUM:123bar` |
| 全局替换 | `echo "a1b2c3" | tc rx replace "\d" "X"` | `aXbXcX` |

### 2.4 过滤行 (rx grep)

| 用例 | 命令 | 预期结果 |
|------|------|----------|
| 匹配行 | `echo -e "ERROR: a\nINFO: b\nERROR: c" | tc rx grep "ERROR"` | `ERROR: a\nERROR: c` |
| 反转匹配 | `echo -e "ERROR: a\nINFO: b" | tc rx grep "ERROR" --invert` | `INFO: b` |
| 多模式 | `echo -e "a\nb\nc" | tc rx grep "a|c"` | `a\nc` |

### 2.5 统计 (rx count)

| 用例 | 命令 | 预期结果 |
|------|------|----------|
| 统计出现次数 | `echo "a1b2c3" | tc rx count "\d"` | `3` |
| 无匹配 | `echo "abc" | tc rx count "\d"` | `0` |

### 2.6 列出所有匹配 (rx findall)

| 用例 | 命令 | 预期结果 |
|------|------|----------|
| 多匹配 | `echo "a1b2c3" | tc rx findall "\d"` | `1\n2\n3` |
| 日期匹配 | `echo "2024-01-15 and 2024-12-25" | tc rx findall "\d{4}-\d{2}-\d{2}"` | 逐行输出日期 |

---

## 3. 通用工具

### 3.1 哈希 (hash)

| 用例 | 命令 | 预期结果 |
|------|------|----------|
| SHA256 | `echo -n "hello" | tc hash sha256` | `2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824` |
| MD5 | `echo -n "hello" | tc hash md5` | `5d41402abc4b2a76b9719d911017c592` |
| 不支持的算法 | `echo "hello" | tc hash sha1` | 报错: "不支持的算法" |

### 3.2 编码/解码 (enc / dec)

| 用例 | 命令 | 预期结果 |
|------|------|----------|
| base64 编码 | `echo -n "hello world" | tc enc base64` | `aGVsbG8gd29ybGQ=` |
| base64 解码 | `echo "aGVsbG8gd29ybGQ=" | tc dec base64` | `hello world` |
| url 编码 | `echo -n "hello world" | tc enc url` | `hello+world` |
| url 解码 | `echo "hello+world" | tc dec url` | `hello world` |

### 3.3 统计 (stats)

| 用例 | 命令 | 预期结果 |
|------|------|----------|
| 文本统计 | `echo -e "hello world\nfoo bar" | tc stats` | `2\t4\t20\t20` (行/词/字/字节) |
| 空输入 | `echo -n "" | tc stats` | `0\t0\t0\t0` |

### 3.4 排序 (sort)

| 用例 | 命令 | 预期结果 |
|------|------|----------|
| 字母排序 | `echo -e "c\na\nb" | tc sort` | `a\nb\nc` |
| 数值排序 | `echo -e "10\n2\n30" | tc sort -n` | `2\n10\n30` |
| 反转排序 | `echo -e "a\nb\nc" | tc sort -r` | `c\nb\na` |

### 3.5 去重 (uniq)

| 用例 | 命令 | 预期结果 |
|------|------|----------|
| 相邻去重 | `echo -e "a\na\nb\nc\nc" | tc uniq` | `a\nb\nc` |
| 计数 | `echo -e "a\na\nb\nc\nc\nc" | tc uniq -c` | `2\ta\n1\tb\n3\tc` |

### 3.6 切片 (head / tail)

| 用例 | 命令 | 预期结果 |
|------|------|----------|
| head 默认 10 行 | `seq 1 20 | tc head` | 输出 1-10 |
| head 指定行数(位置参数) | `seq 1 20 | tc head 3` | `1\n2\n3` |
| head 指定行数(flag) | `seq 1 20 | tc head -n 3` | `1\n2\n3` |
| tail 默认 10 行 | `seq 1 20 | tc tail` | 输出 11-20 |
| tail 指定行数(位置参数) | `seq 1 20 | tc tail 3` | `18\n19\n20` |

```bash
$ echo -e "a\nb\nc\nd\ne\nf" | tc head 3
a
b
c

$ echo -e "a\nb\nc\nd\ne\nf" | tc tail 2
e
f
```

### 3.7 随机采样 (sample)

| 用例 | 命令 | 预期结果 |
|------|------|----------|
| 小于总数 | `seq 1 10 | tc sample 3` | 输出 3 行随机行 |
| 大于等于总数 | `seq 1 3 | tc sample 5` | 输出全部 3 行 |
| 无效参数 | `echo "a\nb" | tc sample 0` | 报错 |
| 随机性验证 | `seq 1 5 | tc sample 2` (执行两次) | 两次结果不一致 |

```bash
$ echo -e "1\n2\n3\n4\n5" | tc sample 2
3
4

$ echo -e "1\n2\n3\n4\n5" | tc sample 2
1
5
```

### 3.8 字段提取 (fields)

| 用例 | 命令 | 预期结果 |
|------|------|----------|
| 默认分隔符(tab) | `echo -e "a\tb\tc" | tc fields 2` | `b` |
| 自定义分隔符 | `echo "a:b:c:d" | tc fields 2 --sep ":"` | `b` |
| 字段超出范围 | `echo "a:b" | tc fields 5 --sep ":"` | 空输出 |

### 3.9 行数统计 (count)

| 用例 | 命令 | 预期结果 |
|------|------|----------|
| 多行 | `echo -e "a\nb\nc" | tc count` | `3` |
| 空输入 | `echo -n "" | tc count` | `0` |

### 3.10 快捷 fmt

| 用例 | 命令 | 预期结果 |
|------|------|----------|
| 根级别 fmt | `echo '{"a":1}' | tc fmt` | 格式化 JSON（同 `tc json fmt`） |

---

## 4. 管道组合（杀手功能）

| 用例 | 命令 | 预期结果 |
|------|------|----------|
| JSON + 正则 | `cat app.log | tc rx extract '(\{.*\})' | tc json get "error_code"` | 从日志提取 JSON 再查字段 |
| 过滤 + 提取 | `cat users.json | tc json filter 'age > 18' | tc json get "*.email"` | 过滤后提取 emails |
| 统计 Top N | `cat access.log | tc rx extract "(\d+\.\d+\.\d+\.\d+)" | tc sort | tc uniq -c | tc sort -n -r | tc head 10` | Top 10 IP |
