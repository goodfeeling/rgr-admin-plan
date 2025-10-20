package redis

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// use lua for redis
const rateLimitScript = `
local key = KEYS[1]
local limit = tonumber(ARGV[1])
local window = tonumber(ARGV[2])
local now = tonumber(ARGV[3])

-- remove expired requests
redis.call('ZREMRANGEBYSCORE', key, 0, now - window)

-- get current count
local current = redis.call('ZCARD', key)

-- if exceeded limit, return 0
if current >= limit then
    return 0
end

-- 添加当前请求
redis.call('ZADD', key, now, now)
redis.call('EXPIRE', key, window)

return 1
`

// RedisLuaRateLimiter
type RedisLuaRateLimiter struct {
	client *redis.Client
	script *redis.Script
	prefix string
	limit  int64
	window int64
}

// NewRedisLuaRateLimiter
func NewRedisLuaRateLimiter(client *redis.Client) *RedisLuaRateLimiter {
	limit, err := strconv.Atoi(os.Getenv("SERVER_LIMIT_RATE"))
	if err != nil {
		limit = 100
	}
	minute, err := strconv.Atoi(os.Getenv("SERVER_LIMIT_MINUTE"))
	if err != nil {
		minute = 1
	}

	window := int64(time.Duration(minute) * time.Minute / time.Second)
	return &RedisLuaRateLimiter{
		client: client,
		script: redis.NewScript(rateLimitScript),
		prefix: "limiter",
		limit:  int64(limit),
		window: int64(window),
	}
}

// Allow
func (r *RedisLuaRateLimiter) Allow(key string) (bool, error) {
	keyName := fmt.Sprintf("%s:%s", r.prefix, key)
	now := time.Now().Unix()

	result, err := r.script.Run(context.Background(), r.client, []string{keyName},
		r.limit, r.window, now).Result()
	if err != nil {
		return false, err
	}

	return result.(int64) == 1, nil
}

// RateLimitMiddleware
func (r *RedisLuaRateLimiter) RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		allowed, err := r.Allow(clientIP)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "rate limiter error",
			})
			return
		}

		if !allowed {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":   "too many requests",
				"message": "rate limit exceeded",
			})
			return
		}

		c.Next()
	}
}
