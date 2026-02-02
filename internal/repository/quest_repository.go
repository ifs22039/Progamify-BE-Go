package repository

import (
	"encoding/json"
	"errors"
	"math"
	"math/rand"
	"strings"
	"time"

	"boysitorus/Progamify-Restful-API/internal/model"
	"boysitorus/Progamify-Restful-API/pkg/utils"

	"gorm.io/gorm"
)




type QuestRepository interface {
	GetQuestByUserID(userID uint) (*model.Quest, error)
	GetQuestByID(id uint) (*model.Quest, error)
	GetAdaptiveQuest(theta float64) (*model.Quest, error)
	GetDifficultyByLevel(levelId uint) string

	GradeQuest(
		user *model.User,
		quest *model.Quest,
		request model.SubmitQuestRequest,
	) (*model.TakeQuest, error)

	CreateTakeQuest(takeQuest *model.TakeQuest) error
}

type questRepository struct {
	db       *gorm.DB
	userRepo UserRepository
}

func NewQuestRepository(db *gorm.DB, userRepo UserRepository) QuestRepository {
	return &questRepository{db: db, userRepo: userRepo}
}

//////////////////////////////
// QUEST SELECTION
//////////////////////////////

func (qr *questRepository) GetQuestByID(id uint) (*model.Quest, error) {
	var quest model.Quest
	if err := qr.db.Preload("Answers").First(&quest, id).Error; err != nil {
		return nil, err
	}
	return &quest, nil
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
	if err := qr.db.
		Preload("Answers").
		Where("difficulty = ?", difficulty).
		Order("RAND()").
		First(&quest).Error; err != nil {
		return nil, errors.New("no suitable quest found")
	}

	return &quest, nil
}
func (qr *questRepository) GetAdaptiveQuest(theta float64) (*model.Quest, error) {
	var quest model.Quest

	err := qr.db.
		Order(gorm.Expr("ABS(beta - ?)", theta)).
		Preload("Answers").
		First(&quest).Error

	if err != nil {
		return nil, err
	}

	return &quest, nil
}

func (qr *questRepository) GetDifficultyByLevel(levelId uint) string {
	rand.Seed(time.Now().UnixNano())
	p := rand.Float64()

	switch {
	case levelId <= 10:
		return "Easy"
	case levelId <= 29:
		if p < 0.7 {
			return "Easy"
		}
		return "Medium"
	case levelId <= 49:
		if p < 0.5 {
			return "Medium"
		} else if p < 0.85 {
			return "Hard"
		}
		return "Easy"
	case levelId <= 79:
		if p < 0.4 {
			return "Medium"
		} else if p < 0.75 {
			return "Hard"
		}
		return "Very Hard"
	default:
		if p < 0.3 {
			return "Medium"
		} else if p < 0.6 {
			return "Hard"
		}
		return "Very Hard"
	}
}

//////////////////////////////
// TAKE QUEST (FULL GRADING + IRT)
//////////////////////////////

func (qr *questRepository) CreateTakeQuest(takeQuest *model.TakeQuest) error {
	return qr.db.Create(takeQuest).Error
}

//////////////////////////////
// IRT UTIL (INLINE, TIDAK DIPISAH)
//////////////////////////////

func raschProbability(theta, beta float64) float64 {
	return math.Exp(theta-beta) / (1 + math.Exp(theta-beta))
}

func updateTheta(theta, p float64, correct bool) float64 {
	var u float64 = 0
	if correct {
		u = 1
	}
	learningRate := 0.3
	return theta + learningRate*(u-p)
}

//////////////////////////////
// LEGACY GRADING (TIDAK DIHAPUS)
//////////////////////////////

func (qr *questRepository) GradeQuest(
	user *model.User,
	quest *model.Quest,
	request model.SubmitQuestRequest,
) (*model.TakeQuest, error) {

	answerJSON := request.Answer
	detail := answerJSON.(map[string]interface{})

	var question model.Quest
	if err := qr.db.Preload("Answers").First(&question, detail["question_id"]).Error; err != nil {
		return nil, err
	}

	rewardExp := 0
	rewardPoint := 0
	isCorrect := false

	exp := 0
	point := 0

	takeQuestAnswer := make(map[string]interface{})

	switch question.Type {

	case "multiple_choice":
		var correct model.QuestAnswer
		var correctIndex int
		for i, a := range question.Answers {
			if a.IsCorrect {
				correct = a
				correctIndex = i
			}
		}

		if int(correct.ID) == int(detail["answer_id"].(float64)) {
			exp += question.Exp
			point += question.Point
			isCorrect = true
		}

		takeQuestAnswer = map[string]interface{}{
			"question_id":          question.ID,
			"user_answer_id":       detail["answer_id"],
			"correct_answer_id":    correct.ID,
			"correct_answer_index": correctIndex,
		}

	case "true_false", "short_answer":
		correct := question.Answers[0]
		userAnswer := strings.ToLower(detail["answer_text"].(string))
		if strings.ToLower(correct.Content) == userAnswer {
			exp += question.Exp
			point += question.Point
			isCorrect = true
		}

	case "essay":
		correct := question.Answers[0]
		userAnswer := detail["answer_text"].(string)

		flag := true
		var similarity float64

		for flag {
			result, err := utils.EssayGrading(correct.Content, userAnswer)
			if err == nil {
				flag = false
				similarity = result
			}
		}

		if similarity >= 50 {
			exp += question.Exp
			point += question.Point
			isCorrect = true
		}
	}

	rewardExp += exp
	rewardPoint += point

	// IRT
	p := raschProbability(user.Theta, quest.Beta)
	thetaBefore := user.Theta
	thetaAfter := updateTheta(user.Theta, p, isCorrect)

	user.Theta = thetaAfter
	if err := qr.userRepo.Update(user); err != nil {
		return nil, err
	}

	answerDetail, _ := json.Marshal(takeQuestAnswer)

	takeQuest := &model.TakeQuest{
		UserID:      user.ID,
		QuestID:    quest.ID,
		Answer:     answerDetail,
		IsCorrect:  isCorrect,
		Score:      func() float64 { if isCorrect { return 100 } else { return 0 } }(),
		RewardExp:  rewardExp,
		RewardPoint: rewardPoint,

		// IRT
		Probability: p,
		ThetaBefore: thetaBefore,
		ThetaAfter:  thetaAfter,
	}

	return takeQuest, nil
}
