package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strings"

	"boysitorus/Progamify-Restful-API/internal/model"
	"boysitorus/Progamify-Restful-API/pkg/utils"

	"gorm.io/gorm"
)

type QuestRepository interface {
	GetQuestByID(id uint) (*model.Quest, error)
	GetQuestByUserID(userID uint) (*model.Quest, error)
	AddTakeQuest(userID uint, request model.SubmitQuestRequest) (*model.TakeQuest, error)
}

type questRepository struct {
	db       *gorm.DB
	userRepo UserRepository
}

func (qr *questRepository) GetQuestByUserID(userID uint) (*model.Quest, error) {
	var user model.User
	if err := qr.db.First(&user, userID).Error; err != nil {
		return nil, errors.New("user not found")
	}

	// Collect all question_ids that the user has answered in exercises
	var takeExercises []model.TakeExercise
	if err := qr.db.Where("user_id = ?", userID).Find(&takeExercises).Error; err == nil && len(takeExercises) > 0 {
		var answeredQuestionIDs []uint
		for _, takeEx := range takeExercises {
			var answers map[string]interface{}
			if err := json.Unmarshal(takeEx.Answers, &answers); err == nil {
				for _, v := range answers {
					if detail, ok := v.(map[string]interface{}); ok {
						if qid, ok := detail["question_id"].(float64); ok {
							answeredQuestionIDs = append(answeredQuestionIDs, uint(qid))
						}
					}
				}
			}
		}
		// Remove duplicates
		uniqueQuestionIDs := make(map[uint]bool)
		var uniqueList []uint
		for _, id := range answeredQuestionIDs {
			if !uniqueQuestionIDs[id] {
				uniqueQuestionIDs[id] = true
				uniqueList = append(uniqueList, id)
			}
		}
		if len(uniqueList) > 0 {
			// Pick a random question_id from answered ones
			randIndex := rand.Intn(len(uniqueList))
			questionID := uniqueList[randIndex]

			// Get the ExQuestion
			var exQuestion model.ExQuestion
			if err := qr.db.Preload("Answers").First(&exQuestion, questionID).Error; err == nil {
				// Convert ExQuestion to Quest
				quest := &model.Quest{
					Model:      exQuestion.Model,
					Title:      exQuestion.Content, // Use Content as Title
					Content:    exQuestion.Content,
					Timer:      0, // Default
					Point:      exQuestion.Point,
					Exp:        exQuestion.Exp,
					Difficulty: "medium", // Default
					Type:       exQuestion.Type,
					Feedback:   exQuestion.Feedback,
					Answers:    make([]model.QuestAnswer, len(exQuestion.Answers)),
				}
				for i, ans := range exQuestion.Answers {
					quest.Answers[i] = model.QuestAnswer{
						Model:     ans.Model,
						QuestID:   quest.ID,
						Content:   ans.Content,
						IsCorrect: ans.IsCorrect,
					}
				}
				return quest, nil
			}
		}
	}

	// If no answered questions, get a random new quest
	var quest model.Quest
	if err := qr.db.Preload("Answers").Order("RAND()").First(&quest).Error; err != nil {
		return nil, errors.New("no suitable quest found")
	}

	return &quest, nil
}

func NewQuestRepository(db *gorm.DB, userRepo UserRepository) QuestRepository {
	return &questRepository{db, userRepo}
}

