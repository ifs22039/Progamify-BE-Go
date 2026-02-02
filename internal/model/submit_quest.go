package model

type SubmitQuest struct {
	QuestID     uint        `json:"quest_id"`
	Answer      any         `json:"answer"`
	IsCorrect   bool        `json:"is_correct"`
	Score       float64     `json:"score"`
	RewardExp   int         `json:"reward_exp"`
	RewardPoint int         `json:"reward_point"`
}
