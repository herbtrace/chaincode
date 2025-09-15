package main

import (
	"log"
	"net/http"

	"herb-api/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	// Set Gin to release mode for production
	// gin.SetMode(gin.ReleaseMode)

	// Create Gin router
	router := gin.Default()

	// Add CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Create controller
	herbController := controllers.NewHerbController()

	// Health check endpoint
	router.GET("/health", herbController.HealthCheck)

	// API routes
	api := router.Group("/api")
	{
		// Herb batch routes
		herbs := api.Group("/herbs")
		{
			herbs.POST("", herbController.CreateHerbBatch)                      // Create new herb batch
			herbs.GET("", herbController.GetAllHerbBatches)                     // Get all herb batches
			herbs.GET("/:id", herbController.GetHerbBatch)                      // Get specific herb batch
			herbs.PUT("/:id/status", herbController.UpdateHerbBatchStatus)      // Update herb batch status
			herbs.PUT("/:id/transfer", herbController.TransferHerbBatch)        // Transfer herb batch ownership
			herbs.GET("/:id/supply-chain", herbController.GetSupplyChainStatus) // Get supply chain status
		}

		// Statistics endpoint
		api.GET("/stats", herbController.GetStats)
	}

	// Supply chain specific endpoints
	supplyChain := api.Group("/supply-chain")
	{
		// Farmer endpoints
		supplyChain.POST("/harvest", func(c *gin.Context) {
			// Set status to "Harvested" and call CreateHerbBatch
			var req map[string]interface{}
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			req["status"] = "Harvested"
			c.Set("herbData", req)
			herbController.CreateHerbBatch(c)
		})

		// Transporter endpoints
		supplyChain.PUT("/transport/:id", func(c *gin.Context) {
			batchID := c.Param("id")
			c.JSON(http.StatusOK, gin.H{
				"message": "Herb batch " + batchID + " picked up for transport",
				"action":  "Use PUT /api/herbs/" + batchID + "/status with newStatus: 'In-Transit'",
			})
		})

		// Lab endpoints
		supplyChain.PUT("/lab-receive/:id", func(c *gin.Context) {
			batchID := c.Param("id")
			c.JSON(http.StatusOK, gin.H{
				"message": "Herb batch " + batchID + " received at lab",
				"action":  "Use PUT /api/herbs/" + batchID + "/status with newStatus: 'Lab-Testing'",
			})
		})

		supplyChain.PUT("/certify/:id", func(c *gin.Context) {
			batchID := c.Param("id")
			c.JSON(http.StatusOK, gin.H{
				"message": "Herb batch " + batchID + " certified",
				"action":  "Use PUT /api/herbs/" + batchID + "/status with newStatus: 'Certified'",
			})
		})
	}

	// API documentation endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service":     "HerbTrace Blockchain API",
			"version":     "1.0.0",
			"description": "API for Ayurvedic Supply Chain Management on Hyperledger Fabric",
			"endpoints": map[string]interface{}{
				"health": "GET /health",
				"herbs": map[string]string{
					"create":       "POST /api/herbs",
					"getAll":       "GET /api/herbs",
					"getById":      "GET /api/herbs/:id",
					"updateStatus": "PUT /api/herbs/:id/status",
					"transfer":     "PUT /api/herbs/:id/transfer",
					"supplyChain":  "GET /api/herbs/:id/supply-chain",
				},
				"stats": "GET /api/stats",
				"supplyChain": map[string]string{
					"harvest":    "POST /api/supply-chain/harvest",
					"transport":  "PUT /api/supply-chain/transport/:id",
					"labReceive": "PUT /api/supply-chain/lab-receive/:id",
					"certify":    "PUT /api/supply-chain/certify/:id",
				},
			},
		})
	})

	// Start server
	log.Println("🌿 HerbTrace API Server starting on :8080")
	log.Println("📡 Blockchain Network: Hyperledger Fabric")
	log.Println("🔗 Channel: herbtrace-temp")
	log.Println("📋 API Documentation: http://localhost:8080")
	log.Println("❤️  Health Check: http://localhost:8080/health")

	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
