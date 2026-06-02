#!/bin/bash

# 测试 contact 模块配置修复

echo "=== 测试 Contact 模块配置修复 ==="
echo ""

# 1. 检查数据库中 contact 模块的当前状态
echo "1. 检查数据库中 contact 模块的当前状态:"
sqlite3 medical.db "SELECT module_name, enabled, extra_data FROM module_configs WHERE module_name = 'contact';"
echo ""

# 2. 模拟保存 contact 配置（使用 curl）
echo "2. 模拟保存 contact 配置:"
echo "发送测试数据到后端 API..."

# 首先登录获取 token
TOKEN=$(curl -s -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' | \
  grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
  echo "登录失败，无法获取 token"
  exit 1
fi

echo "获取到 token: ${TOKEN:0:20}..."

# 保存 contact 配置
curl -s -X POST http://localhost:8080/api/admin/modules \
  -H "Authorization: Bearer $TOKEN" \
  -F "moduleName=contact" \
  -F "enabled=true" \
  -F "email=test@example.com" \
  -F "phone=+86 138 0000 0000" \
  -F "whatsapp=+86 138 0000 0000" \
  -F "address=Industrial Park, Guangdong, China" \
  | python3 -m json.tool 2>/dev/null || echo "响应格式错误"

echo ""
echo "3. 检查保存后的数据库状态:"
sqlite3 medical.db "SELECT module_name, enabled, extra_data FROM module_configs WHERE module_name = 'contact';"
echo ""

# 4. 测试公开 API 是否能获取到 contact 配置
echo "4. 测试公开 API 获取 contact 配置:"
curl -s http://localhost:8080/api/public/modules | \
  python3 -c "import sys, json; data = json.load(sys.stdin); contact = next((m for m in data.get('data', data) if m.get('moduleName') == 'contact'), None); print(json.dumps(contact, indent=2) if contact else '未找到 contact 模块')"

echo ""
echo "=== 测试完成 ==="