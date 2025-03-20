package handler

import (
	"boysitorus/Progamify-Restful-API/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AvatarHandler struct {
	service service.AvatarService
}

func NewAvatarHandler(service service.AvatarService) *AvatarHandler {
	return &AvatarHandler{service: service}
}

func (h *AvatarHandler) GetAll(c *gin.Context) {
	userID := c.MustGet("userId").(uint)

	avatars, err := h.service.GetAvatarsByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, avatars)
}

func (h *AvatarHandler) BuyAvatar(c *gin.Context) {
	userID := c.MustGet("userId").(uint)

	var request struct {
		AvatarID uint `json:"avatar_id"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	haveAvatar, err := h.service.BuyAvatar(uint(userID), uint(request.AvatarID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, haveAvatar)
}
