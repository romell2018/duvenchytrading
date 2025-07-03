package services

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
)

type Candle struct {
	Time   string  `json:"datetime"`
	Open   float64 `json:"open,string"`
	High   float64 `json:"high,string"`
	Low    float64 `json:"low,string"`
	Close  float64 `json:"close,string"`
	Volume float64 `json:"volume,string"`
}

type TimeSeriesResponse struct {
	Values  []Candle `json:"values"`
	Status  string   `json:"status"`
	Message string   `json:"message"`
}

type SignalResult struct {
	Symbol           string  `json:"symbol"`
	RSI              float64 `json:"rsi"`
	MACDHistogram    float64 `json:"macd_histogram"`
	Signal           string  `json:"signal"`
	Reason           string  `json:"reason"`
	Support          float64 `json:"support"`
	Resistance       float64 `json:"resistance"`
	BBUpper          float64 `json:"bb_upper"`
	BBLower          float64 `json:"bb_lower"`
	BBMid            float64 `json:"bb_mid"`
	BollingerSqueeze bool    `json:"bollinger_squeeze"`
	Divergence       string  `json:"divergence"`
}

func AnalyzeSymbol(symbol string) (*SignalResult, error) {
	apiKey := os.Getenv("TWELVE_DATA_API_KEY")
	url := fmt.Sprintf("https://api.twelvedata.com/time_series?symbol=%s&interval=1h&outputsize=100&apikey=%s", symbol, apiKey)

	client := resty.New()
	resp, err := client.R().Get(url)
	if err != nil {
		return nil, err
	}

	var data TimeSeriesResponse
	if err := json.Unmarshal(resp.Body(), &data); err != nil {
		return nil, err
	}

	if data.Status == "error" {
		return nil, fmt.Errorf("Twelve Data error: %s", data.Message)
	}

	closes := make([]float64, len(data.Values))
	for i, c := range data.Values {
		closes[len(data.Values)-1-i] = c.Close
	}

	rsi := CalcRSI(closes, 14)
	macdLine, signalLine, _ := CalcMACD(closes, 12, 26, 9)
	macdHist := macdLine[len(macdLine)-1] - signalLine[len(signalLine)-1]

	support, resistance := DetectSupportResistance(closes)
	bbLower, bbMid, bbUpper := BollingerBands(closes, 20)
	bbSqueeze := (bbUpper-bbLower)/bbMid < 0.02
	divergence := DetectDivergence(closes, []float64{rsi})

	// ⬇️ Signal logic
	signal := "neutral"
	reason := "Indicators are mixed"

	if rsi < 30 && macdHist > 0 {
		signal = "bullish"
		reason = "RSI oversold and MACD histogram turning positive"
	} else if rsi > 70 && macdHist < 0 {
		signal = "bearish"
		reason = "RSI overbought and MACD histogram turning negative"
	} else if divergence == "bullish" {
		signal = "bullish"
		reason = "Bullish divergence detected between RSI and price"
	} else if divergence == "bearish" {
		signal = "bearish"
		reason = "Bearish divergence detected between RSI and price"
	} else {
		// ⬇️ Relaxed confluence logic
		if rsi < 45 && macdHist > 0 {
			signal = "bullish"
			reason = "RSI moderately oversold and MACD turning up"
		} else if rsi > 55 && macdHist < 0 {
			signal = "bearish"
			reason = "RSI weakening and MACD turning down"
		} else if bbSqueeze && macdHist > 0 {
			signal = "bullish"
			reason = "Bollinger Squeeze with MACD breakout"
		} else if closes[len(closes)-1] > resistance {
			signal = "bullish"
			reason = "Price breaking above resistance"
		} else if closes[len(closes)-1] < support {
			signal = "bearish"
			reason = "Price breaking below support"
		}
	}

	return &SignalResult{
		Symbol:           symbol,
		RSI:              rsi,
		MACDHistogram:    macdHist,
		Signal:           signal,
		Reason:           reason,
		Support:          support,
		Resistance:       resistance,
		BBUpper:          bbUpper,
		BBLower:          bbLower,
		BBMid:            bbMid,
		BollingerSqueeze: bbSqueeze,
		Divergence:       divergence,
	}, nil
}
