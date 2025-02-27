package model

import "gorm.io/gorm"

type TakeLesson struct {
	gorm.Model
	LessonID uint `json:"lesson_id" gorm:"foreignKey:LessonID"`
	UserID   uint `json:"user_id" gorm:"foreignKey:UserID"`
	TopicID  uint `json:"topic_id"`
}
