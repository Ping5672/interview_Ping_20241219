package database

import (
	"fmt"
	"interview_Ping_20241219/internal/config"
	"interview_Ping_20241219/internal/models"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	dbConfig := config.GetDatabaseConfig()
	dsn := dbConfig.GetDSN()

	// Add retry logic
	for i := 0; i < 5; i++ {
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		fmt.Printf("Failed to connect to database. Retrying in 5 seconds... (Attempt %d/5)\n", i+1)
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database after 5 attempts: %v", err))
	}

	// Auto Migrate all models
	err = DB.AutoMigrate(
		&models.Player{},
		&models.Level{},
		&models.Room{},
		&models.Reservation{},
		&models.Challenge{},
		&models.ChallengePool{},
		&models.GameLog{},
		&models.Payment{},
	)
	if err != nil {
		panic(fmt.Sprintf("Failed to migrate database: %v", err))
	}

	// Initialize challenge pool if it doesn't exist
	var pool models.ChallengePool
	if err := DB.FirstOrCreate(&pool, models.ChallengePool{Amount: 0}).Error; err != nil {
		panic(fmt.Sprintf("Failed to initialize challenge pool: %v", err))
	}

	fmt.Println("Successfully connected to database")
}
