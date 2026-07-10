# tc 命令速查表

## JSON 处理

| 命令 | 说明 | 示例 |
|------|------|------|
| `tc json fmt` | 格式化 | `echo '{"a":1}' \| tc json fmt` |
| `tc json min` | 压缩 | `cat data.json \| tc json min` |
| `tc json get <path>` | 查询 | `cat d.json \| tc json get "users.0.name"` |
| `tc json filter <expr>` | 过滤 | `cat d.json \| tc json filter 'age > 18'` |
| `tc json keys` | 所有 key | `echo '{"a":1,"b":2}' \| tc json keys` |
| `tc json values` | 所有 value | `echo '{"a":1,"b":2}' \| tc json values` |
| `tc json paths` | 所有路径 | `echo '{"a":{"b":1}}' \| tc json paths` |
| `tc json table <fields>` | 转表格 | `cat d.json \| tc json table "name,email"` |
| `tc json diff <f1> <f2>` | 比较差异 | `tc json diff old.json new.json` |

### 过滤表达式

```
age > 18          数值比较
name == "John"    字符串相等
status != "off"   不等于
email ~ "@gmail"  包含
phone             存在性检查
```

## 正则处理

| 命令 | 说明 | 示例 |
|------|------|------|
| `tc rx match <p>` | 高亮匹配 | `echo "abc123" \| tc rx match "\d+"` |
| `tc rx extract <p>` | 提取捕获组 | `echo "a=1" \| tc rx extract "(\w+)=(\w+)"` |
| `tc rx replace <p> <r>` | 替换 | `echo "foo" \| tc rx replace "foo" "bar"` |
| `tc rx grep <p>` | 过滤行 | `cat log \| tc rx grep "ERROR"` |
| `tc rx grep <p> -v` | 反转过滤 | `cat log \| tc rx grep "DEBUG" -v` |
| `tc rx count <p>` | 统计次数 | `cat log \| tc rx count "ERROR"` |
| `tc rx findall <p>` | 列出所有匹配 | `echo "a1b2" \| tc rx findall "\d"` |

## 通用工具

| 命令 | 说明 | 示例 |
|------|------|------|
| `tc hash <algo>` | 哈希 | `echo -n "hi" \| tc hash sha256` |
| `tc enc <type>` | 编码 | `echo "hi" \| tc enc base64` |
| `tc dec <type>` | 解码 | `echo "aGk=" \| tc dec base64` |
| `tc stats` | 统计 | `cat f.txt \| tc stats` |
| `tc count` | 行数 | `cat f.txt \| tc count` |
| `tc sort` | 排序 | `cat f.txt \| tc sort` |
| `tc sort -n` | 数值排序 | `echo -e "10\n2\n1" \| tc sort -n` |
| `tc sort -r` | 反转排序 | `cat f.txt \| tc sort -r` |
| `tc uniq` | 去重 | `cat f.txt \| tc sort \| tc uniq` |
| `tc head [n]` | 前 N 行 | `cat f.txt \| tc head 20` |
| `tc tail [n]` | 后 N 行 | `cat f.txt \| tc tail 20` |
| `tc sample [n]` | 随机采样 | `cat f.txt \| tc sample 100` |
| `tc fields <n>` | 提取字段 | `echo "a:b:c" \| tc fields 2 -s ":"` |

## 组合技

```bash
# Top 10 IP
cat access.log | tc rx extract "(\d+\.\d+\.\d+\.\d+)" | tc sort | tc uniq -c | tc sort -n -r | tc head 10

# 从日志提取 JSON 错误码
cat app.log | tc rx extract '(\{.*\})' | tc json get "error_code"

# JSON 过滤后提取邮箱
cat users.json | tc json filter 'age > 18' | tc json get "*.email"

# 日志中各状态码计数
cat access.log | tc rx extract 'HTTP/\d\.\d" (\d{3})' | tc sort | tc uniq -c | tc sort -n -r
```
