package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type ChatEntry struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIRequest struct {
	Model    string      `json:"model"`
	Messages []ChatEntry `json:"messages"`
}

type OpenAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

var AIMapping = map[string]string{
	"6E":  "EUR/USD",
	"6A":  "AUD/USD",
	"6B":  "GBP/USD",
	"6J":  "USD/JPY",
	"CL":  "WTI Crude Oil",
	"GC":  "Gold",
	"ES":  "S&P 500",
	"NQ":  "Nasdaq 100",
	"YM":  "Dow Jones",
	"RTY": "Russell 2000",
}

func GetAISignal(c *gin.Context) {
	symbol := strings.ToUpper(c.Param("symbol"))
	mapped := AIMapping[symbol]
	if mapped == "" {
		mapped = symbol
	}

	openaiKey := os.Getenv("OPENAI_API_KEY")
	if openaiKey == "" {
		c.JSON(500, gin.H{"error": "Missing OpenAI API key"})
		return
	}

	// Simple dynamic context — you can enhance it later with live data
	prompt := fmt.Sprintf(`Analyze current market sentiment and risk bias for %s.

Return:
- Sentiment (bullish, bearish, neutral)
- Risk bias (risk-on, risk-off, neutral)
- A 2-sentence professional summary`, mapped)

	reqBody := OpenAIRequest{
		Model: "gpt-3.5-turbo",
		Messages: []ChatEntry{
			{Role: "system", Content: "You are a financial market sentiment analyst."},
			{Role: "user", Content: prompt},
		},
	}

	jsonData, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+openaiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.JSON(500, gin.H{"error": "OpenAI request failed", "details": err.Error()})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		c.JSON(500, gin.H{"error": "OpenAI returned non-200", "details": string(body)})
		return
	}

	var result OpenAIResponse
	_ = json.Unmarshal(body, &result)

	c.JSON(200, gin.H{
		"symbol":   symbol,
		"analyzed": mapped,
		"summary":  result.Choices[0].Message.Content,
		"source":   "OpenAI GPT-3.5",
	})
}
