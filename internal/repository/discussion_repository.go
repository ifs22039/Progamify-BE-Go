package repository

import (
	"boysitorus/Progamify-Restful-API/internal/model"
	"gorm.io/gorm"
)

type DiscussionRepository interface {
	FindDiscussionByLessonID(lessonID uint) ([]map[string]interface{}, error)
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
	err := r.db.Preload("Replies.DetailUser").First(&discussion, discussionID).Error

	if err != nil {
		// Return the error early if the record is not found
		return nil, err
	}

	return &discussion, nil
}

func (r *discussionRepository) FindDiscussionByLessonID(lessonID uint) ([]map[string]interface{}, error) {
	var discussions []map[string]interface{}

	// Fetch discussions with necessary fields
	if err := r.db.Table("discussions").
		Select("users.name, DATE_FORMAT(discussions.created_at, '%d/%m/%y %H:%i') as date, discussions.title, discussions.content, discussions.id").
		Joins("JOIN users ON users.id = discussions.user_id").
		Where("lesson_id = ?", lessonID).
		Order("date DESC, id DESC").
		Find(&discussions).Error; err != nil {
		return nil, err
	}

	for i, discussion := range discussions {
		var replyCount int64
		if err := r.db.Table("disc_replies").
			Where("discussion_id = ?", discussion["id"]).
			Count(&replyCount).Error; err != nil {
			return nil, err
		}
		discussions[i]["replies"] = replyCount
	}

	return discussions, nil
}

func (r *discussionRepository) Create(discussion *model.Discussion) error {
	return r.db.Create(discussion).Error
}
