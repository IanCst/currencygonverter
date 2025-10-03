package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	API_KEY = "cur_live_eUg8Pk4ZIf94e9ESkagJ8TrzG08ShExFJIq0KEvO"
	BASE_URL = "https://api.currencyapi.com/v3/latest"
)

type CurrencyResponse struct {
	Meta struct {
		LastUpdatedAt string `json:"last_updated_at"`
	} `json:"meta"`
	Data map[string]struct {
		Code  string  `json:"code"`
		Value float64 `json:"value"`
	} `json:"data"`
}

type CurrencyData struct {
	USD      float64
	BRL      float64
	LastUpdated string
}

func fetchCurrencyRates() (*CurrencyData, error) {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Build URL with parameters
	url := fmt.Sprintf("%s?apikey=%s&currencies=USD,BRL&base_currency=USD", BASE_URL, API_KEY)

	// Make the request
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	// Parse JSON response
	var currencyResp CurrencyResponse
	if err := json.Unmarshal(body, &currencyResp); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	// Extract currency values
	currencyData := &CurrencyData{
		LastUpdated: currencyResp.Meta.LastUpdatedAt,
	}

	if usd, exists := currencyResp.Data["USD"]; exists {
		currencyData.USD = usd.Value
	}

	if brl, exists := currencyResp.Data["BRL"]; exists {
		currencyData.BRL = brl.Value
	}

	return currencyData, nil
}

func (cd *CurrencyData) String() string {
	if cd == nil {
		return "No currency data available"
	}
	
	return fmt.Sprintf("USD: $%.4f\nBRL: R$%.4f\nLast Updated: %s", 
		cd.USD, cd.BRL, cd.LastUpdated)
}