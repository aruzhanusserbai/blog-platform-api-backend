package main

import (
	"blogPlatform/config"
	"blogPlatform/handlers"
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

	app.POST("/posts", handlers.CreatePost)
	app.GET("/posts", handlers.GetPosts)
	app.GET("/post/:id", handlers.GetPost)
	app.PATCH("/posts/:id", handlers.UpdatePost)
	app.DELETE("/post/:id", handlers.DeletePost)

	app.POST("/posts/:id/comments", handlers.AddCommentToPost)
	app.GET("/posts/:id/comments", handlers.GetComments)
	app.PATCH("/comments/:id", handlers.UpdateComment)
	app.DELETE("/comments/:id", handlers.DeleteComment)

	app.GET("/tags", handlers.GetTags)
	app.POST("/tags", handlers.CreateTags)
	app.POST("/posts/:id/tags", handlers.AddTagsToPost)

	app.GET("/users/:id", handlers.GetUser)
	app.POST("/users", handlers.CreateUser)

	log.Println("Server running on port 8080")
	app.Run(fmt.Sprintf(":%v", PORT))
}
