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
	}

	log.Println("NotificationService running on :8081")
	r.Run(":8081")
}
