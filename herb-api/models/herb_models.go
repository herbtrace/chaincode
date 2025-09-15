package models

import "time"

// HerbBatch represents the herb batch data structure
type HerbBatch struct {
	ID            string `json:"id" binding:"required"`
	BotanicalName string `json:"botanicalName" binding:"required"`
	Farm          string `json:"farm" binding:"required"`
	HarvestDate   string `json:"harvestDate" binding:"required"`
	Owner         string `json:"owner" binding:"required"`
	Status        string `json:"status" binding:"required"`
}

// CreateHerbBatchRequest represents the request payload for creating a herb batch
type CreateHerbBatchRequest struct {
	ID            string `json:"id" binding:"required"`
	BotanicalName string `json:"botanicalName" binding:"required"`
	Farm          string `json:"farm" binding:"required"`
	HarvestDate   string `json:"harvestDate" binding:"required"`
	Owner         string `json:"owner" binding:"required"`
	Status        string `json:"status" binding:"required"`
}

// UpdateStatusRequest represents the request payload for updating herb batch status
type UpdateStatusRequest struct {
	NewStatus string `json:"newStatus" binding:"required"`
}

// TransferRequest represents the request payload for transferring herb batch ownership
type TransferRequest struct {
	NewOwner string `json:"newOwner" binding:"required"`
}

// APIResponse represents a standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// SupplyChainEvent represents a supply chain event
type SupplyChainEvent struct {
	BatchID     string    `json:"batchId"`
	Event       string    `json:"event"`
	Actor       string    `json:"actor"`
	Location    string    `json:"location,omitempty"`
	Timestamp   time.Time `json:"timestamp"`
	Description string    `json:"description,omitempty"`
}

// Supply chain status constants
const (
	StatusHarvested   = "Harvested"
	StatusInTransit   = "In-Transit"
	StatusLabTesting  = "Lab-Testing"
	StatusCertified   = "Certified"
	StatusProcessing  = "Processing"
	StatusPackaged    = "Packaged"
	StatusDistributed = "Distributed"
	StatusDelivered   = "Delivered"
)
