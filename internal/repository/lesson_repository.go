package repository

import (
	"boysitorus/Progamify-Restful-API/internal/model"
	"boysitorus/Progamify-Restful-API/pkg/utils"
	"errors"
	"gorm.io/gorm"
	"strings"
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

		var user model.User
		l.db.Where("id = ?", userID).First(&user)

		var lesson model.Lesson
		l.db.Where("id = ?", lessonID).First(&lesson)

		user.TotalExp = user.TotalExp + lesson.Exp

		err = l.db.Save(&user).Error

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

	if err != nil {
		// Return the error early if the record is not found
		return nil, err
	}

	// Ensure that lesson.Content is not empty before attempting to convert
	if lesson.Content != "" {
		// Convert the HTML content
		html := utils.ConvertHTML(lesson.Content)
		html = strings.ReplaceAll(html, "\\\"", "\"")
		lesson.Content = html
	}

	return &lesson, err
}

func NewLessonRepository(db *gorm.DB) LessonRepository {
	return &lessonRepository{db}
}
