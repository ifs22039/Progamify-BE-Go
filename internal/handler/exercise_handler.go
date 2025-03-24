package handler

import (
	"boysitorus/Progamify-Restful-API/internal/model"
	"boysitorus/Progamify-Restful-API/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ExerciseHandler struct {
	exerciseService service.ExerciseService
	achievementService service.AchievementService
}

func NewExerciseHandler(exerciseService service.ExerciseService, achievementService service.AchievementService) *ExerciseHandler {
	return &ExerciseHandler{exerciseService, achievementService}
}

func (h *ExerciseHandler) GetExercise(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid exercise ID"})
		return
	}

	exercise, err := h.exerciseService.GetExerciseById(uint(id))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exercise not found"})
		return
	}

	c.JSON(http.StatusOK, exercise)
}

func (h *ExerciseHandler) SubmitExercise(c *gin.Context) {
	userId := c.MustGet("userId").(uint)
	var request model.SubmitExerciseRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	takeExercise, err := h.exerciseService.AddTakeExercise(userId, request)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit exercise"})
		return
	}

	mergedExercise := map[string]interface{}{
		"ID":         	takeExercise.ID,
		"CreatedAt":  	takeExercise.CreatedAt,
		"UpdatedAt":  	takeExercise.UpdatedAt,
		"DeletedAt":  	takeExercise.DeletedAt,
		"LessonID":     takeExercise.LessonID,
		"TopicID":    	takeExercise.TopicID,
		"Answers":      takeExercise.Answers,
		"Score": 				takeExercise.Score,
		"TotalCorrect": takeExercise.TotalCorrect,
		"TotalQuestion":takeExercise.TotalQuestion,
		"TotalExp": 		takeExercise.TotalExp,
		"TotalPoint":   takeExercise.TotalPoint,
		"RewardExp": 		takeExercise.RewardExp,
		"RewardPoint": 	takeExercise.RewardPoint,
	}

	achievements, err := h.achievementService.CheckAndAssignAchievement(userId)
	if err != nil {
	}
	if len(achievements) > 0 {
		mergedExercise["achievement"] = achievements
	}
	
	c.JSON(http.StatusOK, mergedExercise)
}
