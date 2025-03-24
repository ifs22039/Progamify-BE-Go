package handler

import (
	"boysitorus/Progamify-Restful-API/internal/model"
	"boysitorus/Progamify-Restful-API/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AchievementHandler struct {
	achievementService service.AchievementService
}
func NewAchievementHandler(achievementService service.AchievementService) *AchievementHandler {
	return &AchievementHandler{achievementService}
}

func (h *AchievementHandler) GetAchievement(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	achievement, err := h.achievementService.GetAchievement(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Achievement not found for this id"})
		return
	}

	c.JSON(http.StatusOK, achievement)
}

func (h *AchievementHandler) AddHaveAchievement(c *gin.Context) {
	userId := c.MustGet("userId").(uint)

	var request model.AddAchievementRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	haveAchievement, err := h.achievementService.AddHaveAchievement(userId, request.AchievementID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add achievement"})
		return
	}

	c.JSON(http.StatusOK, haveAchievement)

}

func (h *AchievementHandler) AssignAchievementIfEligible(c *gin.Context) {

	userId := c.MustGet("userId").(uint)

	achievement, err := h.achievementService.CheckAndAssignAchievement(userId)
	if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign achievement"})
			return
	}

	if achievement == nil {
			c.JSON(http.StatusOK, gin.H{"message": "No achievement assigned"})
			return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Achievement assigned", "achievement": achievement})
}


func (h *AchievementHandler) GetAchievements(c *gin.Context) {
	userId := c.MustGet("userId").(uint)

	achievement, err := h.achievementService.GetAchievements(uint(userId))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Achievements not found for this user"})
		return
	}

	c.JSON(http.StatusOK, achievement)
}


func (h *AchievementHandler) GetAchievementsByUserId(c *gin.Context) {
	userId, err := strconv.ParseUint(c.Param("userId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	achievements, err := h.achievementService.GetAchievements(uint(userId))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Achievements not found for this user"})
		return
	}

	c.JSON(http.StatusOK, achievements)
}
