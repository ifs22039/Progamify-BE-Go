package service

import (
	"boysitorus/Progamify-Restful-API/internal/model"
	"boysitorus/Progamify-Restful-API/internal/repository"
)

type LeaderboardService interface {
	GetLeaderboard(limit int) ([]model.User, error)
}

type leaderboardService struct {
	userRepo repository.UserRepository
}

func NewLeaderboardService(userRepo repository.UserRepository) LeaderboardService {
	return &leaderboardService{userRepo}
}

func (s *leaderboardService) GetLeaderboard(limit int) ([]model.User, error) {
	return s.userRepo.GetTopUsersByExp(limit)
}
