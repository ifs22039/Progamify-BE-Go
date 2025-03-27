package handler

import (
	"boysitorus/Progamify-Restful-API/internal/model"
	"boysitorus/Progamify-Restful-API/internal/service"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
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

	err = h.userService.UpdateLastLogin(userId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "failed to update last login"})
		return
	}

	fmt.Println(user)

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
		"avatar_id":          user.AvatarID,
		"detail_avatar":      user.Avatar,
		"total_lesson_taken": lessonTaken,
	}

	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	userId := c.MustGet("userId").(uint)

	var request struct {
		Name     string `json:"name" binding:"required"`                               // Nama minimal 8 karakter
		Nim      string `json:"nim" binding:"required"`                                // Nim wajib diisi
		Angkatan int    `json:"angkatan" binding:"required,numeric,min=1000,max=9999"` // Angkatan harus 4 digit
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.UpdateUser(userId, request.Name, request.Nim, request.Angkatan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user updated successfully", "user": user})
}

func (h *UserHandler) UpdatePassword(c *gin.Context) {
	userId := c.MustGet("userId").(uint)

	var req model.UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.UpdatePassword(userId, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) ChangeAvatar(c *gin.Context) {
	userID := c.MustGet("userId").(uint)

	var request struct {
		AvatarID uint `json:"avatar_id"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Call the repository function to change avatar
	user, err := h.userService.ChangeAvatar(userID, request.AvatarID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	// Respond with updated user
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) GetUserLessonTaken(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("userId"), 10, 64)

	lessonTaken, err := h.userService.GetLessonTaken(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "error finding lesson taken for user"})
		return
	}

	response := map[string]int{
		"total_lesson_taken": lessonTaken,
	}

	c.JSON(http.StatusOK, response)
}
