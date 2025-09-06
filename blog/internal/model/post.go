package model

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	Title    string `gorm:"not null"`
	Content  string `gorm:"type:text;not null"`
	UserID   uint
	User     User
	Comments []Comment
}
