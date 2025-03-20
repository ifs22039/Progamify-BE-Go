package repository

import (
	"boysitorus/Progamify-Restful-API/internal/model"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type AvatarRepository interface {
	GetAvatarsByUser(userID uint) ([]map[string]interface{}, error)
	BuyAvatar(userID uint, avatarID uint) (*model.HaveAvatar, error)
}

type avatarRepository struct {
	db *gorm.DB
}

func NewAvatarRepository(db *gorm.DB) AvatarRepository {
	return &avatarRepository{db}
}

func (ar *avatarRepository) GetAvatarsByUser(userID uint) ([]map[string]interface{}, error) {
	var avatars []map[string]interface{}

	err := ar.db.
		Table("avatars").
		Select(`
			avatars.id, 
			avatars.title, 
			avatars.price, 
			avatars.picture, 
			CASE 
				WHEN have_avatars.id IS NOT NULL THEN 1 
				ELSE 0 
			END AS is_locked
		`).
		Joins("LEFT JOIN have_avatars ON avatars.id = have_avatars.avatar_id AND have_avatars.user_id = ?", userID).
		Find(&avatars).Error

	return avatars, err
}

func (ar *avatarRepository) BuyAvatar(userID uint, avatarID uint) (*model.HaveAvatar, error) {
	var user model.User
	var avatar model.Avatar
	var haveAvatar model.HaveAvatar

	if err := ar.db.First(&user, userID).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	if err := ar.db.First(&avatar, avatarID).Error; err != nil {
		return nil, fmt.Errorf("avatar not found: %w", err)
	}

	if err := ar.db.Where("user_id = ? AND avatar_id = ?", userID, avatarID).First(&haveAvatar).Error; err == nil {
		return nil, errors.New("user telah membeli avatar ini sebelumnya")
	}

	if user.TotalPoint < int(avatar.Price) {
		return nil, errors.New("Point tidak cukup")
	}

	user.TotalPoint -= int(avatar.Price)
	if err := ar.db.Save(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to update user points: %w", err)
	}

	haveAvatar = model.HaveAvatar{
		UserID:   userID,
		AvatarID: avatarID,
	}

	if err := ar.db.Create(&haveAvatar).Error; err != nil {
		return nil, fmt.Errorf("failed to create HaveAvatar record: %w", err)
	}

	return &haveAvatar, nil
}
