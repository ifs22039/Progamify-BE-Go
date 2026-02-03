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
		err = r.db.Model(&model.User{}).
    Where("id = ?", userID).
    Update("level_id", newLevel.ID).Error
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
	return r.db.Model(&model.User{}).
		Where("id = ?", userID).
		Update("total_exp", gorm.Expr("total_exp + ?", exp)).
		Error
}

func (r *userRepository) AddPoint(userID uint, point int) error {
	return r.db.Model(&model.User{}).
		Where("id = ?", userID).
		Update("total_point", gorm.Expr("total_point + ?", point)).
		Error
}

func (r *userRepository) GetTopUsersByExp(limit int, userId uint) ([]model.UserWithRank, error) {
    var allUsers []model.User
    var result []model.UserWithRank

    // Step 1: Urutkan semua user berdasarkan total_exp DESC
    err := r.db.Order("total_exp DESC, id ASC").Preload("Avatar").Find(&allUsers).Error
    if err != nil {
        return nil, err
    }

    var userInTop bool = false
    var userWithRank model.UserWithRank

    // Step 2: Bangun top N users dan temukan userId rank-nya
    for i, user := range allUsers {
        rank := i + 1
        if rank <= limit {
            // Masukkan ke dalam top N
            result = append(result, model.UserWithRank{
                User: user,
                Rank: rank,
            })
            if user.ID == userId {
                userInTop = true
            }
        }
        if user.ID == userId {
            // Simpan userId jika belum ada di top N
            userWithRank = model.UserWithRank{
                User: user,
                Rank: rank,
            }
        }
    }

    // Step 3: Jika userId tidak ada di top N, tambahkan dia (dengan rank asli)
    if !userInTop {
        result = append(result, userWithRank)
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
