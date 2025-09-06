package model

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	Content string `gorm:"not null"`
	PostID  uint
	UserID  uint
	User    User
}
