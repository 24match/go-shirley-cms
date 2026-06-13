#!/bin/bash

# React 前端构建脚本
# 将构建产物输出到 ../frontend 目录

echo "🚀 开始构建 React 前端..."

# 检查 node_modules 是否存在
if [ ! -d "node_modules" ]; then
    echo "📦 安装依赖..."
    npm install
fi

# 执行构建
echo "🔨 执行构建..."
npm run build

# 验证构建结果
if [ -d "../frontend/assets" ]; then
    echo "✅ 构建成功！"
    echo "📁 构建产物已输出到 ../frontend 目录"
    echo ""
    echo "下一步："
    echo "1. 在项目根目录启动 Go 后端：go run main.go"
    echo "2. 访问 http://localhost:8080"
else
    echo "❌ 构建失败！请检查错误信息。"
    exit 1
fi