package model

import "gorm.io/gorm"

type Lesson struct {
	gorm.Model
	TopicID   uint       `json:"-"`
	Name      string     `json:"name" gorm:"size:255;not null"`
	Exercises []Exercise `json:"exercises,omitempty" gorm:"foreignKey:LessonID"`
	Content   string     `json:"content" gorm:"not null"`
}
