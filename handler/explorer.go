package handlers

import (
	"net/http"
	"strconv"

	services "github.com/DecodeWorms/BorderBlitz/service"
	"github.com/gin-gonic/gin"
)

// ExplorerHandler handles HTTP requests related to transaction exploring
type ExplorerHandler struct {
	transactionService *services.TransactionService
	walletService      *services.WalletService
}

// NewExplorerHandler creates a new explorer handler
func NewExplorerHandler(
	transactionService *services.TransactionService,
	walletService *services.WalletService,
) *ExplorerHandler {
	return &ExplorerHandler{
		transactionService: transactionService,
		walletService:      walletService,
	}
}

// GetWalletOverview gets an overview of a wallet, including balances and USD equivalent
func (h *ExplorerHandler) GetWalletOverview(c *gin.Context) {
	walletIDStr := c.Query("wallet_id")
	walletID, err := strconv.ParseUint(walletIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid wallet ID"})
		return
	}

	wallet, err := h.walletService.GetWallet(uint(walletID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "wallet not found"})
		return
	}

	// Calculate total USD value and prepare pie chart data
	type PieChartItem struct {
		Currency   string  `json:"currency"`
		Amount     float64 `json:"amount"`
		USDValue   float64 `json:"usd_value"`
		Percentage float64 `json:"percentage"`
	}

	pieChartData := make([]PieChartItem, 0, len(wallet.Balances))
	totalUSDValue := wallet.GetTotalUSDValue()

	for _, balance := range wallet.Balances {
		if balance.Amount > 0 {
			usdValue := balance.Stablecoin.GetUSDValue(balance.Amount)
			percentage := 0.0
			if totalUSDValue > 0 {
				percentage = (usdValue / totalUSDValue) * 100
			}

			pieChartData = append(pieChartData, PieChartItem{
				Currency:   balance.Stablecoin.Symbol,
				Amount:     balance.Amount,
				USDValue:   usdValue,
				Percentage: percentage,
			})
		}
	}

	// Get recent transactions
	transactions, _, err := h.transactionService.GetTransactionHistory(uint(walletID), 5, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch transaction history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"wallet":              wallet,
		"total_usd_value":     totalUSDValue,
		"pie_chart_data":      pieChartData,
		"recent_transactions": transactions,
	})
}
