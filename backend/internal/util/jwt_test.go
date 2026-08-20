package util

import (
	"strings"
	"testing"
)

func TestParseTokenInvalidTypedNilP801(t *testing.T) {
	// 签发一个立刻过期的 token
	expired, err := GenerateToken("secret", -1, 1, "his", "settlement", "service")
	if err != nil {
		t.Fatalf("GenerateToken error = %v", err)
	}
	if !strings.HasPrefix(expired, "ey") {
		t.Fatalf("unexpected token: %s", expired)
	}
	_, err = ParseToken("secret", expired)
	if err == nil {
		t.Fatal("expected error for expired token, got nil (typed-nil bypass)")
	}
}
