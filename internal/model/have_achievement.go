package model

import "gorm.io/gorm"

type HaveAchievement struct {
	gorm.Model
	AchievementID   uint     `json:"badge_id"`
	UserID    uint    `json:"user_id"`
	Achievement   	Achievement 	`gorm:"foreignKey:AchievementID"`
}
