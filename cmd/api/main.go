package main

import (
	"boysitorus/Progamify-Restful-API/internal/config"
	"boysitorus/Progamify-Restful-API/internal/handler"
	"boysitorus/Progamify-Restful-API/internal/repository"
	"boysitorus/Progamify-Restful-API/internal/routes"
	"boysitorus/Progamify-Restful-API/internal/service"
	"boysitorus/Progamify-Restful-API/pkg/database"
	"log"
)

func main() {
	// Load config
	cfg := config.LoadConfig()

	// Initialize DB
	db := database.InitDB(cfg)

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	topicRepo := repository.NewTopicRepository(db)
	lessonRepo := repository.NewLessonRepository(db, userRepo)
	exerciseRepo := repository.NewExerciseRepository(db, userRepo)
	questRepo := repository.NewQuestRepository(db, userRepo)
	discussionRepo := repository.NewDiscussionRepository(db)
	levelRepo := repository.NewLevelRepository(db)
	avatarRepo := repository.NewAvatarRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo)
	userService := service.NewUserService(userRepo)
	topicService := service.NewTopicService(topicRepo)
	lessonService := service.NewLessonService(lessonRepo)
	exerciseService := service.NewExerciseService(exerciseRepo)
	leaderboardService := service.NewLeaderboardService(userRepo)
	questService := service.NewQuestService(questRepo)
	discussionService := service.NewDiscussionService(discussionRepo)
	levelService := service.NewLevelService(levelRepo)
	avatarService := service.NewAvatarService(avatarRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	topicHandler := handler.NewTopicHandler(topicService)
	lessonHandler := handler.NewLessonHandler(lessonService)
	exerciseHandler := handler.NewExerciseHandler(exerciseService)
	leaderboardHandler := handler.NewLeaderboardHandler(leaderboardService)
	questHandler := handler.NewQuestHandler(questService)
	discussionHandler := handler.NewDiscussionHandler(discussionService)
	levelHandler := handler.NewLevelHandler(levelService)
	avatarHandler := handler.NewAvatarHandler(avatarService)

	// Initialize Gin router
	r := routes.SetupRoutes(authHandler,
		userHandler,
		topicHandler,
		lessonHandler,
		exerciseHandler,
		leaderboardHandler,
		questHandler,
		discussionHandler,
		levelHandler,
		avatarHandler,
	)

	log.Fatal(r.Run("0.0.0.0:8080"))
}
