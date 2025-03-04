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
	lessonRepo := repository.NewLessonRepository(db)
	exerciseRepo := repository.NewExerciseRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo)
	userService := service.NewUserService(userRepo)
	topicService := service.NewTopicService(topicRepo)
	lessonService := service.NewLessonService(lessonRepo)
	exerciseService := service.NewExerciseService(exerciseRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	topicHandler := handler.NewTopicHandler(topicService)
	lessonHandler := handler.NewLessonHandler(lessonService)
	exerciseHandler := handler.NewExerciseHandler(exerciseService)

	// Initialize Gin router
	r := routes.SetupRoutes(authHandler, userHandler, topicHandler, lessonHandler, exerciseHandler)

	log.Fatal(r.Run(":8080"))
}
