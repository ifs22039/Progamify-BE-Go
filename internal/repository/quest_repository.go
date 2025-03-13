package repository

import (
	"errors"
	"math/rand"
	"time"

	"boysitorus/Progamify-Restful-API/internal/model"
	"gorm.io/gorm"
)

type QuestRepository interface {
	GetQuestByUserID(userID uint) (*model.Quest, error)
	GetDifficultyByLevel(levelId uint) string
}

type questRepository struct {
	db *gorm.DB
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

	// Menggunakan method receiver qr.GetDifficultyByLevel
	difficulty := qr.GetDifficultyByLevel(uint(level.Level))

	var quest model.Quest
	if err := qr.db.Preload("Answers").Where("difficulty = ?", difficulty).Order("RAND()").First(&quest).Error; err != nil {
		return nil, errors.New("no suitable quest found")
	}

	return &quest, nil
}

// Perbaikan tipe parameter agar sesuai dengan interface
func (qr *questRepository) GetDifficultyByLevel(levelId uint) string {
	rand.Seed(time.Now().UnixNano())
	probability := rand.Float64() // Angka acak antara 0.0 - 1.0

	switch {
	case levelId <= 10: // Level 1-10: Selalu Easy
		return "Easy"
	case levelId <= 29: // Level 11-29: Easy mulai berkurang, Medium mulai ada
		if probability < 0.7 {
			return "Easy"
		}
		return "Medium"
	case levelId <= 49: // Level 30-49: Medium lebih sering, Hard mulai ada
		if probability < 0.5 {
			return "Medium"
		} else if probability < 0.85 {
			return "Hard"
		}
		return "Easy"
	case levelId <= 79: // Level 50-79: Hard lebih dominan, Very Hard mulai ada
		if probability < 0.4 {
			return "Medium"
		} else if probability < 0.75 {
			return "Hard"
		}
		return "Very Hard"
	default: // Level 80+: Very Hard dominan, Hard masih ada, Medium sedikit
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
