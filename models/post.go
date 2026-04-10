package models

import "time"

type Post struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	Comments  []Comment `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE;" json:"comments,omitempty"`
	AuthorID  uint      `json:"-"`
	Author    Author    `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	Tags      []Tag     `gorm:"many2many:post_tags;" json:"tags,omitempty"`
}
