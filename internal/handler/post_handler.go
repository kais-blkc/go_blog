package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/kais-blkc/go-blog/internal/service"
	"github.com/kais-blkc/go-blog/internal/shared"
)

type PostHandler struct {
	postService *service.PostService
}

func NewPostHandler(postService *service.PostService) *PostHandler {
	return &PostHandler{
		postService: postService,
	}
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	var req service.CreatePostRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, shared.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Неверный формат запроса"})
		return
	}

	userID := h.getCurrentUserID(c)

	post, err := h.postService.Create(req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Ошибка при создании поста"})
		return
	}

	c.JSON(http.StatusOK, post)
}

func (h *PostHandler) GetAllPosts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	posts, err := h.postService.GetAll(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Ошибка при получении постов"})
		return
	}

	c.JSON(http.StatusOK, posts)
}

func (h *PostHandler) GetPostBySlug(c *gin.Context) {
	slug := c.Param("slug")

	post, err := h.postService.GetPostBySlug(slug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Ошибка при получении поста"})
		return
	}

	c.JSON(http.StatusOK, post)
}

func (h *PostHandler) GetUserPosts(c *gin.Context) {
	userID := h.getCurrentUserID(c)

	posts, err := h.postService.GetUserPosts(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Ошибка при получении постов пользователя"})
		return
	}

	c.JSON(http.StatusOK, posts)
}

func (h *PostHandler) UpdatePost(c *gin.Context) {
	var req service.UpdatePostRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, shared.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Неверный формат запроса"})
		return
	}

	userID := h.getCurrentUserID(c)

	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, shared.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Неверный ID поста"})
		return
	}

	post, err := h.postService.UpdatePost(uint(postID), req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Ошибка при обновлении поста"})
		return
	}

	c.JSON(http.StatusOK, post)
}

func (h *PostHandler) DeletePost(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, shared.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Неверный ID поста"})
		return
	}

	userID := h.getCurrentUserID(c)

	err = h.postService.DeletePost(uint(postID), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Ошибка при удалении поста"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Пост успешно удален"})
}

func (h *PostHandler) getCurrentUserID(c *gin.Context) uint {
	userID, exists := c.Get(shared.ContextUserID)
	if !exists {
		h.sendError(c, http.StatusUnauthorized, "Требуется авторизация")
		return 0
	}

	id, ok := userID.(uint)
	if !ok {
		h.sendError(c, http.StatusInternalServerError, "Некорректный тип ID пользователя")
		return 0
	}

	return id
}

func (h *PostHandler) sendError(c *gin.Context, code int, message string) {
	c.AbortWithStatusJSON(code, shared.ErrorResponse{
		Code:    code,
		Message: message,
	})
}
