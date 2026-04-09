package handlers

import (
	"blogPlatform/config"
	"blogPlatform/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetTags(c *gin.Context) {
	var tags []models.Tag

	if err := config.DB.Find(&tags).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data"})
		return
	}
	c.JSON(http.StatusOK, tags)
}

func AddTagsToPost(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("id"))
	var post models.Post
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := config.DB.First(&post, postID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Post not found"})
		return
	}

	var input struct {
		TagIDs []uint `json:"tag_ids"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var tags []models.Tag
	if err := config.DB.Where("id IN ?", input.TagIDs).Find(&tags).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Tags not found"})
		return
	}

	if err := config.DB.Model(&post).Association("Tags").Append(&tags); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add tags"})
		return
	}

	if err := config.DB.Preload("Tags").First(&post, postID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated post"})
		return
	}

	c.JSON(http.StatusOK, post)
}

func CreateTags(c *gin.Context) {
	var tag models.Tag

	if err := c.ShouldBindJSON(&tag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	if err := config.DB.Create(&tag).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "Tag not created"})
		return
	}
	c.JSON(201, tag)
}
