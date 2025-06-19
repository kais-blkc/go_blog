package model

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	ID       string `json:"id" gorm:"primary_key"`
	Content  string `json:"content" gorm:"not null"`
	AuthorID uint   `json:"author_id"`
	PostID   uint   `json:"post_id" `

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Author User `json:"author" gorm:"foreignKey:AuthorID"`
	Post   Post `json:"post" gorm:"foreignKey:PostID"`
}
