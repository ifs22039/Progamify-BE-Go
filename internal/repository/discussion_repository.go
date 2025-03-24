package repository

import (
	"boysitorus/Progamify-Restful-API/internal/model"
	"gorm.io/gorm"
)

type DiscussionRepository interface {
	FindDiscussionByLessonID(lessonID uint) ([]model.Discussion, error)
	Create(discussion *model.Discussion) error
	FindById(discussionID uint) (*model.Discussion, error)
	AddReply(reply *model.DiscReply) error
}

type discussionRepository struct {
	db *gorm.DB
}

func NewDiscussionRepository(db *gorm.DB) DiscussionRepository {
	return &discussionRepository{db: db}
}

func (r *discussionRepository) AddReply(reply *model.DiscReply) error {
	return r.db.Create(reply).Error
}

func (r *discussionRepository) FindById(discussionID uint) (*model.Discussion, error) {
	var discussion model.Discussion
	err := r.db.Preload("User.Avatar").Preload("Replies.DetailUser.Avatar").First(&discussion, discussionID).Error

	if err != nil {
		// Return the error early if the record is not found
		return nil, err
	}

	return &discussion, nil
}

func (r *discussionRepository) FindDiscussionByLessonID(lessonID uint) ([]model.Discussion, error) {
	var discussions []model.Discussion
	err := r.db.Preload("User.Avatar").Preload("Replies").Where("lesson_id", lessonID).Find(&discussions).Error
	return discussions, err
}

func (r *discussionRepository) Create(discussion *model.Discussion) error {
	return r.db.Create(discussion).Error
}
