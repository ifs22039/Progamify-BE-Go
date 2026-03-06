package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"boysitorus/Progamify-Restful-API/internal/model"
	"boysitorus/Progamify-Restful-API/pkg/utils"

	"gorm.io/gorm"
)

type QuestRepository interface {
	GetQuestByUserID(userID uint) (*model.Quest, error)
	AddTakeQuest(userID uint, request model.SubmitQuestRequest) (*model.TakeQuest, error)
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

	var quest model.Quest
	if err := qr.db.Preload("Answers").Order("RAND()").First(&quest).Error; err != nil {
		return nil, errors.New("no suitable quest found")
	}

	return &quest, nil
}

func NewQuestRepository(db *gorm.DB, userRepo UserRepository) QuestRepository {
	return &questRepository{db, userRepo}
}

func (qr *questRepository) AddTakeQuest(userID uint, request model.SubmitQuestRequest) (*model.TakeQuest, error) {

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

	err = qr.db.Preload("Answers").First(&question, detail["question_id"]).Error

	if err != nil {
		return nil, err
	}

	exp := 0
	point := 0
	// rewardExp += question.Exp
	// rewardPoint += question.Point

	if question.Type == "multiple_choice" {
		fmt.Println("DEBUG: Iterating over multiple choice")
		var correctAnswer model.QuestAnswer
		var correctAnswerIndex int
		for index, item := range question.Answers {
			if item.IsCorrect {
				correctAnswer = item
				correctAnswerIndex = index
			}
		}

		// fmt.Printf("Debug: detail[\"answer_id\"] value: %v, type: %T\n", detail["answer_id"], detail["answer_id"])


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
		jawabanUser := detail["answer_text"].(string)

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
			"user_answer":          detail["answer_text"],
			"correct_answer":       correctAnswer.Content,
		}
	} else if question.Type == "essay" {
		fmt.Println("masuk ke essay kita")
		correctAnswer := question.Answers[0]
		jawabanUser := detail["answer_text"].(string)

		log.Println("Correct Answer:", correctAnswer.Content) 
    log.Println("User Answer:", jawabanUser)     

		var similarity float64 = 0

		// grade the essay once; result contains multiple metrics
		gradeResult, err := utils.EssayGrading(correctAnswer.Content, jawabanUser)
		if err != nil {
			log.Printf("Essay grading failed for question %d: %v", question.ID, err)
		} else {
			// use final score percentage for correctness check
			similarity = gradeResult.FinalScore * 100
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
			"correct_answer_index": 0,
			"user_answer":          detail["answer_text"],
			"correct_answer":       correctAnswer.Content,
		}
	} else if question.Type == "multiple_answer" {
		point := question.Point
			exp := question.Exp

			var userAnswers []int
			for _, val := range detail["answers"].([]interface{}) {
				answer := val.(map[string]interface{})
				if answerID, ok := answer["answer_id"].(float64); ok {
					userAnswers = append(userAnswers, int(answerID))
				}
			}

			var jawabanUserBenar []int

			var correctAnswers []int
			var correctAnswersIndex []int
			for index, item := range question.Answers {
				if item.IsCorrect {
					correctAnswers = append(correctAnswers, int(item.ID))
					correctAnswersIndex = append(correctAnswersIndex, index)
					for _, ans := range userAnswers {
						if ans == int(item.ID) {
							jawabanUserBenar = append(jawabanUserBenar, int(item.ID))
						}
					}
				}
			}

			var countJawabanBenar = len(jawabanUserBenar)
			var countJawabanSalah = len(userAnswers) - len(jawabanUserBenar)

			var expGained int = 0
			var pointGained int = 0

			if countJawabanBenar == len(correctAnswers) && len(correctAnswers) == len(userAnswers) {
				expGained = exp
				pointGained = point
				is_correct = true
			} else {
				expEachAns := exp / len(correctAnswers)
				pointEachAns := point / len(correctAnswers)
				if countJawabanBenar == countJawabanSalah || countJawabanSalah > countJawabanBenar {
					fmt.Println("Condition 1")
					expGained = 0
					pointGained = 0
				} else if countJawabanBenar > countJawabanSalah {
					fmt.Println("Condition 2")
					expGained = (expEachAns * countJawabanBenar) - (expEachAns * countJawabanSalah)
					pointGained = (pointEachAns * countJawabanBenar) - (pointEachAns * countJawabanSalah)
					is_correct = true
				}
			}

			expGained += expGained
			expGained += pointGained

			takeQuestAnswer = map[string]interface{}{
				"question_id":            question.ID,
				"feedback":               question.Feedback,
				"exp_gained":             expGained,
				"point_gained":           pointGained,
				"user_answer_id":         detail["answer_id"],
				"user_correct_answer_id": jawabanUserBenar,
				"correct_answer_index":   0,
				"correct_answers_index" : correctAnswersIndex,
				"user_answer_index":      detail["index_jawaban"],
				"type":                   question.Type,
				"user_answer":          detail["answer_text"],
			}
	} else {
		fmt.Println("DEBUG: Question type didn't detect")
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
		Answer:       answerDetail,
		Score:         score,
		IsCorrect:  	 is_correct,
		RewardExp:     rewardExp,
		RewardPoint:   rewardPoint,
	}

	// fmt.Printf("Debug:\n QuestID: %v,\n UserID: %v,\n Answers: %v,\n Score:%v,\n IsCorrect:%v,\n RewardExp:%v,\n RewardPoint:%v,\n", quest.ID, userID, answerDetail, score, is_correct, rewardExp, rewardPoint)
	// fmt.Printf("takeQuestAnswer: %+v\n", takeQuestAnswer)
	if qr.userRepo == nil {
    fmt.Printf("userRepo nil")
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