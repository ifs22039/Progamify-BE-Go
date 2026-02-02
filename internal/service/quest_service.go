package service

import (
	"boysitorus/Progamify-Restful-API/internal/model"
	"boysitorus/Progamify-Restful-API/internal/repository"
)

type QuestService interface {
	GetQuestById(id uint) (*model.Quest, error)
	SubmitQuest(userID uint, request model.SubmitQuestRequest) (*model.TakeQuest, error)
}

type questService struct {
	questRepo repository.QuestRepository
	userRepo  repository.UserRepository
}

func NewQuestService(
	questRepo repository.QuestRepository,
	userRepo repository.UserRepository,
) QuestService {
	return &questService{
		questRepo: questRepo,
		userRepo:  userRepo,
	}
}

func (s *questService) GetQuestById(id uint) (*model.Quest, error) {
	return s.questRepo.GetQuestByID(id)
}

func (s *questService) SubmitQuest(
	userID uint,
	request model.SubmitQuestRequest,
) (*model.TakeQuest, error) {

	// 1️⃣ Ambil user (theta)
	user, err := s.userRepo.FindById(userID)
	if err != nil {
		return nil, err
	}

	// 2️⃣ Ambil quest (beta)
	quest, err := s.questRepo.GetQuestByID(request.QuestID)
	if err != nil {
		return nil, err
	}

	// 3️⃣ Grading + IRT (SUDAH BENAR ADA DI REPOSITORY)
	takeQuest, err := s.questRepo.GradeQuest(user, quest, request)
	if err != nil {
		return nil, err
	}

	// 4️⃣ Simpan hasil
	if err := s.questRepo.CreateTakeQuest(takeQuest); err != nil {
		return nil, err
	}

	return takeQuest, nil
}
