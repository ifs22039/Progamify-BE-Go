package repository

import (
	"boysitorus/Progamify-Restful-API/internal/model"
	"gorm.io/gorm"
)

type TopicProgressDTO struct {
	ID                 uint   `json:"id"`
	Name               string `json:"name"`
	TotalLessons       int64  `json:"total_lessons"`
	TotalTakeLessons   int64  `json:"total_take_lessons"`
	TotalExercises     int64  `json:"total_exercises"`
	TotalTakeExercises int64  `json:"total_take_exercises"`
}

type TopicRepository interface {
	FindAll() ([]model.Topic, error)
	FindById(id uint) (*model.Topic, error)
	FindTopicWithLessons(id uint) (*model.Topic, error)
	FindAllByUserId(userId uint) ([]TopicProgressDTO, error)
	LessonsTakenByUserId(topicId uint, userId uint) ([]model.TakeLesson, error)
	ExercisesTakenByUserId(topicId uint, userId uint) ([]model.TakeExercise, error)
}

type topicRepository struct {
	db *gorm.DB
}

func (r *topicRepository) ExercisesTakenByUserId(topicId uint, userId uint) ([]model.TakeExercise, error) {
	var takeExercise []model.TakeExercise
	err := r.db.Where("topic_id = ? AND user_id = ?", topicId, userId).Find(&takeExercise).Error
	return takeExercise, err
}

func NewTopicRepository(db *gorm.DB) TopicRepository {
	return &topicRepository{db}
}

func (r *topicRepository) FindById(id uint) (*model.Topic, error) {
	var topic model.Topic
	err := r.db.Preload("Lessons.Exercises.Questions").First(&topic, id).Error
	return &topic, err
}

func (r *topicRepository) LessonsTakenByUserId(topicId uint, userId uint) ([]model.TakeLesson, error) {
	var takeLessons []model.TakeLesson
	err := r.db.Where("topic_id = ? AND user_id = ?", topicId, userId).Find(&takeLessons).Error
	return takeLessons, err
}

func (r *topicRepository) FindAll() ([]model.Topic, error) {
	var topics []model.Topic
	err := r.db.Find(&topics).Error
	return topics, err
}

func (r *topicRepository) FindTopicWithLessons(id uint) (*model.Topic, error) {
	var topic model.Topic
	err := r.db.Preload("Lessons").First(&topic, id).Error
	return &topic, err
}

func (r *topicRepository) FindAllByUserId(userId uint) ([]TopicProgressDTO, error) {
	var results []TopicProgressDTO

	query := `
    SELECT 
        topics.id, 
        topics.name, 
        COUNT(DISTINCT lessons.id) AS total_lessons,
        COUNT(DISTINCT take_lessons.id) AS total_take_lessons,
        COUNT(DISTINCT exercises.id) AS total_exercises,
        COUNT(DISTINCT take_exercises.id) AS total_take_exercises
    FROM topics
		LEFT JOIN lessons ON topics.id = lessons.topic_id
		LEFT JOIN take_lessons ON topics.id = take_lessons.topic_id AND take_lessons.user_id = ?
		LEFT JOIN exercises ON topics.id = exercises.topic_id
		LEFT JOIN take_exercises ON topics.id = take_exercises.topic_id AND take_exercises.user_id = ?
    GROUP BY topics.id, topics.name
    `

	err := r.db.Raw(query, userId, userId).Scan(&results).Error

	return results, err
}
