package service

import (
	"boysitorus/Progamify-Restful-API/internal/model"
	"boysitorus/Progamify-Restful-API/internal/repository"
)

type GiftService interface {
	GetAll() ([]model.Gift, error)
	BuyGift(userID uint, giftID uint) (*model.HaveGift, error)
	GetUserGift(userID uint) ([]model.HaveGift, error)
}

type giftService struct {
	repo repository.GiftRepository
}

func (g *giftService) GetUserGift(userID uint) ([]model.HaveGift, error) {
	return g.repo.GetUserGift(userID)
}

func (g *giftService) GetAll() ([]model.Gift, error) {
	return g.repo.GetAll()
}

func (g *giftService) BuyGift(userID uint, giftID uint) (*model.HaveGift, error) {
	return g.repo.BuyGift(userID, giftID)
}

func NewGiftService(repo repository.GiftRepository) GiftService {
	return &giftService{repo: repo}
}
