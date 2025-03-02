package model

import "gorm.io/gorm"

type ExQuestion struct {
	gorm.Model
	Point      int  `json:"point" gorm:"not null"`
	Exp        int  `json:"exp" gorm:"not null"`
	ExerciseID uint `json:"-"`
}
