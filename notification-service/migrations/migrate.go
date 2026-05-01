package migrations

import (
	"log"
	"notification-service/config"
	"notification-service/models"
)

func Run() {
	err := config.DB.AutoMigrate(&models.Notification{})
	if err != nil {
		log.Fatal("Migration failed: ", err)
	}
	log.Println("Migration completed!")
}
