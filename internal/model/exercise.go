package model

import "gorm.io/gorm"

type Exercise struct {
	gorm.Model
	Title     string       `json:"title" gorm:"size:255;not null"`
	LessonID  uint         `json:"lesson_id"`
	Questions []ExQuestion `json:"questions" gorm:"foreignKey:ExerciseID"`

	Beta float64 `json:"beta" gorm:"column:beta;default:0"`
}
