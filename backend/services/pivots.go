package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

type Pivots struct {
	Support1    float64 `json:"support1"`
	Resistance1 float64 `json:"resistance1"`
}

func parseFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

func GetPivots(symbol string) (*Pivots, error) {
	apiKey := os.Getenv("TWELVE_DATA_API_KEY")
	url := fmt.Sprintf("https://api.twelvedata.com/pivot_points?symbol=%s&interval=1day&apikey=%s", symbol, apiKey)

	fmt.Println("🔍 [Pivot] Request:", url)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("❌ Pivot HTTP error:", err)
		return fallbackToCalculatedPivots(symbol)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("📦 [Pivot] Raw response:", string(body))

	var result struct {
		Values []struct {
			S1 string `json:"support1"`
			R1 string `json:"resistance1"`
		} `json:"values"`
	}

	if err := json.Unmarshal(body, &result); err != nil || len(result.Values) == 0 {
		fmt.Println("⚠️ Pivot JSON decode failed or empty — falling back")
		return fallbackToCalculatedPivots(symbol)
	}

	s1, err1 := parseFloat(result.Values[0].S1)
	r1, err2 := parseFloat(result.Values[0].R1)
	if err1 != nil || err2 != nil {
		fmt.Println("⚠️ Pivot parse error — falling back")
		return fallbackToCalculatedPivots(symbol)
	}

	fmt.Println("✅ [Pivot] Success:", s1, r1)
	return &Pivots{Support1: s1, Resistance1: r1}, nil
}

func fallbackToCalculatedPivots(symbol string) (*Pivots, error) {
	apiKey := os.Getenv("TWELVE_DATA_API_KEY")
	url := fmt.Sprintf("https://api.twelvedata.com/time_series?symbol=%s&interval=1day&outputsize=1&apikey=%s", symbol, apiKey)

	fmt.Println("🔁 [Fallback] OHLC Request:", url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fallback OHLC HTTP error: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("📦 [Fallback] Raw OHLC:", string(body))

	var ohlc struct {
		Values []struct {
			High  string `json:"high"`
			Low   string `json:"low"`
			Close string `json:"close"`
		} `json:"values"`
	}
	if err := json.Unmarshal(body, &ohlc); err != nil || len(ohlc.Values) == 0 {
		return nil, fmt.Errorf("fallback OHLC JSON error or empty: %v", err)
	}

	high, err1 := parseFloat(ohlc.Values[0].High)
	low, err2 := parseFloat(ohlc.Values[0].Low)
	closePrice, err3 := parseFloat(ohlc.Values[0].Close)
	if err1 != nil || err2 != nil || err3 != nil {
		return nil, fmt.Errorf("fallback OHLC parse error")
	}

	pivot := (high + low + closePrice) / 3
	s1 := (2 * pivot) - high
	r1 := (2 * pivot) - low

	fmt.Printf("✅ [Fallback] Calculated: S1=%.4f, R1=%.4f\n", s1, r1)
	return &Pivots{Support1: s1, Resistance1: r1}, nil
}
