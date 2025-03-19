package handler

import (
	"boysitorus/Progamify-Restful-API/internal/model"
	"boysitorus/Progamify-Restful-API/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService}
}

func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	userId := c.MustGet("userId").(uint)

	user, err := h.userService.GetUser(userId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	lessonTaken, err := h.userService.GetLessonTaken(userId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "error finding lesson taken for user"})
		return
	}

	response := map[string]interface{}{
		"ID":                 user.ID,
		"CreatedAt":          user.CreatedAt,
		"UpdatedAt":          user.UpdatedAt,
		"DeletedAt":          user.DeletedAt,
		"name":               user.Name,
		"email":              user.Email,
		"nim":                user.Nim,
		"angkatan":           user.Angkatan,
		"total_point":        user.TotalPoint,
		"total_exp":          user.TotalExp,
		"level_id":           user.LevelId,
		"total_lesson_taken": lessonTaken,
	}

	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	userId := c.MustGet("userId").(uint)

	var req model.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.UpdateUser(userId, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		return
	}

	c.JSON(http.StatusOK, user)
}
