package config

import (
	"context"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis(logger *slog.Logger) *redis.Client {
	redisAddr := os.Getenv("REDIS_ADDR")     // пример: "localhost:6379"
	redisPass := os.Getenv("REDIS_PASSWORD") // если нет пароля, можно оставить ""
	redisDBStr := os.Getenv("REDIS_DB")      // номер БД, по умолчанию 0

	dbNum := 0
	if redisDBStr != "" {
		n, err := strconv.Atoi(redisDBStr)
		if err != nil {
			logger.Warn("invalid REDIS_DB, using 0", "value", redisDBStr)
		} else {
			dbNum = n
		}
	}

	var client *redis.Client
	maxAttempts := 12
	backoff := 2 * time.Second

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		client = redis.NewClient(&redis.Options{
			Addr:     redisAddr,
			Password: redisPass,
			DB:       dbNum,
		})

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		_, err := client.Ping(ctx).Result()
		if err == nil {
			logger.Info("connected to Redis", "addr", redisAddr)
			return client
		}

		logger.Warn("Redis connect attempt failed", "attempt", attempt, "error", err)
		time.Sleep(backoff)
		if backoff < 10*time.Second {
			backoff *= 2
			if backoff > 10*time.Second {
				backoff = 10 * time.Second
			}
		}
	}

	logger.Error("failed to connect to Redis after retries")
	os.Exit(1)
	return nil
}
