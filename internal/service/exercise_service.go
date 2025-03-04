package service

import (
	"boysitorus/Progamify-Restful-API/internal/model"
	"boysitorus/Progamify-Restful-API/internal/repository"
)

type ExerciseService interface {
	GetExerciseById(id uint) (*model.Exercise, error)
}

type exerciseService struct {
	exerciseRepo repository.ExerciseRepository
}

func (e *exerciseService) GetExerciseById(id uint) (*model.Exercise, error) {
	return e.exerciseRepo.FindById(id)
}

func NewExerciseService(exerciseRepo repository.ExerciseRepository) ExerciseService {
	return &exerciseService{exerciseRepo}
}
