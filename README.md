
# Dcard Backend Demo
Implement the middleware including:
- Limit the number of requests from the same IP per hour not to exceed 1000
- Add the remaining number of requests (X-RateLimit-Remaining) and the time to zero the rate limit (X-RateLimit-Reset) in the response headers
- If it exceeds the limit, it will return 429 (Too Many Requests)
- Can be achieved using various databases

## How to use 
```go
import (
	"github.com/Hank-Kuo/dcard-backend/pkg/middleware/limiterMemory"
	"github.com/Hank-Kuo/dcard-backend/pkg/middleware/limiterRedis"
)

// memory mode
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

// redis mode
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

```

## How to start 
```script
# run 
make run  # run go server 

# test script
make test ./pkg/middleware/limiterMemory # test middleware
make test ./pkg/middleware/limiterRedis # test middleware

# build
make build # build binary file 
make build-darwin
make build-window
make build-linux

# docker 
docker-run # run docker 
make docker-compose # run docker compose 
```

## Architecture 
In this project, I use two database to achieve this homework. One is Cache, another is redis. Developer can use the one of them in different case. 
- ![](https://i.imgur.com/bmOIZaD.png)
- 
### Algorithm 
- **Simple Counter Algorithm**: When someone send request to server, the server will take the number of times in database( or memory), then compare whether more than we excepted times. If the counter doesn't exceed, it will minus 1, then return current number, otherwise, it will return -1. 

![](https://i.imgur.com/551jdNt.png)



### Database 
- **Cache**
    - If the server don't need auto scale or not using cluster, I don't need to consider the data share problem. So in here, I save rate limit's information in cache, I don't need to install external database.
    - When the many request come, the serve will suffer the race condition, so I use lock way to avoid this situation. 
- **Redis**
    -  Consider the multiple server will share data, so I use redis to solve this problem. 
    -  When the many request come, the serve will suffer the race condition, to avoid this problem, I use lua script to solve. Because the redis is based on single thread, so it doesn't cause race condition.

## Demo
In here, I run the simple server with middleware I made. The entry point is in ```cmd/main.go```.

You can select the cache(```ServerWithRedisLimiter```) or redis(```ServerWithLocalLimiter```) type. 

- Script: 
    - ```make run``` ( If you want to redis mode, you need to set the configuration in main.go or ```make docker-compose```) 

- Success case
    - the limit rate doesn't exceed the number of we expected.
    - ![](https://i.imgur.com/KptdBq1.png)
- Error case 
    - the limit rate exceeds the number of we expected.
    - ![](https://i.imgur.com/4jV1g6i.png)

## Docker
- Script
```
make docker-build-image # Step 1 
make docker-run # only run server, not redis ( Step 2)
make docker-compose # (Step 2)
```  

## Test
In the test case, I simulate the server can only accept 10 requests in an hour, then the script will send parael request to server. 

Then the 3 test case. First is after 1 request, the request status code should be 200, remaned request should be 9, reset time should be after 1 hour. Second is after 5 request, the request status code should be 200, remaned request should be 4, reset time should be after 1 hour. Last is after 10 request, the request status code should be 200, remaned request should be 4, reset time should be after 1 hour. 

- Script: 
    - ```make test ./pkg/middleware/limiterMemory # test middleware```
    - ```make test ./pkg/middleware/limiterRedis # test middleware```
        - make sure you redis host is 127.0.0.1:6379 and no password.     

- Result:
    - ![](https://i.imgur.com/Yxpe5WM.png)

## Extra Test
In Go offically provide cli to detect race condite problem (```go xxx -race```) 

- Script: ```go run main.go -race```

```go
// test race condition
func main() {
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
```
![](https://i.imgur.com/shkotz8.png)