package repository

import (
	"boysitorus/Progamify-Restful-API/internal/model"
	"boysitorus/Progamify-Restful-API/pkg/utils"
	"encoding/json"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type ExerciseRepository interface {
	FindById(id uint) (*model.Exercise, error)
	AddTakeExercise(userID uint, request model.SubmitExerciseRequest) (*model.TakeExercise, error)
}

type exerciseRepository struct {
	db       *gorm.DB
	userRepo UserRepository
}

func NewExerciseRepository(db *gorm.DB, userRepo UserRepository) ExerciseRepository {
	return &exerciseRepository{db, userRepo}
}

func (e *exerciseRepository) AddTakeExercise(userID uint, request model.SubmitExerciseRequest) (*model.TakeExercise, error) {
	var exercise model.Exercise

	err := e.db.Preload("Questions.Answers").First(&exercise, request.ExerciseID).Error

	if err != nil {
		return nil, err
	}

	var lesson model.Lesson

	err = e.db.First(&lesson, exercise.LessonID).Error

	if err != nil {
		return nil, err
	}

	// Count previous attempts
	var attemptCount int64
	e.db.Model(&model.TakeExercise{}).Where("user_id = ? AND exercise_id = ?", userID, request.ExerciseID).Count(&attemptCount)
	attemptNumber := int(attemptCount) + 1

	answersJSON := request.Answers

	takeExerciseAnswer := make(map[int]interface{})
	totalExp := 0
	totalPoint := 0
	rewardExp := 0
	rewardPoint := 0
	var totalCorrect float64 = 0

	//Grading
	for key, value := range answersJSON {
		detail := value.(map[string]interface{})

		var question model.ExQuestion

		err = e.db.Preload("Answers").First(&question, detail["question_id"]).Error

		if err != nil {
			return nil, err
		}

		exp := 0
		point := 0
		rewardExp += question.Exp
		rewardPoint += question.Point

		if question.Type == "multiple_choice" {
			var correctAnswer model.ExAnswer
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

				totalExp += exp
				totalPoint += point
				totalCorrect += 1
			}

			takeExerciseAnswer[key] = map[string]interface{}{
				"question_id":          question.ID,
				"feedback":             question.Feedback,
				"exp_gained":           exp,
				"point_gained":         point,
				"user_answer_id":       detail["answer_id"],
				"user_answer_index":    detail["index_jawaban"],
				"correct_answer_id":    correctAnswer.ID,
				"correct_answer_index": correctAnswerIndex,
				"type":                 question.Type,
			}
		} else if question.Type == "true_false" {
			correctAnswer := question.Answers[0]
			correctAnswerIndex := 0
			var jawabanUser string
			
			if val, ok := detail["answer_text"].(string); ok {
				jawabanUser = val
			} else if val, ok := detail["index_jawaban"].(string); ok {
				jawabanUser = val
			}

			if strings.TrimSpace(strings.ToLower(correctAnswer.Content)) == strings.TrimSpace(strings.ToLower(jawabanUser)) {
				correctAnswerIndex = int(detail["index_jawaban"].(float64))
				exp = exp + question.Exp
				point = point + question.Point

				totalExp += exp
				totalPoint += point
				totalCorrect += 1
			} else {
				if int(detail["index_jawaban"].(float64)) == 0 {
					correctAnswerIndex = 1
				} else {
					correctAnswerIndex = 0
				}
			}

			takeExerciseAnswer[key] = map[string]interface{}{
				"question_id":          question.ID,
				"feedback":             question.Feedback,
				"exp_gained":           exp,
				"point_gained":         point,
				"user_answer":          jawabanUser,
				"correct_answer":       correctAnswer.Content,
				"user_answer_index":    detail["index_jawaban"],
				"correct_answer_index": correctAnswerIndex,
				"type":                 question.Type,
			}
		} else if question.Type == "short_answer" {
			correctAnswer := question.Answers[0]
			var jawabanUser string
			
			if val, ok := detail["answer_text"].(string); ok {
				jawabanUser = val
			} else if val, ok := detail["index_jawaban"].(string); ok {
				jawabanUser = val
			}

			if strings.TrimSpace(strings.ToLower(correctAnswer.Content)) == strings.TrimSpace(strings.ToLower(jawabanUser)) {
				exp = exp + question.Exp
				point = point + question.Point

				totalExp += exp
				totalPoint += point
				totalCorrect += 1
			}

			takeExerciseAnswer[key] = map[string]interface{}{
				"question_id":          question.ID,
				"feedback":             question.Feedback,
				"exp_gained":           exp,
				"point_gained":         point,
				"user_answer":          jawabanUser,
				"correct_answer":       correctAnswer.Content,
				"type":                 question.Type,
			}
		} else if question.Type == "essay" {
			fmt.Println("Ada soal essay nih")
			correctAnswer := question.Answers[0]
			var jawabanUser string
			
			if val, ok := detail["answer_text"].(string); ok {
				jawabanUser = val
			} else if val, ok := detail["index_jawaban"].(string); ok {
				jawabanUser = val
			}

			flag := true
			var similarity float64 = 0

			for flag {
				result, err := utils.EssayGrading(correctAnswer.Content, jawabanUser)
				fmt.Println(result)
				fmt.Println(err)
				if err == nil {
					flag = false
				}
				similarity = result
			}

			if similarity >= 50 {
				exp = exp + question.Exp
				point = point + question.Point

				totalExp += exp
				totalPoint += point
				totalCorrect += 1
			}

			takeExerciseAnswer[key] = map[string]interface{}{
				"question_id":          question.ID,
				"feedback":             question.Feedback,
				"exp_gained":           exp,
				"point_gained":         point,
				"user_answer":          jawabanUser,
				"correct_answer":       correctAnswer.Content,
				"similarity_score":     similarity,
				"type":                 question.Type,
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
				totalCorrect += 1
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
					totalCorrect += 1
				}
			}

			totalExp += expGained
			totalPoint += pointGained

			takeExerciseAnswer[key] = map[string]interface{}{
				"question_id":            question.ID,
				"feedback":               question.Feedback,
				"exp_gained":             expGained,
				"point_gained":           pointGained,
				"correct_answer_id":      correctAnswers,
				"user_answer_id":         userAnswers,
				"user_correct_answer_id": jawabanUserBenar,
				"correct_answer_index":   correctAnswersIndex,
				"user_answer_index":      detail["index_jawaban"],
				"type":                   question.Type,
			}
		} else if question.Type == "matching" {
			// Matching question type
			// Format: each answer contains keyword|explanation (separated by |)
			// is_correct = 1 indicates valid matching pair
			// Example: "Variable|Container untuk menyimpan data"
			
			var userMatchings []map[string]interface{}
			if matchings, ok := detail["matchings"].([]interface{}); ok {
				for _, m := range matchings {
					userMatchings = append(userMatchings, m.(map[string]interface{}))
				}
			}

			// Parse all answers to extract keywords and explanations
			var matchingPairs []map[string]interface{}
			var keywordsList []map[string]interface{}
			
			for _, answer := range question.Answers {
				if answer.IsCorrect {
					// Parse content: "keyword|explanation"
					parts := strings.Split(answer.Content, "|")
					keyword := strings.TrimSpace(parts[0])
					explanation := ""
					
					if len(parts) > 1 {
						explanation = strings.TrimSpace(parts[1])
					}
					
					pair := map[string]interface{}{
						"id":          answer.ID,
						"keyword":     keyword,
						"explanation": explanation,
					}
					
					matchingPairs = append(matchingPairs, pair)
					keywordsList = append(keywordsList, map[string]interface{}{
						"id":      answer.ID,
						"keyword": keyword,
					})
				}
			}

			// Count total matching pairs
			totalMatchPairs := len(matchingPairs)
			var userCorrectCount int = 0

			// Validate user matchings
			for _, userMatch := range userMatchings {
				var matchID int
				if mID, ok := userMatch["match_id"].(float64); ok {
					matchID = int(mID)
				} else if mID, ok := userMatch["match_id"].(int); ok {
					matchID = mID
				}
				
				// Check if this match ID exists in our valid pairs
				for _, pair := range matchingPairs {
					if pID, ok := pair["id"].(uint); ok && matchID == int(pID) {
						userCorrectCount++
						break
					}
				}
			}

			var expGained int = 0
			var pointGained int = 0

			// Calculate score based on correct matches
			if totalMatchPairs > 0 {
				if userCorrectCount == totalMatchPairs && len(userMatchings) == totalMatchPairs {
					// Perfect match - all pairs matched correctly
					expGained = question.Exp
					pointGained = question.Point
					totalCorrect += 1
				} else if userCorrectCount > 0 {
					// Partial match - give proportional score
					percentage := float64(userCorrectCount) / float64(totalMatchPairs)
					expGained = int(float64(question.Exp) * percentage)
					pointGained = int(float64(question.Point) * percentage)
					
					// Count as correct if more than 50% matched
					if userCorrectCount > totalMatchPairs/2 {
						totalCorrect += 1
					}
				}
			}

			totalExp += expGained
			totalPoint += pointGained

			takeExerciseAnswer[key] = map[string]interface{}{
				"question_id":       question.ID,
				"feedback":          question.Feedback,
				"exp_gained":        expGained,
				"point_gained":      pointGained,
				"user_matchings":    userMatchings,
				"correct_count":     userCorrectCount,
				"total_matches":     totalMatchPairs,
				"keywords":          keywordsList,
				"all_pairs":         matchingPairs,
				"user_answer_index": detail["index_jawaban"],
				"type":              question.Type,
			}
		}
	}

	totalQuestions := float64(len(answersJSON))
	// If answersJSON is empty for some reason, fallback to exercise pool size or a minimum of 1
	if totalQuestions == 0 {
		totalQuestions = float64(len(exercise.Questions))
	}
	if totalQuestions == 0 {
		totalQuestions = 1
	}
	score := totalCorrect / totalQuestions

	answerDetail, err := json.Marshal(takeExerciseAnswer)
	if err != nil {
		return nil, err
	}

	newTakeExercise := model.TakeExercise{
		ExerciseID:    exercise.ID,
		LessonID:      lesson.ID,
		UserID:        userID,
		TopicID:       lesson.TopicID,
		AttemptNumber: attemptNumber,
		Answers:       answerDetail,
		Score:         score,
		TotalQuestion: int(totalQuestions),
		TotalCorrect:  int(totalCorrect),
		TotalExp:      totalExp,
		TotalPoint:    totalPoint,
		RewardExp:     rewardExp,
		RewardPoint:   rewardPoint,
	}

	err = e.userRepo.AddPoint(userID, totalPoint)

	if err != nil {
		return nil, err
	}

	err = e.userRepo.AddExp(userID, totalExp)

	if err != nil {
		return nil, err
	}

	err = e.userRepo.CheckLevel(userID)

	if err != nil {
		return nil, err
	}

	err = e.db.Create(&newTakeExercise).Error
	if err != nil {
		return nil, err
	}

	return &newTakeExercise, err
}

func (e *exerciseRepository) FindById(id uint) (*model.Exercise, error) {
	var exercise model.Exercise

	err := e.db.
		Preload("Questions", func(db *gorm.DB) *gorm.DB {
			return db.Order("RAND()").Limit(5)
		}).
		Preload("Questions.Answers").
		First(&exercise, id).Error

	if err != nil {
		return nil, err
	}

	return &exercise, nil
}

