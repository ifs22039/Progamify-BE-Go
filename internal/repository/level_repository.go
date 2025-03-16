package repository

import (
	"boysitorus/Progamify-Restful-API/internal/model"
	"gorm.io/gorm"
)

type LevelRepository interface {
	GetLevel(id uint) (*model.Level, error)
}

type levelRepository struct {
	db *gorm.DB
}

func NewLevelRepository(db *gorm.DB) LevelRepository {
	return &levelRepository{db}
}

func (r *levelRepository) GetLevel(id uint) (*model.Level, error) {
	var level model.Level
	err := r.db.First(&level, id).Error
	return &level, err
}