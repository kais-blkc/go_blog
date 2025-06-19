package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/kais-blkc/go-blog/internal/handler"
)

func SetupAuthRoutes(
	api *gin.RouterGroup,
	authHandler *handler.AuthHandler,
) {
	auth := api.Group("/auth/")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}
}
