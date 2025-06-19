package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kais-blkc/go-blog/internal/app"
	"github.com/kais-blkc/go-blog/internal/middleware"
)

func SetupRoutes(
	router *gin.Engine,
	app *app.Application,
) {
	// Middleware
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.RequestBodyLogger())

	root := router.Group("/")
	{
		root.GET("/", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "Its root directory, please use /api/",
			})
		})
	}

	api := router.Group("/api/")
	{
		SetupAuthRoutes(api, app.AuthHandler)
		SetupPostRoutes(api, app.PostHandler, app.AuthHandler, app.AuthService)
	}
}
