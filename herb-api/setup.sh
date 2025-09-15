#!/bin/bash

echo "🌿 Setting up HerbTrace API..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21+ first."
    exit 1
fi

echo "✅ Go is installed: $(go version)"

# Initialize Go module if not exists
if [ ! -f "go.mod" ]; then
    echo "📦 Initializing Go module..."
    go mod init herb-api
fi

# Install dependencies
echo "📥 Installing dependencies..."
go mod tidy

# Check if Fabric network is running
echo "🔍 Checking Fabric network..."
cd ../test-network
if ! docker ps | grep -q "peer0.org1.example.com"; then
    echo "⚠️  Fabric network is not running. Starting network..."
    ./network.sh up createChannel -c herbtrace-temp
    ./network.sh deployCC -ccn herbbatch -ccp ../herb-asset/chaincode-go -ccl go -ccv 1.0 -ccs 1 -c herbtrace-temp
else
    echo "✅ Fabric network is running"
fi

cd ../herb-api

echo ""
echo "🎉 Setup complete!"
echo ""
echo "To start the API server:"
echo "  go run main.go"
echo ""
echo "API will be available at: http://localhost:8080"
echo "Health check: http://localhost:8080/health"
echo "Documentation: http://localhost:8080/"
echo ""
