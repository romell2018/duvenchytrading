package handlers

import (
	"backend/services"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// type OpenAIRequest struct {
// 	Model    string      `json:"model"`
// 	Messages []ChatEntry `json:"messages"`
// }

// type ChatEntry struct {
// 	Role    string `json:"role"`
// 	Content string `json:"content"`
// }

// type OpenAIResponse struct {
// 	Choices []struct {
// 		Message struct {
// 			Content string `json:"content"`
// 		} `json:"message"`
// 	} `json:"choices"`
// }

type TwelveNewsResponse struct {
	Data []struct {
		Title string `json:"title"`
	} `json:"data"`
}

// symbol → result cache and timestamp
var newsCache = make(map[string]string)
var lastFetch = make(map[string]int64)
var mutex sync.Mutex

var symbolMap = map[string]string{
	"6E":  "EUR/USD",
	"6A":  "AUD/USD",
	"6B":  "GBP/USD",
	"CL":  "WTI",
	"GC":  "GOLD",
	"SI":  "SILVER",
	"ES":  "SPY",
	"NQ":  "QQQ",
	"RTY": "IWM",
	"ZB":  "TLT",
	"ZN":  "IEF",
	"MES": "SPY",
	"MNQ": "QQQ",
}

func GetNewsBias(c *gin.Context) {
	symbol := strings.ToUpper(c.Param("symbol"))
	mapped := symbolMap[symbol]
	if mapped == "" {
		mapped = symbol
	}

	mutex.Lock()
	lastTime := lastFetch[mapped]
	cached := newsCache[mapped]
	mutex.Unlock()

	if time.Now().Unix()-lastTime < 60 && cached != "" {
		c.JSON(http.StatusOK, gin.H{
			"symbol":    symbol,
			"sentiment": "cached",
			"bias":      "cached",
			"summary":   cached,
			"source":    "Cached GPT",
		})
		return
	}

	openaiKey := os.Getenv("OPENAI_API_KEY")
	if openaiKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Missing OpenAI API key"})
		return
	}

	// Get live price
	quote, err := services.GetQuote(mapped)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch price"})
		return
	}
	entry, err := strconv.ParseFloat(quote.Price, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse price"})
		return
	}

	// Get pivot levels
	pivots, err := services.GetPivots(mapped)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pivots"})
		return
	}

	// Dummy headlines
	headlines := []string{
		fmt.Sprintf("%s sees elevated volatility amid macro uncertainty.", mapped),
		fmt.Sprintf("Investors reevaluate their positions in %s.", mapped),
		fmt.Sprintf("Recent developments bring %s into focus.", mapped),
	}

	// Prompt with price + pivots
	prompt := fmt.Sprintf(`Analyze the following market conditions for %s and return:
- Sentiment (bullish, bearish, neutral)
- Market bias (risk-on, risk-off, neutral)
- A 2-sentence professional summary
- (Optional) Trade idea based on price action

Headlines:
%s

Current market data:
Price: %.4f
Support (S1): %.4f
Resistance (R1): %.4f`, mapped, strings.Join(headlines, "\n"), entry, pivots.Support1, pivots.Resistance1)

	// Send to OpenAI
	reqBody := OpenAIRequest{
		Model: "gpt-3.5-turbo",
		Messages: []ChatEntry{
			{Role: "system", Content: "You are a financial market analyst."},
			{Role: "user", Content: prompt},
		},
	}

	jsonData, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+openaiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "OpenAI request failed", "details": string(body)})
		return
	}
	defer resp.Body.Close()

	var result OpenAIResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode OpenAI response"})
		return
	}

	raw := result.Choices[0].Message.Content

	mutex.Lock()
	newsCache[mapped] = raw
	lastFetch[mapped] = time.Now().Unix()
	mutex.Unlock()

	c.JSON(http.StatusOK, gin.H{
		"symbol":     symbol,
		"sentiment":  "ai-generated",
		"bias":       "ai-generated",
		"summary":    raw,
		"price":      fmt.Sprintf("%.4f", entry),
		"support":    fmt.Sprintf("%.4f", pivots.Support1),
		"resistance": fmt.Sprintf("%.4f", pivots.Resistance1),
		"source":     "OpenAI GPT-3.5 with price + pivots",
	})
}
