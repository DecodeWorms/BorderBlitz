package models

// CreateWalletRequest is the request body for creating a wallet
type CreateWalletRequest struct {
	UserType string `json:"user_type" binding:"required,oneof=email mobile"`
	Email    string `json:"email" binding:"required"`
}

// DepositRequest is the request body for depositing funds
type DepositRequest struct {
	StablecoinID uint    `json:"stablecoin_id" binding:"required"`
	Amount       float64 `json:"amount" binding:"required,gt=0"`
}

// TransferRequest is the request body for transferring funds
type TransferRequest struct {
	SenderWalletID   uint    `json:"sender_wallet_id" binding:"required"`
	ReceiverWalletID uint    `json:"receiver_wallet_id" binding:"required"`
	SourceCoinID     uint    `json:"source_coin_id" binding:"required"`
	DestCoinID       uint    `json:"dest_coin_id" binding:"required"`
	Amount           float64 `json:"amount" binding:"required,gt=0"`
}

// DepositTransactionRequest is the request body for depositing funds
type DepositTransactionRequest struct {
	WalletID     uint    `json:"wallet_id" binding:"required"`
	StablecoinID uint    `json:"stablecoin_id" binding:"required"`
	Amount       float64 `json:"amount" binding:"required,gt=0"`
}

// SwapRequest is the request body for swapping stablecoins
type SwapRequest struct {
	WalletID     uint    `json:"wallet_id" binding:"required"`
	SourceCoinID uint    `json:"source_coin_id" binding:"required"`
	DestCoinID   uint    `json:"dest_coin_id" binding:"required"`
	Amount       float64 `json:"amount" binding:"required,gt=0"`
}

// GetExchangeRateRequest is the request body for getting exchange rates
type GetExchangeRateRequest struct {
	SourceCoinID uint `json:"source_coin_id" binding:"required"`
	DestCoinID   uint `json:"dest_coin_id" binding:"required"`
}

// CreateUserRequest is the request body for creating request
type CreateUserRequest struct {
	UserType string `json:"user_type" binding:"required,oneof=email mobile"`
}
