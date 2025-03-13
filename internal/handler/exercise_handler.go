package handler

import (
	"boysitorus/Progamify-Restful-API/internal/model"
	"boysitorus/Progamify-Restful-API/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ExerciseHandler struct {
	exerciseService service.ExerciseService
}

func NewExerciseHandler(exerciseService service.ExerciseService) *ExerciseHandler {
	return &ExerciseHandler{exerciseService}
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
	// Return response ke client
	c.JSON(http.StatusOK, takeExercise)
}
