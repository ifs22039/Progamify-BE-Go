package repository

import (
	"boysitorus/Progamify-Restful-API/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByEmail(email string) (*model.User, error)
	FindById(id uint) (*model.User, error)
	Update(user *model.User) error
	Create(user *model.User) error
	ExistsByEmail(email string) (bool, error)
	GetTopUsersByExp(limit int) ([]model.User, error)
	AddPoint(userID uint, point int) error
	AddExp(userID uint, exp int) error
	CheckLevel(userID uint) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
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
	err := r.db.First(&user, id).Error
	return &user, err
}

func (r *userRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
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

func (r *userRepository) GetTopUsersByExp(limit int) ([]model.User, error) {
	var users []model.User
	err := r.db.Order("total_exp DESC").Limit(limit).Find(&users).Error
	return users, err
}
