package handlers

import (
	"blogPlatform/config"
	"blogPlatform/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddCommentToPost(c *gin.Context) {
	var comment models.Comment
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	comment.PostID = uint(postID)
	comment.AuthorID = c.GetUint("user_id")

	if err := config.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	var post models.Post
	if err := config.DB.Preload("Author").First(&post, postID).Error; err == nil {
		var commenter models.Author
		if err := config.DB.First(&commenter, comment.AuthorID).Error; err == nil {
			go SendNotification(NotifyPayload{
				AuthorID:      post.AuthorID,
				AuthorEmail:   post.Author.Email,
				AuthorName:    post.Author.Username,
				PostTitle:     post.Title,
				CommenterName: commenter.Username,
			})
		}
	}

	c.JSON(http.StatusCreated, comment)
}

func GetComments(c *gin.Context) {
	var comments []models.Comment

	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := config.DB.
		Where("post_id = ?", uint(postID)).
		Preload("Author").
		Find(&comments).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data"})
		return
	}

	c.JSON(http.StatusOK, comments)
}

func UpdateComment(c *gin.Context) {
	var comment models.Comment
	id := c.Param("id")

	if err := config.DB.First(&comment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	userID := c.GetUint("user_id")
	if userID != comment.AuthorID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not your comment"})
		return
	}

	var input struct {
		Content *string `json:"content"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	if input.Content != nil {
		comment.Content = *input.Content
	}

	if err := config.DB.Save(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update"})
		return
	}
	c.JSON(http.StatusOK, comment)
}

func DeleteComment(c *gin.Context) {
	id := c.Param("id")
	var comment models.Comment

	if err := config.DB.First(&comment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	userID := c.GetUint("user_id")
	if userID != comment.AuthorID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not your comment"})
		return
	}

	if err := config.DB.Delete(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted"})
}
