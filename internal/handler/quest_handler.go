package handler

import (
	"boysitorus/Progamify-Restful-API/internal/model"
	"boysitorus/Progamify-Restful-API/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type QuestHandler struct {
	questService service.QuestService
	badgeService service.BadgeService
}

func NewQuestHandler(questService service.QuestService, badgeService service.BadgeService) *QuestHandler {
	return &QuestHandler{questService, badgeService}
}

func (h *QuestHandler) GetQuest(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quest ID"})
		return
	}

	quest, err := h.questService.GetQuestById(uint(id))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Quest not found"})
		return
	}

	c.JSON(http.StatusOK, quest)
}

func (h *QuestHandler) SubmitQuest(c *gin.Context) {
	userId := c.MustGet("userId").(uint)
	var request model.SubmitQuestRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	takeQuest, err := h.questService.AddTakeQuest(userId, request)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit exercise"})
		return
	}

	badges, err := h.badgeService.CheckAndAssignBadge(userId)
	if err != nil {
	}

	response := gin.H{
		"quest": takeQuest,
	}

	if len(badges) > 0 {
		response["badge"] = badges
	}
	// Return response ke client
	c.JSON(http.StatusOK, response)
}
