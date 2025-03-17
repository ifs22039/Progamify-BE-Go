package service

import (
	"boysitorus/Progamify-Restful-API/internal/model"
	"boysitorus/Progamify-Restful-API/internal/repository"
)

type LevelService interface {
	GetLevel(id uint) (*repository.LevelWithNext, error) 
	GetLevelByUserId(userId uint) (*model.Level, error)
}

type levelService struct {
	levelRepo repository.LevelRepository
}

func (s *levelService) GetLevel(id uint) (*repository.LevelWithNext, error) {
	return s.levelRepo.GetLevel(id)
}

func (s *levelService) GetLevelByUserId(userId uint) (*model.Level, error) {
	return s.levelRepo.GetLevelByUserId(userId)
}

func NewLevelService(levelRepo repository.LevelRepository) LevelService {
	return &levelService{levelRepo}
}
