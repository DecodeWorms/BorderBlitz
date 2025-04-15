package handlers

import (
	"net/http"
	"strconv"

	"github.com/DecodeWorms/BorderBlitz/models"
	services "github.com/DecodeWorms/BorderBlitz/service"
	"github.com/gin-gonic/gin"
)

// WalletHandler handles HTTP requests related to wallets
type WalletHandler struct {
	walletService *services.WalletService
}

// NewWalletHandler creates a new wallet handler
func NewWalletHandler(walletService *services.WalletService) *WalletHandler {
	return &WalletHandler{
		walletService: walletService,
	}
}

// CreateWallet creates a new wallet
func (h *WalletHandler) CreateWallet(c *gin.Context) {
	userID := c.Query("user_id")
	var req models.CreateWalletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wallet, err := h.walletService.CreateWallet(userID, req.Email, req.UserType)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, wallet)
}

// GetWallet gets a wallet by ID
func (h *WalletHandler) GetWallet(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid wallet ID"})
		return
	}

	wallet, err := h.walletService.GetWallet(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "wallet not found"})
		return
	}

	c.JSON(http.StatusOK, wallet)
}

// GetWalletByUser gets a wallet by user ID and type
func (h *WalletHandler) GetWalletByUser(c *gin.Context) {
	userID := c.Query("user_id")
	userType := c.Query("user_type")

	if userID == "" || userType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id and user_type are required"})
		return
	}

	wallet, err := h.walletService.GetWalletByUserID(userID, userType)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "wallet not found"})
		return
	}

	c.JSON(http.StatusOK, wallet)
}

// Deposit adds funds to a wallet
func (h *WalletHandler) Deposit(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid wallet ID"})
		return
	}

	var req models.DepositRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	balance, err := h.walletService.DepositFunds(uint(id), req.StablecoinID, req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, balance)
}
