package model

import (
	"encoding/json"

	"gorm.io/gorm"
)

type TakeExercise struct {
	gorm.Model
	ExerciseID    uint            `json:"exercise_id" gorm:"foreignKey:ExerciseID"`
	LessonID      uint            `json:"lesson_id" gorm:"foreignKey:LessonID"`
	UserID        uint            `json:"user_id" gorm:"foreignKey:UserID"`
	TopicID       uint            `json:"topic_id" gorm:"foreignKey:TopicID"`
	AttemptNumber int             `json:"attempt_number"`
	Answers       json.RawMessage `json:"answers"`
	Score         float64         `json:"score"`
	TotalCorrect  int             `json:"total_correct"`
	TotalQuestion int             `json:"total_question"`
	TotalExp      int             `json:"total_exp"`
	TotalPoint    int             `json:"total_point"`
	RewardExp     int             `json:"reward_exp"`
	RewardPoint   int             `json:"reward_point"`
}
