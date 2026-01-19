package main

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func proxyToService(baseURL string) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Убираем префикс /api
		// /api/notifications/123/read -> /notifications/123/read
		targetPath := strings.TrimPrefix(c.Request.URL.Path, "/api")

		// Собираем итоговый URL
		targetURL := baseURL + targetPath

		// Читаем тело запроса
		body, _ := io.ReadAll(c.Request.Body)

		// Создаём новый HTTP-запрос
		req, err := http.NewRequest(
			c.Request.Method,
			targetURL,
			bytes.NewReader(body),
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Копируем заголовки
		req.Header = c.Request.Header.Clone()

		// Отправляем запрос в микросервис
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "service unavailable"})
			return
		}
		defer resp.Body.Close()

		// Читаем ответ
		respBody, _ := io.ReadAll(resp.Body)

		// Возвращаем ответ клиенту
		c.Data(
			resp.StatusCode,
			resp.Header.Get("Content-Type"),
			respBody,
		)
	}
}
