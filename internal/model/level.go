package model

import "gorm.io/gorm"

type Level struct {
	gorm.Model
	Level      int `json:"level" gorm:"not null"`
	LevelId    uint `json:"-"`
}
