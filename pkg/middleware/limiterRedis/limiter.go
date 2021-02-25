package limiterRedis

import (
	constant "dcardBackend/pkg/constant"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

func RateLimiter(param func() (*redis.Client, int, time.Duration)) gin.HandlerFunc {
	return func(c *gin.Context) {
		redisClient, ipLimit, restTime := param()
		ip := c.ClientIP()
		now := time.Now().Unix()
		script := redis.NewScript(constant.Script)
		config := []interface{}{now, ipLimit, int((restTime).Seconds())}
		value, err := script.Run(redisClient, []string{ip}, config...).Result()

		if err != nil {
			c.AbortWithStatus(500)
			return
		}

		result := value.([]interface{})
		remaining := result[0].(int64)
		resetTime := time.Unix(result[1].(int64), 0)
		if remaining == -1 {
			c.Header("X-RateLimit-Remaining", strconv.FormatInt(0, 10))
			c.Header("X-RateLimit-Reset", strconv.FormatInt(resetTime.Unix(), 10))

			c.AbortWithStatus(429)
			return
		}
		c.Header("X-RateLimit-Remaining", strconv.FormatInt(remaining, 10))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(resetTime.Unix(), 10))
		c.Next()
	}
}
