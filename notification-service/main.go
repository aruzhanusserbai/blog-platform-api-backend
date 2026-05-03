package main

import (
	"log"
	"net/http"
	"notification-service/config"
	"notification-service/middleware"
	"notification-service/migrations"
	"notification-service/models"

	"github.com/gin-gonic/gin"
)

type NotifyRequest struct {
	AuthorID      uint   `json:"author_id"`
	AuthorEmail   string `json:"author_email"`
	AuthorName    string `json:"author_name"`
	PostTitle     string `json:"post_title"`
	CommenterName string `json:"commenter_name"`
}

func main() {
	config.ConnectDB()
	migrations.Run()

	r := gin.Default()

	r.POST("/notify", func(c *gin.Context) {
		var req NotifyRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
			return
		}

		notification := models.Notification{
			AuthorID:      req.AuthorID,
			AuthorEmail:   req.AuthorEmail,
			AuthorName:    req.AuthorName,
			PostTitle:     req.PostTitle,
			CommenterName: req.CommenterName,
		}

		if err := config.DB.Create(&notification).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save notification"})
			return
		}

		log.Printf("📬 Уведомление сохранено: %s прокомментировал пост '%s'",
			req.CommenterName, req.PostTitle)

		c.JSON(http.StatusOK, gin.H{
			"status":       "notification saved",
			"notification": notification,
		})
	})

	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/notifications", func(c *gin.Context) {
			userID := c.GetUint("user_id")

			var notifications []models.Notification
			if err := config.DB.
				Where("author_id = ?", userID).
				Order("created_at desc").
				Find(&notifications).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch"})
				return
			}

			c.JSON(http.StatusOK, notifications)
		})

		protected.GET("/notifications/unread", func(c *gin.Context) {
			userID := c.GetUint("user_id")

			var notifications []models.Notification
			if err := config.DB.
				Where("author_id = ? AND is_read = false", userID).
				Order("created_at desc").
				Find(&notifications).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch"})
				return
			}
			c.JSON(http.StatusOK, notifications)
		})

		protected.GET("/notifications/:id", func(c *gin.Context) {
			userID := c.GetUint("user_id")
			id := c.Param("id")

			var notification models.Notification
			if err := config.DB.
				Where("id = ? AND author_id = ?", id, userID).
				First(&notification).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "notification not found"})
				return
			}
			c.JSON(http.StatusOK, notification)
		})

		protected.DELETE("/notifications/:id", func(c *gin.Context) {
			userID := c.GetUint("user_id")
			id := c.Param("id")

			var notification models.Notification
			if err := config.DB.
				Where("id = ? AND author_id = ?", id, userID).
				First(&notification).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "notification not found"})
				return
			}

			if err := config.DB.Delete(&notification).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "notification deleted"})
		})

	}

	log.Println("NotificationService running on :8081")
	r.Run(":8081")
}
