package tests

import (
	"interview_Ping_20241219/internal/api"
	"interview_Ping_20241219/internal/config"
	"interview_Ping_20241219/internal/database"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupTestEnvironment(t *testing.T) *gin.Engine {
	// Set test environment flag
	config.IsTestEnvironment = true

	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Initialize test database
	database.InitDB()

	// Clean up database
	cleanupDatabase()

	// Create new server
	server := api.NewServer()
	return server.Router()
}

func cleanupDatabase() {
	// Delete in correct order to respect foreign key constraints
	db := database.DB

	// Delete in correct order to respect foreign key constraints
	// First, delete all dependent tables
	db.Exec("DELETE FROM payments")   // Delete payments first
	db.Exec("DELETE FROM game_logs")  // Then logs
	db.Exec("DELETE FROM challenges") // Then challenges
	db.Exec("DELETE FROM challenge_pools")
	db.Exec("DELETE FROM reservations") // Then reservations
	db.Exec("DELETE FROM rooms")        // Then rooms
	db.Exec("DELETE FROM players")      // Then players
	db.Exec("DELETE FROM levels")       // Finally levels
}
