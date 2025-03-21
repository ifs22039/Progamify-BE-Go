package repository

import (
	"boysitorus/Progamify-Restful-API/internal/model"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type GiftRepository interface {
	GetAll() ([]model.Gift, error)
	BuyGift(userID uint, giftID uint) (*model.HaveGift, error)
}

type giftRepository struct {
	db *gorm.DB
}

func (g *giftRepository) GetAll() ([]model.Gift, error) {
	var gifts []model.Gift
	err := g.db.Find(&gifts).Error
	return gifts, err
}

func (g *giftRepository) BuyGift(userID uint, giftID uint) (*model.HaveGift, error) {
	var user model.User
	var gift model.Gift

	if err := g.db.First(&user, userID).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	if err := g.db.First(&gift, giftID).Error; err != nil {
		return nil, fmt.Errorf("gift not found: %w", err)
	}

	if user.TotalPoint < int(gift.Price) {
		return nil, errors.New("Point tidak cukup")
	}

	user.TotalPoint -= int(gift.Price)
	if err := g.db.Save(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to update user points: %w", err)
	}

	haveGift := model.HaveGift{
		UserID:   userID,
		GiftID:   giftID,
		IsActive: true,
	}

	if err := g.db.Create(&haveGift).Error; err != nil {
		return nil, fmt.Errorf("failed to create HaveGift record: %w", err)
	}

	return &haveGift, nil
}

func NewGiftRepository(db *gorm.DB) GiftRepository {
	return &giftRepository{db}
}
