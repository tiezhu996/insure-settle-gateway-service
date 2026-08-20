package middleware

import (
	"log/slog"
	"time"

	"github.com/blueship581/gbinsureapi/internal/constants"
	"github.com/gin-gonic/gin"
)

// RequestLogger 请求日志：request_id、method、path、status、latency、剩余配额。
func RequestLogger(log *slog.Logger, rl *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		requestID, _ := c.Get(RequestIDKey)
		log.InfoContext(c.Request.Context(), constants.LOG_REQUEST_START,
			"request_id", requestID, "method", c.Request.Method, "path", c.Request.URL.Path)
		c.Next()
		clientID := uint(0)
		if id, ok := c.Get(ClientIDKey); ok {
			if v, ok := id.(uint); ok {
				clientID = v
			}
		}
		remaining := rl.Remaining(clientID)
		log.InfoContext(c.Request.Context(), constants.LOG_REQUEST_FINISH,
			"request_id", requestID, "method", c.Request.Method, "path", c.Request.URL.Path,
			"status", c.Writer.Status(), "latency_ms", time.Since(start).Milliseconds(),
			"rate_remaining", remaining)
	}
}
