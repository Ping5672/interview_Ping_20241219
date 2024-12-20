package models

import (
    "time"
)

type Player struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    Name      string    `gorm:"uniqueIndex;not null" json:"name"`
    Level     uint      `json:"level"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}