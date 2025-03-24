package service

import (
	"boysitorus/Progamify-Restful-API/internal/model"
	"boysitorus/Progamify-Restful-API/internal/repository"
)

type DiscussionService interface {
	GetDiscussionsByLessonID(lessonID uint) ([]model.Discussion, error)
	CreateDiscussion(discussion *model.Discussion) error
	GetById(discussionID uint) (*model.Discussion, error)
	AddReply(reply *model.DiscReply) error
}

type discussionService struct {
	repo repository.DiscussionRepository
}

func NewDiscussionService(repo repository.DiscussionRepository) DiscussionService {
	return &discussionService{repo: repo}
}

func (s *discussionService) AddReply(reply *model.DiscReply) error {
	return s.repo.AddReply(reply)
}

func (s *discussionService) GetById(discussionID uint) (*model.Discussion, error) {
	return s.repo.FindById(discussionID)
}

func (s *discussionService) GetDiscussionsByLessonID(lessonID uint) ([]model.Discussion, error) {
	return s.repo.FindDiscussionByLessonID(lessonID)
}

func (s *discussionService) CreateDiscussion(discussion *model.Discussion) error {
	return s.repo.Create(discussion)
}
