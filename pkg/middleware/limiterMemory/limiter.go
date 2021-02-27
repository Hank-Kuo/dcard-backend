package limiterMemory

import (
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type LimiterInfo struct {
	count  int
	expire time.Time
}

type Limiter struct {
	mux        sync.RWMutex
	ipLimit    int
	duration   time.Duration
	collection map[string]LimiterInfo
}

func NewLimiter(r int, t time.Duration) *Limiter {
	return &Limiter{
		ipLimit:    r,
		duration:   t,
		collection: make(map[string]LimiterInfo),
	}
}

func (l *Limiter) Allow(key string) bool {
	l.mux.Lock()
	defer l.mux.Unlock()
	data, ok := l.collection[key]

	var count int
	var resetTime time.Time
	if !ok || data.expire.Before(time.Now()) {
		count = l.ipLimit - 1
		resetTime = time.Now().Add(l.duration)
	} else {
		// refresh count
		count = data.count - 1
		resetTime = data.expire
		if count < 0 {
			count = -1
			l.collection[key] = LimiterInfo{count, resetTime}
			return false
		}
	}
	l.collection[key] = LimiterInfo{count, resetTime}
	return true
}

func (l *Limiter) Get(key string) (int, int64, error) {
	l.mux.RLock()
	defer l.mux.RUnlock()
	data, ok := l.collection[key]
	if !ok {
		return 0, 0, errors.New("can't find record")
	}
	if data.count == -1 {
		return 0, (data.expire).Unix(), nil
	}
	return data.count, (data.expire).Unix(), nil
}

func RateLimiter(param func() (int, time.Duration)) gin.HandlerFunc {
	ipLimit, restTime := param()
	limiter := NewLimiter(ipLimit, restTime)
	return func(c *gin.Context) {
		ip := c.ClientIP()
		ok := limiter.Allow(ip)
		remaining, t, e := limiter.Get(ip)
		if e != nil {
			c.AbortWithStatus(500)
			return
		}

		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(t, 10))
		if !ok {
			c.AbortWithStatus(429)
			return
		}
		c.Next()
	}
}

/*
// test race condition
func main() {
	var wg sync.WaitGroup
	a := NewLimiter(10, time.Hour)

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
*/
