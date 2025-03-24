package service

import (
	"boysitorus/Progamify-Restful-API/internal/model"
	"boysitorus/Progamify-Restful-API/internal/repository"
)

type AchievementService interface {
	GetAchievement(id uint) (*model.Achievement, error) 
	AddHaveAchievement(userID uint, achievementId uint) (*model.HaveAchievement, error)
	CheckAndAssignAchievement(userID uint) ([]model.HaveAchievement, error)
	GetAchievements(id uint) ([]model.UserAchievement, error)
}

type achievementService struct {
	achievementRepo repository.AchievementRepository
}

func (s *achievementService) GetAchievement(id uint) (*model.Achievement, error) {
	return s.achievementRepo.FindAchievement(id)
}

func (s *achievementService) AddHaveAchievement (userId uint, achievementId uint) (*model.HaveAchievement, error) {
	return s.achievementRepo.AddHaveAchievement(userId, achievementId)
}

func (s *achievementService) CheckAndAssignAchievement (userID uint) ([]model.HaveAchievement, error) {
	return s.achievementRepo.AssignAchievementIfEligible(userID)
}

func (s *achievementService) GetAchievements(id uint) ([]model.UserAchievement, error) {
	return s.achievementRepo.GetAchievements(id)
}

func NewAchievementService(achievementRepo repository.AchievementRepository) AchievementService {
	return &achievementService{achievementRepo}
}
