// tech analysis, trading signals, and market data
package handlers

import (
	"backend/config"
	"backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetSignal(c *gin.Context) {
	symbol := c.Param("symbol")

	if !config.IsSupported(symbol) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported symbol"})
		return
	}

	mapped, ok := config.SymbolMap[symbol]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Symbol not mapped"})
		return
	}

	signal, err := services.AnalyzeSymbol(mapped)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, signal)
}
