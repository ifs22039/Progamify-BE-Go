package model

import (
	"encoding/json"

	"gorm.io/gorm"
)

type TakeExercise struct {
	gorm.Model

	ExerciseID    uint            `json:"exercise_id"`
	LessonID      uint            `json:"lesson_id"`
	UserID        uint            `json:"user_id"`
	TopicID       uint            `json:"topic_id"`
	AttemptNumber int             `json:"attempt_number"`

	Answers       json.RawMessage `json:"answers"`
	Score         float64         `json:"score"`
	TotalCorrect  int             `json:"total_correct"`
	TotalQuestion int             `json:"total_question"`
	TotalExp      int             `json:"total_exp"`
	TotalPoint    int             `json:"total_point"`
	RewardExp     int             `json:"reward_exp"`
	RewardPoint   int             `json:"reward_point"`

	ThetaBefore float64 `json:"theta_before" gorm:"column:theta_before"`
	ThetaAfter  float64 `json:"theta_after" gorm:"column:theta_after"`
	BetaBefore  float64 `json:"beta_before" gorm:"column:beta_before"`
	BetaAfter   float64 `json:"beta_after" gorm:"column:beta_after"`
}
