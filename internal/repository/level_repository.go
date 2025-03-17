package repository

import (
	"boysitorus/Progamify-Restful-API/internal/model"

	"gorm.io/gorm"
)

type LevelRepository interface {
	GetLevel(id uint) (*LevelWithNext, error) 
	GetLevelByUserId(userId uint) (*model.Level, error)
}

type LevelWithNext struct {
	model.Level
	NextLevel *model.Level `json:"next_level,omitempty"` 
}

type levelRepository struct {
	db *gorm.DB
}

func NewLevelRepository(db *gorm.DB) LevelRepository {
	return &levelRepository{db}
}

func (r *levelRepository) GetLevel(id uint) (*LevelWithNext, error)  {
	var level model.Level
	err := r.db.First(&level, id).Error
	if err != nil {
		return nil, err
	}

	var nextLevel model.Level
	var nextLevelPtr *model.Level

	err = r.db.Where("exp_needed > ?", level.ExpNeeded).Order("exp_needed ASC").First(&nextLevel).Error
	if err == nil { 
		nextLevelPtr = &nextLevel
	}
	
	levelWithNext := &LevelWithNext{
		Level:     level,
		NextLevel: nextLevelPtr,
	}

	return levelWithNext, nil
}

func (r *levelRepository) GetLevelByUserId(userId uint) (*model.Level, error) {
	var user model.User
	err := r.db.First(&user, userId).Error
	if err != nil {
		return nil, err
	}

	var level model.Level
	err = r.db.First(&level, user.LevelId).Error
	if err != nil {
		return nil, err
	}

	return &level, nil
}