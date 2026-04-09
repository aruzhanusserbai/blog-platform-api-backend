package models

type Tag struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`
}
