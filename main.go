package main

import (
	"blogPlatform/config"
	"blogPlatform/handlers"
	"blogPlatform/middleware"
	"blogPlatform/models"
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	PORT := os.Getenv("PORT")
	config.ConnectDB()

	config.DB.AutoMigrate(
		&models.Author{},
		&models.Post{},
		&models.Comment{},
		&models.Tag{},
	)

	app := gin.Default()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	app.POST("/register", handlers.Register)
	app.POST("/login", handlers.LogIn)
	app.GET("/posts", handlers.GetPosts)
	app.GET("/posts/:id", handlers.GetPost)
	app.GET("/tags", handlers.GetTags)
	app.GET("/posts/:id/comments", handlers.GetComments)

	protected := app.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/posts", handlers.CreatePost)

		protected.PATCH("/posts/:id", handlers.UpdatePost)
		protected.DELETE("/posts/:id", handlers.DeletePost)

		protected.POST("/posts/:id/comments", handlers.AddCommentToPost)
		protected.PATCH("/comments/:id", handlers.UpdateComment)
		protected.DELETE("/comments/:id", handlers.DeleteComment)

		protected.POST("/tags", handlers.CreateTags)
		protected.POST("/posts/:id/tags", handlers.AddTagsToPost)

		protected.GET("/users/:id", handlers.GetUser)
		protected.POST("/users", handlers.CreateUser)
	}

	log.Println("Server running on port 8080")
	app.Run(fmt.Sprintf(":%v", PORT))
}
