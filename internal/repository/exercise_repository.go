package repository

import (
	"boysitorus/Progamify-Restful-API/internal/model"
	"gorm.io/gorm"
)

type ExerciseRepository interface {
	FindById(id uint) (*model.Exercise, error)
}

type exerciseRepository struct {
	db *gorm.DB
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

func NewExerciseRepository(db *gorm.DB) ExerciseRepository {
	return &exerciseRepository{db}
}
