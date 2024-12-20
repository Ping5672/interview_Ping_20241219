package services

import (
	"interview_Ping_20241219/internal/database"
	"interview_Ping_20241219/internal/models"
	"net/http"
	"strconv" // Add this import

	"github.com/gin-gonic/gin"
)

func RegisterLogRoutes(router *gin.Engine) {
	logs := router.Group("/logs")
	{
		logs.GET("", GetLogs)
		logs.POST("", CreateLog)
	}
}

// GetLogs handles GET /logs with query parameters
func GetLogs(c *gin.Context) {
	playerID := c.Query("player_id")
	action := c.Query("action")
	startTime := c.Query("start_time")
	endTime := c.Query("end_time")
	limitStr := c.DefaultQuery("limit", "50")

	// Convert limit string to integer
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid limit parameter",
			"details": err.Error(),
		})
		return
	}

	query := database.DB.Model(&models.GameLog{}).Preload("Player")

	// Apply filters
	if playerID != "" {
		query = query.Where("player_id = ?", playerID)
	}
	if action != "" {
		query = query.Where("action = ?", action)
	}
	if startTime != "" {
		query = query.Where("created_at >= ?", startTime)
	}
	if endTime != "" {
		query = query.Where("created_at <= ?", endTime)
	}

	// Apply limit and order
	query = query.Order("created_at DESC").Limit(limit)

	var logs []models.GameLog
	if err := query.Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch logs",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, logs)
}

// CreateLog handles POST /logs
func CreateLog(c *gin.Context) {
	var log models.GameLog
	if err := c.ShouldBindJSON(&log); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input",
			"details": err.Error(),
		})
		return
	}

	// Validate required fields
	if log.PlayerID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "player_id is required",
		})
		return
	}

	if log.Action == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "action is required",
		})
		return
	}

	// Validate action type
	validActions := map[models.LogActionType]bool{
		models.ActionRegister:      true,
		models.ActionLogin:         true,
		models.ActionLogout:        true,
		models.ActionEnterRoom:     true,
		models.ActionLeaveRoom:     true,
		models.ActionJoinChallenge: true,
		models.ActionChallengeEnd:  true,
	}

	if !validActions[log.Action] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid action type",
		})
		return
	}

	// Validate player exists
	var player models.Player
	if err := database.DB.First(&player, log.PlayerID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Player not found",
		})
		return
	}

	// Create log entry
	if err := database.DB.Create(&log).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create log",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, log)
}
