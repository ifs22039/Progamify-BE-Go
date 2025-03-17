package service

import (
	"boysitorus/Progamify-Restful-API/internal/model"
	"boysitorus/Progamify-Restful-API/internal/repository"
)

type DiscussionService interface {
	GetDiscussionsByLessonID(lessonID uint) ([]map[string]interface{}, error)
	CreateDiscussion(discussion *model.Discussion) error
}

type discussionService struct {
	repo repository.DiscussionRepository
}

func NewDiscussionService(repo repository.DiscussionRepository) DiscussionService {
	return &discussionService{repo: repo}
}

func (s *discussionService) GetDiscussionsByLessonID(lessonID uint) ([]map[string]interface{}, error) {
	return s.repo.FindDiscussionByLessonID(lessonID)
}

func (s *discussionService) CreateDiscussion(discussion *model.Discussion) error {
	return s.repo.Create(discussion)
}
