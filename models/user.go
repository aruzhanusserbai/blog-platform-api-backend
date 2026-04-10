package models

import "time"

type Author struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"unique" json:"username"`
	Email     string    `gorm:"unique" json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	Posts     []Post    `gorm:"foreignKey:AuthorID;constraint:OnDelete:CASCADE;"`
	Comments  []Comment `gorm:"foreignKey:AuthorID;constraint:OnDelete:CASCADE;"`
}
