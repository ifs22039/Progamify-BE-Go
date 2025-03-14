package service

import (
	"boysitorus/Progamify-Restful-API/internal/config"
	"boysitorus/Progamify-Restful-API/internal/model"
	"boysitorus/Progamify-Restful-API/internal/repository"
	"boysitorus/Progamify-Restful-API/pkg/utils"
	"fmt"
)

type AuthService interface {
	Login(email, password string) (string, *model.User, error)
	Register(req *model.RegisterRequest) (*model.User, error)
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo}
}

func (s *authService) Login(email, password string) (string, *model.User, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", nil, err
	}

	if !utils.CheckPassword(password, user.Password) {
		return "", nil, err
	}

	token, err := config.GenerateToken(user.ID)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

func (s *authService) Register(req *model.RegisterRequest) (*model.User, error) {
	exists, err := s.userRepo.ExistsByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("error checking email: %v", err)
	}
	if exists {
		return nil, fmt.Errorf("email already exists")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %v", err)
	}

	user := &model.User{
		Name:       req.Name,
		Email:      req.Email,
		Password:   hashedPassword,
		Nim:        req.Nim,
		Angkatan:   req.Angkatan,
		TotalPoint: 0,
		TotalExp:   0,
		LevelId:    1,
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %v", err)
	}

	return user, nil
}
