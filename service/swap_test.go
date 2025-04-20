package services

import (
	"testing"

	"github.com/DecodeWorms/BorderBlitz/mocks"
	"github.com/DecodeWorms/BorderBlitz/models"
	"go.uber.org/mock/gomock"
)

func TestSwapService_GetExchangeRate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock stablecoin repository
	stablecoinRepo := mocks.NewMockStableCoin(ctrl)

	// Define source and destination coins
	sourceCoin := &models.Stablecoin{
		ID:      1,
		Symbol:  "USDT",
		Name:    "Tether",
		USDRate: 1.0,
	}
	destCoin := &models.Stablecoin{
		ID:      2,
		Symbol:  "DAI",
		Name:    "DAI",
		USDRate: 0.5,
	}

	// Setup expectations
	stablecoinRepo.EXPECT().FindByID(uint(1)).Return(sourceCoin, nil)
	stablecoinRepo.EXPECT().FindByID(uint(2)).Return(destCoin, nil)

	// No DB involved
	service := NewSwapService(nil, nil, stablecoinRepo, nil)

	// Execute
	rate, err := service.GetExchangeRate(1, 2)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedRate := 2.0 // 1.0 / 0.5
	if rate != expectedRate {
		t.Errorf("expected exchange rate %.2f, got %.2f", expectedRate, rate)
	}
}
