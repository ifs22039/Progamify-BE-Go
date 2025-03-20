package repository

import (
	"boysitorus/Progamify-Restful-API/internal/model"
	"gorm.io/gorm"
)

type AvatarRepository interface {
	GetAvatarsByUser(userID uint) ([]map[string]interface{}, error)
	BuyAvatar(haveAvatar *model.HaveAvatar) (model.HaveAvatar, error)
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

func (ar *avatarRepository) BuyAvatar(haveAvatar *model.HaveAvatar) (model.HaveAvatar, error) {
	if err := ar.db.Create(haveAvatar).Error; err != nil {
		return model.HaveAvatar{}, err
	}

	return *haveAvatar, nil
}
