package service

import (
	"boysitorus/Progamify-Restful-API/internal/model"
	"boysitorus/Progamify-Restful-API/internal/repository"
)

type ExerciseService interface {
	// userID may be zero for anonymous callers; when non-zero we'll look up
	// the user's theta and use it to drive adaptive question selection.  Any
	// previously correctly answered questions will be omitted from the result.
	GetExerciseById(id uint, userID uint) (*model.Exercise, error)
	AddTakeExercise(userID uint, request model.SubmitExerciseRequest) (*model.TakeExercise, error)
}

type exerciseService struct {
	exerciseRepo repository.ExerciseRepository
	userRepo     repository.UserRepository
}

func (e *exerciseService) AddTakeExercise(userID uint, request model.SubmitExerciseRequest) (*model.TakeExercise, error) {
	return e.exerciseRepo.AddTakeExercise(userID, request)
}

func (e *exerciseService) GetExerciseById(id uint, userID uint) (*model.Exercise, error) {
	// determine theta based on userID (if provided)
	theta := 0.0
	if userID != 0 {
		if u, err := e.userRepo.FindById(userID); err == nil && u != nil {
			theta = u.Theta
		} else {
			// if for some reason the user isn't found we treat them as anonymous
			// so we won't attempt to filter questions later.  This mirrors how we
			// handled the legacy behaviour before introducing the userID param.
			userID = 0
		}
	}
	return e.exerciseRepo.FindById(id, theta, userID)
}

func NewExerciseService(exerciseRepo repository.ExerciseRepository, userRepo repository.UserRepository) ExerciseService {
	return &exerciseService{exerciseRepo, userRepo}
}
