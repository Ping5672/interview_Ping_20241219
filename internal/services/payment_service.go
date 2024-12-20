package services

import (
	"fmt"
	"interview_Ping_20241219/internal/database"
	"interview_Ping_20241219/internal/models"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func RegisterPaymentRoutes(router *gin.Engine) {
	payments := router.Group("/payments")
	{
		payments.POST("", ProcessPayment)
		payments.GET("/:id", GetPayment)
	}
}

type PaymentRequest struct {
	PlayerID uint                 `json:"player_id" binding:"required"`
	Amount   float64              `json:"amount" binding:"required,gt=0"`
	Method   models.PaymentMethod `json:"method" binding:"required"`
	Details  string               `json:"details"`
}

// Mock payment processing for different methods
func processPaymentByMethod(method models.PaymentMethod, amount float64) (string, error) {
	// Simulate processing time
	time.Sleep(time.Millisecond * 500)

	// Simulate random failure (10% chance)
	if rand.Float64() < 0.1 {
		return "", fmt.Errorf("payment failed: transaction declined")
	}

	// Generate mock transaction ID
	transactionID := fmt.Sprintf("%s_%d", method, time.Now().UnixNano())

	switch method {
	case models.PaymentMethodCreditCard:
		// Mock credit card processing
		return fmt.Sprintf("CC_%s", transactionID), nil
	case models.PaymentMethodBank:
		// Mock bank transfer processing
		return fmt.Sprintf("BT_%s", transactionID), nil
	case models.PaymentMethodThirdParty:
		// Mock third party payment processing
		return fmt.Sprintf("TP_%s", transactionID), nil
	case models.PaymentMethodBlockchain:
		// Mock blockchain payment processing
		return fmt.Sprintf("BC_%s", transactionID), nil
	default:
		return "", fmt.Errorf("unsupported payment method")
	}
}

func ProcessPayment(c *gin.Context) {
	var req PaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input",
			"details": err.Error(),
		})
		return
	}

	// Validate player exists
	var player models.Player
	if err := database.DB.First(&player, req.PlayerID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Player not found",
		})
		return
	}

	// Get payment processor
	processor := CreatePaymentProcessor(req.Method)
	if processor == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid payment method",
		})
		return
	}

	// Create initial payment record
	payment := models.Payment{
		PlayerID: req.PlayerID,
		Amount:   req.Amount,
		Method:   req.Method,
		Status:   models.PaymentStatusPending,
		Details:  req.Details,
	}

	// Start transaction
	tx := database.DB.Begin()

	if err := tx.Create(&payment).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create payment record",
			"details": err.Error(),
		})
		return
	}

	// Process payment
	transactionID, err := processor.Process(req.Amount)
	if err != nil {
		payment.Status = models.PaymentStatusFailed
		payment.ErrorMessage = err.Error()
		if err := tx.Save(&payment).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to update payment record",
				"details": err.Error(),
			})
			return
		}
		tx.Commit()
		c.JSON(http.StatusBadRequest, gin.H{
			"error":      "Payment failed",
			"details":    err.Error(),
			"payment_id": payment.ID,
		})
		return
	}

	// Update successful payment
	payment.Status = models.PaymentStatusSuccess
	payment.TransactionID = transactionID

	if err := tx.Save(&payment).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update payment record",
			"details": err.Error(),
		})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message":        "Payment processed successfully",
		"payment_id":     payment.ID,
		"transaction_id": transactionID,
		"status":         payment.Status,
		"amount":         payment.Amount,
		"method":         payment.Method,
	})
}

func GetPayment(c *gin.Context) {
	id := c.Param("id")
	var payment models.Payment

	if err := database.DB.Preload("Player").First(&payment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Payment not found",
		})
		return
	}

	c.JSON(http.StatusOK, payment)
}
