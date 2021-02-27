package main

import (
	"time"

	"github.com/Hank-Kuo/dcard-backend/pkg/middleware/limiterMemory"
	"github.com/Hank-Kuo/dcard-backend/pkg/middleware/limiterRedis"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

func NewClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
	return client
}

func ServerWithLocalLimiter() {
	r := gin.Default()
	r.Use(limiterMemory.RateLimiter(func() (int, time.Duration) {
		return 1000, time.Hour * 1
	}))
	r.GET("/", func(c *gin.Context) {})

	r.Run(":8888")
}

func ServerWithRedisLimiter() {
	r := gin.Default()
	client := NewClient()
	r.Use(limiterRedis.RateLimiter(func() (*redis.Client, int, time.Duration) {
		return client, 10, time.Hour * 1
	}))

	r.GET("/", func(c *gin.Context) {})
	r.Run(":8888")
}

func main() {
	ServerWithLocalLimiter()
}
