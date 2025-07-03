package services

import (
	"math"
)

func DetectSupportResistance(data []float64) (support, resistance float64) {
	n := len(data)
	if n < 3 {
		return 0, 0
	}
	support = data[0]
	resistance = data[0]
	for _, price := range data {
		if price < support {
			support = price
		}
		if price > resistance {
			resistance = price
		}
	}
	return support, resistance
}

func BollingerBands(data []float64, period int) (float64, float64, float64) {
	if len(data) < period {
		return 0, 0, 0
	}

	sum := 0.0
	for i := len(data) - period; i < len(data); i++ {
		sum += data[i]
	}
	ma := sum / float64(period)

	variance := 0.0
	for i := len(data) - period; i < len(data); i++ {
		variance += math.Pow(data[i]-ma, 2)
	}
	stddev := math.Sqrt(variance / float64(period))

	upper := ma + 2*stddev
	lower := ma - 2*stddev
	return lower, ma, upper
}

func DetectDivergence(price []float64, rsi []float64) string {
	if len(price) < 2 || len(rsi) < 2 {
		return "none"
	}

	if price[len(price)-1] > price[len(price)-2] && rsi[len(rsi)-1] < rsi[len(rsi)-2] {
		return "bearish"
	}

	if price[len(price)-1] < price[len(price)-2] && rsi[len(rsi)-1] > rsi[len(rsi)-2] {
		return "bullish"
	}

	return "none"
}
