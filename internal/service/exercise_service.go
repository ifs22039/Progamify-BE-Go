package service

import (
	"boysitorus/Progamify-Restful-API/internal/model"
	"boysitorus/Progamify-Restful-API/internal/repository"
)

type ExerciseService interface {
	GetExerciseById(id uint) (*model.Exercise, error)
	AddTakeExercise(userID uint, request model.SubmitExerciseRequest) (*model.TakeExercise, error)
}

type exerciseService struct {
	exerciseRepo repository.ExerciseRepository
}

func (e *exerciseService) AddTakeExercise(userID uint, request model.SubmitExerciseRequest) (*model.TakeExercise, error) {
	return e.exerciseRepo.AddTakeExercise(userID, request)
}

func (e *exerciseService) GetExerciseById(id uint) (*model.Exercise, error) {
	return e.exerciseRepo.FindById(id)
}

func NewExerciseService(exerciseRepo repository.ExerciseRepository) ExerciseService {
	return &exerciseService{exerciseRepo}
}
