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
) *gin.Engine {
	// Initialize Gin router
	r := gin.Default()

	// Public routes
	setupPublicRoutes(r, authHandler)

	// Protected routes
	setupProtectedRoutes(r, userHandler, authHandler)

	return r
}

// setupPublicRoutes defines routes that don't require authentication
func setupPublicRoutes(r *gin.Engine, authHandler *handler.AuthHandler) {
	r.POST("/login", authHandler.Login)
	r.POST("/register", authHandler.Register)
}

// setupProtectedRoutes defines routes that require authentication
func setupProtectedRoutes(
	r *gin.Engine,
	userHandler *handler.UserHandler,
	authHandler *handler.AuthHandler,
) {
	// Create a route group with authentication middleware
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		// User routes
		api.GET("/me", userHandler.GetCurrentUser)
		api.PUT("/me", userHandler.UpdateUser)

		// Auth routes
		api.POST("/logout", authHandler.Logout)
	}
}
