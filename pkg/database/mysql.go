package database

import (
	"boysitorus/Progamify-Restful-API/internal/config"
	"boysitorus/Progamify-Restful-API/internal/model"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB(cfg *config.Config) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Auto migrate
	db.AutoMigrate(
		&model.User{},
		&model.TakeExercise{},
		&model.Exercise{},
		&model.Lesson{},
		&model.Topic{},
		&model.Quest{},
		&model.TakeQuest{},
		&model.TakeLesson{},
		&model.Level{},
		&model.Avatar{},
		&model.Badge{},
		&model.Gift{},
		&model.Achievement{},
		&model.HaveAchievement{},
		&model.HaveAvatar{},
		&model.HaveBadge{},
		&model.HaveGift{},
		&model.Discussion{},
		&model.DiscReply{},
		&model.ExQuestion{},
		&model.ExAnswer{},
		&model.QuestAnswer{},
	)

	return db
}
