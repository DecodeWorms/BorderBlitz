package handlers

import (
	"net/http"

	"github.com/DecodeWorms/BorderBlitz/models"
	services "github.com/DecodeWorms/BorderBlitz/service"
	"github.com/gin-gonic/gin"
)

// SwapHandler handles HTTP requests related to swapping stablecoins
type SwapHandler struct {
	swapService *services.SwapService
}

// NewSwapHandler creates a new swap handler
func NewSwapHandler(swapService *services.SwapService) *SwapHandler {
	return &SwapHandler{
		swapService: swapService,
	}
}

// Swap swaps one stablecoin for another within the same wallet
func (h *SwapHandler) Swap(c *gin.Context) {
	var req models.SwapRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, err := h.swapService.Swap(
		req.WalletID,
		req.SourceCoinID,
		req.DestCoinID,
		req.Amount,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tx)
}

// GetExchangeRate gets the exchange rate between two stablecoins
func (h *SwapHandler) GetExchangeRate(c *gin.Context) {
	var req models.GetExchangeRateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rate, err := h.swapService.GetExchangeRate(req.SourceCoinID, req.DestCoinID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"source_coin_id": req.SourceCoinID,
		"dest_coin_id":   req.DestCoinID,
		"exchange_rate":  rate,
	})
}
