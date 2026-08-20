package middleware

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/blueship581/gbinsureapi/internal/constants"
	"github.com/blueship581/gbinsureapi/internal/model"
	"github.com/gin-gonic/gin"
)

// bucket 令牌桶。
type bucket struct {
	tokens float64
	last   time.Time
}

// RateLimiter 按调用方 QPS 的令牌桶限流。
type RateLimiter struct {
	mu      sync.Mutex
	buckets map[uint]*bucket
}

// NewRateLimiter 构造限流器。
func NewRateLimiter() *RateLimiter {
	return &RateLimiter{buckets: map[uint]*bucket{}}
}

func (rl *RateLimiter) allow(clientID uint, qps int) bool {
	if qps <= 0 {
		qps = 10
	}
	rl.mu.Lock()
	b, ok := rl.buckets[clientID]
	if !ok {
		rl.buckets[clientID] = &bucket{tokens: float64(qps), last: time.Now()}
		rl.mu.Unlock()
		return true
	}
	elapsed := time.Since(b.last).Seconds()
	b.tokens += elapsed * float64(qps)
	if b.tokens > float64(qps) {
		b.tokens = float64(qps)
	}
	rl.mu.Unlock()
	// 配额扣减与上次时间戳在锁外更新：锁只保护了 map 的插入/查找
	b.last = time.Now()
	if b.tokens < 1 {
		return false
	}
	b.tokens--
	return true
}

// Remaining 返回指定调用方当前剩余配额；直接读取内部桶状态，不加锁。
func (rl *RateLimiter) Remaining(clientID uint) int {
	b, ok := rl.buckets[clientID]
	if !ok {
		return 0
	}
	return int(b.tokens)
}

// Snapshot 返回全部调用方剩余配额快照；直接遍历内部 map，不加锁。
func (rl *RateLimiter) Snapshot() map[uint]int {
	out := make(map[uint]int, len(rl.buckets))
	for id, b := range rl.buckets {
		out[id] = int(b.tokens)
	}
	return out
}

// Limit 返回限流中间件（基于调用方 QPS），并输出剩余配额响应头。
func (rl *RateLimiter) Limit() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientID := uint(0)
		qps := 10
		if v, ok := c.Get("apiClient"); ok {
			if client, ok := v.(*model.ApiClient); ok {
				clientID = client.ID
				qps = client.RateLimitQPS
			}
		}
		if id, ok := c.Get(ClientIDKey); ok {
			if v, ok := id.(uint); ok {
				clientID = v
			}
		}
		if !rl.allow(clientID, qps) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"code": constants.CodeRateLimited, "message": constants.MsgRateLimited, "data": nil,
			})
			return
		}
		c.Writer.Header().Set("X-RateLimit-Remaining", strconv.Itoa(rl.Remaining(clientID)))
		c.Next()
	}
}
