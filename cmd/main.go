package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/Hank-Kuo/dcard-backend/pkg/middleware/limiterMemory"
	"github.com/Hank-Kuo/dcard-backend/pkg/middleware/limiterRedis"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
)

func NewClient() *redis.Client {
	redisURL := os.Getenv("REDIS_URL")
	client := redis.NewClient(&redis.Options{
		Addr:     redisURL, // "redis:6379"
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
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"狀態": "ok",
		})
	})

	r.Run(":8080")
}

func ServerWithRedisLimiter() {
	r := gin.Default()
	client := NewClient()
	r.Use(limiterRedis.RateLimiter(func() (*redis.Client, int, time.Duration) {
		return client, 1000, time.Hour * 1
	}))

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"狀態": "ok",
		})
	})
	r.Run(":8080")
}

func TestRace() {
	var wg sync.WaitGroup
	a := limiterMemory.NewLimiter(10, time.Hour)

	n := 1000
	wg.Add(n)

	for i := 1; i <= n; i++ {
		go func() {
			allow := a.Allow("127.0.0.1")
			if allow {
				fmt.Println(i)
			}
			wg.Done()
		}()
	}
}

func main() {
	godotenv.Load()
	ServerWithRedisLimiter() // run redis type
	// ServerWithLocalLimiter() // run local type

	// TestRace()
}
