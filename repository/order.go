package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rustacean-dev/possystem/model"
)

// GetAllOrders retrieves all order records from PocketBase,
// using the `expand=user_id` query to also fetch related user info.
// Requires "List/Search" access rule: @request.auth.id != "
func GetAllOrders(token string) ([]model.Order, error) {
	url := "http://127.0.0.1:8090/api/collections/orders/records?expand=user_id"

	// Create request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Set Authorization header with Bearer
	req.Header.Set("Authorization", "Bearer "+token)
	fmt.Println("🔐 Token:", token)

	// Perform request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// TEMP: Dump full raw body before decoding
	raw, _ := io.ReadAll(resp.Body)
	//fmt.Println("\n RAW JSON FROM POCKETBASE ↓\n", string(raw), "\n")

	// Reuse raw body for decoding
	resp.Body = io.NopCloser(bytes.NewBuffer(raw))

	// Check for non-200 response
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch orders: %s", string(raw))
	}

	// Decode response
	var res struct {
		Items []model.Order `json:"items"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	// Print createdAt from each order to verify
	for _, order := range res.Items {
		fmt.Printf("✅ Order %s created at %s\n", order.ID, order.CreatedAt)
	}

	return res.Items, nil
}

// CreateOrder sends a new order to PocketBase for storage.
// Requires "Create" rule on the 'orders' collection: @request.auth.id != ""
func CreateOrder(order model.Order, token string) error {
	data, err := json.Marshal(order)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "http://127.0.0.1:8090/api/collections/orders/records", bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to create order (%d): %s", resp.StatusCode, string(body))
	}
	return nil
}

// GetAllItems fetches all item records from PocketBase.
// Requires "List/Search" access on 'items': @request.auth.id != ""
func GetAllItems(token string) ([]model.Item, error) {
	req, err := http.NewRequest("GET", "http://127.0.0.1:8090/api/collections/items/records", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to fetch items: %s", string(body))
	}

	// Decode JSON response
	var res struct {
		Items []model.Item `json:"items"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}
	return res.Items, nil
}

// GetItemByName fetches an item from PocketBase using a filter on the name field.
// Requires "List/Search" access: @request.auth.id != ""
func GetItemByName(name, token string) (model.Item, error) {
	url := fmt.Sprintf("http://127.0.0.1:8090/api/collections/items/records?filter=name=\"%s\"", name)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return model.Item{}, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return model.Item{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return model.Item{}, fmt.Errorf("failed to get item: %s", string(body))
	}

	var res struct {
		Items []model.Item `json:"items"`
	}
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return model.Item{}, err
	}
	if len(res.Items) == 0 {
		return model.Item{}, fmt.Errorf("item not found")
	}
	return res.Items[0], nil
}
