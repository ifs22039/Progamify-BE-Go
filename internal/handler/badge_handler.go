package handler

import (
	"boysitorus/Progamify-Restful-API/internal/model"
	"boysitorus/Progamify-Restful-API/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BadgeHandler struct {
	badgeService service.BadgeService
}
func NewBadgeHandler(badgeService service.BadgeService) *BadgeHandler {
	return &BadgeHandler{badgeService}
}

func (h *BadgeHandler) GetBadge(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	badge, err := h.badgeService.GetBadge(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Badge not found for this id"})
		return
	}

	c.JSON(http.StatusOK, badge)
}

func (h *BadgeHandler) AddHaveBadge(c *gin.Context) {
	userId := c.MustGet("userId").(uint)

	var request model.AddBadgeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	haveBadge, err := h.badgeService.AddHaveBadge(userId, request.BadgeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add badge"})
		return
	}

	c.JSON(http.StatusOK, haveBadge)

}

func (h *BadgeHandler) AssignBadgeIfEligible(c *gin.Context) {

	userId := c.MustGet("userId").(uint)

	badge, err := h.badgeService.CheckAndAssignBadge(userId)
	if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign badge"})
			return
	}

	if badge == nil {
			c.JSON(http.StatusOK, gin.H{"message": "No badge assigned"})
			return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Badge assigned", "badge": badge})
}


func (h *BadgeHandler) GetBadges(c *gin.Context) {
	userId := c.MustGet("userId").(uint)

	badge, err := h.badgeService.GetBadges(uint(userId))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Badges not found for this user"})
		return
	}

	c.JSON(http.StatusOK, badge)
}
