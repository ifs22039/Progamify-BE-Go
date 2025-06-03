package repository

import (
	"boysitorus/Progamify-Restful-API/internal/model"
	"boysitorus/Progamify-Restful-API/pkg/utils"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindByEmail(email string) (*model.User, error)
	FindById(id uint) (*model.User, error)
	Update(user *model.User) error
	Create(user *model.User) error
	ExistsByEmail(email string) (bool, error)
	GetTopUsersByExp(limit int, userId uint) ([]model.UserWithRank, error)
	AddPoint(userID uint, point int) error
	AddExp(userID uint, exp int) error
	CheckLevel(userID uint) error
	LessonTaken(userID uint) (int, error)
	ChangeAvatar(userID uint, avatarID uint) (*model.User, error)
	UpdatePassword(id uint, req *model.UpdatePasswordRequest) (*model.User, error)
	UpdateLastLogin(id uint) error
	UpdateUser(user *model.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) LessonTaken(userID uint) (int, error) {
	var count int64

	if err := r.db.Model(&model.TakeLesson{}).
		Where("user_id = ?", userID).
		Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}

func (r *userRepository) CheckLevel(userID uint) error {
	var user model.User
	err := r.db.First(&user, userID).Error

	if err != nil {
		return err
	}

	currentExp := &user.TotalExp

	var newLevel model.Level
	err = r.db.Where("exp_needed <= ?", currentExp).Order("exp_needed DESC").First(&newLevel).Error

	if err != nil {
		return err
	}

	if user.LevelId < newLevel.ID {
		err = r.db.Model(&user).Update("level_id", newLevel.ID).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *userRepository) FindById(id uint) (*model.User, error) {
	var user model.User
	err := r.db.Preload("Avatar").First(&user, id).Error
	return &user, err
}

func (r *userRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) UpdateUser(user *model.User) error {
	return r.db.Model(user).Select("Name", "Nim", "Angkatan").Updates(user).Error
}

func (r *userRepository) Create(user *model.User) error {
	err := r.db.Create(user).Error

	haveAvatar := model.HaveAvatar{
		UserID:   user.ID,
		AvatarID: 1,
	}

	err = r.db.Create(&haveAvatar).Error
	if err != nil {
		return err
	}

	return err
}

func (r *userRepository) ExistsByEmail(email string) (bool, error) {
	var count int64
	err := r.db.Model(&model.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

func (r *userRepository) AddExp(userID uint, exp int) error {
	var user model.User
	err := r.db.First(&user, userID).Error

	if err != nil {
		return err
	}

	user.TotalExp += exp

	err = r.db.Save(&user).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) AddPoint(userID uint, point int) error {
	var user model.User
	err := r.db.First(&user, userID).Error

	if err != nil {
		return err
	}

	user.TotalPoint += point

	err = r.db.Save(&user).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) GetTopUsersByExp(limit int, userId uint) ([]model.UserWithRank, error) {
	var topUsers []model.User
	var result []model.UserWithRank

	err := r.db.Order("total_exp DESC").Limit(limit).Preload("Avatar").Find(&topUsers).Error
	if err != nil {
		return nil, err
	}

	var rank int64
	err = r.db.Model(&model.User{}).
		Where("total_exp > (?)", r.db.Model(&model.User{}).Select("total_exp").Where("id = ?", userId)).
		Count(&rank).Error
	if err != nil {
		return nil, err
	}
	rank++

	userIncluded := false
	for i, user := range topUsers {
		result = append(result, model.UserWithRank{
			User: user,
			Rank: i + 1,
		})
		if user.ID == userId {
			userIncluded = true
		}
	}

	if !userIncluded {
		var currentUser model.User
		err = r.db.Where("id = ?", userId).Preload("Avatar").First(&currentUser).Error
		if err != nil {
			return nil, err
		}
		result = append(result, model.UserWithRank{
			User: currentUser,
			Rank: int(rank),
		})
	}

	return result, nil
}

func (r *userRepository) ChangeAvatar(userID uint, avatarID uint) (*model.User, error) {
	var user model.User
	var count int64

	err := r.db.Table("have_avatars").
		Where("user_id = ? AND avatar_id = ?", userID, avatarID).
		Count(&count).Error

	if err != nil {
		return &model.User{}, err
	}

	if count == 0 {
		return &model.User{}, errors.New("avatar not owned by user")
	}

	err = r.db.First(&user, userID).Error
	if err != nil {
		return &model.User{}, err
	}

	err = r.db.Model(&user).Update("avatar_id", avatarID).Error
	if err != nil {
		return &model.User{}, err
	}

	return &user, nil
}

func (r *userRepository) UpdatePassword(id uint, req *model.UpdatePasswordRequest) (*model.User, error) {
	var user model.User

	err := r.db.First(&user, id).Error
	if err != nil {
		return &model.User{}, err
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, fmt.Errorf("invalid current password")
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return nil, fmt.Errorf("error hashing new password: %v", err)
	}

	err = r.db.Model(&user).Update("password", hashedPassword).Error
	if err != nil {
		return &model.User{}, err
	}

	return &user, nil
}

func (r *userRepository) UpdateLastLogin(id uint) error {
	var user model.User

	err := r.db.First(&user, id).Error
	if err != nil {
		return err
	}

	err = r.db.Model(&user).Update("last_login", time.Now()).Error
	if err != nil {
		return err
	}

	return nil
}
