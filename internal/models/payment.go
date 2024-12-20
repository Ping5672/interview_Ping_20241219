package models

import (
	"time"
)

type PaymentStatus string
type PaymentMethod string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusSuccess   PaymentStatus = "success"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusCancelled PaymentStatus = "cancelled"

	PaymentMethodCreditCard PaymentMethod = "credit_card"
	PaymentMethodBank       PaymentMethod = "bank_transfer"
	PaymentMethodThirdParty PaymentMethod = "third_party"
	PaymentMethodBlockchain PaymentMethod = "blockchain"
)

type Payment struct {
	ID            uint          `gorm:"primaryKey" json:"id"`
	Amount        float64       `json:"amount"`
	Method        PaymentMethod `json:"method"`
	Status        PaymentStatus `json:"status"`
	TransactionID string        `json:"transaction_id"`
	PlayerID      uint          `json:"player_id"`
	Player        Player        `gorm:"foreignKey:PlayerID" json:"player"`
	Details       string        `json:"details"`
	ErrorMessage  string        `json:"error_message,omitempty"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}
