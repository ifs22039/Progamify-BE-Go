package model

import "gorm.io/gorm"

type Level struct {
	gorm.Model
	Level     int   `json:"level"`
	ExpNeeded int64 `json:"exp_needed"`
}
