package service

import (
	"boysitorus/Progamify-Restful-API/internal/model"
	"boysitorus/Progamify-Restful-API/internal/repository"
)

type LessonService interface {
	GetLessonById(id uint) (*model.Lesson, error)
	AddTakeLesson(topicID uint, lessonID uint, userID uint) (*model.TakeLesson, error)
}

type lessonService struct {
	lessonRepo repository.LessonRepository
}

func (l *lessonService) AddTakeLesson(topicID uint, lessonID uint, userID uint) (*model.TakeLesson, error) {
	return l.lessonRepo.AddTakeLesson(topicID, lessonID, userID)
}

func (l *lessonService) GetLessonById(id uint) (*model.Lesson, error) {
	return l.lessonRepo.FindById(id)
}

func NewLessonService(lessonRepo repository.LessonRepository) LessonService {
	return &lessonService{lessonRepo}
}
