package services

import (
	"interview_Ping_20241219/internal/database"
	"interview_Ping_20241219/internal/models"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	CHALLENGE_COST     = 20.01
	WIN_PROBABILITY    = 0.01 // 1%
	CHALLENGE_DURATION = 30   // seconds
	COOLDOWN_DURATION  = 60   // seconds
)

func RegisterChallengeRoutes(router *gin.Engine) {
	challenges := router.Group("/challenges")
	{
		challenges.POST("", JoinChallenge)
		challenges.GET("/results", GetChallengeResults)
	}
}

type JoinChallengeRequest struct {
	PlayerID uint    `json:"player_id" binding:"required"`
	Amount   float64 `json:"amount" binding:"required,eq=20.01"`
}

func JoinChallenge(c *gin.Context) {
	var req JoinChallengeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input",
			"details": err.Error(),
		})
		return
	}

	// Check if player exists
	var player models.Player
	if err := database.DB.First(&player, req.PlayerID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Player not found",
		})
		return
	}

	// Check cooldown period
	var lastChallenge models.Challenge
	database.DB.Where("player_id = ?", req.PlayerID).
		Order("created_at DESC").
		First(&lastChallenge)

	if !lastChallenge.CreatedAt.IsZero() &&
		time.Since(lastChallenge.CreatedAt).Seconds() < COOLDOWN_DURATION {
		c.JSON(http.StatusTooEarly, gin.H{
			"error":     "Please wait one minute between challenges",
			"wait_time": COOLDOWN_DURATION - time.Since(lastChallenge.CreatedAt).Seconds(),
		})
		return
	}

	// Start transaction
	tx := database.DB.Begin()

	// Update pool
	var pool models.ChallengePool
	if err := tx.First(&pool).Error; err != nil {
		pool = models.ChallengePool{Amount: 0}
	}
	pool.Amount += CHALLENGE_COST
	if err := tx.Save(&pool).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update pool"})
		return
	}

	// Determine if player wins
	isWinner := rand.Float64() < WIN_PROBABILITY

	challenge := models.Challenge{
		PlayerID:  req.PlayerID,
		Amount:    CHALLENGE_COST,
		IsWinner:  isWinner,
		StartTime: time.Now(),
		EndTime:   time.Now().Add(time.Second * CHALLENGE_DURATION),
	}

	if err := tx.Create(&challenge).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create challenge"})
		return
	}

	// If player wins, empty the pool
	if isWinner {
		challenge.Amount = pool.Amount
		pool.Amount = 0
		if err := tx.Save(&pool).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update pool"})
			return
		}
		if err := tx.Save(&challenge).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update challenge"})
			return
		}
	}

	tx.Commit()

	c.JSON(http.StatusCreated, gin.H{
		"challenge_id": challenge.ID,
		"is_winner":    challenge.IsWinner,
		"amount":       challenge.Amount,
		"pool_amount":  pool.Amount,
	})
}

func GetChallengeResults(c *gin.Context) {
	var challenges []models.Challenge
	if err := database.DB.Preload("Player").
		Order("created_at DESC").
		Limit(10).
		Find(&challenges).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch results"})
		return
	}

	var pool models.ChallengePool
	if err := database.DB.First(&pool).Error; err != nil {
		pool = models.ChallengePool{Amount: 0}
	}

	c.JSON(http.StatusOK, gin.H{
		"challenges":  challenges,
		"pool_amount": pool.Amount,
	})
}
