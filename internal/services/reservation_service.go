package services

import (
	"interview_Ping_20241219/internal/database"
	"interview_Ping_20241219/internal/models"
	"net/http"
	"strconv"
	"time" // Add this import for string to int conversion

	"github.com/gin-gonic/gin"
)

type ReservationRequest struct {
	RoomID    uint   `json:"room_id" binding:"required"`
	PlayerID  uint   `json:"player_id" binding:"required"`
	Date      string `json:"date" binding:"required"`       // Format: "2006-01-02"
	StartTime string `json:"start_time" binding:"required"` // Format: "15:04"
	EndTime   string `json:"end_time" binding:"required"`   // Format: "15:04"
}

func RegisterReservationRoutes(router *gin.Engine) {
	reservations := router.Group("/reservations")
	{
		reservations.GET("", ListReservations)
		reservations.POST("", CreateReservation)
	}
}

// ListReservations handles GET /reservations with optional query parameters
func ListReservations(c *gin.Context) {
	// Get query parameters
	roomID := c.Query("room_id")
	date := c.Query("date")
	limitStr := c.DefaultQuery("limit", "10")

	// Convert limit to integer
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid limit parameter",
			"details": err.Error(),
		})
		return
	}

	var reservations []models.Reservation
	query := database.DB.Model(&models.Reservation{}).Preload("Room").Preload("Player")

	// Apply filters
	if roomID != "" {
		query = query.Where("room_id = ?", roomID)
	}
	if date != "" {
		query = query.Where("date = ?", date)
	}

	// Apply limit
	query = query.Limit(limit)

	if err := query.Find(&reservations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch reservations",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, reservations)
}

// CreateReservation handles POST /reservations
func CreateReservation(c *gin.Context) {
	var req ReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input",
			"details": err.Error(),
		})
		return
	}

	// Parse date and times
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid date format. Use YYYY-MM-DD",
		})
		return
	}

	startTime, err := time.Parse("15:04", req.StartTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid start time format. Use HH:mm",
		})
		return
	}

	endTime, err := time.Parse("15:04", req.EndTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid end time format. Use HH:mm",
		})
		return
	}

	// Check if room exists
	var room models.Room
	if err := database.DB.First(&room, req.RoomID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Room not found",
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

	// Check for conflicting reservations
	var conflictCount int64
	database.DB.Model(&models.Reservation{}).
		Where("room_id = ? AND date = ? AND ((start_time <= ? AND end_time > ?) OR (start_time < ? AND end_time >= ?))",
			req.RoomID, date, endTime, startTime, endTime, startTime).
		Count(&conflictCount)

	if conflictCount > 0 {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Room is already reserved for this time period",
		})
		return
	}

	// Create reservation
	reservation := models.Reservation{
		RoomID:    req.RoomID,
		PlayerID:  req.PlayerID,
		Date:      date,
		StartTime: startTime,
		EndTime:   endTime,
	}

	if err := database.DB.Create(&reservation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create reservation",
			"details": err.Error(),
		})
		return
	}

	// Load related data
	database.DB.Preload("Room").Preload("Player").First(&reservation, reservation.ID)

	c.JSON(http.StatusCreated, reservation)
}
