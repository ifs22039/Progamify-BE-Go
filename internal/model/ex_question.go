package model

import "gorm.io/gorm"

type ExQuestion struct {
	gorm.Model
	ExerciseID uint `json:"-"`
}
