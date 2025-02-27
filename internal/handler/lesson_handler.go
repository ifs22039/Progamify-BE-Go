package handler

import (
	"boysitorus/Progamify-Restful-API/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type LessonHandler struct {
	lessonService service.LessonService
}

func NewLessonHandler(lessonService service.LessonService) *LessonHandler {
	return &LessonHandler{lessonService}
}

func (h *LessonHandler) GetLesson(c *gin.Context) {
	userID := c.MustGet("userId").(uint)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lesson ID"})
		return
	}

	lesson, err := h.lessonService.GetLessonById(uint(id))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Lesson not found"})
		return
	}

	_, err = h.lessonService.AddTakeLesson(lesson.TopicID, lesson.ID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to record take lesson"})
		return
	}

	c.JSON(http.StatusOK, lesson)
}
