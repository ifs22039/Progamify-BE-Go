package model

import "gorm.io/gorm"

type HaveAvatar struct {
	gorm.Model
	UserID   uint `json:"user_id"`
	AvatarID uint `json:"avatar_id"`
}
