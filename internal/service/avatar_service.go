package service

import (
	"boysitorus/Progamify-Restful-API/internal/model"
	"boysitorus/Progamify-Restful-API/internal/repository"
)

type AvatarService interface {
	GetAvatarsByUser(userID uint) ([]map[string]interface{}, error)
	BuyAvatar(haveAvatar *model.HaveAvatar) (model.HaveAvatar, error)
}

type avatarService struct {
	repo repository.AvatarRepository
}

func (as *avatarService) GetAvatarsByUser(userID uint) ([]map[string]interface{}, error) {
	return as.repo.GetAvatarsByUser(userID)
}

func (as *avatarService) BuyAvatar(haveAvatar *model.HaveAvatar) (model.HaveAvatar, error) {
	return as.repo.BuyAvatar(haveAvatar)
}

func NewAvatarService(repo repository.AvatarRepository) AvatarService {
	return &avatarService{repo: repo}
}
