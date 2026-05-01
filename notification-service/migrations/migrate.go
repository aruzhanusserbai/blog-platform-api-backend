package migrations

import (
	"blogPlatform/notification-service/config"
	"blogPlatform/notification-service/models"
	"log"
)

func Run() {
	err := config.DB.AutoMigrate(&models.Notification{})
	if err != nil {
		log.Fatal("Migration failed: ", err)
	}
	log.Println("Migration completed!")
}
