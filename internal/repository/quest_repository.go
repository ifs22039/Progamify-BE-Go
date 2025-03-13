package repository

import (
	"encoding/json"
	"errors"
	"math/rand"
	"strings"
	"time"

	"boysitorus/Progamify-Restful-API/internal/model"
	"boysitorus/Progamify-Restful-API/pkg/utils"

	"gorm.io/gorm"
)

type QuestRepository interface {
	GetQuestByUserID(userID uint) (*model.Quest, error)
	GetDifficultyByLevel(levelId uint) string
}

type questRepository struct {
	db *gorm.DB
	userRepo UserRepository
}

func (qr *questRepository) GetQuestByUserID(userID uint) (*model.Quest, error) {
	var user model.User
	if err := qr.db.First(&user, userID).Error; err != nil {
		return nil, errors.New("user not found")
	}

	var level model.Level
	if err := qr.db.First(&level, user.LevelId).Error; err != nil {
		return nil, errors.New("level not found")
	}

	difficulty := qr.GetDifficultyByLevel(uint(level.Level))

	var quest model.Quest
	if err := qr.db.Preload("Answers").Where("difficulty = ?", difficulty).Order("RAND()").First(&quest).Error; err != nil {
		return nil, errors.New("no suitable quest found")
	}

	return &quest, nil
}

func (qr *questRepository) GetDifficultyByLevel(levelId uint) string {
	rand.Seed(time.Now().UnixNano())
	probability := rand.Float64() 

	switch {
	case levelId <= 10: 
		return "Easy"
	case levelId <= 29: 
		if probability < 0.7 {
			return "Easy"
		}
		return "Medium"
	case levelId <= 49:
		if probability < 0.5 {
			return "Medium"
		} else if probability < 0.85 {
			return "Hard"
		}
		return "Easy"
	case levelId <= 79:
		if probability < 0.4 {
			return "Medium"
		} else if probability < 0.75 {
			return "Hard"
		}
		return "Very Hard"
	default: 
		if probability < 0.3 {
			return "Medium"
		} else if probability < 0.6 {
			return "Hard"
		}
		return "Very Hard"
	}
}

func NewQuestRepository(db *gorm.DB) QuestRepository {
	return &questRepository{db: db}
}

func (qr *questRepository) submitQuest(userID uint, request model.SubmitQuestRequest) (*model.TakeQuest, error) {

	var quest model.Quest

	err := qr.db.Preload("Answers").First(&quest, request.QuestID).Error

	if err != nil {
		return nil, err
	}

	answerJSON := request.Answer

	takeQuestAnswer := make(map[string]interface{})

	rewardExp := 0
	rewardPoint := 0
	is_correct := false

	//Grading
	
	detail := answerJSON.(map[string]interface{})

	var question model.Quest

	err = qr.db.Preload("Answers").First(&quest, detail["question_id"]).Error

	if err != nil {
		return nil, err
	}

	exp := 0
	point := 0
	rewardExp += question.Exp
	rewardPoint += question.Point

	if question.Type == "multiple_choice" {
		var correctAnswer model.QuestAnswer
		var correctAnswerIndex int
		for index, item := range question.Answers {
			if item.IsCorrect {
				correctAnswer = item
				correctAnswerIndex = index
			}
		}

		if int(correctAnswer.ID) == int(detail["answer_id"].(float64)) {
			exp = exp + question.Exp
			point = point + question.Point

			rewardExp += exp
			rewardPoint += point
			is_correct = true
		}

		takeQuestAnswer = map[string]interface{}{
			"question_id":          question.ID,
			"feedback":             question.Feedback,
			"exp_gained":           exp,
			"point_gained":         point,
			"user_answer_id":       detail["answer_id"],
			"user_answer_index":    detail["index_jawaban"],
			"correct_answer_id":    correctAnswer.ID,
			"correct_answer_index": correctAnswerIndex,
		}
	} else if question.Type == "true_false" {
		correctAnswer := question.Answers[0]
		correctAnswerIndex := 0
		jawabanUser := detail["answer_text"].(string)

		if strings.ToLower(correctAnswer.Content) == strings.ToLower(jawabanUser) {
			correctAnswerIndex = int(detail["index_jawaban"].(float64))
			exp = exp + question.Exp
			point = point + question.Point

			rewardExp += exp
			rewardPoint += point
			is_correct = true
		} else {
			if int(detail["index_jawaban"].(float64)) == 0 {
				correctAnswerIndex = 1
			} else {
				correctAnswerIndex = 0
			}
		}

		takeQuestAnswer = map[string]interface{}{
			"question_id":          question.ID,
			"feedback":             question.Feedback,
			"exp_gained":           exp,
			"point_gained":         point,
			"user_answer":          detail["answer_text"],
			"correct_answer":       correctAnswer.Content,
			"user_answer_index":    detail["index_jawaban"],
			"correct_answer_index": correctAnswerIndex,
		}
	} else if question.Type == "short_answer" {
		correctAnswer := question.Answers[0]
		jawabanUser := detail["index_jawaban"].(string)

		if strings.ToLower(correctAnswer.Content) == strings.ToLower(jawabanUser) {
			exp = exp + question.Exp
			point = point + question.Point

			rewardExp += exp
			rewardPoint += point
			is_correct = true
		}

		takeQuestAnswer = map[string]interface{}{
			"question_id":          question.ID,
			"feedback":             question.Feedback,
			"exp_gained":           exp,
			"point_gained":         point,
			"user_answer_index":    detail["index_jawaban"],
			"correct_answer_index": correctAnswer.Content,
		}
	} else if question.Type == "essay" {
		correctAnswer := question.Answers[0]
		jawabanUser := detail["index_jawaban"].(string)

		flag := true

		var similarity float64 = 0

		for flag {
			result, err := utils.EssayGrading(correctAnswer.Content, jawabanUser)
			if err == nil {
				flag = false
			}
			similarity = result
		}

		if similarity >= 50 {
			exp = exp + question.Exp
			point = point + question.Point

			rewardExp += exp
			rewardPoint += point
			is_correct = true
		}

		takeQuestAnswer = map[string]interface{}{
			"question_id":          question.ID,
			"feedback":             question.Feedback,
			"exp_gained":           exp,
			"point_gained":         point,
			"user_answer_index":    detail["index_jawaban"],
			"correct_answer_index": correctAnswer.Content,
		}
	} else if question.Type == "multiple_answer" {

	}

	var score float64
	if is_correct {
			score = 100
	} else {
			score = 0
	}

	answerDetail, err := json.Marshal(takeQuestAnswer)
	if err != nil {
		return nil, err
	}

	newTakeQuest := model.TakeQuest{
		QuestID:    	 quest.ID,
		UserID:        userID,
		Answers:       answerDetail,
		Score:         score,
		IsCorrect:  	 is_correct,
		RewardExp:     rewardExp,
		RewardPoint:   rewardPoint,
	}

	err = qr.userRepo.AddPoint(userID, rewardPoint)

	if err != nil {
		return nil, err
	}

	err = qr.userRepo.AddExp(userID, rewardExp)

	if err != nil {
		return nil, err
	}

	err = qr.userRepo.CheckLevel(userID)

	if err != nil {
		return nil, err
	}

	err =qr.db.Create(&newTakeQuest).Error
	if err != nil {
		return nil, err
	}

	return &newTakeQuest, err
}