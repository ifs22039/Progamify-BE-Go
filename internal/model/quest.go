package model

import "gorm.io/gorm"

type Quest struct {
	gorm.Model
	Title     	string       	`json:"title" gorm:"size:255;not null"`
	Content  		string 			 	`json:"content" gorm:"size:255;not null"`
	Point     	int        	 	`json:"point" gorm:"not null"`
	Exp       	int        	 	`json:"exp" gorm:"not null"`
	Difficulty 	string 				`json:"difficulty" gorm:"size:255;not null"`
	Type      	string      	`json:"type" gorm:"not null"`
	Feedback  	string      	`json:"feedback"`
	Answers 		[]QuestAnswer `json:"answers" gorm:"foreignKey:QuestID"`
}
