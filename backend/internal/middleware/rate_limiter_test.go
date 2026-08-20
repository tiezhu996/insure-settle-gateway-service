package middleware

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"sync"
	"testing"

	"github.com/blueship581/gbinsureapi/internal/constants"
	"github.com/gin-gonic/gin"
)

// startBarrier 返回一个 start channel 同步屏障：所有参与者就绪后 close 同时释放。
func startBarrier() (chan struct{}, *sync.WaitGroup) {
	start := make(chan struct{})
	var wg sync.WaitGroup
	return start, &wg
}

func TestRateLimiterConcurrentSnapshotP101(t *testing.T) {
	rl := NewRateLimiter()
	start, wg := startBarrier()
	const workers = 6
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()
			<-start
			for j := 0; j < 200; j++ {
				if w%2 == 0 {
					rl.allow(uint(1+w%3), 10)
				} else {
					_ = rl.Remaining(uint(1 + w%3))
					_ = rl.Snapshot()
				}
			}
		}(i)
	}
	close(start)
	wg.Wait()
	if got := rl.Remaining(1); got < 0 || got > 10 {
		t.Fatalf("Remaining(1) = %d, want in [0,10]", got)
	}
}

func TestRateLimiterRemainingHeaderStableP102(t *testing.T) {
	gin.SetMode(gin.TestMode)
	log := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	rl := NewRateLimiter()
	engine := gin.New()
	engine.Use(RequestLogger(log, rl))
	engine.Use(rl.Limit())
	engine.GET("/api/v1/settlements", func(c *gin.Context) {
		c.Set(ClientIDKey, uint(9))
		c.JSON(http.StatusOK, gin.H{"code": constants.CodeOK})
	})
	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			for j := 0; j < 80; j++ {
				req := httptest.NewRequest(http.MethodGet, "/api/v1/settlements", nil)
				rec := httptest.NewRecorder()
				engine.ServeHTTP(rec, req)
				if rec.Code == http.StatusTooManyRequests {
					continue
				}
				h := rec.Header().Get("X-RateLimit-Remaining")
				v, err := strconv.Atoi(h)
				if err != nil {
					t.Errorf("bad remaining header %q: %v", h, err)
					return
				}
				if v < 0 || v > 10 {
					t.Errorf("remaining header %d out of range", v)
					return
				}
			}
		}()
	}
	close(start)
	wg.Wait()
}

func TestRateLimiterSnapshotImmutableP103(t *testing.T) {
	rl := NewRateLimiter()
	start, wg := startBarrier()
	const workers = 6
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()
			<-start
			for j := 0; j < 150; j++ {
				snap := rl.Snapshot()
				// 快照副本被改写不得影响限流器内部状态
				for id := range snap {
					snap[id] = -999
					delete(snap, id)
				}
				rl.allow(1, 10)
			}
		}(i)
	}
	close(start)
	wg.Wait()
	if got := rl.Remaining(1); got < 0 || got > 10 {
		t.Fatalf("Remaining(1) polluted by snapshot = %d", got)
	}
}
