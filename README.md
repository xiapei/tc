# tc — ToolCraft

> 程序员的文本/数据瑞士军刀。JSON + 正则 + 文本处理，一个命令搞定。

[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

---

## 为什么用 tc？

| 痛点 | tc 的解决方式 |
|------|-------------|
| jq 语法记不住 | `tc json get "users.0.name"` — 直觉语法 |
| grep 不懂 JSON | tc 理解 JSON 结构，能按字段过滤 |
| 正则+JSON 组合要写脚本 | 管道一行搞定 |
| 工具太多记不住 | 一个 `tc` 全搞定 |

## 安装

```bash
# Go（推荐）
go install github.com/user/tc@latest

# Homebrew
brew install user/tc/tc

# 手动下载
# https://github.com/user/tc/releases
```

## 快速上手

### JSON 处理

```bash
# 格式化
echo '{"name":"John","age":30}' | tc json fmt

# 查询
cat data.json | tc json get "users.0.email"
cat data.json | tc json get "users.*.name"

# 过滤
cat users.json | tc json filter 'age > 18'
cat users.json | tc json filter 'status == "active"'
cat users.json | tc json filter 'email ~ "@gmail"'

# 转换
cat data.json | tc json keys       # 所有 key
cat data.json | tc json values     # 所有 value
cat data.json | tc json paths      # 所有路径

# 表格输出
cat users.json | tc json table "name,email,age"
cat users.json | tc json table "id,title" --csv

# Diff
tc json diff old.json new.json
```

### 正则处理

```bash
# 高亮匹配
echo "abc123def456" | tc rx match "\d+"
cat access.log | tc rx match "5\d{2}"

# 提取捕获组
echo "2024-01-15 error occurred" | tc rx extract "(\d{4}-\d{2}-\d{2})"
echo "name=John age=30" | tc rx extract "(\w+)=(\w+)"

# 替换
echo "hello world" | tc rx replace "world" "Go"
cat file.txt | tc rx replace "(\d+)" "NUM:$1"

# 过滤行
cat access.log | tc rx grep "ERROR|WARN"
cat access.log | tc rx grep "DEBUG" --invert  # 反转

# 统计
cat app.log | tc rx count "ERROR"

# 列出所有匹配
echo "a1b2c3" | tc rx findall "\d"
```

### 组合技（杀手功能）

```bash
# 从日志提取 JSON 再查字段
cat app.log | tc rx extract '(\{.*\})' | tc json get "error_code"

# JSON 数组过滤后提取特定字段
cat users.json | tc json filter 'age > 18' | tc json get "*.email"

# 日志中统计各状态码出现次数
cat access.log | tc rx extract "HTTP/\d\.\d\" (\d{3})" | tc sort | tc uniq --count | tc sort -n -r

# Top 10 IP
cat access.log | tc rx extract "(\d+\.\d+\.\d+\.\d+)" | tc sort | tc uniq --count | tc sort -n -r | tc head 10
```

### 通用工具

```bash
# 哈希
echo -n "hello" | tc hash sha256
echo -n "hello" | tc hash md5

# 编码/解码
echo "hello world" | tc enc base64
echo "aGVsbG8gd29ybGQ=" | tc dec base64
echo "hello world" | tc enc url

# Base64 编解码（直接传参）
tc b64encode "hello world"           # → aGVsbG8gd29ybGQ=
tc b64decode "aGVsbG8gd29ybGQ="      # → hello world
tc b64encode --url "abc 123?"         # URL-safe 模式
echo -n "你好世界" | tc b64encode      # 支持中文
tc b64decode "5L2g5aW95LiW55WM"      # 中文解码

# 统计
cat file.txt | tc stats        # 行数 词数 字符数 字节数
cat file.txt | tc count        # 行数

# 排序/去重
cat file.txt | tc sort
cat file.txt | tc sort -n      # 数值排序
cat file.txt | tc sort | tc uniq

# 采样
cat bigfile.txt | tc sample 100

# 切片
cat file.txt | tc head 10
cat file.txt | tc tail 10

# 字段提取
echo "a:b:c:d" | tc fields 2 --sep ":"

### 进制转换

```bash
# 任意进制互转
tc base 16 FF              # 16→10: 255
tc base 10 255             # 10→16: 0xFF
tc base 2 1010             # 2→10: 10
tc base 8 77               # 8→10: 63

# 批量转换
tc base 16 "FF" "1A" "2B"

# 显示所有进制
tc base 16 ff --all        # 2/8/10/16/36
```

### 端口检查

```bash
# 检查端口是否被占用（Windows 下显示进程名）
tc port 8080
# ✗ 端口 8080 被占用 (PID: 14776, 进程: node.exe)

tc port 3000
# ✓ 端口 3000 可用
```

### 密码生成

```bash
# 16 位随机密码（默认：大小写+数字+符号）
tc gen password

# 自定义长度和字符集
tc gen password 20                    # 20 位
tc gen password --no-sym              # 不含符号
tc gen password --no-upper --no-sym   # 仅小写+数字
tc gen password 12 --digit-only       # 仅数字 12 位
```

### JWT 解码

```bash
# 解码 JWT token（不验证签名，查看 payload）
tc jwt decode "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U"

# 支持管道
cat token.txt | tc jwt decode
```

### JSON / YAML 互转

```bash
# JSON → YAML
cat data.json | tc convert j2y

# YAML → JSON
cat config.yaml | tc convert y2j

# 保存到文件
cat data.json | tc convert j2y > config.yaml
```

## 命令速查表

```
tc json fmt          格式化 JSON
tc json min          压缩 JSON
tc json get <path>   查询字段
tc json filter <expr> 过滤数组
tc json keys         提取所有 key
tc json values       提取所有 value
tc json paths        列出所有路径
tc json table <f>    转为表格
tc json diff <f1> <f2> 比较差异

tc rx match <p>      高亮匹配
tc rx extract <p>    提取捕获组
tc rx replace <p> <r> 替换
tc rx grep <p>       过滤行
tc rx count <p>      统计次数
tc rx findall <p>    列出所有匹配

tc hash <algo>       计算哈希
tc enc <type>        编码
tc dec <type>        解码
tc stats             统计信息
tc count             行数
tc sort              排序
tc uniq              去重
tc head [n]          前 N 行
tc tail [n]          后 N 行
tc sample [n]        随机采样
tc fields <n>        提取字段
tc b64encode [data]   Base64 编码
tc b64decode [data]   Base64 解码

tc base <from> <val>  进制转换
tc port <n>           端口检查
tc gen password       密码生成
tc jwt decode <token> JWT 解码
tc convert <j2y|y2j>  JSON/YAML 互转

tc version           版本信息
```

## 设计原则

- **管道友好** — stdin → 处理 → stdout
- **零配置** — 开箱即用
- **快速** — Go 编译，启动 < 50ms
- **小体积** — 单二进制，< 10MB
- **错误友好** — 告诉你怎么修，不是堆栈

## 技术栈

- Go 1.22+
- [cobra](https://github.com/spf13/cobra) — CLI 框架
- [gjson](https://github.com/tidwall/gjson) — JSON 查询
- [sjson](https://github.com/tidwall/sjson) — JSON 修改
- [yaml.v3](https://gopkg.in/yaml.v3) — YAML 解析
- [fatih/color](https://github.com/fatih/color) — 终端颜色

## License

MIT
