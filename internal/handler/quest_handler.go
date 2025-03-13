package handler

import (
	"boysitorus/Progamify-Restful-API/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type QuestHandler struct {
	questService service.QuestService
}

func NewQuestHandler(questService service.QuestService) *QuestHandler {
	return &QuestHandler{questService}
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
