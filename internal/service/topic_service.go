package service

import (
	"boysitorus/Progamify-Restful-API/internal/model"
	"boysitorus/Progamify-Restful-API/internal/repository"
)

type TopicService interface {
	GetTopicById(id uint) (*model.Topic, error)
	GetAllTopics() ([]model.Topic, error)
	GetTopicWithLessons(id uint) (*model.Topic, error)
	GetAllTopicsByUserId(userId uint) ([]repository.TopicProgressDTO, error)
	GetLessonsTakenByUserId(topicId uint, userId uint) ([]model.TakeLesson, error)
}

// topicService struct
type topicService struct {
	topicRepo repository.TopicRepository
}

func (s *topicService) GetLessonsTakenByUserId(topicId uint, userId uint) ([]model.TakeLesson, error) {
	return s.topicRepo.LessonsTakenByUserId(topicId, userId)
}

func NewTopicService(topicRepo repository.TopicRepository) TopicService {
	return &topicService{topicRepo}
}

func (s *topicService) GetTopicById(id uint) (*model.Topic, error) {
	return s.topicRepo.FindById(id)
}

func (s *topicService) GetAllTopics() ([]model.Topic, error) {
	return s.topicRepo.FindAll()
}

func (s *topicService) GetTopicWithLessons(id uint) (*model.Topic, error) {
	return s.topicRepo.FindTopicWithLessons(id)
}

func (s *topicService) GetAllTopicsByUserId(userId uint) ([]repository.TopicProgressDTO, error) {
	return s.topicRepo.FindAllByUserId(userId)
}
