package models

import (
	"time"
)

type LogActionType string

const (
	ActionRegister      LogActionType = "註冊"
	ActionLogin         LogActionType = "登入"
	ActionLogout        LogActionType = "登出"
	ActionEnterRoom     LogActionType = "進入房間"
	ActionLeaveRoom     LogActionType = "退出房間"
	ActionJoinChallenge LogActionType = "參加挑戰"
	ActionChallengeEnd  LogActionType = "挑戰結果"
)

type GameLog struct {
	ID        uint          `gorm:"primaryKey" json:"id"`
	PlayerID  uint          `json:"player_id"`
	Player    Player        `gorm:"foreignKey:PlayerID" json:"player"`
	Action    LogActionType `json:"action"`
	Details   string        `json:"details"`
	CreatedAt time.Time     `json:"created_at"`
}
