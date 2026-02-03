package repository

import (
	"boysitorus/Progamify-Restful-API/internal/model"
	"boysitorus/Progamify-Restful-API/pkg/utils"
	"encoding/json"
	"fmt"

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

func (e *exerciseRepository) AddTakeExercise(
    userID uint,
    request model.SubmitExerciseRequest,
) (*model.TakeExercise, error) {

    fmt.Println("✅ AddTakeExercise DIPANGGIL | userID:", userID)
	// 1. Ambil Exercise
	var exercise model.Exercise
	if err := e.db.Preload("Questions.Answers").
		First(&exercise, request.ExerciseID).Error; err != nil {
		return nil, err
	}

	// 2. Ambil Lesson
	var lesson model.Lesson
	if err := e.db.First(&lesson, exercise.LessonID).Error; err != nil {
		return nil, err
	}

	// 3. Ambil User (Theta)
	user, err := e.userRepo.FindById(userID)
	if err != nil {
		return nil, err
	}

	thetaBefore := user.Theta
	betaBefore := exercise.Beta

	// 4. Hitung attempt
	var attemptCount int64
	e.db.Model(&model.TakeExercise{}).
		Where("user_id = ? AND exercise_id = ?", userID, request.ExerciseID).
		Count(&attemptCount)
	attemptNumber := int(attemptCount) + 1

	// =========================
	// 5️⃣ GRADING ASLI KAMU
	// =========================

	answersJSON := request.Answers
	takeExerciseAnswer := make(map[int]interface{})

	totalCorrect := 0.0
	totalExp := 0
	totalPoint := 0
	rewardExp := 0
	rewardPoint := 0

	for key, value := range answersJSON {
		detail := value.(map[string]interface{})

		var question model.ExQuestion
		if err := e.db.Preload("Answers").
			First(&question, detail["question_id"]).Error; err != nil {
			return nil, err
		}

		exp := 0
		point := 0
		rewardExp += question.Exp
		rewardPoint += question.Point

		if question.Type == "multiple_choice" {
			var correct model.ExAnswer
			var correctIndex int

			for i, a := range question.Answers {
				if a.IsCorrect {
					correct = a
					correctIndex = i
				}
			}

			if int(correct.ID) == int(detail["answer_id"].(float64)) {
				exp = question.Exp
				point = question.Point
				totalCorrect++
			}

			totalExp += exp
			totalPoint += point

			takeExerciseAnswer[key] = map[string]interface{}{
				"question_id":          question.ID,
				"type":                 question.Type,
				"exp_gained":           exp,
				"point_gained":         point,
				"correct_answer_index": correctIndex,
			}
		}
	}

	totalQuestion := float64(len(exercise.Questions))
	score := totalCorrect / totalQuestion

	// =========================
	// 6️⃣ IRT (RASCH 1PL)
	// =========================

	isCorrect := score >= 0.6
	p := utils.RaschProbability(thetaBefore, betaBefore)
	thetaAfter := utils.UpdateTheta(thetaBefore, p, isCorrect)
	betaAfter := betaBefore // beta statis (Rasch 1PL)

	// =========================
	// 7️⃣ UPDATE USER.THETA
	// =========================

	if err := e.db.Model(&model.User{}).
		Where("id = ?", userID).
		Update("theta", thetaAfter).Error; err != nil {
		return nil, err
	}

	// =========================
	// 8️⃣ SIMPAN TAKE_EXERCISE
	// =========================

	answerDetail, _ := json.Marshal(takeExerciseAnswer)

	newTake := model.TakeExercise{
		ExerciseID:    exercise.ID,
		LessonID:      lesson.ID,
		UserID:        userID,
		TopicID:       lesson.TopicID,
		AttemptNumber: attemptNumber,
		Answers:       answerDetail,

		Score:         score,
		TotalCorrect:  int(totalCorrect),
		TotalQuestion: int(totalQuestion),
		TotalExp:      totalExp,
		TotalPoint:    totalPoint,
		RewardExp:     rewardExp,
		RewardPoint:   rewardPoint,

		ThetaBefore: thetaBefore,
		ThetaAfter:  thetaAfter,
		BetaBefore:  betaBefore,
		BetaAfter:   betaAfter,
	}

	if err := e.db.Create(&newTake).Error; err != nil {
		return nil, err
	}

	// =========================
	// 9️⃣ GAMIFICATION
	// =========================

	_ = e.userRepo.AddPoint(userID, totalPoint)
	_ = e.userRepo.AddExp(userID, totalExp)
	_ = e.userRepo.CheckLevel(userID)

	fmt.Println("✅ IRT OK | theta:", thetaBefore, "→", thetaAfter)

	return &newTake, nil
}

/* =========================
   FIND EXERCISE
========================= */

func (e *exerciseRepository) FindById(id uint) (*model.Exercise, error) {
	var exercise model.Exercise
	if err := e.db.
		Preload("Questions", func(db *gorm.DB) *gorm.DB {
			return db.Order("RAND()").Limit(5)
		}).
		Preload("Questions.Answers").
		First(&exercise, id).Error; err != nil {
		return nil, err
	}
	return &exercise, nil
}
