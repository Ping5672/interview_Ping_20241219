package models

import (
    "time"
)

type Room struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    Name        string    `gorm:"uniqueIndex;not null" json:"name"`
    Description string    `json:"description"`
    Status      string    `gorm:"default:'available'" json:"status"` // available, occupied, maintenance
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type Reservation struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    RoomID    uint      `json:"room_id"`
    Room      Room      `gorm:"foreignKey:RoomID" json:"room"`
    PlayerID  uint      `json:"player_id"`
    Player    Player    `gorm:"foreignKey:PlayerID" json:"player"`
    Date      time.Time `json:"date"`
    StartTime time.Time `json:"start_time"`
    EndTime   time.Time `json:"end_time"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}