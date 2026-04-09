package models

type Comment struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Content  string `json:"content"`
	PostID   uint   `json:"post_id"`
	AuthorID uint   `json:"author_id"`
}
