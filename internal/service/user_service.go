package service

import (
	"boysitorus/Progamify-Restful-API/internal/model"
	"boysitorus/Progamify-Restful-API/internal/repository"
)

type UserService interface {
	GetUser(id uint) (*model.User, error)
	UpdateUser(id uint, req *model.UpdateUserRequest) (*model.User, error)
	GetLessonTaken(userID uint) (int, error)
	ChangeAvatar(userID uint, avatarID uint) (*model.User, error)
	UpdatePassword(id uint, req *model.UpdatePasswordRequest) (*model.User, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func (s *userService) ChangeAvatar(userID uint, avatarID uint) (*model.User, error) {
	return s.userRepo.ChangeAvatar(userID, avatarID)
}

func (s *userService) GetLessonTaken(userID uint) (int, error) {
	return s.userRepo.LessonTaken(userID)
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo}
}

func (s *userService) GetUser(id uint) (*model.User, error) {
	return s.userRepo.FindById(id)
}

func (s *userService) UpdateUser(id uint, req *model.UpdateUserRequest) (*model.User, error) {
	user, err := s.userRepo.FindById(id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Nim != "" {
		user.Nim = req.Nim
	}
	if req.Angkatan != 0 {
		user.Angkatan = req.Angkatan
	}

	err = s.userRepo.Update(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) UpdatePassword(id uint, req *model.UpdatePasswordRequest) (*model.User, error) {
	return s.userRepo.UpdatePassword(id, req)
}
