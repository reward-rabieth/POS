package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rustacean-dev/possystem/model"
)

// CreateItem sends a POST request to PocketBase to create a new item record.
// It marshals the item data into JSON and includes the bearer token if provided.
// Requires the "Create" API rule in PocketBase to allow authenticated users:
// @request.auth.id != ""
func CreateItem(item model.Item, token string) error {
	itemJSON, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("failed to encode item: %w", err)
	}

	req, err := http.NewRequest("POST", "http://127.0.0.1:8090/api/collections/items/records", bytes.NewBuffer(itemJSON))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Add token if present
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("PocketBase error: %s", string(body))
	}

	return nil
}

// GetItemByID fetches a specific item record from PocketBase by ID.
// Requires "View" API rule: @request.auth.id != ""
func GetItemByID(id, token string) (model.Item, error) {
	url := fmt.Sprintf("http://127.0.0.1:8090/api/collections/items/records/%s", id)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return model.Item{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return model.Item{}, fmt.Errorf("item lookup failed (%d)", resp.StatusCode)
	}

	var rec struct {
		ID          string  `json:"id"`
		Name        string  `json:"name"`
		Price       float64 `json:"price"`
		Description string  `json:"description"`
		Quantity    int     `json:"quantity"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&rec); err != nil {
		return model.Item{}, err
	}

	return model.Item{
		ID:          rec.ID,
		Name:        rec.Name,
		Price:       rec.Price,
		Description: rec.Description,
		Quantity:    rec.Quantity,
	}, nil
}

// UpdateItemStock sends a PATCH request to PocketBase to update the quantity
// of an item by its ID. The request uses JSON and includes an authorization token.
//
// Returns an error if the update fails.
// UpdateItemStock modifies the quantity field of an existing item.
// Requires "Update" API rule: @request.auth.id != ""
func UpdateItemStock(id string, newQty int, token string) error {
	data := map[string]any{
		"quantity": newQty,
	}

	body, _ := json.Marshal(data)
	req, _ := http.NewRequest("PATCH", fmt.Sprintf("http://127.0.0.1:8090/api/collections/items/records/%s", id), bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("failed to update stock (%d)", resp.StatusCode)
	}

	return nil
}
