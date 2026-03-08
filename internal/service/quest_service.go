package service

import (
	"boysitorus/Progamify-Restful-API/internal/model"
	"boysitorus/Progamify-Restful-API/internal/repository"
)

type QuestService interface {
	GetQuestById(id uint) (*model.Quest, error)
	AddTakeQuest(userID uint, request model.SubmitQuestRequest) (*model.TakeQuest, error)
}

type questService struct {
	questRepo repository.QuestRepository
}

func (e *questService) GetQuestById(id uint) (*model.Quest, error) {
	return e.questRepo.GetQuestByID(id)
}

func (e *questService) AddTakeQuest(userID uint, request model.SubmitQuestRequest) (*model.TakeQuest, error) {
	return e.questRepo.AddTakeQuest(userID, request)
}

func NewQuestService(questRepo repository.QuestRepository) QuestService {
	return &questService{questRepo}
}
