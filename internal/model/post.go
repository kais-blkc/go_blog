package model

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	Title     string `json:"title" gorm:"not null"`
	Content   string `json:"content" gorm:"not null"`
	Excerpt   string `json:"excerpt"`
	Slug      string `json:"slug" gorm:"uniqueIndex;not null"`
	Published bool   `json:"publisher" gorm:"default:true"`
	AuthorID  uint   `json:"author_id"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Author   User      `json:"author" gorm:"foreignKey:AuthorID"`
	Comments []Comment `json:"comments" gorm:"foreignKey:PostID"`
}
