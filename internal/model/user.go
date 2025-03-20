package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name       string `json:"name" gorm:"size:255;not null"`
	Email      string `json:"email" gorm:"size:255;unique;not null"`
	Password   string `json:"-" gorm:"not null"`
	Nim        string `json:"nim" gorm:"size:20;unique;not null"`
	Angkatan   int    `json:"angkatan" gorm:"not null"`
	TotalPoint int    `json:"total_point" gorm:"default:0"`
	TotalExp   int    `json:"total_exp" gorm:"default:0"`
	LevelId    uint   `json:"level_id" gorm:"not null"`
	AvatarID   uint   `json:"avatar_id"`
	Avatar     Avatar `json:"avatar" gorm:"foreignKey:AvatarID;"`
}
