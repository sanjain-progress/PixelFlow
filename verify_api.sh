#!/bin/bash

# 1. Login to get Token
echo "--- Logging in ---"
# Note: We use the Auth Service client to get a token first
# In a real app, the frontend would do this
TOKEN=$(go run apps/auth/client/main.go | grep "Token:" | awk '{print $2}')
echo "Token: $TOKEN"

if [ -z "$TOKEN" ]; then
    echo "Failed to get token. Is Auth Service running?"
    exit 1
fi

# 2. Test Health Check
echo -e "\n--- Health Check ---"
curl -s http://localhost:8080/health

# 3. Test Upload (Create Task)
echo -e "\n\n--- Uploading Image ---"
curl -X POST http://localhost:8080/api/upload \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"image_url": "https://example.com/image.jpg"}'

# 4. Test List Tasks
echo -e "\n\n--- Listing Tasks ---"
curl -s -X GET http://localhost:8080/api/tasks \
  -H "Authorization: Bearer $TOKEN"
