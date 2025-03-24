package model

import "gorm.io/gorm"

type Achievement struct {
	gorm.Model
	Title     	string   	`json:"title"`
	Description string 		`json:"description"`
	Picture			string 		`json:"picture"`
}

type UserAchievement struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Picture     string `json:"picture"`
	Count       int    `json:"count"`
}
