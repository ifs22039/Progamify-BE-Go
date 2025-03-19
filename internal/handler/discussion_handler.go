package handler

import (
	"boysitorus/Progamify-Restful-API/internal/model"
	"boysitorus/Progamify-Restful-API/internal/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type DiscussionHandler struct {
	service service.DiscussionService
}

func NewDiscussionHandler(service service.DiscussionService) *DiscussionHandler {
	return &DiscussionHandler{service: service}
}

func (h *DiscussionHandler) GetById(c *gin.Context) {
	discussionID, err := strconv.Atoi(c.Param("discussionID"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid discussion ID"})
		return
	}

	discussion, err := h.service.GetById(uint(discussionID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, discussion)
}

func (h *DiscussionHandler) GetDiscussionsByLessonID(c *gin.Context) {
	lessonID, err := strconv.Atoi(c.Param("lessonID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lesson ID"})
		return
	}
	discussions, err := h.service.GetDiscussionsByLessonID(uint(lessonID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, discussions)
}

func (h *DiscussionHandler) CreateDiscussion(c *gin.Context) {
	userID := c.MustGet("userId").(uint)

	var discussion model.Discussion

	if err := c.ShouldBindJSON(&discussion); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Println(err.Error())
		return
	}

	discussion.UserID = userID

	if err := h.service.CreateDiscussion(&discussion); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err.Error())
		return
	}
	c.JSON(http.StatusCreated, discussion)
}

func (h *DiscussionHandler) CreateReply(c *gin.Context) {
	userID := c.MustGet("userId").(uint)

	var reply model.DiscReply

	if err := c.ShouldBindJSON(&reply); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Println(err.Error())
		return
	}

	reply.UserID = userID

	if err := h.service.AddReply(&reply); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err.Error())
		return
	}
	
	c.JSON(http.StatusCreated, reply)
}
