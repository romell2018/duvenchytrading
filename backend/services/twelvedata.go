package services

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
)

const baseURL = "https://api.twelvedata.com"

type QuoteResponse struct {
	Symbol        string      `json:"symbol"`
	Name          string      `json:"name"`
	Exchange      string      `json:"exchange"`
	Currency      string      `json:"currency"`
	Price         string      `json:"price"`
	PreviousClose string      `json:"previous_close"`
	Change        string      `json:"change"`
	PercentChange string      `json:"percent_change"`
	Timestamp     json.Number `json:"timestamp"` // ✅ Fix here
	Status        string      `json:"status"`
	Message       string      `json:"message"`
}

// ✅ Make sure this function is called GetQuote (uppercase to export)
func GetQuote(symbol string) (*QuoteResponse, error) {
	apiKey := os.Getenv("TWELVE_DATA_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("Twelve Data API key not set in .env")
	}

	client := resty.New()
	url := fmt.Sprintf("%s/quote?symbol=%s&apikey=%s", baseURL, symbol, apiKey)

	resp, err := client.R().Get(url)
	if err != nil {
		return nil, err
	}

	var quote QuoteResponse
	err = json.Unmarshal(resp.Body(), &quote)
	if err != nil {
		return nil, err
	}

	if quote.Status == "error" {
		return nil, fmt.Errorf("API error: %s", quote.Message)
	}

	return &quote, nil
}
