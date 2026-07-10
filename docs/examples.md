# tc 使用示例

## 场景 1：分析 API 响应

```bash
# 调用 API 并格式化
curl -s https://api.example.com/users | tc json fmt

# 提取特定字段
curl -s https://api.example.com/users | tc json get "data.*.email"

# 过滤活跃用户
curl -s https://api.example.com/users | tc json filter 'status == "active"'
```

## 场景 2：分析 Nginx 日志

```bash
# 统计 5xx 错误数
cat access.log | tc rx grep " 5\d{2} " | tc count

# Top 10 访问 IP
cat access.log | tc rx extract "^(\d+\.\d+\.\d+\.\d+)" | tc sort | tc uniq -c | tc sort -n -r | tc head 10

# 统计各 HTTP 方法
cat access.log | tc rx extract '"(GET|POST|PUT|DELETE)' | tc sort | tc uniq -c | tc sort -n -r

# 找出慢请求（响应时间 > 5s）
cat access.log | tc rx extract '(\d+\.\d{3})$' | tc rx grep '^[5-9]\.'
```

## 场景 3：处理配置文件

```bash
# 查看 JSON 配置
cat config.json | tc json fmt

# 修改前备份，然后 diff
cp config.json config.json.bak
# ... 修改 config.json ...
tc json diff config.json.bak config.json

# 提取所有环境变量名
cat .env | tc rx extract "^(\w+)=" | tc sort
```

## 场景 4：数据清洗

```bash
# 从 CSV 提取第 2 列
cat data.csv | tc fields 2 -s ","

# 去除空行
cat messy.txt | tc rx grep "^.+$"

# 提取所有邮箱
cat contacts.txt | tc rx findall "[\w.]+@[\w.]+\.\w+"

# 提取所有 URL
cat page.html | tc rx findall "https?://[^\s\"'<>]+"
```

## 场景 5：快速哈希/编码

```bash
# 计算文件哈希
tc hash sha256 important.tar.gz

# 编码 JWT payload
echo '{"sub":"12345","name":"John"}' | tc enc base64

# 解码 URL 编码
echo "hello%20world%21" | tc dec url
```

## 场景 6：开发调试

```bash
# 格式化 API 响应
curl -s localhost:3000/api/users | tc json fmt

# 查看数据库导出
cat dump.json | tc json get "rows.*.id" | tc count

# 对比两个环境的配置
tc json diff dev.json prod.json

# 快速生成测试数据摘要
cat test-data.json | tc json table "id,name,email" --csv > test-report.csv
```
