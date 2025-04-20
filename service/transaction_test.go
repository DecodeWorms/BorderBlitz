package services

import (
	"testing"

	"github.com/DecodeWorms/BorderBlitz/mocks"
	"github.com/DecodeWorms/BorderBlitz/models"
	"go.uber.org/mock/gomock"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestTransactionService_Transfer_ValidAmount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create real GORM DB (won't be used due to early return)
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect test database: %v", err)
	}

	// Create mocked repositories
	txRepo := mocks.NewMockTransaction(ctrl)
	walletRepo := mocks.NewMockWallet(ctrl)
	coinRepo := mocks.NewMockStableCoin(ctrl)

	// Initialize service
	service := NewTransactionService(db, txRepo, walletRepo, coinRepo)

	// Test invalid amount
	_, err = service.Transfer(1, 2, 1, 2, -100)

	// Assert
	expectedErr := "transfer amount must be positive"
	if err == nil || err.Error() != expectedErr {
		t.Errorf("expected error %q, got %v", expectedErr, err)
	}
}

func TestTransactionService_Transfer_InValidAmount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create real GORM DB (won't be used due to early return)
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect test database: %v", err)
	}

	// Create mocked repositories
	txRepo := mocks.NewMockTransaction(ctrl)
	walletRepo := mocks.NewMockWallet(ctrl)
	coinRepo := mocks.NewMockStableCoin(ctrl)

	// Initialize service
	service := NewTransactionService(db, txRepo, walletRepo, coinRepo)

	// Test invalid amount
	_, err = service.Transfer(1, 2, 1, 2, -100)

	// Assert
	expectedErr := "transfer amount must be positive"
	if err == nil || err.Error() != expectedErr {
		t.Errorf("expected error %q, got %v", expectedErr, err)
	}
}

func TestTransactionService_GetTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mocked repository
	txRepo := mocks.NewMockTransaction(ctrl)

	// Create a mock transaction
	mockTx := &models.Transaction{
		ID:                uint(1),
		Type:              models.TransactionTransfer,
		SenderWalletID:    uintPtr(1),
		ReceiverWalletID:  uintPtr(2),
		SourceCoinID:      uint(1),
		SourceAmount:      100,
		DestinationCoinID: uint(2),
		DestinationAmount: 100,
		ExchangeRate:      1.0,
		Status:            "completed",
		Reference:         "tx123",
	}

	// Mock the FindByID method to return the mock transaction
	txRepo.EXPECT().FindByID(uint(1)).Return(mockTx, nil)

	// Initialize the service with the mocked repository
	service := NewTransactionService(nil, txRepo, nil, nil)

	// Test GetTransaction method with valid ID
	tx, err := service.GetTransaction(uint(1))

	// Assert that no error occurred
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Assert that the returned transaction matches the expected mock
	if tx.ID != mockTx.ID || tx.Reference != mockTx.Reference {
		t.Errorf("expected transaction %+v, got %+v", mockTx, tx)
	}
}

// Helper function to create pointers to uint
func uintPtr(i uint) *uint {
	return &i
}

func TestTransactionService_GetTransactionHistory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mocked repository
	txRepo := mocks.NewMockTransaction(ctrl)

	// Create a mock transaction list
	mockTxs := []models.Transaction{
		{
			ID:                uint(1),
			Type:              models.TransactionTransfer,
			SenderWalletID:    uintPtr(1),
			ReceiverWalletID:  uintPtr(2),
			SourceCoinID:      uint(1),
			SourceAmount:      100,
			DestinationCoinID: uint(2),
			DestinationAmount: 100,
			ExchangeRate:      1.0,
			Status:            "completed",
			Reference:         "tx123",
		},
		{
			ID:                uint(2),
			Type:              models.TransactionTransfer,
			SenderWalletID:    uintPtr(1),
			ReceiverWalletID:  uintPtr(3),
			SourceCoinID:      uint(1),
			SourceAmount:      200,
			DestinationCoinID: uint(2),
			DestinationAmount: 200,
			ExchangeRate:      1.0,
			Status:            "completed",
			Reference:         "tx124",
		},
	}

	// Mock the FindByWalletID method to return the mock transactions
	txRepo.EXPECT().FindByWalletID(uint(1), 10, 0).Return(mockTxs, nil)

	// Mock the CountByWalletID method to return the count of transactions
	txRepo.EXPECT().CountByWalletID(uint(1)).Return(int64(2), nil)

	// Initialize the service with the mocked repository
	service := NewTransactionService(nil, txRepo, nil, nil)

	// Test GetTransactionHistory method with valid parameters
	transactions, count, err := service.GetTransactionHistory(uint(1), 10, 0)

	// Assert no error occurred
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Assert that the returned transactions match the mock
	if len(transactions) != len(mockTxs) {
		t.Errorf("expected %d transactions, got %d", len(mockTxs), len(transactions))
	}

	// Assert that the count matches the expected count
	if count != 2 {
		t.Errorf("expected count %d, got %d", 2, count)
	}
}
