package model

import "gorm.io/gorm"

type Quest struct {
	gorm.Model
	Title      string        `json:"title" gorm:"size:255;not null"`
	Content    string        `json:"content" gorm:"size:255;not null"`
	Timer 		 int 					 `json:"timer" gorm:"not null;"`
	Point      int           `json:"point" gorm:"not null"`
	Exp        int           `json:"exp" gorm:"not null"`
	DifficultyLabel string  `json:"difficulty_label" gorm:"size:50;not null"`
	DifficultyIRT   float64 `json:"difficulty_irt" gorm:"default:0"`
	Type       string        `json:"type" gorm:"not null"`
	Feedback   string        `json:"feedback"`
	Answers    []QuestAnswer `json:"answers" gorm:"foreignKey:QuestID"`
	Beta  float64 					 `gorm:"default:0"`
}
