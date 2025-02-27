package main

import (
	"boysitorus/Progamify-Restful-API/internal/config"
	"boysitorus/Progamify-Restful-API/internal/handler"
	"boysitorus/Progamify-Restful-API/internal/middleware"
	"boysitorus/Progamify-Restful-API/internal/repository"
	"boysitorus/Progamify-Restful-API/internal/service"
	"boysitorus/Progamify-Restful-API/pkg/database"
	"github.com/gin-gonic/gin"
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

	// Initialize services
	authService := service.NewAuthService(userRepo)
	userService := service.NewUserService(userRepo)
	topicService := service.NewTopicService(topicRepo)
	lessonService := service.NewLessonService(lessonRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	topicHandler := handler.NewTopicHandler(topicService)
	lessonHandler := handler.NewLessonHandler(lessonService)

	// Initialize Gin router
	r := gin.Default()

	api := r.Group("/api")

	api.POST("/users/login", authHandler.Login)
	api.POST("/users/register", authHandler.Register)

	api.Use(middleware.AuthMiddleware())
	{
		api.GET("/users/current", userHandler.GetCurrentUser)
		api.PUT("/users/current", userHandler.UpdateUser)
		api.DELETE("/users/logout", authHandler.Logout)

		api.GET("/topics", topicHandler.ListTopics)
		api.GET("/topics/:id", topicHandler.GetTopic)

		api.GET("/lessons/:id", lessonHandler.GetLesson)
	}

	log.Fatal(r.Run(":8080"))
}
