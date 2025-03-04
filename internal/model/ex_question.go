package model

import "gorm.io/gorm"

type ExQuestion struct {
	gorm.Model
	Point      int        `json:"point" gorm:"not null"`
	Exp        int        `json:"exp" gorm:"not null"`
	ExerciseID uint       `json:"-"`
	Content    string     `json:"content" gorm:"not null"`
	Type       string     `json:"type" gorm:"not null"`
	Feedback   string     `json:"feedback"`
	Answers    []ExAnswer `json:"answers" gorm:"foreignKey:ExQuestionID"`
}
