# HerbTrace API - Blockchain Backend for Ayurvedic Supply Chain

A Go-based REST API server that interfaces with Hyperledger Fabric blockchain for tracking Ayurvedic herbs through the supply chain.

## 🏗️ Architecture

```
Frontend → REST API (Gin) → Fabric Network → Smart Contracts → Blockchain Ledger
```

## 🚀 Quick Start

### Prerequisites
- Go 1.21+
- Docker & Docker Compose (for Fabric network)
- Hyperledger Fabric network running (from test-network)

### 1. Start Blockchain Network
```bash
cd ../test-network
./network.sh up createChannel -c herbtrace-temp
./network.sh deployCC -ccn herbbatch -ccp ../herb-asset/chaincode-go -ccl go -ccv 1.0 -ccs 1 -c herbtrace-temp
```

### 2. Install Dependencies
```bash
cd herb-api
go mod tidy
```

### 3. Run API Server
```bash
go run main.go
```

The API will start on `http://localhost:8080`

## 📋 API Endpoints

### Health Check
```bash
GET /health
```

### Herb Batch Management
```bash
# Create new herb batch
POST /api/herbs
{
  "id": "batch8",
  "botanicalName": "Neem",
  "farm": "Delhi Organics",
  "harvestDate": "2025-09-15",
  "owner": "Farmer Name",
  "status": "Harvested"
}

# Get all herb batches
GET /api/herbs

# Get specific herb batch
GET /api/herbs/{id}

# Update herb batch status
PUT /api/herbs/{id}/status
{
  "newStatus": "In-Transit"
}

# Transfer ownership
PUT /api/herbs/{id}/transfer
{
  "newOwner": "New Owner Name"
}

# Get supply chain timeline
GET /api/herbs/{id}/supply-chain
```

### Supply Chain Workflows
```bash
# Farmer harvests herbs
POST /api/supply-chain/harvest

# Transporter picks up
PUT /api/supply-chain/transport/{id}

# Lab receives herbs
PUT /api/supply-chain/lab-receive/{id}

# Lab certifies herbs
PUT /api/supply-chain/certify/{id}
```

### Statistics
```bash
GET /api/stats
```

## 🌿 Supply Chain Statuses

- `Harvested` - Herbs harvested from farm
- `In-Transit` - Herbs being transported
- `Lab-Testing` - Quality testing in progress
- `Certified` - Lab certification completed
- `Processing` - Manufacturing/processing
- `Packaged` - Ready for distribution
- `Distributed` - Sent to retailers
- `Delivered` - Delivered to end consumer

## 📝 Example Usage

### Create a Herb Batch
```bash
curl -X POST http://localhost:8080/api/herbs \
  -H "Content-Type: application/json" \
  -d '{
    "id": "batch9",
    "botanicalName": "Withania somnifera",
    "farm": "Kerala Ayurveda Farms",
    "harvestDate": "2025-09-15",
    "owner": "Ravi Kumar",
    "status": "Harvested"
  }'
```

### Get All Herbs
```bash
curl http://localhost:8080/api/herbs
```

### Update Status
```bash
curl -X PUT http://localhost:8080/api/herbs/batch9/status \
  -H "Content-Type: application/json" \
  -d '{"newStatus": "In-Transit"}'
```

### Transfer Ownership
```bash
curl -X PUT http://localhost:8080/api/herbs/batch9/transfer \
  -H "Content-Type: application/json" \
  -d '{"newOwner": "Transport Company"}'
```

## 🔧 Configuration

Update the paths in `services/fabric_service.go` if your network is in a different location:

```go
NetworkPath:   "../test-network", // Adjust this path
ChannelName:   "herbtrace-temp",
ChaincodeName: "herbbatch",
```

## 🐛 Troubleshooting

### API Server Issues
1. Check if Fabric network is running: `docker ps`
2. Verify chaincode is deployed: `./network.sh cc query -ccn herbbatch -c herbtrace-temp -ccqc '{"function":"GetAllHerbBatches","Args":[]}'`
3. Check API logs for detailed error messages

### Blockchain Connection Issues
1. Ensure test-network is in correct relative path
2. Verify channel name and chaincode name match your deployment
3. Check if CLI commands work manually

## 🏆 For Hackathon Demo

### Demo Scenarios
1. **Farmer Harvest**: Create new herb batch
2. **Transport**: Update status to "In-Transit"
3. **Lab Testing**: Update status to "Lab-Testing"
4. **Certification**: Update status to "Certified"
5. **Traceability**: Show complete supply chain history

### Demo Commands
```bash
# Show API documentation
curl http://localhost:8080/

# Health check
curl http://localhost:8080/health

# Statistics
curl http://localhost:8080/api/stats
```

## 🌟 Features

- ✅ REST API interface to Hyperledger Fabric
- ✅ Complete CRUD operations for herb batches
- ✅ Supply chain status tracking
- ✅ Ownership transfer capabilities
- ✅ Statistics and analytics
- ✅ CORS enabled for frontend integration
- ✅ Comprehensive error handling
- ✅ Production-ready architecture

## 📚 Next Steps

1. Add authentication middleware
2. Implement database caching layer
3. Add WebSocket support for real-time updates
4. Create frontend integration
5. Add comprehensive logging
6. Implement rate limiting

Made with ❤️ for sustainable Ayurvedic supply chains! 🌿
