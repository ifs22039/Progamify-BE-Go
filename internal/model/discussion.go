package model

import "gorm.io/gorm"

type Discussion struct {
	gorm.Model
	LessonID uint        `json:"lesson_id" gorm:";not null"`
	UserID   uint        `json:"user_id" gorm:"size:255;not null"`
	Title    string      `json:"title" gorm:"not null"`
	Content  string      `json:"content" gorm:"not null"`
	Replies  []DiscReply `json:"replies,omitempty" gorm:"foreignKey:DiscussionID"`
}
