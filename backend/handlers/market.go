package handlers

import (
	"backend/config"
	"backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetQuote handles the request to fetch the live price and chart data for a given symbol.
func GetQuote(c *gin.Context) {
	symbol := c.Param("symbol")

	if !config.IsSupported(symbol) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Symbol not supported"})
		return
	}

	mapped, ok := config.SymbolMap[symbol]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Symbol not mapped to Twelve Data"})
		return
	}

	quote, err := services.GetQuote(mapped)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, quote)
}
