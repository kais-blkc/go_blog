package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/kais-blkc/go-blog/internal/handler"
	"github.com/kais-blkc/go-blog/internal/middleware"
	"github.com/kais-blkc/go-blog/internal/service"
)

func SetupPostRoutes(
	api *gin.RouterGroup,
	postHandler *handler.PostHandler,
	authHandler *handler.AuthHandler,
	authService *service.AuthService,
) {
	posts := api.Group("/posts")
	{
		posts.GET("/", postHandler.GetAllPosts)
		posts.GET("/:slug", postHandler.GetPostBySlug)
	}

	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(authService))
	{
		posts := protected.Group("/posts")
		{
			posts.GET("/my", postHandler.GetUserPosts)
			posts.POST("/", postHandler.CreatePost)
			posts.PUT("/:id", postHandler.UpdatePost)
			posts.DELETE("/:slug", postHandler.DeletePost)
		}
	}
}
