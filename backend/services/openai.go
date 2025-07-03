package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type NewsBiasResult struct {
	Symbol    string `json:"symbol"`
	Sentiment string `json:"sentiment"`
	Bias      string `json:"bias"`
	Summary   string `json:"summary"`
	Source    string `json:"source"`
}

func AnalyzeNewsBias(symbol string) (*NewsBiasResult, error) {
	dummyHeadline := fmt.Sprintf("Recent economic data suggests growing uncertainty around %s", symbol)

	prompt := fmt.Sprintf(`Analyze the following news for trading sentiment and bias:
Headline: %s

Return a JSON with:
- sentiment (bullish, bearish, neutral)
- bias (e.g. risk-on, hawkish, dovish, inflation-fear, etc.)
- short summary`, dummyHeadline)

	payload := map[string]interface{}{
		"model": "gpt-4",
		"messages": []map[string]string{
			{"role": "system", "content": "You are a financial news sentiment analyst."},
			{"role": "user", "content": prompt},
		},
		"temperature": 0.3,
	}

	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY"))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	raw, _ := io.ReadAll(res.Body)

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, err
	}

	var parsed NewsBiasResult
	if err := json.Unmarshal([]byte(result.Choices[0].Message.Content), &parsed); err != nil {
		return nil, err
	}

	parsed.Symbol = symbol
	parsed.Source = "GPT-4 (simulated)"
	return &parsed, nil
}
