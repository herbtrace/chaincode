package services

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"herb-api/models"
)

type FabricService struct {
	NetworkPath   string
	ChannelName   string
	ChaincodeName string
}

// NewFabricService creates a new instance of FabricService
func NewFabricService() *FabricService {
	return &FabricService{
		NetworkPath:   "../test-network", // Adjust path relative to your API server
		ChannelName:   "herbtrace-temp",
		ChaincodeName: "herbbatch",
	}
}

// CreateHerbBatch creates a new herb batch on the blockchain
func (fs *FabricService) CreateHerbBatch(herb models.CreateHerbBatchRequest) error {
	// Clean the arguments by replacing spaces with underscores to avoid shell parsing issues
	cleanID := strings.ReplaceAll(herb.ID, " ", "_")
	cleanBotanicalName := strings.ReplaceAll(herb.BotanicalName, " ", "_")
	cleanFarm := strings.ReplaceAll(herb.Farm, " ", "_")
	cleanHarvestDate := strings.ReplaceAll(herb.HarvestDate, " ", "_")
	cleanOwner := strings.ReplaceAll(herb.Owner, " ", "_")
	cleanStatus := strings.ReplaceAll(herb.Status, " ", "_")

	// Construct the chaincode invoke command
	args := fmt.Sprintf(`{"function":"CreateHerbBatch","Args":["%s","%s","%s","%s","%s","%s"]}`,
		cleanID, cleanBotanicalName, cleanFarm, cleanHarvestDate, cleanOwner, cleanStatus)

	cmd := exec.Command("./network.sh", "cc", "invoke",
		"-ccn", fs.ChaincodeName,
		"-c", fs.ChannelName,
		"-d", "3", // CLI delay
		"-ccic", args)

	cmd.Dir = fs.NetworkPath
	output, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("failed to create herb batch: %v, output: %s", err, string(output))
	}

	// Check if the output contains success indicators
	if !strings.Contains(string(output), "Invoke successful") {
		return fmt.Errorf("chaincode invocation failed: %s", string(output))
	}

	return nil
} // ReadHerbBatch retrieves a herb batch from the blockchain
func (fs *FabricService) ReadHerbBatch(batchID string) (*models.HerbBatch, error) {
	args := fmt.Sprintf(`{"function":"ReadHerbBatch","Args":["%s"]}`, batchID)

	cmd := exec.Command("./network.sh", "cc", "query",
		"-ccn", fs.ChaincodeName,
		"-c", fs.ChannelName,
		"-ccqc", args)

	cmd.Dir = fs.NetworkPath
	output, err := cmd.CombinedOutput()

	if err != nil {
		return nil, fmt.Errorf("failed to read herb batch: %v, output: %s", err, string(output))
	}

	// Extract JSON from the output
	outputStr := string(output)
	lines := strings.Split(outputStr, "\n")
	var jsonLine string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "{") && strings.HasSuffix(line, "}") {
			jsonLine = line
			break
		}
	}

	if jsonLine == "" {
		return nil, fmt.Errorf("no valid JSON found in output: %s", outputStr)
	}

	var herbBatch models.HerbBatch
	if err := json.Unmarshal([]byte(jsonLine), &herbBatch); err != nil {
		return nil, fmt.Errorf("failed to parse herb batch JSON: %v", err)
	}

	return &herbBatch, nil
}

// GetAllHerbBatches retrieves all herb batches from the blockchain
func (fs *FabricService) GetAllHerbBatches() ([]models.HerbBatch, error) {
	args := `{"function":"GetAllHerbBatches","Args":[]}`

	cmd := exec.Command("./network.sh", "cc", "query",
		"-ccn", fs.ChaincodeName,
		"-c", fs.ChannelName,
		"-ccqc", args)

	cmd.Dir = fs.NetworkPath
	output, err := cmd.CombinedOutput()

	if err != nil {
		return nil, fmt.Errorf("failed to get all herb batches: %v, output: %s", err, string(output))
	}

	// Extract JSON from the output
	outputStr := string(output)
	lines := strings.Split(outputStr, "\n")
	var jsonLine string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			jsonLine = line
			break
		}
	}

	if jsonLine == "" {
		return nil, fmt.Errorf("no valid JSON array found in output: %s", outputStr)
	}

	var herbBatches []models.HerbBatch
	if err := json.Unmarshal([]byte(jsonLine), &herbBatches); err != nil {
		return nil, fmt.Errorf("failed to parse herb batches JSON: %v", err)
	}

	return herbBatches, nil
}

// UpdateHerbBatchStatus updates the status of a herb batch
func (fs *FabricService) UpdateHerbBatchStatus(batchID, newStatus string) error {
	args := fmt.Sprintf(`{"function":"UpdateHerbBatchStatus","Args":["%s","%s"]}`,
		batchID, newStatus)

	cmd := exec.Command("./network.sh", "cc", "invoke",
		"-ccn", fs.ChaincodeName,
		"-c", fs.ChannelName,
		"-d", "3",
		"-ccic", args)

	cmd.Dir = fs.NetworkPath
	output, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("failed to update herb batch status: %v, output: %s", err, string(output))
	}

	if !strings.Contains(string(output), "Invoke successful") {
		return fmt.Errorf("chaincode invocation failed: %s", string(output))
	}

	return nil
}

// TransferHerbBatch transfers ownership of a herb batch
func (fs *FabricService) TransferHerbBatch(batchID, newOwner string) (string, error) {
	args := fmt.Sprintf(`{"function":"TransferHerbBatch","Args":["%s","%s"]}`,
		batchID, newOwner)

	cmd := exec.Command("./network.sh", "cc", "invoke",
		"-ccn", fs.ChaincodeName,
		"-c", fs.ChannelName,
		"-d", "3",
		"-ccic", args)

	cmd.Dir = fs.NetworkPath
	output, err := cmd.CombinedOutput()

	if err != nil {
		return "", fmt.Errorf("failed to transfer herb batch: %v, output: %s", err, string(output))
	}

	if !strings.Contains(string(output), "Invoke successful") {
		return "", fmt.Errorf("chaincode invocation failed: %s", string(output))
	}

	// For now, return empty string as old owner (would need to parse from chaincode response in production)
	return "Previous Owner", nil
}

// HerbBatchExists checks if a herb batch exists on the blockchain
func (fs *FabricService) HerbBatchExists(batchID string) (bool, error) {
	args := fmt.Sprintf(`{"function":"HerbBatchExists","Args":["%s"]}`, batchID)

	cmd := exec.Command("./network.sh", "cc", "query",
		"-ccn", fs.ChaincodeName,
		"-c", fs.ChannelName,
		"-ccqc", args)

	cmd.Dir = fs.NetworkPath
	output, err := cmd.CombinedOutput()

	if err != nil {
		return false, fmt.Errorf("failed to check herb batch existence: %v, output: %s", err, string(output))
	}

	// Parse boolean result
	outputStr := strings.TrimSpace(string(output))
	lines := strings.Split(outputStr, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "true" {
			return true, nil
		} else if line == "false" {
			return false, nil
		}
	}

	return false, fmt.Errorf("unexpected output format: %s", outputStr)
}
