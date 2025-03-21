package model

import "gorm.io/gorm"

type Gift struct {
	gorm.Model
	Title   string `json:"title"`
	Picture string `json:"picture_url"`
	Price   int64  `json:"price"`
}
