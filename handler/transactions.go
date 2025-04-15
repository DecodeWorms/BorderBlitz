package handlers

import (
	"net/http"
	"strconv"

	"github.com/DecodeWorms/BorderBlitz/models"
	services "github.com/DecodeWorms/BorderBlitz/service"
	"github.com/gin-gonic/gin"
)

// TransactionHandler handles HTTP requests related to transactions
type TransactionHandler struct {
	transactionService *services.TransactionService
}

// NewTransactionHandler creates a new transaction handler
func NewTransactionHandler(transactionService *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
	}
}

// Transfer transfers funds from one wallet to another
func (h *TransactionHandler) Transfer(c *gin.Context) {
	var req models.TransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, err := h.transactionService.Transfer(
		req.SenderWalletID,
		req.ReceiverWalletID,
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

// Deposit adds funds to a wallet
func (h *TransactionHandler) Deposit(c *gin.Context) {
	var req models.DepositTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, err := h.transactionService.Deposit(req.WalletID, req.StablecoinID, req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tx)
}

// GetTransaction gets a transaction by ID
func (h *TransactionHandler) GetTransaction(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction ID"})
		return
	}

	tx, err := h.transactionService.GetTransaction(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
		return
	}

	c.JSON(http.StatusOK, tx)
}

// GetTransactionHistory gets the transaction history for a wallet
func (h *TransactionHandler) GetTransactionHistory(c *gin.Context) {
	walletIDStr := c.Query("wallet_id")
	walletID, err := strconv.ParseUint(walletIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid wallet ID"})
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	offsetStr := c.DefaultQuery("offset", "0")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	transactions, total, err := h.transactionService.GetTransactionHistory(uint(walletID), limit, offset)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transactions": transactions,
		"total":        total,
		"limit":        limit,
		"offset":       offset,
	})
}
