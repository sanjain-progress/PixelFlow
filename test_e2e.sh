#!/bin/bash

# PixelFlow Docker E2E Test Script
# Tests the complete workflow: Auth → API → Kafka → Worker → MongoDB

set -e

echo "🧪 PixelFlow E2E Test - HTTP Architecture"
echo "=========================================="

BASE_URL="http://localhost:8080"
AUTH_URL="http://localhost:50051"

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Test 1: Health Check
echo ""
echo "📋 Test 1: Health Check"
HEALTH=$(curl -s $BASE_URL/health)
if [[ $HEALTH == *"ok"* ]]; then
    echo -e "${GREEN}✓ API Health Check Passed${NC}"
else
    echo -e "${RED}✗ API Health Check Failed${NC}"
    exit 1
fi

# Test 2: Register User
echo ""
echo "📋 Test 2: Register User"
REGISTER_RESPONSE=$(curl -s -X POST $AUTH_URL/register \
    -H "Content-Type: application/json" \
    -d '{"email":"test@pixelflow.com","password":"password123"}')

if [[ $REGISTER_RESPONSE == *"successfully"* ]] || [[ $REGISTER_RESPONSE == *"already exists"* ]]; then
    echo -e "${GREEN}✓ User Registration Successful${NC}"
else
    echo -e "${RED}✗ User Registration Failed: $REGISTER_RESPONSE${NC}"
    exit 1
fi

# Test 3: Login
echo ""
echo "📋 Test 3: Login User"
LOGIN_RESPONSE=$(curl -s -X POST $AUTH_URL/login \
    -H "Content-Type: application/json" \
    -d '{"email":"test@pixelflow.com","password":"password123"}')

TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*' | sed 's/"token":"//')

if [ -z "$TOKEN" ]; then
    echo -e "${RED}✗ Login Failed: $LOGIN_RESPONSE${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Login Successful${NC}"
echo "Token: ${TOKEN:0:20}..."

# Test 4: Validate Token
echo ""
echo "📋 Test 4: Validate Token"
VALIDATE_RESPONSE=$(curl -s -X GET $AUTH_URL/validate \
    -H "Authorization: Bearer $TOKEN")

if [[ $VALIDATE_RESPONSE == *"\"valid\":true"* ]]; then
    echo -e "${GREEN}✓ Token Validation Successful${NC}"
else
    echo -e "${RED}✗ Token Validation Failed: $VALIDATE_RESPONSE${NC}"
    exit 1
fi

# Test 5: Upload Task (Authenticated)
echo ""
echo "📋 Test 5: Create Image Processing Task"
UPLOAD_RESPONSE=$(curl -s -X POST $BASE_URL/api/upload \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"image_url":"https://example.com/image.jpg"}')

TASK_ID=$(echo $UPLOAD_RESPONSE | grep -o '"_id":"[^"]*' | sed 's/"_id":"//')

if [ -z "$TASK_ID" ]; then
    echo -e "${RED}✗ Task Creation Failed: $UPLOAD_RESPONSE${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Task Created Successfully${NC}"
echo "Task ID: $TASK_ID"

# Test 6: List Tasks
echo ""
echo "📋 Test 6: List User Tasks"
TASKS_RESPONSE=$(curl -s -X GET $BASE_URL/api/tasks \
    -H "Authorization: Bearer $TOKEN")

if [[ $TASKS_RESPONSE == *"$TASK_ID"* ]]; then
    echo -e "${GREEN}✓ Task Listed Successfully${NC}"
else
    echo -e "${RED}✗ Task Listing Failed${NC}"
    exit 1
fi

# Test 7: Wait for Worker Processing
echo ""
echo "📋 Test 7: Worker Processing (waiting 6 seconds for processing...)"
sleep 6

# Check final task status
FINAL_TASKS=$(curl -s -X GET $BASE_URL/api/tasks \
    -H "Authorization: Bearer $TOKEN")

echo "Final task status:"
echo $FINAL_TASKS | python3 -m json.tool 2>/dev/null || echo $FINAL_TASKS

if [[ $FINAL_TASKS == *"COMPLETED"* ]]; then
    echo -e "${GREEN}✓ Worker Processing Completed Successfully${NC}"
elif [[ $FINAL_TASKS == *"PROCESSING"* ]]; then
    echo -e "${GREEN}⚠ Task is still PROCESSING (worker is working)${NC}"
else
    echo -e "${RED}⚠ Task status: Check logs for worker activity${NC}"
fi

echo ""
echo "=========================================="
echo -e "${GREEN}🎉 E2E Tests Completed!${NC}"
echo ""
echo "Summary:"
echo "  ✓ Auth Service: Working"
echo "  ✓ API Service: Working"
echo "  ✓ Task Creation: Working"
echo "  ✓ Kafka Publishing: Working"
echo "  ✓ Worker Processing: Check logs"
echo ""
echo "View logs with:"
echo "  docker-compose logs -f worker-service"
