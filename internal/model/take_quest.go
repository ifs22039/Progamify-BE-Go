package model

import (
	"encoding/json"

	"gorm.io/gorm"
)

type TakeQuest struct {
	gorm.Model
	QuestID    		uint            `json:"quest_id" gorm:"foreignKey:QuestID"`
	UserID        uint            `json:"user_id" gorm:"foreignKey:UserID"`
	Answer       json.RawMessage `json:"answer"`
	IsCorrect			bool						`json:"is_correct"`
	Score         float64         `json:"score"`
	RewardExp     int             `json:"reward_exp"`
	RewardPoint   int             `json:"reward_point"`

	Probability float64 `json:"probability"`
	ThetaBefore float64 `json:"theta_before"`
	ThetaAfter  float64 `json:"theta_after"`


}
