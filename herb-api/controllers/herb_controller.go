package controllers

import (
	"net/http"

	"herb-api/models"
	"herb-api/services"

	"github.com/gin-gonic/gin"
)

type HerbController struct {
	fabricService *services.FabricService
}

// NewHerbController creates a new instance of HerbController
func NewHerbController() *HerbController {
	return &HerbController{
		fabricService: services.NewFabricService(),
	}
}

// CreateHerbBatch handles POST /api/herbs
func (hc *HerbController) CreateHerbBatch(c *gin.Context) {
	var req models.CreateHerbBatchRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request payload",
			Error:   err.Error(),
		})
		return
	}

	// Check if herb batch already exists
	exists, err := hc.fabricService.HerbBatchExists(req.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to check if herb batch exists",
			Error:   err.Error(),
		})
		return
	}

	if exists {
		c.JSON(http.StatusConflict, models.APIResponse{
			Success: false,
			Message: "Herb batch with this ID already exists",
		})
		return
	}

	// Create herb batch on blockchain
	if err := hc.fabricService.CreateHerbBatch(req); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to create herb batch on blockchain",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Herb batch created successfully",
		Data: map[string]string{
			"batchId": req.ID,
			"status":  "Created on blockchain",
		},
	})
}

// GetHerbBatch handles GET /api/herbs/:id
func (hc *HerbController) GetHerbBatch(c *gin.Context) {
	batchID := c.Param("id")

	if batchID == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Batch ID is required",
		})
		return
	}

	herbBatch, err := hc.fabricService.ReadHerbBatch(batchID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "Herb batch not found",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Herb batch retrieved successfully",
		Data:    herbBatch,
	})
}

// GetAllHerbBatches handles GET /api/herbs
func (hc *HerbController) GetAllHerbBatches(c *gin.Context) {
	herbBatches, err := hc.fabricService.GetAllHerbBatches()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to retrieve herb batches",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Herb batches retrieved successfully",
		Data: map[string]interface{}{
			"batches": herbBatches,
			"count":   len(herbBatches),
		},
	})
}

// UpdateHerbBatchStatus handles PUT /api/herbs/:id/status
func (hc *HerbController) UpdateHerbBatchStatus(c *gin.Context) {
	batchID := c.Param("id")
	var req models.UpdateStatusRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request payload",
			Error:   err.Error(),
		})
		return
	}

	// Validate status
	validStatuses := []string{
		models.StatusHarvested,
		models.StatusInTransit,
		models.StatusLabTesting,
		models.StatusCertified,
		models.StatusProcessing,
		models.StatusPackaged,
		models.StatusDistributed,
		models.StatusDelivered,
	}

	isValidStatus := false
	for _, status := range validStatuses {
		if req.NewStatus == status {
			isValidStatus = true
			break
		}
	}

	if !isValidStatus {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid status. Valid statuses are: Harvested, In-Transit, Lab-Testing, Certified, Processing, Packaged, Distributed, Delivered",
		})
		return
	}

	// Update status on blockchain
	if err := hc.fabricService.UpdateHerbBatchStatus(batchID, req.NewStatus); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to update herb batch status",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Herb batch status updated successfully",
		Data: map[string]string{
			"batchId":   batchID,
			"newStatus": req.NewStatus,
		},
	})
}

// TransferHerbBatch handles PUT /api/herbs/:id/transfer
func (hc *HerbController) TransferHerbBatch(c *gin.Context) {
	batchID := c.Param("id")
	var req models.TransferRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request payload",
			Error:   err.Error(),
		})
		return
	}

	oldOwner, err := hc.fabricService.TransferHerbBatch(batchID, req.NewOwner)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to transfer herb batch",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Herb batch transferred successfully",
		Data: map[string]string{
			"batchId":  batchID,
			"oldOwner": oldOwner,
			"newOwner": req.NewOwner,
		},
	})
}

// GetSupplyChainStatus handles GET /api/herbs/:id/supply-chain
func (hc *HerbController) GetSupplyChainStatus(c *gin.Context) {
	batchID := c.Param("id")

	herbBatch, err := hc.fabricService.ReadHerbBatch(batchID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "Herb batch not found",
			Error:   err.Error(),
		})
		return
	}

	// Create a supply chain timeline based on current status
	timeline := createSupplyChainTimeline(herbBatch)

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Supply chain status retrieved successfully",
		Data: map[string]interface{}{
			"batchInfo": herbBatch,
			"timeline":  timeline,
		},
	})
}

// HealthCheck handles GET /health
func (hc *HerbController) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "HerbTrace API is running",
		Data: map[string]string{
			"status":  "healthy",
			"service": "HerbTrace Blockchain API",
			"version": "1.0.0",
		},
	})
}

// GetStats handles GET /api/stats
func (hc *HerbController) GetStats(c *gin.Context) {
	herbBatches, err := hc.fabricService.GetAllHerbBatches()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to retrieve statistics",
			Error:   err.Error(),
		})
		return
	}

	stats := calculateStats(herbBatches)

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Statistics retrieved successfully",
		Data:    stats,
	})
}

// Helper function to create supply chain timeline
func createSupplyChainTimeline(herb *models.HerbBatch) []map[string]interface{} {
	timeline := []map[string]interface{}{
		{
			"stage":       "Farming",
			"status":      "Completed",
			"actor":       herb.Owner,
			"location":    herb.Farm,
			"date":        herb.HarvestDate,
			"description": "Herbs harvested from certified organic farm",
		},
	}

	// Add stages based on current status
	switch herb.Status {
	case models.StatusInTransit:
		timeline = append(timeline, map[string]interface{}{
			"stage":       "Transportation",
			"status":      "In Progress",
			"actor":       herb.Owner,
			"description": "Herbs in transit to next stage",
		})
	case models.StatusLabTesting:
		timeline = append(timeline, map[string]interface{}{
			"stage":       "Quality Testing",
			"status":      "In Progress",
			"actor":       herb.Owner,
			"description": "Quality testing and certification in progress",
		})
	case models.StatusCertified:
		timeline = append(timeline, map[string]interface{}{
			"stage":       "Certification",
			"status":      "Completed",
			"actor":       herb.Owner,
			"description": "Quality certification completed",
		})
	}

	return timeline
}

// Helper function to calculate statistics
func calculateStats(herbs []models.HerbBatch) map[string]interface{} {
	statusCount := make(map[string]int)
	farmCount := make(map[string]int)

	for _, herb := range herbs {
		statusCount[herb.Status]++
		farmCount[herb.Farm]++
	}

	return map[string]interface{}{
		"totalBatches":    len(herbs),
		"statusBreakdown": statusCount,
		"farmBreakdown":   farmCount,
		"activeSupplyChain": statusCount[models.StatusInTransit] +
			statusCount[models.StatusLabTesting] +
			statusCount[models.StatusProcessing],
	}
}
