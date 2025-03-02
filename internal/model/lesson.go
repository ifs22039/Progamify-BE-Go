package model

import "gorm.io/gorm"

type Lesson struct {
	gorm.Model
	TopicID   uint       `json:"-"`
	Name      string     `json:"name,omitempty" gorm:"size:255;not null"`
	Content   string     `json:"content,omitempty" gorm:"not null"`
	Exp       int        `json:"exp,omitempty"`
	Exercises []Exercise `json:"exercises,omitempty" gorm:"foreignKey:LessonID"`
}
