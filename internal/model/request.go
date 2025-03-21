package model

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UpdateUserRequest struct {
	Name       string `json:"name"`
	Nim        string `json:"nim"`
	Angkatan   int    `json:"angkatan"`
	TotalExp   int    `json:"total_exp"`
	TotalPoint int    `json:"total_point"`
}

type RegisterRequest struct {
	Name                 string `json:"name" binding:"required"`
	Email                string `json:"email" binding:"required,email"`
	Nim                  string `json:"nim" binding:"required"`
	Angkatan             int    `json:"angkatan" binding:"required"`
	Password             string `json:"password" binding:"required,min=8"`
	PasswordConfirmation string `json:"password_confirmation" binding:"required,eqfield=Password"`
	AvatarID             uint    `json:"avatar_id" binding:"required"`
}

type SubmitExerciseRequest struct {
	ExerciseID uint                `json:"exercise_id" binding:"required"`
	Answers    map[int]interface{} `json:"answers" binding:"required"`
}

type SubmitQuestRequest struct {
	QuestID uint        `json:"quest_id" binding:"required"`
	Answer  interface{} `json:"answer" binding:"required"`
}

type AddBadgeRequest struct {
	BadgeID uint `json:"badge_id" binding:"required"`
}
