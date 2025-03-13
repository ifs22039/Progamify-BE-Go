package model

import "gorm.io/gorm"

type QuestAnswer struct {
	gorm.Model
	QuestID 				uint   `json:"-"`
	Content      		string `json:"content" gorm:"not null"`
}
