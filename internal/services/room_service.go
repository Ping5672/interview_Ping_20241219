package services

import (
	"interview_Ping_20241219/internal/database"
	"interview_Ping_20241219/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RoomRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

func RegisterRoomRoutes(router *gin.Engine) {
	rooms := router.Group("/rooms")
	{
		rooms.GET("", ListRooms)
		rooms.POST("", CreateRoom)
		rooms.GET("/:id", GetRoom)
		rooms.PUT("/:id", UpdateRoom)
		rooms.DELETE("/:id", DeleteRoom)
	}
}

// ListRooms godoc
// @Summary List all game rooms
// @Description Get a list of all available game rooms
// @Tags rooms
// @Accept json
// @Produce json
// @Success 200 {array} models.Room
// @Failure 500 {object} ErrorResponse
// @Router /rooms [get]
func ListRooms(c *gin.Context) {
	var rooms []models.Room
	if err := database.DB.Find(&rooms).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch rooms",
			"details": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, rooms)
}

// CreateRoom godoc
// @Summary Create a new game room
// @Description Create a new game room with the provided information
// @Tags rooms
// @Accept json
// @Produce json
// @Param room body RoomRequest true "Room information"
// @Success 201 {object} models.Room
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /rooms [post]
func CreateRoom(c *gin.Context) {
	var req RoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input",
			"details": err.Error(),
		})
		return
	}

	// Check for existing room with same name
	var existingRoom models.Room
	if err := database.DB.Where("name = ?", req.Name).First(&existingRoom).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Room with this name already exists",
		})
		return
	}

	room := models.Room{
		Name:        req.Name,
		Description: req.Description,
		Status:      "available",
	}

	if err := database.DB.Create(&room).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create room",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, room)
}

// GetRoom godoc
// @Summary Get a specific game room
// @Description Get detailed information about a specific game room
// @Tags rooms
// @Accept json
// @Produce json
// @Param id path int true "Room ID"
// @Success 200 {object} models.Room
// @Failure 404 {object} ErrorResponse
// @Router /rooms/{id} [get]
func GetRoom(c *gin.Context) {
	id := c.Param("id")
	var room models.Room

	if err := database.DB.First(&room, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Room not found",
		})
		return
	}

	c.JSON(http.StatusOK, room)
}

// UpdateRoom godoc
// @Summary Update a game room
// @Description Update the information of a specific game room
// @Tags rooms
// @Accept json
// @Produce json
// @Param id path int true "Room ID"
// @Param room body RoomRequest true "Room information"
// @Success 200 {object} models.Room
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /rooms/{id} [put]
func UpdateRoom(c *gin.Context) {
	id := c.Param("id")
	var room models.Room

	if err := database.DB.First(&room, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Room not found",
		})
		return
	}

	var req RoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input",
			"details": err.Error(),
		})
		return
	}

	room.Name = req.Name
	room.Description = req.Description

	if err := database.DB.Save(&room).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update room",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, room)
}

// DeleteRoom godoc
// @Summary Delete a game room
// @Description Delete a specific game room
// @Tags rooms
// @Accept json
// @Produce json
// @Param id path int true "Room ID"
// @Success 200 {object} MessageResponse
// @Failure 500 {object} ErrorResponse
// @Router /rooms/{id} [delete]
func DeleteRoom(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&models.Room{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete room",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Room deleted successfully",
	})
}
