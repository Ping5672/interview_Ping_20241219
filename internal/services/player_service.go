package services

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "interview_Ping_20241219/internal/database"
    "interview_Ping_20241219/internal/models"
)

// PlayerRequest represents the request body for creating/updating a player
type PlayerRequest struct {
    Name  string `json:"name" binding:"required"`
    Level uint   `json:"level" binding:"required"`
}

// RegisterPlayerRoutes registers all player-related routes
// @Summary Register player routes
// @Description Sets up all player-related API endpoints

func RegisterPlayerRoutes(router *gin.Engine) {
    players := router.Group("/players")
    {
        players.GET("", ListPlayers)
        players.POST("", CreatePlayer)
        players.GET("/:id", GetPlayer)
        players.PUT("/:id", UpdatePlayer)
        players.DELETE("/:id", DeletePlayer)
    }
}

// ListPlayers godoc
// @Summary List all players
// @Description Get a list of all registered players
// @Tags players
// @Accept json
// @Produce json
// @Success 200 {array} models.Player
// @Failure 500 {object} ErrorResponse
// @Router /players [get]
func ListPlayers(c *gin.Context) {
    var players []models.Player
    result := database.DB.Find(&players)
    
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to fetch players",
            "details": result.Error.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, players)
}

// CreatePlayer godoc
// @Summary Create a new player
// @Description Create a new player with the provided information
// @Tags players
// @Accept json
// @Produce json
// @Param player body PlayerRequest true "Player information"
// @Success 201 {object} models.Player
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /players [post]
func CreatePlayer(c *gin.Context) {
    var req PlayerRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid input",
            "details": err.Error(),
        })
        return
    }

    player := models.Player{
        Name:  req.Name,
        Level: req.Level,
    }

    if err := database.DB.Create(&player).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to create player",
            "details": err.Error(),
        })
        return
    }

    c.JSON(http.StatusCreated, player)
}

func GetPlayer(c *gin.Context) {
    id := c.Param("id")
    var player models.Player
    
    if err := database.DB.First(&player, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Player not found"})
        return
    }
    
    c.JSON(http.StatusOK, player)
}

func UpdatePlayer(c *gin.Context) {
    id := c.Param("id")
    var player models.Player
    
    if err := database.DB.First(&player, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Player not found"})
        return
    }
    
    if err := c.ShouldBindJSON(&player); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    database.DB.Save(&player)
    c.JSON(http.StatusOK, player)
}

func DeletePlayer(c *gin.Context) {
    id := c.Param("id")
    if err := database.DB.Delete(&models.Player{}, id).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"message": "Player deleted successfully"})
}