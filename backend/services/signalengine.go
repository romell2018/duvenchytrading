package services

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
)

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
	Entry            float64 `json:"entry"`
	StopLoss         float64 `json:"stop_loss"`
	TakeProfit       float64 `json:"take_profit"`
	Confidence       string  `json:"confidence"`
}

type TimeSeriesResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Values  []struct {
		Close float64 `json:"close,string"`
	} `json:"values"`
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
		return nil, fmt.Errorf("twelve data error: %s", data.Message)
	}

	closes := make([]float64, len(data.Values))
	for i, c := range data.Values {
		closes[len(data.Values)-1-i] = c.Close
	}

	if len(closes) < 26 {
		return nil, fmt.Errorf("not enough data points for analysis")
	}

	rsi := CalcRSI(closes, 14)
	macdLine, signalLine, _ := CalcMACD(closes, 12, 26, 9)
	macdHist := macdLine[len(macdLine)-1] - signalLine[len(signalLine)-1]

	support, resistance := DetectSupportResistance(closes)
	bbLower, bbMid, bbUpper := BollingerBands(closes, 20)
	bbSqueeze := (bbUpper-bbLower)/bbMid < 0.02
	divergence := DetectDivergence(closes, []float64{rsi})
	latestPrice := closes[len(closes)-1]

	// 🔁 Generate pivot-based trade idea
	pivotSignal, pivotReason, entry, sl, tp := GenerateTradeSetup(latestPrice, support, resistance)

	// 🧠 Combine with indicator confluence
	signal := pivotSignal
	reason := pivotReason

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
	} else if rsi < 45 && macdHist > 0 {
		signal = "bullish"
		reason = "RSI moderately oversold and MACD turning up"
	} else if rsi > 55 && macdHist < 0 {
		signal = "bearish"
		reason = "RSI weakening and MACD turning down"
	} else if bbSqueeze && macdHist > 0 {
		signal = "bullish"
		reason = "Bollinger Squeeze with MACD breakout"
	}

	// 🧠 GPT Sentiment-Based Bias Lookup
	newsBias, err := AnalyzeNewsBias(symbol)
	gptSentiment := ""
	if err == nil {
		gptSentiment = newsBias.Sentiment
	} else {
		gptSentiment = "unknown"
	}

	// ✅ Confidence scoring
	confidence := "low"
	if gptSentiment == signal {
		confidence = "high"
	} else if gptSentiment == "neutral" {
		confidence = "medium"
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
		Entry:            entry,
		StopLoss:         sl,
		TakeProfit:       tp,
		Confidence:       confidence,
	}, nil

}

// GenerateTradeSetup provides a basic pivot-based trade signal and reason.
func GenerateTradeSetup(price, support, resistance float64) (string, string, float64, float64, float64) {
	rangeSize := resistance - support
	if rangeSize == 0 {
		rangeSize = price * 0.01 // fallback to 1% range if support == resistance
	}
	entryBuffer := rangeSize * 0.1
	stopBuffer := rangeSize * 0.05

	entry := 0.0
	sl := 0.0
	tp := 0.0
	signal := ""
	reason := ""

	if price > resistance+entryBuffer {
		entry = resistance + entryBuffer
		sl = resistance - stopBuffer
		tp = entry + rangeSize
		signal = "📈 Go long above resistance"
		reason = fmt.Sprintf("Breakout long. Entry: %.5f | SL: %.5f | TP: %.5f", entry, sl, tp)
	} else if price < support-entryBuffer {
		entry = support - entryBuffer
		sl = support + stopBuffer
		tp = entry - rangeSize
		signal = "📉 Short breakdown below support"
		reason = fmt.Sprintf("Breakdown short. Entry: %.5f | SL: %.5f | TP: %.5f", entry, sl, tp)
	} else {
		// 🔄 NEW: still return a potential trade idea (anticipatory setup)
		entry = resistance + entryBuffer
		sl = support - stopBuffer
		tp = entry + rangeSize
		signal = "🕵 Watching for breakout"
		reason = fmt.Sprintf("Price in range. Watching for breakout. Suggested long: Entry %.5f | SL %.5f | TP %.5f", entry, sl, tp)
	}

	return signal, reason, entry, sl, tp
}
