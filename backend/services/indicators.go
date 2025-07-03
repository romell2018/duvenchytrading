package services

func EMA(data []float64, period int) []float64 {
	ema := make([]float64, len(data))
	k := 2.0 / float64(period+1)
	ema[0] = data[0]
	for i := 1; i < len(data); i++ {
		ema[i] = data[i]*k + ema[i-1]*(1-k)
	}
	return ema
}

func CalcMACD(data []float64, shortPeriod, longPeriod, signalPeriod int) ([]float64, []float64, []float64) {
	shortEMA := EMA(data, shortPeriod)
	longEMA := EMA(data, longPeriod)
	macd := make([]float64, len(data))
	for i := 0; i < len(data); i++ {
		macd[i] = shortEMA[i] - longEMA[i]
	}
	signal := EMA(macd, signalPeriod)
	histogram := make([]float64, len(data))
	for i := 0; i < len(data); i++ {
		histogram[i] = macd[i] - signal[i]
	}
	return macd, signal, histogram
}

func CalcRSI(closes []float64, period int) float64 {
	var gains, losses float64
	for i := 1; i <= period; i++ {
		diff := closes[i] - closes[i-1]
		if diff > 0 {
			gains += diff
		} else {
			losses -= diff
		}
	}

	if losses == 0 {
		return 100
	}

	rs := gains / losses
	rsi := 100 - (100 / (1 + rs))
	return rsi
}
