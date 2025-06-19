package repository

import (
	"gorm.io/gorm"

	"github.com/kais-blkc/go-blog/internal/model"
)

type PostRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) Create(post *model.Post) error {
	return r.db.Create(post).Error
}

func (r *PostRepository) GetAll(limit, offset int) ([]model.Post, error) {
	var posts []model.Post
	err := r.db.Preload("Author").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&posts).Error

	return posts, err
}

func (r *PostRepository) GetByID(id uint) (*model.Post, error) {
	var post model.Post
	err := r.db.
		Preload("Author").
		Preload("Comments.Author").
		First(&post, id).Error

	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (r *PostRepository) GetBySlug(slug string) (*model.Post, error) {
	var post model.Post
	err := r.db.
		Preload("Author").
		Preload("Comments.Author").
		Where("slug = ? AND published = ?", slug, true).
		First(&post).Error

	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (r *PostRepository) GetByAuthor(authorID uint) ([]model.Post, error) {
	var posts []model.Post
	err := r.db.
		Where("author_id = ?", authorID).
		Order("created_at DESC").
		Find(&posts).Error

	return posts, err
}

func (r *PostRepository) Update(post *model.Post) error {
	return r.db.Save(post).Error
}

func (r *PostRepository) Delete(id uint) error {
	return r.db.Delete(&model.Post{}, id).Error
}
