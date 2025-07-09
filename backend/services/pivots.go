package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

func parseFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

type Pivots struct {
	Support1    float64
	Resistance1 float64
}

func GetPivots(symbol string) (*Pivots, error) {
	apiKey := os.Getenv("TWELVE_DATA_API_KEY")
	url := fmt.Sprintf("https://api.twelvedata.com/pivot_points?symbol=%s&interval=1day&apikey=%s", symbol, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Values []struct {
			S1 string `json:"support1"`
			R1 string `json:"resistance1"`
		} `json:"values"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if len(result.Values) == 0 {
		return nil, fmt.Errorf("no pivot data returned")
	}

	s1, err := parseFloat(result.Values[0].S1)
	if err != nil {
		return nil, err
	}
	r1, err := parseFloat(result.Values[0].R1)
	if err != nil {
		return nil, err
	}

	return &Pivots{
		Support1:    s1,
		Resistance1: r1,
	}, nil
}
