package models

import "time"

type Notification struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	AuthorID      uint      `json:"author_id"`
	AuthorEmail   string    `json:"author_email"`
	AuthorName    string    `json:"author_name"`
	PostTitle     string    `json:"post_title"`
	CommenterName string    `json:"commenter_name"`
	CommentText   string    `json:"comment_text"`
	IsRead        bool      `gorm:"default:false" json:"is_read"`
	CreatedAt     time.Time `json:"created_at"`
}
