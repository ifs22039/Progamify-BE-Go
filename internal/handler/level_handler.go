package handler

import (
	"boysitorus/Progamify-Restful-API/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LevelHandler struct {
	levelService service.LevelService
}
func NewLevelHandler(levelService service.LevelService) *LevelHandler {
	return &LevelHandler{levelService}
}

func (h *LevelHandler) GetLevel(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid level ID"})
		return
	}
	level, err := h.levelService.GetLevel(uint(id))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "level not found"})
		return
	}

	c.JSON(http.StatusOK, level)
}

func (h *LevelHandler) GetLevelByUserId(c *gin.Context) {
	userId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	level, err := h.levelService.GetLevelByUserId(uint(userId))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Level not found for this user"})
		return
	}

	c.JSON(http.StatusOK, level)
}