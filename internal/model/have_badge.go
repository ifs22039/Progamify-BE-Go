package model

import "gorm.io/gorm"

type HaveBadge struct {
	gorm.Model
	BadgeID   uint     `json:"badge_id"`
	UserID    uint    `json:"user_id"`
	Badge   	Badge 	`gorm:"foreignKey:BadgeID"`
}
