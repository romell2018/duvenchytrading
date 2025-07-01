package services

// import (
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"os"
// 	"strconv"
// )

// func FetchDailyCandles(symbol string) (float64, float64, float64, error) {
// 	apiKey := os.Getenv("ALPHA_VANTAGE_API_KEY")
// 	if apiKey == "" {
// 		return 0, 0, 0, fmt.Errorf("ALPHA_VANTAGE_API_KEY environment variable not set")
// 	}

// 	if len(symbol) != 6 {
// 		return 0, 0, 0, fmt.Errorf("invalid symbol format")
// 	}

// 	from := symbol[:3]
// 	to := symbol[3:]

// 	url := fmt.Sprintf("https://www.alphavantage.co/query?function=FX_DAILY&from_symbol=%s&to_symbol=%s&apikey=%s", from, to, apiKey)

// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return 0, 0, 0, err
// 	}
// 	defer resp.Body.Close()

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return 0, 0, 0, err
// 	}

// 	fmt.Println("🔎 ALPHA RESPONSE:")
// 	fmt.Println(string(body)) // debug print

// 	var result struct {
// 		TimeSeries map[string]map[string]string `json:"Time Series FX (Daily)"`
// 	}

// 	err = json.Unmarshal(body, &result)
// 	if err != nil {
// 		return 0, 0, 0, err
// 	}

// 	// Get the most recent candle (the first key)
// 	for _, candle := range result.TimeSeries {
// 		highStr := candle["2. high"]
// 		lowStr := candle["3. low"]
// 		closeStr := candle["4. close"]

// 		high, err1 := strconv.ParseFloat(highStr, 64)
// 		low, err2 := strconv.ParseFloat(lowStr, 64)
// 		close, err3 := strconv.ParseFloat(closeStr, 64)

// 		if err1 != nil || err2 != nil || err3 != nil {
// 			return 0, 0, 0, fmt.Errorf("parse error")
// 		}

// 		return high, low, close, nil
// 	}

// 	return 0, 0, 0, fmt.Errorf("no data found")
// }
// type TimeSeriesResponse struct {
// 	Values []struct {
// 		Close string `json:"close"`
// 	} `json:"values"`
// 	Status string `json:"status"`
// }

// func FetchRecentCloses(symbol string) ([]float64, error) {
// 	apiKey := os.Getenv("TWELVE_DATA_API_KEY")
// 	if apiKey == "" {
// 		return nil, fmt.Errorf("TWELVE_DATA_API_KEY environment variable not set")
// 	}

// 	url := fmt.Sprintf("https://api.twelvedata.com/time_series?symbol=%s&interval=1day&outputsize=50&apikey=%s", symbol, apiKey)
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var data TimeSeriesResponse
// 	err = json.Unmarshal(body, &data)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if len(data.Values) == 0 {
// 		return nil, fmt.Errorf("no data returned")
// 	}

// 	closes := []float64{}
// 	for i := len(data.Values) - 1; i >= 0; i-- { // oldest to newest
// 		closePrice, err := strconv.ParseFloat(data.Values[i].Close, 64)
// 		if err != nil {
// 			return nil, err
// 		}
// 		closes = append(closes, closePrice)
// 	}

// 	return closes, nil
// }