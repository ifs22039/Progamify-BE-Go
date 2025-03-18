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
	db       *gorm.DB
	userRepo UserRepository
}

func NewLessonRepository(db *gorm.DB, userRepo UserRepository) LessonRepository {
	return &lessonRepository{db, userRepo}
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

		var lesson model.Lesson
		l.db.Where("id = ?", lessonID).First(&lesson)

		err = l.userRepo.AddExp(userID, lesson.Exp)

		if err != nil {

			return nil, err
		}

		err = l.db.Where("lesson_id = ? AND user_id = ?", lessonID, userID).First(&takeLesson).Error

		return &takeLesson, err
	}

	//tidak mengembalikan apa apa jika sudah mengambil lesson
	return nil, nil
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
