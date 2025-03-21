package model

import "gorm.io/gorm"

type HaveGift struct {
	gorm.Model
	UserID   uint `json:"user_id"`
	GiftID   uint `json:"gift_id"`
	IsActive bool `json:"is_active"`
}
