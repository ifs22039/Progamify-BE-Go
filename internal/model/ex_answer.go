package model

import "gorm.io/gorm"

type ExAnswer struct {
	gorm.Model
	ExQuestionID uint   `json:"-"`
	Content      string `json:"content" gorm:"not null"`
	IsCorrect    bool   `json:"is_correct" gorm:"not null"`
}
