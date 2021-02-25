# Dcard Backend Demo
Implement the middleware including:
- Limit the number of requests from the same IP per hour not to exceed 1000
- Add the remaining number of requests (X-RateLimit-Remaining) and the time to zero the rate limit (X-RateLimit-Reset) in the response headers
- If it exceeds the limit, it will return 429 (Too Many Requests)
- Can be achieved using various databases

## How to use 


## Architecture 
In this project, I implement two way to achieve this homework. 

## Algorithm 


### Database 
- Cache
    - If the server don't need auto scale or not using cluster, I don't need to consider the data share problem. So in here, I save rate limit's information in cache, I don't need to install external database.
    - In here, I use go-cache library as our database. In go-cache, the author use atomic way to fix race condition problem, so I don't need to consider this problem, then I import the time/rate standard library, using the rate limit that have been created (token-bucket algorithm).
- Redis 
    -  f    
## 
## 