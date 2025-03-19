package model

import "gorm.io/gorm"

type DiscReply struct {
	gorm.Model
	DiscussionID uint   `json:"discussion_id" gorm:""`
	UserID       uint   `json:"user_id" gorm:"size:255;not null"`
	Content      string `json:"content" gorm:"not null"`
	DetailUser   User   `json:"detail_user,omitempty" gorm:"foreignkey:UserID"`
}

func (DiscReply) TableName() string {
	return "disc_replies"
}
