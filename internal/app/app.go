package app

import (
	"gorm.io/gorm"

	"github.com/kais-blkc/go-blog/internal/config"
	"github.com/kais-blkc/go-blog/internal/handler"
	"github.com/kais-blkc/go-blog/internal/repository"
	"github.com/kais-blkc/go-blog/internal/service"
)

type Application struct {
	DB          *gorm.DB
	cfg         *config.Config
	AuthHandler *handler.AuthHandler
	PostHandler *handler.PostHandler
	AuthService *service.AuthService
}

func NewApp(db *gorm.DB, cfg *config.Config) *Application {

	// Инициализируем репозитории
	userRepo := repository.NewUserRepository(db)
	postRepo := repository.NewPostRepository(db)

	// Инициализируем сервисы
	authService := service.NewAuthService(userRepo, cfg.JwtSecret)
	postService := service.NewPostService(postRepo)

	// Инициализируем обработчики
	authHandler := handler.NewAuthHandler(authService)
	postHandler := handler.NewPostHandler(postService)

	return &Application{
		DB:          db,
		cfg:         cfg,
		AuthHandler: authHandler,
		PostHandler: postHandler,
	}
}
