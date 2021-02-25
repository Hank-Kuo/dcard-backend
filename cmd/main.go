package main

import (
	"github.com/Hank-Kuo/dcard-backend/pkg/middleware/limiterMemory"
	"github.com/gin-gonic/gin"
	"time"
)

func main() {
	r := gin.New()

	r.Use(limiterMemory.RateLimiter(func() (int, time.Duration) {
		return 10, time.Hour * 1
	}))

	r.GET("/", func(c *gin.Context) {})
	r.Run(":8888")
}

/*
import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type Bank struct {
	balance int
	mux     sync.Mutex
}

func (b *Bank) Deposit(amount int) {
	b.mux.Lock()
	b.balance += amount
	b.mux.Unlock()
}
func (b *Bank) Balance() int {
	b.mux.Lock()
	balance := b.balance
	b.mux.Unlock()
	return balance
}

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
	// when key is not exist or already expires, refresh map
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
	return data.count, (data.expire).Unix(), nil
}

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
			//b.Deposit(1000)
			wg.Done()
		}()
	}
	fmt.Println(a.Get("127.0.0.1"))
	// fmt.Println(b.Balance())
}


func RateLimiter(param func() (time.Duration, int),
	abort func(*gin.Context)) gin.HandlerFunc {

	expire, bucket := param()
	limiter := NewLimiter(bucket, expire)
	return func(c *gin.Context) {
		ip := c.ClientIP()
		ok := limiter.Allow(ip)
		remaining, resetTime, e := limiter.Get(ip)
		if e != nil {
			fmt.Println("error")
		}

		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(resetTime, 10))
		if !ok {
			c.AbortWithStatus(429)
			return
		}
		c.Next()
	}
}
*/
