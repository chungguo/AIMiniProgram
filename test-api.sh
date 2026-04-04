#!/bin/bash

# API 测试脚本

BASE_URL="http://localhost:8000/api"

echo "=== AIMiniProgram API Test ==="
echo ""

# 测试健康检查
echo "1. Testing health endpoint..."
curl -s ${BASE_URL}/health | jq .
echo ""

# 测试模型列表
echo "2. Testing models list..."
curl -s "${BASE_URL}/models?page=1&limit=5" | jq '.data[:2]'
echo ""

# 测试模型家族
echo "3. Testing model families..."
curl -s ${BASE_URL}/models/families | jq .
echo ""

# 测试单个模型
echo "4. Testing single model (gpt-4o)..."
curl -s ${BASE_URL}/models/detail/gpt-4o | jq '.data | {id, name, family}'
echo ""

# 测试论文列表
echo "5. Testing papers list..."
curl -s "${BASE_URL}/papers?page=1&limit=3" | jq '.data[:1]'
echo ""

# 测试评测数据
echo "6. Testing analysis list..."
curl -s "${BASE_URL}/analysis/artificialanalysis?page=1&limit=3" | jq '.data[:1]'
echo ""

echo "=== Test Complete ==="
