package model

import "gorm.io/gorm"

type Badge struct {
	gorm.Model
	Title     	string   	`json:"title"`
	Description string 		`json:"description"`
	Picture			string 		`json:"picture"`
}

type UserBadge struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Picture     string `json:"picture"`
	Count       int    `json:"count"`
}
