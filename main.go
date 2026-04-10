package main

import (
	"blogPlatform/config"
	"blogPlatform/handlers"
	"blogPlatform/middleware"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	PORT := os.Getenv("PORT")
	config.ConnectDB()

	app := gin.Default()

	app.POST("/register", handlers.Register)
	app.POST("/login", handlers.LogIn)
	app.GET("/posts", handlers.GetPosts)
	app.GET("/post/:id", handlers.GetPost)
	app.GET("/tags", handlers.GetTags)

	protected := app.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/posts", handlers.CreatePost)

		protected.PATCH("/posts/:id", handlers.UpdatePost)
		protected.DELETE("/post/:id", handlers.DeletePost)

		protected.POST("/posts/:id/comments", handlers.AddCommentToPost)
		protected.GET("/posts/:id/comments", handlers.GetComments)
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
