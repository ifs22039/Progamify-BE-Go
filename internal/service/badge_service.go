package service

import (
	"boysitorus/Progamify-Restful-API/internal/model"
	"boysitorus/Progamify-Restful-API/internal/repository"
)

type BadgeService interface {
	GetBadge(id uint) (*model.Badge, error) 
	AddHaveBadge(userID uint, badgeId uint) (*model.HaveBadge, error)
	CheckAndAssignBadge(userID uint) ([]model.HaveBadge, error)
	GetBadges(id uint) ([]model.UserBadge, error)
}

type badgeService struct {
	badgeRepo repository.BadgeRepository
}

func (s *badgeService) GetBadge(id uint) (*model.Badge, error) {
	return s.badgeRepo.FindBadge(id)
}

func (s *badgeService) AddHaveBadge (userId uint, badgeId uint) (*model.HaveBadge, error) {
	return s.badgeRepo.AddHaveBadge(userId, badgeId)
}

func (s *badgeService) CheckAndAssignBadge (userID uint) ([]model.HaveBadge, error) {
	return s.badgeRepo.AssignBadgeIfEligible(userID)
}

func (s *badgeService) GetBadges(id uint) ([]model.UserBadge, error) {
	return s.badgeRepo.GetBadges(id)
}

func NewBadgeService(badgeRepo repository.BadgeRepository) BadgeService {
	return &badgeService{badgeRepo}
}
