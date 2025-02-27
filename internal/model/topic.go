package model

import "gorm.io/gorm"

type Topic struct {
	gorm.Model
	Name    string   `json:"name"`
	Lessons []Lesson `json:"lessons" gorm:"foreignKey:TopicID"`
}
