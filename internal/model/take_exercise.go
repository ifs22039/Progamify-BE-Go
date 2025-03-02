package model

import "gorm.io/gorm"

type TakeExercise struct {
	gorm.Model
	ExerciseID uint `json:"exercise_id" gorm:"foreignKey:ExerciseID"`
	LessonID   uint `json:"lesson_id" gorm:"foreignKey:LessonID"`
	UserID     uint `json:"user_id" gorm:"foreignKey:UserID"`
	TopicID    uint `json:"topic_id"`
}
