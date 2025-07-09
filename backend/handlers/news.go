package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type OpenAIRequest struct {
	Model    string      `json:"model"`
	Messages []ChatEntry `json:"messages"`
}

type ChatEntry struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

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
			"source":    "Cached GPT + Twelve Data",
		})
		return
	}

	tdKey := os.Getenv("TWELVE_DATA_API_KEY")
	openaiKey := os.Getenv("OPENAI_API_KEY")
	if tdKey == "" || openaiKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Missing API keys"})
		return
	}

	newsURL := fmt.Sprintf("https://api.twelvedata.com/news?symbol=%s&apikey=%s", mapped, tdKey)
	res, err := http.Get(newsURL)
	if err != nil || res.StatusCode != 200 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Twelve Data fetch failed"})
		return
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	var newsResp TwelveNewsResponse
	_ = json.Unmarshal(body, &newsResp)

	headlines := []string{}
	for _, item := range newsResp.Data {
		if len(headlines) < 5 {
			headlines = append(headlines, item.Title)
		}
	}

	if len(headlines) == 0 {
		headlines = []string{
			fmt.Sprintf("%s volatility increases amid macro pressures", mapped),
			fmt.Sprintf("Risk sentiment shifting against %s", mapped),
		}
	}

	prompt := fmt.Sprintf(`Analyze the following headlines for %s and return:
- Sentiment (bullish, bearish, neutral)
- Market bias (risk-on, risk-off, neutral)
- A 2-sentence professional summary

Headlines:
%s`, mapped, strings.Join(headlines, "\n"))

	reqBody := OpenAIRequest{
		Model: "gpt-3.5-turbo",
		Messages: []ChatEntry{
			{Role: "system", Content: "You are a financial news sentiment analyst."},
			{Role: "user", Content: prompt},
		},
	}

	jsonData, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+openaiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != 200 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "OpenAI error"})
		return
	}
	defer resp.Body.Close()

	var result OpenAIResponse
	_ = json.NewDecoder(resp.Body).Decode(&result)
	raw := result.Choices[0].Message.Content

	// Cache result
	mutex.Lock()
	newsCache[mapped] = raw
	lastFetch[mapped] = time.Now().Unix()
	mutex.Unlock()

	c.JSON(http.StatusOK, gin.H{
		"symbol":    symbol,
		"sentiment": "ai-generated",
		"bias":      "ai-generated",
		"summary":   raw,
		"source":    "GPT-3.5 + Twelve Data (cached)",
	})
}
