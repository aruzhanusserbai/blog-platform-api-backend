package handlers

import (
	"blogPlatform/config"
	dto2 "blogPlatform/dto"
	"blogPlatform/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func CreatePost(c *gin.Context) {
	var dto dto2.CreatePostDTO

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if strings.TrimSpace(dto.Title) == "" || strings.TrimSpace(dto.Content) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Title and content cannot be empty",
		})
		return
	}

	post := models.Post{
		Title:    dto.Title,
		Content:  dto.Content,
		AuthorID: c.GetUint("user_id"),
	}

	if err := config.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Post is not created"})
		return
	}

	if len(dto.Tags) > 0 {
		var tags []models.Tag
		config.DB.Where("name IN ?", dto.Tags).Find(&tags)
		config.DB.Model(&post).Association("Tags").Replace(&tags)
	}

	config.DB.Preload("Author").Preload("Tags").First(&post, post.ID)

	c.JSON(http.StatusCreated, post)
}

func GetPosts(c *gin.Context) {
	var posts []models.Post

	if err := config.DB.Preload("Tags").Preload("Comments").Preload("Author").Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Can't fetch the data"})
		return
	}

	c.JSON(http.StatusOK, posts)
}

func GetPost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid credentials"})
		return
	}

	var post models.Post
	if err := config.DB.Preload("Tags").Preload("Comments").Preload("Author").First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}
	c.JSON(http.StatusOK, post)
}

func UpdatePost(c *gin.Context) {
	var post models.Post
	id := c.Param("id")

	if err := config.DB.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	userID := c.GetUint("user_id")
	if userID != post.AuthorID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not your post"})
		return
	}

	var input struct {
		Title   *string  `json:"title"`
		Content *string  `json:"content"`
		Tags    []string `json:"tags"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if input.Title != nil {
		post.Title = *input.Title
	}

	if input.Content != nil {
		post.Content = *input.Content
	}

	if input.Tags != nil {
		var tags []models.Tag

		for _, t := range input.Tags {
			var tag models.Tag

			err := config.DB.Where("name = ?", t).First(&tag).Error
			if err != nil {
				tag = models.Tag{Name: t}
				config.DB.Create(&tag)
			}

			tags = append(tags, tag)
		}

		config.DB.Model(&post).Association("Tags").Replace(tags)
	}

	if err := config.DB.Save(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
		return
	}

	c.JSON(http.StatusOK, post)
}

func DeletePost(c *gin.Context) {
	id := c.Param("id")
	var post models.Post

	if err := config.DB.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	userID := c.GetUint("user_id")
	if userID != post.AuthorID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not your post"})
		return
	}

	if err := config.DB.Delete(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Post not deleted"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Post deleted"})
}
