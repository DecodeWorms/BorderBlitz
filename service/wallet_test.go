package services

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/DecodeWorms/BorderBlitz/mocks"
	"github.com/DecodeWorms/BorderBlitz/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestWalletService_CreateWallet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWalletRepo := mocks.NewMockWallet(ctrl)
	mockStableCoinRepo := mocks.NewMockStableCoin(ctrl)
	mockUserRepo := mocks.NewMockUser(ctrl)

	service := NewWalletService(mockWalletRepo, mockStableCoinRepo, mockUserRepo)

	tests := []struct {
		name           string
		userID         string
		email          string
		userType       string
		mockUser       *models.Users
		mockWallet     *models.Wallet
		mockFindUser   error
		mockFindWallet error
		mockCreateErr  error
		expectedErr    string
	}{
		{
			name:           "successfully creates wallet",
			userID:         "user-123",
			email:          "user@example.com",
			userType:       "email",
			mockUser:       &models.Users{ID: "user-123"},
			mockFindUser:   nil,
			mockFindWallet: fmt.Errorf("wallet not found"),
			mockCreateErr:  nil,
		},
		{
			name:         "user not found",
			userID:       "user-123",
			userType:     "email",
			mockFindUser: errors.New("not found"),
			expectedErr:  "user's record is not found",
		},
		{
			name:         "invalid user type",
			userID:       "user-123",
			userType:     "social",
			mockUser:     &models.Users{ID: "user-123"},
			mockFindUser: nil,
			expectedErr:  "user type must be 'email' or 'mobile'",
		},
		{
			name:           "wallet already exists",
			userID:         "user-123",
			userType:       "mobile",
			mockUser:       &models.Users{ID: "user-123"},
			mockFindUser:   nil,
			mockFindWallet: nil,
			expectedErr:    "wallet already exists for this mobile",
		},
		{
			name:           "failed to create wallet",
			userID:         "user-123",
			userType:       "email",
			mockUser:       &models.Users{ID: "user-123"},
			mockFindUser:   nil,
			mockFindWallet: fmt.Errorf("wallet not found"),
			mockCreateErr:  errors.New("db error"),
			expectedErr:    "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockFindUser == nil {
				mockUserRepo.EXPECT().
					FindByUserID(tt.userID).
					Return(tt.mockUser, nil)
			} else {
				mockUserRepo.EXPECT().
					FindByUserID(tt.userID).
					Return(nil, tt.mockFindUser)
			}

			if tt.mockFindUser == nil && tt.expectedErr == "" || strings.HasPrefix(tt.expectedErr, "wallet") || tt.mockCreateErr != nil {
				mockWalletRepo.EXPECT().
					FindByUserID(tt.userID, strings.ToLower(tt.userType)).
					Return(tt.mockWallet, tt.mockFindWallet)
			}

			if tt.mockCreateErr != nil || tt.expectedErr == "" && tt.mockFindWallet != nil {
				mockWalletRepo.EXPECT().
					Create(gomock.Any()).
					Return(tt.mockCreateErr)
			}

			result, err := service.CreateWallet(tt.userID, tt.email, tt.userType)

			if tt.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				require.Equal(t, strings.ToLower(tt.userType), result.UserType)
			}
		})
	}
}

func TestWalletService_GetWallet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successfully_gets_wallet", func(t *testing.T) {
		mockWalletRepo := mocks.NewMockWallet(ctrl)
		mockStableCoinRepo := mocks.NewMockStableCoin(ctrl)
		mockUserRepo := mocks.NewMockUser(ctrl)

		expectedWallet := &models.Wallet{
			ID:       1,
			UserID:   "user-uuid",
			UserType: "email",
		}

		mockWalletRepo.
			EXPECT().
			FindByID(uint(1)).
			Return(expectedWallet, nil)

		service := NewWalletService(mockWalletRepo, mockStableCoinRepo, mockUserRepo)

		wallet, err := service.GetWallet(1)
		assert.NoError(t, err)
		assert.Equal(t, expectedWallet, wallet)
	})

	t.Run("wallet_not_found", func(t *testing.T) {
		mockWalletRepo := mocks.NewMockWallet(ctrl)
		mockStableCoinRepo := mocks.NewMockStableCoin(ctrl)
		mockUserRepo := mocks.NewMockUser(ctrl)

		mockWalletRepo.
			EXPECT().
			FindByID(uint(2)).
			Return(nil, errors.New("not found"))

		service := NewWalletService(mockWalletRepo, mockStableCoinRepo, mockUserRepo)

		wallet, err := service.GetWallet(2)
		assert.Nil(t, wallet)
		assert.EqualError(t, err, "error wallet not available")
	})
}

func TestWalletService_GetWalletByUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successfully_gets_wallet_by_user_id", func(t *testing.T) {
		mockWalletRepo := mocks.NewMockWallet(ctrl)
		mockStableCoinRepo := mocks.NewMockStableCoin(ctrl)
		mockUserRepo := mocks.NewMockUser(ctrl)

		userID := "user-uuid"
		userType := "email"

		expectedWallet := &models.Wallet{
			UserID:   userID,
			UserType: userType,
		}

		mockWalletRepo.
			EXPECT().
			FindByUserID(userID, userType).
			Return(expectedWallet, nil)

		service := NewWalletService(mockWalletRepo, mockStableCoinRepo, mockUserRepo)

		wallet, err := service.GetWalletByUserID(userID, userType)
		assert.NoError(t, err)
		assert.Equal(t, expectedWallet, wallet)
	})

	t.Run("wallet_not_found_by_user_id", func(t *testing.T) {
		mockWalletRepo := mocks.NewMockWallet(ctrl)
		mockStableCoinRepo := mocks.NewMockStableCoin(ctrl)
		mockUserRepo := mocks.NewMockUser(ctrl)

		userID := "non-existent-user"
		userType := "email"

		mockWalletRepo.
			EXPECT().
			FindByUserID(userID, userType).
			Return(nil, errors.New("not found"))

		service := NewWalletService(mockWalletRepo, mockStableCoinRepo, mockUserRepo)

		wallet, err := service.GetWalletByUserID(userID, userType)
		assert.Nil(t, wallet)
		assert.EqualError(t, err, "error wallet not available")
	})
}

