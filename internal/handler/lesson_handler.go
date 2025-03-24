package handler

import (
	"boysitorus/Progamify-Restful-API/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type LessonHandler struct {
	lessonService service.LessonService
	achievementService service.AchievementService
}

func NewLessonHandler(lessonService service.LessonService, achievementService service.AchievementService) *LessonHandler {
	return &LessonHandler{lessonService, achievementService}
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

	takeLesson, err := h.lessonService.AddTakeLesson(lesson.TopicID, lesson.ID, userID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to record take lesson"})
		return
	}

	if takeLesson != nil {
		mergedLesson := map[string]interface{}{
			"ID":         lesson.ID,
			"CreatedAt":  lesson.CreatedAt,
			"UpdatedAt":  lesson.UpdatedAt,
			"DeletedAt":  lesson.DeletedAt,
			"name":       lesson.Name,
			"content":    lesson.Content,
			"exp":        lesson.Exp,
			"takeLesson": takeLesson,
		}

		achievements, err := h.achievementService.CheckAndAssignAchievement(userID)
		if err != nil {
		}
		if len(achievements) > 0 {
			mergedLesson["achievement"] = achievements
		}

		c.JSON(http.StatusOK, mergedLesson)
		return
	}

	c.JSON(http.StatusOK, lesson)
}