func (qr *questRepository) GetQuestByID(id uint) (*model.Quest, error) {
	var quest model.Quest
	if err := qr.db.Preload("Answers").First(&quest, id).Error; err != nil {
		return nil, errors.New("quest not found")
	}
	return &quest, nil
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
	ratio := 0.0

	if strings.EqualFold(question.Type, "multiple_choice") {
		fmt.Println("DEBUG: Iterating over multiple choice")
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
	} else if strings.EqualFold(question.Type, "true_false") {
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
	} else if strings.EqualFold(question.Type, "short_answer") {
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
			"question_id":       question.ID,
			"feedback":          question.Feedback,
			"exp_gained":        exp,
			"point_gained":      point,
			"user_answer_index": detail["index_jawaban"],
			"user_answer":       detail["answer_text"],
			"correct_answer":    correctAnswer.Content,
		}
	} else if strings.EqualFold(question.Type, "essay") {
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
	} else if strings.EqualFold(question.Type, "multiple_answer") {
		// BUG FIX: use = not := to avoid shadowing outer exp/point variables
		point = question.Point
		exp = question.Exp

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

		// BUG FIX: was `expGained += expGained` (doubled) and `expGained += pointGained` (wrong var)
		rewardExp += expGained
		rewardPoint += pointGained

		takeQuestAnswer = map[string]interface{}{
			"question_id":            question.ID,
			"feedback":               question.Feedback,
			"exp_gained":             expGained,
			"point_gained":           pointGained,
			"user_answer_id":         detail["answer_id"],
			"user_correct_answer_id": jawabanUserBenar,
			"correct_answer_index":   0,
			"correct_answers_index":  correctAnswersIndex,
			"user_answer_index":      detail["index_jawaban"],
			"type":                   question.Type,
			"user_answer":            detail["answer_text"],
		}
	} else if strings.EqualFold(question.Type, "matching") {
		// build correct pairs from data stored in question.Content first,
		// then fall back to question.Answers if necessary.
		type pair struct{ keyword, explanation string }
		var correctPairs []pair
		// try parsing JSON array from Content field
		if strings.TrimSpace(question.Content) != "" {
			var arr []map[string]interface{}
			if err := json.Unmarshal([]byte(question.Content), &arr); err == nil {
				for _, item := range arr {
					kw, _ := item["keyword"].(string)
					exp, _ := item["explanation"].(string)
					correctPairs = append(correctPairs, pair{
						keyword:     strings.TrimSpace(kw),
						explanation: strings.TrimSpace(exp),
					})
				}
			}
		}
		// fallback to Answers table if Content had nothing useful
		if len(correctPairs) == 0 {
			for _, ans := range question.Answers {
				content := strings.TrimSpace(ans.Content)
				var parsed map[string]string
				if err := json.Unmarshal([]byte(content), &parsed); err == nil {
					correctPairs = append(correctPairs, pair{
						keyword:     strings.TrimSpace(parsed["keyword"]),
						explanation: strings.TrimSpace(parsed["explanation"]),
					})
				} else {
					correctPairs = append(correctPairs, pair{keyword: content})
				}
			}
		}
		if len(correctPairs) == 0 {
			log.Printf("⚠️ Matching question %d has no stored answer pairs", question.ID)
		}

		submittedPairs, ok := detail["answers"].([]interface{})
		if !ok || len(submittedPairs) == 0 {
			// no answers submitted
			takeQuestAnswer = map[string]interface{}{
				"question_id":   question.ID,
				"feedback":      question.Feedback,
				"exp_gained":    0,
				"point_gained":  0,
				"user_matches":  []map[string]interface{}{},
				"correct_pairs": correctPairs,
			}
		} else {
			correctCount := 0
			totalPairs := len(submittedPairs)

			var userMatches []map[string]interface{}

			for _, pairAny := range submittedPairs {
				p, ok := pairAny.(map[string]interface{})
				if !ok {
					continue
				}

				kw, ok1 := p["keyword"].(string)
				ex, ok2 := p["explanation"].(string)
				if !ok1 || !ok2 {
					continue
				}

				submittedKeyword := strings.TrimSpace(kw)
				submittedExplanation := strings.TrimSpace(ex)

				userMatches = append(userMatches, map[string]interface{}{
					"explanation": submittedExplanation,
					"keyword":     submittedKeyword,
				})

				// look for an exact match in the correctPairs slice
				for _, cp := range correctPairs {
					if strings.EqualFold(cp.keyword, submittedKeyword) &&
						(cp.explanation == "" || strings.EqualFold(cp.explanation, submittedExplanation)) {
						correctCount++
						break
					}
				}
			}

			// BUG FIX: use = not := so outer `ratio` variable is updated for score calculation below
			if totalPairs > 0 {
				ratio = float64(correctCount) / float64(totalPairs)
			}

			// award proportional exp/point
			exp = int(ratio * float64(question.Exp))
			point = int(ratio * float64(question.Point))

			rewardExp += exp
			rewardPoint += point

			is_correct = correctCount == totalPairs && totalPairs > 0

			takeQuestAnswer = map[string]interface{}{
				"question_id":   question.ID,
				"feedback":      question.Feedback,
				"exp_gained":    exp,
				"point_gained":  point,
				"user_matches":  userMatches,
				"correct_pairs": correctPairs,
				"correct_count": correctCount,
				"total_pairs":   totalPairs,
				"is_correct":    is_correct,
			}
		}
	} else {
		fmt.Println("DEBUG: Question type didn't detect")
	}

	var score float64
	if strings.EqualFold(question.Type, "matching") {
		score = ratio * 100
	} else if is_correct {
		score = 100
	} else {
		score = 0
	}

	answerDetail, err := json.Marshal(takeQuestAnswer)
	if err != nil {
		return nil, err
	}

	newTakeQuest := model.TakeQuest{
		QuestID:     quest.ID,
		UserID:      userID,
		Answer:      answerDetail,
		Score:       score,
		IsCorrect:   is_correct,
		RewardExp:   rewardExp,
		RewardPoint: rewardPoint,
	}

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

	err = qr.db.Create(&newTakeQuest).Error
	if err != nil {
		return nil, err
	}

	return &newTakeQuest, err
}