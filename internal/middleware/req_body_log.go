package middleware

import (
	"bytes"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequestBodyLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Читаем тело запроса
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			log.Printf("Ошибка чтения тела запроса: %v", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Не удалось прочитать тело запроса"})
			return
		}

		// Возвращаем тело обратно в Request.Body
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		// Выводим в консоль
		log.Printf("Тело запроса (%s %s): %s", c.Request.Method, c.Request.URL, string(body))

		c.Next()
	}
}
