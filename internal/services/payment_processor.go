package services

import (
	"fmt"
	"interview_Ping_20241219/internal/models"
	"math/rand"
	"time"
)

// PaymentProcessor interface
type PaymentProcessor interface {
	Process(amount float64) (string, error)
}

// CreditCardProcessor implements credit card payment processing
type CreditCardProcessor struct{}

func (p *CreditCardProcessor) Process(amount float64) (string, error) {
	// Simulate credit card gateway call
	time.Sleep(time.Millisecond * 800)

	if rand.Float64() < 0.1 { // 10% chance of failure
		return "", fmt.Errorf("credit card payment failed: insufficient funds")
	}

	// Generate mock transaction ID with CC prefix
	return fmt.Sprintf("CC_%s_%d", generateCardToken(), time.Now().UnixNano()), nil
}

func generateCardToken() string {
	return fmt.Sprintf("CARD_%d", rand.Intn(10000))
}

// BankTransferProcessor implements bank transfer processing
type BankTransferProcessor struct{}

func (p *BankTransferProcessor) Process(amount float64) (string, error) {
	// Simulate bank API call
	time.Sleep(time.Second * 1)

	if rand.Float64() < 0.05 { // 5% chance of failure
		return "", fmt.Errorf("bank transfer failed: invalid bank account")
	}

	// Generate mock bank transaction ID
	return fmt.Sprintf("BT_%s_%d", generateBankToken(), time.Now().UnixNano()), nil
}

func generateBankToken() string {
	return fmt.Sprintf("BANK_%d", rand.Intn(10000))
}

// ThirdPartyProcessor implements third-party payment processing
type ThirdPartyProcessor struct{}

func (p *ThirdPartyProcessor) Process(amount float64) (string, error) {
	// Simulate third-party API call
	time.Sleep(time.Millisecond * 600)

	if rand.Float64() < 0.08 { // 8% chance of failure
		return "", fmt.Errorf("third-party payment failed: service unavailable")
	}

	// Generate mock third-party transaction ID
	return fmt.Sprintf("TP_%s_%d", generateTPToken(), time.Now().UnixNano()), nil
}

func generateTPToken() string {
	return fmt.Sprintf("3RDPARTY_%d", rand.Intn(10000))
}

// BlockchainProcessor implements blockchain payment processing
type BlockchainProcessor struct{}

func (p *BlockchainProcessor) Process(amount float64) (string, error) {
	// Simulate blockchain transaction
	time.Sleep(time.Second * 2) // Blockchain transactions typically take longer

	if rand.Float64() < 0.15 { // 15% chance of failure
		return "", fmt.Errorf("blockchain payment failed: network congestion")
	}

	// Generate mock blockchain transaction hash
	return fmt.Sprintf("BC_%s_%d", generateBlockchainHash(), time.Now().UnixNano()), nil
}

func generateBlockchainHash() string {
	const charset = "abcdef0123456789"
	hash := make([]byte, 32)
	for i := range hash {
		hash[i] = charset[rand.Intn(len(charset))]
	}
	return string(hash)
}

// PaymentFactory creates the appropriate payment processor
func CreatePaymentProcessor(method models.PaymentMethod) PaymentProcessor {
	switch method {
	case models.PaymentMethodCreditCard:
		return &CreditCardProcessor{}
	case models.PaymentMethodBank:
		return &BankTransferProcessor{}
	case models.PaymentMethodThirdParty:
		return &ThirdPartyProcessor{}
	case models.PaymentMethodBlockchain:
		return &BlockchainProcessor{}
	default:
		return nil
	}
}
