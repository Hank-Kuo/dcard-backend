package limiterRedis

import (
	"strconv"
	"time"

	constant "github.com/Hank-Kuo/dcard-backend/pkg/constant"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

func RateLimiter(param func() (*redis.Client, int, time.Duration)) gin.HandlerFunc {
	return func(c *gin.Context) {
		redisClient, ipLimit, restTime := param()
		ip := c.ClientIP()
		now := time.Now().Unix()
		script := redis.NewScript(constant.LuaScript)
		config := []interface{}{now, ipLimit, int((restTime).Seconds())}
		data, err := script.Run(redisClient, []string{ip}, config...).Result()

		if err != nil {
			c.AbortWithStatus(500)
			return
		}

		result := data.([]interface{})
		count := result[0].(int64)
		resetTime := time.Unix(result[1].(int64), 0)
		if count == -1 {
			c.Header("X-RateLimit-Remaining", strconv.FormatInt(0, 10))
			c.Header("X-RateLimit-Reset", strconv.FormatInt(resetTime.Unix(), 10))

			c.AbortWithStatus(429)
			return
		}
		c.Header("X-RateLimit-Remaining", strconv.FormatInt(count, 10))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(resetTime.Unix(), 10))
		c.Next()
	}
}
