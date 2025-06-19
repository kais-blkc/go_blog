package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kais-blkc/go-blog/internal/service"
	"github.com/kais-blkc/go-blog/internal/shared"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req service.RegisterRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Println("Register ShouldBindJSON error")
		log.Println(err)
		log.Println(c.Request.Body)

		c.JSON(http.StatusBadRequest, shared.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Неверный формат запроса"})
		return
	}

	response, err := h.authService.Register(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Ошибка при регистрации"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, shared.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Неверный формат запроса"})
		return
	}

	response, err := h.authService.Login(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Ошибка при входе"})
		return
	}

	c.JSON(http.StatusOK, response)
}
