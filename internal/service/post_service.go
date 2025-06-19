package service

import (
	"errors"
	"math"
	"regexp"
	"strings"

	"github.com/mozillazg/go-unidecode"

	"github.com/kais-blkc/go-blog/internal/model"
	"github.com/kais-blkc/go-blog/internal/repository"
)

const (
	ErrPostNotFound           = "пост не найден"
	ErrNotEnoughRights        = "недостаточно прав для редактирования поста"
	ErrPostAlreadyExistsTitle = "пост с таким заголовком уже существует"
	ErrPostAlreadyExistsSlug  = "пост с таким slug уже существует"
)

type PostService struct {
	postRepo *repository.PostRepository
}

func NewPostService(postRepo *repository.PostRepository) *PostService {
	return &PostService{
		postRepo: postRepo,
	}
}

type CreatePostRequest struct {
	Title   string `json:"title" binding:"required,min=3,max=100"`
	Content string `json:"content" binding:"required,min=10"`
	Excerpt string `json:"excerpt"`
}

type UpdatePostRequest struct {
	Title     string `json:"title" binding:"required,min=3,max=100"`
	Content   string `json:"content" binding:"required,min=10"`
	Excerpt   string `json:"excerpt"`
	Published bool   `json:"published"`
}

func (s *PostService) Create(req CreatePostRequest, authorID uint) (*model.Post, error) {
	slug := s.generateSlug(req.Title)

	existingPost, _ := s.postRepo.GetBySlug(slug)
	if existingPost != nil {
		return nil, errors.New(ErrPostAlreadyExistsTitle)
	}

	excerpt := req.Excerpt
	if excerpt == "" {
		excerpt = s.generateExcerpt(req.Content, 200)
	}

	post := &model.Post{
		Title:    req.Title,
		Content:  req.Content,
		Excerpt:  excerpt,
		Slug:     slug,
		AuthorID: authorID,
	}

	err := s.postRepo.Create(post)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (s *PostService) generateSlug(title string) string {
	slug := unidecode.Unidecode(title)
	slug = strings.ToLower(slug)
	reg := regexp.MustCompile(`[^\a-z0-9\s-]`)
	slug = reg.ReplaceAllString(slug, "")
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = regexp.MustCompile(`-+`).ReplaceAllString(slug, "-")
	return strings.Trim(slug, "-")
}

func (s *PostService) generateExcerpt(content string, maxLen int) string {
	cleanContent := strings.TrimSpace(content)
	if len(cleanContent) <= maxLen {
		return cleanContent
	}

	exerpt := string([]rune(cleanContent)[:maxLen])
	lastSpace := strings.LastIndex(exerpt, " ")
	if lastSpace != -1 {
		exerpt = exerpt[:lastSpace]
	}

	return exerpt + "..."
}

func (s *PostService) GetAll(page, limit int) ([]model.Post, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit
	offset = int(math.Max(0, float64(offset)))

	return s.postRepo.GetAll(limit, offset)
}

func (s *PostService) GetPostBySlug(slug string) (*model.Post, error) {
	return s.postRepo.GetBySlug(slug)
}

func (s *PostService) UpdatePost(postID uint, req UpdatePostRequest, authorId uint) (*model.Post, error) {
	post, err := s.postRepo.GetByID(postID)
	if err != nil {
		return nil, errors.New(ErrPostNotFound)
	}

	if post.AuthorID != authorId {
		return nil, errors.New(ErrNotEnoughRights)
	}

	post.Title = req.Title
	post.Content = req.Content
	post.Excerpt = req.Excerpt
	post.Published = req.Published

	err = s.postRepo.Update(post)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (s *PostService) DeletePost(postID uint, authorID uint) error {
	post, err := s.postRepo.GetByID(postID)
	if err != nil {
		return errors.New(ErrPostNotFound)
	}

	if post.AuthorID != authorID {
		return errors.New(ErrNotEnoughRights)
	}

	return s.postRepo.Delete(postID)
}

func (s *PostService) GetUserPosts(authorID uint) ([]model.Post, error) {
	return s.postRepo.GetByAuthor(authorID)
}
