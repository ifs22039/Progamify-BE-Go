package handler

import (
	"boysitorus/Progamify-Restful-API/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type LeaderboardHandler struct {
	leaderboardService service.LeaderboardService
}

func NewLeaderboardHandler(leaderboardService service.LeaderboardService) *LeaderboardHandler {
	return &LeaderboardHandler{leaderboardService}
}

func (h *LeaderboardHandler) GetLeaderboard(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
		return
	}

	users, err := h.leaderboardService.GetLeaderboard(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch leaderboard"})
		return
	}

	c.JSON(http.StatusOK, users)
}