package models

import "gorm.io/gorm"

// Wallet represents a user wallet
type Wallet struct {
	gorm.Model
	Address string  `json:"address" gorm:"uniqueIndex"`
	Balance float64 `json:"balance"`
}

// Transaction represents a transaction record
type Transaction struct {
	gorm.Model
	WalletID        uint    `json:"wallet_id"`
	Amount          float64 `json:"amount"`
	TransactionType string  `json:"transaction_type"` // e.g., "credit" or "debit"
}
