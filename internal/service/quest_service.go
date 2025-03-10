package service
import (
	"boysitorus/Progamify-Restful-API/internal/model"
	"boysitorus/Progamify-Restful-API/internal/repository"
)

type QuestService interface {
	GetQuestById(id uint) (*model.Quest, error)
}


type questService struct {
	questRepo repository.QuestRepository
}

func (e *questService) GetQuestById(id uint) (*model.Quest, error) {
	return e.questRepo.GetQuestByUserID(id)
}

func NewQuestService(questRepo repository.QuestRepository) QuestService {
	return &questService{questRepo}
}
