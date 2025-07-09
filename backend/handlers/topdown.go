package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetTopDown(c *gin.Context) {
	// TODO: Replace with dynamic sector data (from ETF APIs or sector map)
	topSectors := []string{"Technology", "Energy", "Industrials"}
	bottomSectors := []string{"Utilities", "Real Estate"}

	comment := "Tech and Energy leading based on price momentum. Defensive sectors underperforming."

	c.JSON(http.StatusOK, gin.H{
		"market_trend":    "bullish",
		"leading_sectors": topSectors,
		"lagging_sectors": bottomSectors,
		"commentary":      comment,
		"source":          "Mocked sector ranks (replace with live ETF data)",
	})
}
