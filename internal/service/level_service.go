package service

import (
	"boysitorus/Progamify-Restful-API/internal/model"
	"boysitorus/Progamify-Restful-API/internal/repository"
)

type LevelService interface {
	GetLevel(id uint) (*model.Level, error)
}

type levelService struct {
	levelRepo repository.LevelRepository
}

func NewLevelService(levelRepo repository.LevelRepository) LevelService {
	return &levelService{levelRepo}
}

func (s *levelService) GetLevel(id uint) (*model.Level, error) {
	return s.levelRepo.GetLevel(id)
}
