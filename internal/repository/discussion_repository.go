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
	//if err := r.db.Table("discussions").
	//	Select("users.name, DATE_FORMAT(discussions.created_at, '%d/%m/%y %H:%i') as date, discussions.title, discussions.content").
	//	Joins("join users on users.id = discussions.user_id").
	//	Where("lesson_id = ?", lessonID).
	//	Order("date desc").
	//	Scan(&discussions).Error; err != nil {
	//	return nil, err
	//}
	//return discussions, nil

	if err := r.db.Table("discussions").
		Select("users.name, DATE_FORMAT(discussions.created_at, '%d/%m/%y %H:%i') as date, discussions.title, discussions.content, discussions.id").
		Joins("JOIN users ON users.id = discussions.user_id").
		Where("lesson_id = ?", lessonID).
		Order("date DESC, id DESC").
		Find(&discussions).Error; err != nil {
		return nil, err
	}

	// Load replies for each discussion
	for i, discussion := range discussions {
		var replies []map[string]interface{}
		if err := r.db.Table("disc_replies").
			Select("users.name, DATE_FORMAT(disc_replies.created_at, '%d/%m/%y %H:%i') as date, disc_replies.content").
			Joins("JOIN users ON users.id = disc_replies.user_id").
			Where("discussion_id = ?", discussion["id"]).
			Order("date ASC").
			Scan(&replies).Error; err != nil {
			return nil, err
		}
		discussions[i]["replies"] = replies
	}

	return discussions, nil
}

func (r *discussionRepository) Create(discussion *model.Discussion) error {
	return r.db.Create(discussion).Error
}
