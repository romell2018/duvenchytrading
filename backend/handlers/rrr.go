package handlers

import (
	"backend/config"
	"backend/services"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetRRR(c *gin.Context) {
	symbol := strings.ToUpper(c.Param("symbol"))

	if !config.IsSupported(symbol) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported symbol"})
		return
	}

	mapped, ok := config.SymbolMap[symbol]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Symbol not mapped"})
		return
	}

	// 🔥 Get live price (as string)
	quote, err := services.GetQuote(mapped)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch quote"})
		return
	}

	// ✅ Convert price string to float64
	entry, err := strconv.ParseFloat(quote.Price, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse price"})
		return
	}

	// ✅ Get pivots (float64 already)
	pivots, err := services.GetPivots(mapped)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pivots"})
		return
	}
	target := pivots.Resistance1
	stop := pivots.Support1

	// ✅ Check structure
	if entry <= stop || target <= entry {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid price structure for RRR"})
		return
	}

	rrr := (target - entry) / (entry - stop)

	// ✅ Send response
	response := gin.H{
		"symbol":       symbol,
		"entry_price":  fmt.Sprintf("%.2f", entry),
		"target_price": fmt.Sprintf("%.2f", target),
		"stop_loss":    fmt.Sprintf("%.2f", stop),
		"rr_ratio":     fmt.Sprintf("%.2f", rrr),
		"notes":        "RRR based on S1/R1 pivots. Adjust based on timeframe and volatility.",
	}

	c.JSON(http.StatusOK, response)
}
