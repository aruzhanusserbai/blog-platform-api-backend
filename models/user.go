package models

type Author struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Username string `json:"username"`
	Email    string `gorm:"unique" json:"email"`
	Posts    []Post `gorm:"foreignKey:AuthorID;constraint:OnDelete:CASCADE;"`
}
