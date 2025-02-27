package repository

import (
	"boysitorus/Progamify-Restful-API/internal/model"
	"errors"
	"gorm.io/gorm"
)

type LessonRepository interface {
	FindById(id uint) (*model.Lesson, error)
	AddTakeLesson(topicID uint, lessonID uint, userID uint) (*model.TakeLesson, error)
}

type lessonRepository struct {
	db *gorm.DB
}

func (l *lessonRepository) AddTakeLesson(topicID uint, lessonID uint, userID uint) (*model.TakeLesson, error) {
	var takeLesson model.TakeLesson
	err := l.db.Where("lesson_id = ? AND user_id = ?", lessonID, userID).First(&takeLesson).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		newTakeLesson := model.TakeLesson{
			LessonID: lessonID,
			TopicID:  topicID,
			UserID:   userID,
		}
		err = l.db.Create(&newTakeLesson).Error
		if err != nil {
			return nil, err
		}
		err = l.db.Where("lesson_id = ? AND user_id = ?", lessonID, userID).First(&takeLesson).Error
		return &takeLesson, err
	}
	return &takeLesson, err
}

func (l *lessonRepository) FindById(id uint) (*model.Lesson, error) {
	var lesson model.Lesson
	err := l.db.First(&lesson, id).Error
	return &lesson, err
}

func NewLessonRepository(db *gorm.DB) LessonRepository {
	return &lessonRepository{db}
}
