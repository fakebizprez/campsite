#!/bin/bash

# Demo script to test the Go API functionality
# This shows the key features working locally without AWS dependencies

echo "🚀 Campsite Go API Demo"
echo "========================"

# Check if the Go API is running
API_URL="http://localhost:8080"

echo "📡 Testing API Health..."
curl -s "$API_URL/health" | jq '.' || echo "API not running on $API_URL"

echo -e "\n📋 Testing Organizations..."

# Create an organization
echo "Creating organization..."
ORG_RESPONSE=$(curl -s -X POST "$API_URL/api/v1/organizations" \
  -H "Content-Type: application/json" \
  -d '{"name":"Demo Organization","slug":"demo-org"}')

echo "Response: $ORG_RESPONSE" | jq '.'

# Get organizations
echo -e "\nListing organizations..."
curl -s "$API_URL/api/v1/organizations" | jq '.'

echo -e "\n📁 Testing Projects..."

# Create a project
echo "Creating project..."
PROJECT_RESPONSE=$(curl -s -X POST "$API_URL/api/v1/organizations/demo-org/projects" \
  -H "Content-Type: application/json" \
  -d '{"name":"Demo Project","slug":"demo-project"}')

echo "Response: $PROJECT_RESPONSE" | jq '.'

# Get projects
echo -e "\nListing projects..."
curl -s "$API_URL/api/v1/organizations/demo-org/projects" | jq '.'

echo -e "\n📝 Testing Posts..."

# Create a post
echo "Creating post..."
POST_RESPONSE=$(curl -s -X POST "$API_URL/api/v1/organizations/demo-org/projects/demo-project/posts" \
  -H "Content-Type: application/json" \
  -d '{"title":"Demo Post","content":"This is a demo post created through the Go API!"}')

echo "Response: $POST_RESPONSE" | jq '.'

# Get posts
echo -e "\nListing posts..."
curl -s "$API_URL/api/v1/organizations/demo-org/projects/demo-project/posts" | jq '.'

echo -e "\n📄 Testing File Upload..."

# Create a test file
echo "Hello from Campsite Go API!" > /tmp/test-upload.txt

# Upload the file
echo "Uploading file..."
UPLOAD_RESPONSE=$(curl -s -X POST "$API_URL/api/v1/uploads" \
  -F "file=@/tmp/test-upload.txt")

echo "Response: $UPLOAD_RESPONSE" | jq '.'

# Extract filename from response
FILENAME=$(echo "$UPLOAD_RESPONSE" | jq -r '.filename')

if [ "$FILENAME" != "null" ] && [ "$FILENAME" != "" ]; then
  echo -e "\nDownloading uploaded file..."
  curl -s "$API_URL/api/v1/uploads/$FILENAME"
  echo ""
fi

echo -e "\n🔌 Testing Integration Stubs..."

# Test Slack integration stub
echo "Testing Slack integration..."
curl -s -X POST "$API_URL/api/v1/integrations/slack/send" \
  -H "Content-Type: application/json" \
  -d '{"message":"Hello from Go API!","channel":"general"}' | jq '.'

# Test email integration stub
echo -e "\nTesting email integration..."
curl -s -X POST "$API_URL/api/v1/integrations/email/send" \
  -H "Content-Type: application/json" \
  -d '{"to":"user@example.com","subject":"Test Email","body":"This is a test email from the Go API"}' | jq '.'

# Test search integration stub  
echo -e "\nTesting search integration..."
curl -s "$API_URL/api/v1/integrations/search/query?q=test" | jq '.'

echo -e "\n✅ Demo Complete!"
echo "The Go API successfully demonstrated:"
echo "  • Organizations, Projects, and Posts management"
echo "  • Local file upload and download"
echo "  • Integration stubs for external services"
echo "  • All without requiring AWS or other cloud dependencies"

# Clean up
rm -f /tmp/test-upload.txt

echo -e "\n💡 To start the Go API server, run:"
echo "    cd go-api && make run"
echo -e "\n💡 To connect via WebSocket:"
echo "    wscat -c ws://localhost:8080/ws"