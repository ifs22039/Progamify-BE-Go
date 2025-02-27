package handler

import (
	"boysitorus/Progamify-Restful-API/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type TopicHandler struct {
	topicService service.TopicService
}

func NewTopicHandler(topicService service.TopicService) *TopicHandler {
	return &TopicHandler{topicService}
}

func (h *TopicHandler) GetTopic(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid topic ID"})
		return
	}

	topic, err := h.topicService.GetTopicById(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Topic not found"})
		return
	}

	c.JSON(http.StatusOK, topic)
}

func (h *TopicHandler) GetTopicWithLessons(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid topic ID"})
		return
	}

	topic, err := h.topicService.GetTopicWithLessons(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Topic not found"})
		return
	}

	c.JSON(http.StatusOK, topic)
}

//func (h *TopicHandler) ListTopics(c *gin.Context) {
//	topics, err := h.topicService.GetAllTopics()
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	c.JSON(http.StatusOK, topics)
//}

func (h *TopicHandler) ListTopics(c *gin.Context) {
	userId := c.MustGet("userId").(uint)

	// Get topics with progress
	topics, err := h.topicService.GetAllTopicsByUserId(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, topics)
}
