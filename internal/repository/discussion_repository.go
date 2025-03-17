package repository

import (
	"boysitorus/Progamify-Restful-API/internal/model"
	"gorm.io/gorm"
)

type DiscussionRepository interface {
	FindDiscussionByLessonID(lessonID uint) ([]map[string]interface{}, error)
	Create(discussion *model.Discussion) error
}

type discussionRepository struct {
	db *gorm.DB
}

func NewDiscussionRepository(db *gorm.DB) DiscussionRepository {
	return &discussionRepository{db: db}
}

func (r *discussionRepository) FindDiscussionByLessonID(lessonID uint) ([]map[string]interface{}, error) {
	var discussions []map[string]interface{}
	if err := r.db.Table("discussions").
		Select("users.name, DATE_FORMAT(discussions.created_at, '%d/%m/%y %H:%i') as date, discussions.title, discussions.content").
		Joins("join users on users.id = discussions.user_id").
		Where("lesson_id = ?", lessonID).
		Order("date desc").
		Scan(&discussions).Error; err != nil {
		return nil, err
	}
	return discussions, nil
}

func (r *discussionRepository) Create(discussion *model.Discussion) error {
	return r.db.Create(discussion).Error
}
