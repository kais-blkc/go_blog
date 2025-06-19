package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/kais-blkc/go-blog/internal/service"
	"github.com/kais-blkc/go-blog/internal/shared"
)

const (
	AuthHeader   = "Authorization"
	BearerPrefix = "Bearer"
	ErrorKey     = "error"
)

func AuthMiddleware(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader(AuthHeader)
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, shared.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "Отсутствует токен авторизации"})
			return
		}

		tokenParts := strings.Split(authHeader, " ")

		if len(tokenParts) != 2 || tokenParts[0] != BearerPrefix {
			c.AbortWithStatusJSON(http.StatusUnauthorized, shared.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "Неверный формат токена"})
			return
		}

		userID, err := authService.ValidateToken(tokenParts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, shared.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "Недействительный токен"})
			return
		}

		c.Set(shared.ContextUserID, userID)
		c.Next()
	}
}
