package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

// SmartContract provides functions for managing HerbBatch assets
type SmartContract struct {
	contractapi.Contract
}

// HerbBatch describes the core asset in our supply chain
// Insert struct field in alphabetic order => to achieve determinism across languages
// golang keeps the order when marshal to json but doesn't order automatically
type HerbBatch struct {
	ID            string `json:"ID"`
	BotanicalName string `json:"botanicalName"`
	Farm          string `json:"farm"`
	HarvestDate   string `json:"harvestDate"`
	Owner         string `json:"owner"`
	Status        string `json:"status"` // e.g., "Harvested", "In-Transit", "Certified"
}

// InitLedger adds a base set of herb batches to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	herbBatches := []HerbBatch{
		{ID: "batch1", BotanicalName: "Withania somnifera", Farm: "Kerala Ayurveda Farms", HarvestDate: "2024-08-15", Owner: "Ravi Sharma", Status: "Harvested"},
		{ID: "batch2", BotanicalName: "Curcuma longa", Farm: "Tamil Nadu Spice Co", HarvestDate: "2024-08-20", Owner: "Priya Patel", Status: "In-Transit"},
		{ID: "batch3", BotanicalName: "Ocimum tenuiflorum", Farm: "Maharashtra Herbs", HarvestDate: "2024-07-30", Owner: "Suresh Kumar", Status: "Certified"},
		{ID: "batch4", BotanicalName: "Bacopa monnieri", Farm: "Uttarakhand Organics", HarvestDate: "2024-08-10", Owner: "Anjali Singh", Status: "Harvested"},
		{ID: "batch5", BotanicalName: "Centella asiatica", Farm: "Karnataka Medicinals", HarvestDate: "2024-08-25", Owner: "Vikram Joshi", Status: "In-Transit"},
		{ID: "batch6", BotanicalName: "Tinospora cordifolia", Farm: "Rajasthan Herb Gardens", HarvestDate: "2024-08-12", Owner: "Meera Gupta", Status: "Certified"},
	}

	for _, herbBatch := range herbBatches {
		herbBatchJSON, err := json.Marshal(herbBatch)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(herbBatch.ID, herbBatchJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

// CreateHerbBatch issues a new herb batch to the world state with given details.
func (s *SmartContract) CreateHerbBatch(ctx contractapi.TransactionContextInterface, id string, botanicalName string, farm string, harvestDate string, owner string, status string) error {
	exists, err := s.HerbBatchExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the herb batch %s already exists", id)
	}

	herbBatch := HerbBatch{
		ID:            id,
		BotanicalName: botanicalName,
		Farm:          farm,
		HarvestDate:   harvestDate,
		Owner:         owner,
		Status:        status,
	}
	herbBatchJSON, err := json.Marshal(herbBatch)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, herbBatchJSON)
}

// ReadHerbBatch returns the herb batch stored in the world state with given id.
func (s *SmartContract) ReadHerbBatch(ctx contractapi.TransactionContextInterface, id string) (*HerbBatch, error) {
	herbBatchJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if herbBatchJSON == nil {
		return nil, fmt.Errorf("the herb batch %s does not exist", id)
	}

	var herbBatch HerbBatch
	err = json.Unmarshal(herbBatchJSON, &herbBatch)
	if err != nil {
		return nil, err
	}

	return &herbBatch, nil
}

// UpdateHerbBatch updates an existing herb batch in the world state with provided parameters.
func (s *SmartContract) UpdateHerbBatch(ctx contractapi.TransactionContextInterface, id string, botanicalName string, farm string, harvestDate string, owner string, status string) error {
	exists, err := s.HerbBatchExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the herb batch %s does not exist", id)
	}

	// overwriting original herb batch with new herb batch
	herbBatch := HerbBatch{
		ID:            id,
		BotanicalName: botanicalName,
		Farm:          farm,
		HarvestDate:   harvestDate,
		Owner:         owner,
		Status:        status,
	}
	herbBatchJSON, err := json.Marshal(herbBatch)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, herbBatchJSON)
}

// DeleteHerbBatch deletes a given herb batch from the world state.
func (s *SmartContract) DeleteHerbBatch(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := s.HerbBatchExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the herb batch %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
}

// HerbBatchExists returns true when herb batch with given ID exists in world state
func (s *SmartContract) HerbBatchExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	herbBatchJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return herbBatchJSON != nil, nil
}

// TransferHerbBatch updates the owner field of herb batch with given id in world state, and returns the old owner.
func (s *SmartContract) TransferHerbBatch(ctx contractapi.TransactionContextInterface, id string, newOwner string) (string, error) {
	herbBatch, err := s.ReadHerbBatch(ctx, id)
	if err != nil {
		return "", err
	}

	oldOwner := herbBatch.Owner
	herbBatch.Owner = newOwner

	herbBatchJSON, err := json.Marshal(herbBatch)
	if err != nil {
		return "", err
	}

	err = ctx.GetStub().PutState(id, herbBatchJSON)
	if err != nil {
		return "", err
	}

	return oldOwner, nil
}

// GetAllHerbBatches returns all herb batches found in world state
func (s *SmartContract) GetAllHerbBatches(ctx contractapi.TransactionContextInterface) ([]*HerbBatch, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all herb batches in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var herbBatches []*HerbBatch
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var herbBatch HerbBatch
		err = json.Unmarshal(queryResponse.Value, &herbBatch)
		if err != nil {
			return nil, err
		}
		herbBatches = append(herbBatches, &herbBatch)
	}

	return herbBatches, nil
}

// UpdateHerbBatchStatus updates only the status field of a herb batch with given id in world state
func (s *SmartContract) UpdateHerbBatchStatus(ctx contractapi.TransactionContextInterface, id string, newStatus string) error {
	herbBatch, err := s.ReadHerbBatch(ctx, id)
	if err != nil {
		return err
	}

	herbBatch.Status = newStatus

	herbBatchJSON, err := json.Marshal(herbBatch)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, herbBatchJSON)
}
