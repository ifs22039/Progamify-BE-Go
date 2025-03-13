package repository

import (
	"boysitorus/Progamify-Restful-API/internal/model"
	"boysitorus/Progamify-Restful-API/pkg/utils"
	"encoding/json"
	"gorm.io/gorm"
	"strings"
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
	//var takeExercise model.TakeExercise
	//
	//err := e.db.Where("user_id = ? AND exercise_id = ?", userID, request.ExerciseID).First(&takeExercise).Error
	//
	//if err == nil {
	//	return &takeExercise, errors.New("This user already take the exercise")
	//}

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
			}
		} else if question.Type == "true_false" {
			correctAnswer := question.Answers[0]
			correctAnswerIndex := 0
			jawabanUser := detail["answer_text"].(string)

			if strings.ToLower(correctAnswer.Content) == strings.ToLower(jawabanUser) {
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

				totalExp += exp
				totalPoint += point
				totalCorrect += 1
			}

			takeExerciseAnswer[key] = map[string]interface{}{
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

				totalExp += exp
				totalPoint += point
				totalCorrect += 1
			}

			takeExerciseAnswer[key] = map[string]interface{}{
				"question_id":          question.ID,
				"feedback":             question.Feedback,
				"exp_gained":           exp,
				"point_gained":         point,
				"user_answer_index":    detail["index_jawaban"],
				"correct_answer_index": correctAnswer.Content,
			}
		} else if question.Type == "multiple_answer" {

		}
	}

	totalQuestions := float64(len(exercise.Questions))
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

	err = e.userRepo.AddExp(userID, totalPoint)

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
	err := e.db.Preload("Questions.Answers").First(&exercise, id).Error

	if err != nil {
		// Return the error early if the record is not found
		return nil, err
	}

	return &exercise, err
}
