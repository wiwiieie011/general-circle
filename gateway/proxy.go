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

		targetPath := strings.TrimPrefix(c.Request.URL.Path, "/api")

		targetURL := baseURL + targetPath
		if c.Request.URL.RawQuery != "" {
			targetURL += "?" + c.Request.URL.RawQuery
		}

		body, _ := io.ReadAll(c.Request.Body)

		req, err := http.NewRequest(
			c.Request.Method,
			targetURL,
			bytes.NewReader(body),
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		req.Header = c.Request.Header.Clone()

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "service unavailable"})
			return
		}
		defer resp.Body.Close()

		respBody, _ := io.ReadAll(resp.Body)

		c.Data(
			resp.StatusCode,
			resp.Header.Get("Content-Type"),
			respBody,
		)
	}
}
