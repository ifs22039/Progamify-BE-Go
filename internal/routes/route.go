package routes

import (
	"boysitorus/Progamify-Restful-API/internal/handler"
	"boysitorus/Progamify-Restful-API/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all routes for the application
func SetupRoutes(
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	topicHandler *handler.TopicHandler,
	lessonHandler *handler.LessonHandler,
	exerciseHandler *handler.ExerciseHandler,
	leaderboardHandler *handler.LeaderboardHandler,
	questHandler *handler.QuestHandler,
	discussionHandler *handler.DiscussionHandler,
) *gin.Engine {
	// Initialize Gin router
	r := gin.Default()

	api := r.Group("/api")

	api.POST("/users/login", authHandler.Login)
	api.POST("/users/register", authHandler.Register)
	api.GET("/customrender", authHandler.CustomRender)

	api.Use(middleware.AuthMiddleware())
	{
		api.GET("/users/current", userHandler.GetCurrentUser)
		api.PUT("/users/current", userHandler.UpdateUser)
		api.DELETE("/users/logout", authHandler.Logout)

		api.GET("/topics", topicHandler.ListTopics)
		api.GET("/topics/:id", topicHandler.GetTopic)

		api.GET("/lessons/:id", lessonHandler.GetLesson)

		api.GET("/exercises/:id", exerciseHandler.GetExercise)
		api.POST("/exercises/submit", exerciseHandler.SubmitExercise)

		api.GET("/leaderboard", leaderboardHandler.GetLeaderboard)

		api.GET("/quest/:id", questHandler.GetQuest)
		api.POST("/quest/submit", questHandler.SubmitQuest)

		api.POST("/discussions", discussionHandler.CreateDiscussion)
		api.GET("/discussions/:lessonID", discussionHandler.GetDiscussionsByLessonID)
	}

	return r
}
