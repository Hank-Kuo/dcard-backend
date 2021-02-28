package limiterRedis

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func NewClient() *redis.Client {
	godotenv.Load()
	redisURL := os.Getenv("REDIS_URL_TEST")
	client := redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Password: "",
		DB:       0,
	})
	return client
}

func TestLimiter(t *testing.T) {
	r := gin.New()
	assert := assert.New(t)

	client := NewClient()
	r.Use(RateLimiter(func() (*redis.Client, int, time.Duration) {
		return client, 10, time.Hour * 1
	}))

	r.GET("/", func(c *gin.Context) {})

	go func() {
		err := r.Run(":8888")
		if err != nil {
			t.Fatal("Error run http server", err.Error())
		}
	}()

	runtime.Gosched()

	t.Run("MemoryLimiter After 1 request", func(t *testing.T) {
		resp, err := http.DefaultClient.Get("http://127.0.0.1:8888")
		if err != nil {
			t.Error("Error during requests", err.Error())
			return
		}
		remaining, _ := strconv.Atoi(resp.Header.Get("X-RateLimit-Remaining"))
		assert.Equal(resp.StatusCode, 200, "It should be 200")
		assert.Equal(remaining, 9, "It should be 9")
		fmt.Printf("Status code: %d \t remain: %s \t reset time: %s \n", resp.StatusCode, resp.Header.Get("X-RateLimit-Remaining"), resp.Header.Get("X-RateLimit-Reset"))
	})

	var wg sync.WaitGroup
	n := 4
	wg.Add(n)
	for i := 1; i <= n; i++ {
		go func() {
			_, err := http.DefaultClient.Get("http://127.0.0.1:8888")
			if err != nil {
				t.Error("Error during requests", err.Error())
			}
			wg.Done()
		}()
	}
	wg.Wait()

	t.Run("MemoryLimiter After 5 request", func(t *testing.T) {
		resp, err := http.DefaultClient.Get("http://127.0.0.1:8888")
		if err != nil {
			t.Error("Error during requests", err.Error())
		}
		remaining, _ := strconv.Atoi(resp.Header.Get("X-RateLimit-Remaining"))
		assert.Equal(resp.StatusCode, 200, "It should be 200")
		assert.Equal(remaining, 4, "It should be 4")
		fmt.Printf("Status code: %d \t remain: %s \t reset time: %s \n", resp.StatusCode, resp.Header.Get("X-RateLimit-Remaining"), resp.Header.Get("X-RateLimit-Reset"))
	})

	n = 9
	wg.Add(n)
	for i := 1; i <= n; i++ {
		go func() {
			_, err := http.DefaultClient.Get("http://127.0.0.1:8888")
			if err != nil {
				t.Error("Error during requests", err.Error())
			}
			wg.Done()
		}()
	}
	wg.Wait()

	t.Run("MemoryLimiter After 10 request", func(t *testing.T) {
		resp, err := http.DefaultClient.Get("http://127.0.0.1:8888")
		if err != nil {
			t.Error("Error during requests", err.Error())
		}

		remaining, _ := strconv.Atoi(resp.Header.Get("X-RateLimit-Remaining"))
		assert.Equal(resp.StatusCode, 429, "It should be 429")
		assert.Equal(remaining, 0, "It should be 0")
		fmt.Printf("Status code: %d \t remain: %s \t reset time: %s \n", resp.StatusCode, resp.Header.Get("X-RateLimit-Remaining"), resp.Header.Get("X-RateLimit-Reset"))
	})
}
