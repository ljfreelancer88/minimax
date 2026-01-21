#!/bin/bash

echo "🔍 Checking setup..."

# Check PHP server
echo -n "PHP server (3000): "
curl -s -o /dev/null -w "%{http_code}" http://localhost:3000/ && echo "✅ Running" || echo "❌ Not running"

# Check Annotation proxy
echo -n "Annotation proxy (8080): "
curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/ && echo "✅ Running" || echo "❌ Not running"

# Check assets
echo -n "JS asset: "
curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/annotation-assets/annotation.js && echo "✅ Available" || echo "❌ Not found"

echo -n "CSS asset: "
curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/annotation-assets/annotation.css && echo "✅ Available" || echo "❌ Not found"

# Check API
echo -n "API endpoint: "
curl -s "http://localhost:8080/api/annotations?url=/" | grep -q "annotations" && echo "✅ Working" || echo "❌ Not working"

# Check injection
echo -n "Script injection: "
curl -s http://localhost:8080/ | grep -q "annotation.js" && echo "✅ Injected" || echo "❌ Not injected"
