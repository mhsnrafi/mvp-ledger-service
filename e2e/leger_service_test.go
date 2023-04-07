package e2e

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"ledger-service/services"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

type addFundsResponse struct {
	Success bool `json:"success"`
	Data    struct {
		ID      string `json:"Id"`
		Message string `json:"Message"`
	} `json:"data"`
}

type getBalanceResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Balance float64 `json:"Balance"`
	} `json:"data"`
}

func TestAddFunds(t *testing.T) {
	services.LoadConfig()
	services.ConnectDB()

	// Prepare request data for adding funds
	data := `{"amount": 100}`

	// Use an existing user ID for testing
	userID := "9f3a1d82c5e74e2b"

	// Make an HTTP request to the API endpoint
	resp, err := http.Post(fmt.Sprintf("http://nginx:4000/v1/users/%s/add", userID), "application/json", strings.NewReader(data))
	if err != nil {
		t.Fatalf("Failed to make HTTP request: %s", err.Error())
	}
	defer resp.Body.Close()

	// Check the status code of the HTTP response
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("HTTP request returned status code %d", resp.StatusCode)
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %s", err.Error())
	}

	// Unmarshal the response JSON
	var result addFundsResponse
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("Failed to unmarshal response JSON: %s", err.Error())
	}

	// Check the response content
	if !result.Success {
		t.Fatalf("Unexpected success status: %v", result.Success)
	}
	if result.Data.ID != "9f3a1d82c5e74e2b" {
		t.Fatalf("Unexpected user ID: %s", result.Data.ID)
	}
	if result.Data.Message != "Funds added successfully" {
		t.Fatalf("Unexpected message: %s", result.Data.Message)
	}
}

// Test case for getting the balance of a user

func TestGetBalance(t *testing.T) {
	services.LoadConfig()
	services.ConnectDB()

	// Use an existing user ID for testing
	userID := "9f3a1d82c5e74e2b"

	// Make an HTTP request to the API endpoint
	resp, err := http.Get(fmt.Sprintf("http://nginx:4000/v1/users/%s/balance", userID))
	if err != nil {
		t.Fatalf("Failed to make HTTP request: %s", err.Error())
	}
	defer resp.Body.Close()

	// Check the status code of the HTTP response
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("HTTP request returned status code %d", resp.StatusCode)
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %s", err.Error())
	}

	// Unmarshal the response JSON
	var balance getBalanceResponse
	if err := json.Unmarshal(body, &balance); err != nil {
		t.Fatalf("Failed to unmarshal response JSON: %s", err.Error())
	}

	// Check the response content
	if balance.Data.Balance <= 0 {
		t.Fatalf("Unexpected balance: %f", balance.Data.Balance)
	}
}

// Test case for getting the transaction history of a user
func TestGetTransactionHistory(t *testing.T) {
	services.LoadConfig()
	services.ConnectDB()
	userID := "9f3a1d82c5e74e2b"

	resp, err := http.Get(fmt.Sprintf("http://nginx:4000/v1/users/%s/history", userID))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	var history map[string]interface{}
	err = json.Unmarshal(body, &history)
	require.NoError(t, err)

	require.NotNil(t, history)
	require.NotEmpty(t, history)
}

func TestAddFundsInvalidAmount(t *testing.T) {
	services.LoadConfig()
	services.ConnectDB()

	resp, err := http.PostForm("http://nginx:4000/v1/users/9f3a1d82c5e74e2b/add", url.Values{"amount": {"-100"}})
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	var errorResponse map[string]interface{}
	err = json.Unmarshal(body, &errorResponse)
	require.NoError(t, err)

	require.NotNil(t, errorResponse)
	require.NotEmpty(t, errorResponse["error"])
}

func TestGetBalanceNonExistentUser(t *testing.T) {
	services.LoadConfig()
	services.ConnectDB()

	resp, err := http.Get("http://nginx:4000/v1/users/9999/balance")
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	var errorResponse map[string]interface{}
	err = json.Unmarshal(body, &errorResponse)
	require.NoError(t, err)

	expectedMessage := "User not found"
	require.NotNil(t, errorResponse)
	require.Equal(t, expectedMessage, errorResponse["message"], "Expected error message: %s, got: %s", expectedMessage, errorResponse["message"])
}
