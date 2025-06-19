package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/kais-blkc/go-blog/internal/app"
	"github.com/kais-blkc/go-blog/internal/config"
	"github.com/kais-blkc/go-blog/internal/routes"
	"github.com/kais-blkc/go-blog/pkg/database"
)

func main() {
	// Загружаем конфигурацию
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки конфигурации:", err)
	}

	// Подключаемся к базе данных
	db, err := database.Connect(cfg.DatabaseUrl)
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных:", err)
	}

	// Initialize app
	app := app.NewApp(db, cfg)

	// Настраиваем роутер
	router := gin.Default()
	routes.SetupRoutes(router, app)

	// Запускаем сервер
	log.Printf("Сервер запущен http://localhost:%s", cfg.Port)

	err = router.Run(":" + cfg.Port)
	if err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}

}
