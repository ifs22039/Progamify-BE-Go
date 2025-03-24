package handler

import (
	"boysitorus/Progamify-Restful-API/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type GiftHandler struct {
	service service.GiftService
}

func NewGiftHandler(service service.GiftService) *GiftHandler {
	return &GiftHandler{service: service}
}

func (h *GiftHandler) GetAll(c *gin.Context) {
	gifts, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gifts)
}

func (h *GiftHandler) BuyGift(c *gin.Context) {
	userID := c.MustGet("userId").(uint)

	var request struct {
		GiftID uint `json:"gift_id"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	haveAvatar, err := h.service.BuyGift(uint(userID), uint(request.GiftID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, haveAvatar)
}

func (h *GiftHandler) GetUserGift(c *gin.Context) {
	userID := c.MustGet("userId").(uint)
	userGifts, err := h.service.GetUserGift(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, userGifts)
}
