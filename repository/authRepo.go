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

// AuthAPI returns the base URL for interacting with the PocketBase users collection API.
func AuthAPI() string {
	return "http://127.0.0.1:8090/api/collections/users"
}

func LoginUser(login model.LoginRequest) (*model.LoginResponse, error) {
	// Validate input fields
	if login.Identity == "" || login.Password == "" {
		return nil, fmt.Errorf("identity and password are required")
	}

	// Prepare the request payload by converting it to JSON
	loginJSON, err := json.Marshal(map[string]interface{}{
		"identity": login.Identity,
		"password": login.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("error marshaling login data: %w", err)
	}
	// Create a POST request to PocketBase auth endpoint
	req, err := http.NewRequest("POST", AuthAPI()+"/auth-with-password", bytes.NewBuffer(loginJSON))

	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Log the request details for debugging
	fmt.Println(" Sending request to:", req.URL.String())
	fmt.Println(" Payload:", string(loginJSON))

	// Set up the HTTP client with a timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("login failed: %s", string(body))
	}

	// Parse the successful response into a temporary struct
	var response struct {
		Token  string `json:"token"`
		Record struct {
			ID              string `json:"id"`
			Username        string `json:"username"`
			Email           string `json:"email"`
			EmailVisibility bool   `json:"emailVisibility"`
			Created         string `json:"created"`
			Updated         string `json:"updated"`
			Verified        bool   `json:"verified"`
			Avatar          string `json:"avatar"`
		} `json:"record"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	// Convert to the internal LoginResponse format and return
	loginResponse := &model.LoginResponse{
		Token: response.Token,
		User: model.User{
			ID:              response.Record.ID,
			Username:        response.Record.Username,
			Email:           response.Record.Email,
			EmailVisibility: response.Record.EmailVisibility,
			Verified:        response.Record.Verified,
			Avatar:          response.Record.Avatar,
			CreatedAt:       response.Record.Created,
			UpdatedAt:       response.Record.Updated,
		},
	}

	return loginResponse, nil

}