func TestWalletService_GetBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successfully_gets_wallet_balance", func(t *testing.T) {
		mockWalletRepo := mocks.NewMockWallet(ctrl)
		mockStableCoinRepo := mocks.NewMockStableCoin(ctrl)
		mockUserRepo := mocks.NewMockUser(ctrl)

		walletID := uint(1)
		stablecoinID := uint(100)

		expectedBalance := &models.Balance{
			WalletID:     walletID,
			StablecoinID: stablecoinID,
			Amount:       150.75,
		}

		mockWalletRepo.
			EXPECT().
			GetBalance(walletID, stablecoinID).
			Return(expectedBalance, nil)

		service := NewWalletService(mockWalletRepo, mockStableCoinRepo, mockUserRepo)

		balance, err := service.GetBalance(walletID, stablecoinID)
		assert.NoError(t, err)
		assert.Equal(t, expectedBalance, balance)
	})

	t.Run("balance_not_found_or_error", func(t *testing.T) {
		mockWalletRepo := mocks.NewMockWallet(ctrl)
		mockStableCoinRepo := mocks.NewMockStableCoin(ctrl)
		mockUserRepo := mocks.NewMockUser(ctrl)

		walletID := uint(999)
		stablecoinID := uint(404)

		mockWalletRepo.
			EXPECT().
			GetBalance(walletID, stablecoinID).
			Return(nil, errors.New("balance not found"))

		service := NewWalletService(mockWalletRepo, mockStableCoinRepo, mockUserRepo)

		balance, err := service.GetBalance(walletID, stablecoinID)
		assert.Nil(t, balance)
		assert.EqualError(t, err, "balance not found")
	})
}

func TestWalletService_DepositFunds(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successfully_deposits_funds", func(t *testing.T) {
		mockWalletRepo := mocks.NewMockWallet(ctrl)
		mockStableCoinRepo := mocks.NewMockStableCoin(ctrl)
		mockUserRepo := mocks.NewMockUser(ctrl)

		walletID := uint(1)
		stablecoinID := uint(101)
		initialAmount := 50.0
		depositAmount := 25.0
		expectedNewAmount := 75.0

		balance := &models.Balance{
			WalletID:     walletID,
			StablecoinID: stablecoinID,
			Amount:       initialAmount,
		}

		updatedBalance := &models.Balance{
			WalletID:     walletID,
			StablecoinID: stablecoinID,
			Amount:       expectedNewAmount,
		}

		mockWalletRepo.
			EXPECT().
			GetBalance(walletID, stablecoinID).
			Return(balance, nil)

		mockWalletRepo.
			EXPECT().
			UpdateBalance(updatedBalance).
			Return(nil)

		service := NewWalletService(mockWalletRepo, mockStableCoinRepo, mockUserRepo)

		result, err := service.DepositFunds(walletID, stablecoinID, depositAmount)
		assert.NoError(t, err)
		assert.Equal(t, expectedNewAmount, result.Amount)
	})

	t.Run("fails_due_to_negative_or_zero_deposit", func(t *testing.T) {
		mockWalletRepo := mocks.NewMockWallet(ctrl)
		mockStableCoinRepo := mocks.NewMockStableCoin(ctrl)
		mockUserRepo := mocks.NewMockUser(ctrl)

		service := NewWalletService(mockWalletRepo, mockStableCoinRepo, mockUserRepo)

		// Test zero and negative deposit
		invalidAmounts := []float64{0, -10}

		for _, amount := range invalidAmounts {
			result, err := service.DepositFunds(1, 101, amount)
			assert.Nil(t, result)
			assert.EqualError(t, err, "deposit amount must be positive")
		}
	})
}

func TestWalletService_ListAllWallets(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successfully_lists_all_wallets", func(t *testing.T) {
		mockWalletRepo := mocks.NewMockWallet(ctrl)
		mockStableCoinRepo := mocks.NewMockStableCoin(ctrl)
		mockUserRepo := mocks.NewMockUser(ctrl)

		expectedWallets := []models.Wallet{
			{ID: 1, UserID: "user1", UserType: "email"},
			{ID: 2, UserID: "user2", UserType: "mobile"},
		}

		mockWalletRepo.
			EXPECT().
			ListAll().
			Return(expectedWallets, nil)

		service := NewWalletService(mockWalletRepo, mockStableCoinRepo, mockUserRepo)

		result, err := service.ListAllWallets()
		assert.NoError(t, err)
		assert.Equal(t, expectedWallets, result)
	})

	t.Run("fails_to_list_wallets_due_to_repo_error", func(t *testing.T) {
		mockWalletRepo := mocks.NewMockWallet(ctrl)
		mockStableCoinRepo := mocks.NewMockStableCoin(ctrl)
		mockUserRepo := mocks.NewMockUser(ctrl)

		mockWalletRepo.
			EXPECT().
			ListAll().
			Return(nil, errors.New("failed to fetch wallets"))

		service := NewWalletService(mockWalletRepo, mockStableCoinRepo, mockUserRepo)

		result, err := service.ListAllWallets()
		assert.Nil(t, result)
		assert.EqualError(t, err, "failed to fetch wallets")
	})
}
