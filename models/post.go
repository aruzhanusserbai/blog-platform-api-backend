package models

import "time"

type Post struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	Comments  []Comment `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE;"`
	AuthorID  uint      `json:"author_id"`
	Tags      []Tag     `gorm:"many2many:post_tags;"`
}
