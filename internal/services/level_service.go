package services

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "interview_Ping_20241219/internal/database"
    "interview_Ping_20241219/internal/models"
)

// RegisterLevelRoutes registers all level-related routes
func RegisterLevelRoutes(router *gin.Engine) {
    levels := router.Group("/levels")
    {
        levels.GET("", ListLevels)
        levels.POST("", CreateLevel)
    }
}

// ListLevels godoc
// @Summary List all levels
// @Description Get a list of all available levels
// @Tags levels
// @Accept json
// @Produce json
// @Success 200 {array} models.Level
// @Failure 500 {object} ErrorResponse
// @Router /levels [get]
func ListLevels(c *gin.Context) {
    var levels []models.Level
    result := database.DB.Find(&levels)
    
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to fetch levels",
            "details": result.Error.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, levels)
}

// CreateLevel godoc
// @Summary Create a new level
// @Description Create a new level with the provided name
// @Tags levels
// @Accept json
// @Produce json
// @Param level body LevelRequest true "Level information"
// @Success 201 {object} models.Level
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /levels [post]
func CreateLevel(c *gin.Context) {
    var level models.Level
    
    // Validate input
    if err := c.ShouldBindJSON(&level); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid input",
            "details": err.Error(),
        })
        return
    }

    // Validate required fields
    if level.Name == "" {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Name is required",
        })
        return
    }

    // Create level in database
    result := database.DB.Create(&level)
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to create level",
            "details": result.Error.Error(),
        })
        return
    }

    // Return created level with 201 status
    c.JSON(http.StatusCreated, level)
}

type ErrorResponse struct {
    Error   string `json:"error"`
    Details string `json:"details,omitempty"`
}

type LevelRequest struct {
    Name string `json:"name" binding:"required"`
}