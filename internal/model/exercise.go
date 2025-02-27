package model

import "gorm.io/gorm"

type Exercise struct {
	gorm.Model
	Title     string       `json:"title" gorm:"size:255;not null"`
	LessonID  uint         `json:"-"`
	Questions []ExQuestion `json:"questions" gorm:"foreignKey:ExerciseID"`
}
