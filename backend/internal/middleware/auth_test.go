package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blueship581/gbinsureapi/internal/util"
	"github.com/gin-gonic/gin"
)

func TestJwtAuthInvalidTokenNoPanicP802(t *testing.T) {
	gin.SetMode(gin.TestMode)
	expired, err := util.GenerateToken("secret", -1, 1, "his", "settlement", "service")
	if err != nil {
		t.Fatalf("GenerateToken error = %v", err)
	}
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(JwtAuthRequired("secret"))
	r.GET("/api/v1/settlements/:settlement_no", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/settlements/S1", nil)
	req.Header.Set("Authorization", "Bearer "+expired)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401 for expired token", rec.Code)
	}
}

func TestAdminRequiredNilClaimsP804(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(func(c *gin.Context) {
		// 直接注入 typed-nil claims，模拟解析层返回 nil 无错误
		c.Set(ClientKey, (*util.Claims)(nil))
		c.Next()
	})
	r.Use(AdminRequired())
	r.GET("/api/v1/clients", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/clients", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403 for nil claims (no panic)", rec.Code)
	}
}
