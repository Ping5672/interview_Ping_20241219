package models

import (
	"time"
)

type Challenge struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	PlayerID  uint      `json:"player_id"`
	Player    Player    `gorm:"foreignKey:PlayerID" json:"player"`
	Amount    float64   `json:"amount"`
	IsWinner  bool      `json:"is_winner"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ChallengePool struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Amount    float64   `json:"amount" gorm:"default:0"`
	UpdatedAt time.Time `json:"updated_at"`
}
