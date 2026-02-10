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

		// 1) Try fetch by question_id when provided
		if qid, ok := detail["question_id"].(float64); ok {
			if err := e.db.Preload("Answers").First(&question, uint(qid)).Error; err != nil {
				return nil, err
			}
		} else if kw, ok := detail["keyword"].(string); ok {
			// 2) Try to find question by matching keyword with question.Content from preloaded exercise
			found := false
			for _, q := range exercise.Questions {
				if q.Content == kw {
					question = q
					// ensure answers are loaded
					if len(question.Answers) == 0 {
						if err := e.db.Preload("Answers").First(&question, question.ID).Error; err != nil {
							return nil, err
						}
					}
					found = true
					break
				}
			}
			if !found {
				if err := e.db.Preload("Answers").Where("content = ? AND exercise_id = ?", kw, exercise.ID).First(&question).Error; err != nil {
					return nil, err
				}
			}
		} else {
			return nil, fmt.Errorf("invalid answer payload: missing question_id or keyword")
		}

		exp := 0
		point := 0
		rewardExp += question.Exp
		rewardPoint += question.Point

		// MULTIPLE CHOICE (existing behaviour)
		if question.Type == "multiple_choice" {
			var correct model.ExAnswer
			var correctIndex int

			for i, a := range question.Answers {
				if a.IsCorrect {
					correct = a
					correctIndex = i
				}
			}

			if ansid, ok := detail["answer_id"].(float64); ok {
				if int(correct.ID) == int(ansid) {
					exp = question.Exp
					point = question.Point
					totalCorrect++
				}
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

		// MATCHING: accept either `answer_id` or `explanation` string, or match by `keyword` -> question.Content
		} else if question.Type == "matching" {

    // Ambil semua keyword yang benar (harusnya pure string, bukan JSON)
    var correctKeywords []string
    for _, ans := range question.Answers {
        if ans.IsCorrect {
            // Pastikan Content adalah keyword murni, bukan JSON string
            content := strings.TrimSpace(ans.Content)
            // Jika Content sudah JSON, parse dulu (jika struktur database salah)
            var parsed map[string]string
            if err := json.Unmarshal([]byte(content), &parsed); err == nil {
                if kw, ok := parsed["keyword"]; ok {
                    correctKeywords = append(correctKeywords, strings.TrimSpace(kw))
                }
            } else {
                correctKeywords = append(correctKeywords, content)
            }
        }
    }

    submittedPairs, ok := detail["answers"].([]interface{})
    if !ok || len(submittedPairs) == 0 {
        // handle error
        continue
    }

    correctCount := 0
    totalPairs := len(submittedPairs)  // jumlah yang user isi, atau gunakan jumlah explanations

    var userMatches []map[string]interface{}

    for _, pairAny := range submittedPairs {
        pair, ok := pairAny.(map[string]interface{})
        if !ok { continue }

        submittedKeyword := strings.TrimSpace(pair["keyword"].(string))
        submittedExplanation := strings.TrimSpace(pair["explanation"].(string))

        userMatches = append(userMatches, map[string]interface{}{
            "explanation": submittedExplanation,
            "keyword":     submittedKeyword,
        })

        for _, correctKw := range correctKeywords {
            if submittedKeyword == correctKw {
                correctCount++
                break
            }
        }
    }

			// Skor per soal matching = jumlah pasangan benar / total pasangan
			isFullyCorrect := correctCount == totalPairs && totalPairs > 0
			if isFullyCorrect {
				exp = question.Exp
				point = question.Point
				totalCorrect += 1 // atau += float64(correctCount)/float64(totalPairs)
			}

			totalExp += exp
			totalPoint += point

			takeExerciseAnswer[key] = map[string]interface{}{
				"question_id":       question.ID,
				"type":              question.Type,
				"exp_gained":        exp,
				"point_gained":      point,
				"correct_keywords":  correctKeywords,
				"user_matches":      userMatches,
				"correct_count":     correctCount,
				"total_pairs":       totalPairs,
				"is_fully_correct":  isFullyCorrect,
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