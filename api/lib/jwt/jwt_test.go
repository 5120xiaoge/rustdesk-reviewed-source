package jwt

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const testKey = "non-secret-test-key"

func TestGenerateToken(t *testing.T) {
	if NewJwt(testKey, time.Hour).GenerateToken(1) == "" {
		t.Fatal("token generation failed")
	}
}

func TestParseToken(t *testing.T) {
	service := NewJwt(testKey, time.Hour)
	uid, err := service.ParseToken(service.GenerateToken(999))
	if err != nil || uid != 999 {
		t.Fatalf("unexpected token result: user=%d, error=%v", uid, err)
	}
}

func TestRejectInvalidClaims(t *testing.T) {
	future := jwt.NewNumericDate(time.Now().Add(time.Hour))
	cases := []struct {
		name   string
		method jwt.SigningMethod
		claims UserClaims
	}{
		{"missing expiry", jwt.SigningMethodHS256, UserClaims{UserId: 1}},
		{"expired", jwt.SigningMethodHS256, UserClaims{UserId: 1, RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour))}}},
		{"zero user", jwt.SigningMethodHS256, UserClaims{RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: future}}},
		{"wrong algorithm", jwt.SigningMethodHS384, UserClaims{UserId: 1, RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: future}}},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			token, err := jwt.NewWithClaims(test.method, test.claims).SignedString([]byte(testKey))
			if err != nil {
				t.Fatal(err)
			}
			if _, err := NewJwt(testKey, time.Hour).ParseToken(token); err == nil {
				t.Fatal("invalid token was accepted")
			}
		})
	}
}

func TestRejectEmptyKey(t *testing.T) {
	claims := UserClaims{UserId: 1, RegisteredClaims: jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
	}}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte{})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := NewJwt("", time.Hour).ParseToken(token); err == nil {
		t.Fatal("token with an empty signing key was accepted")
	}
}

func BenchmarkJwtService_GenerateToken(b *testing.B) {
	service := NewJwt(testKey, time.Hour)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.GenerateToken(999)
	}
}

func BenchmarkJwtService_ParseToken(b *testing.B) {
	service := NewJwt(testKey, time.Hour)
	token := service.GenerateToken(999)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.ParseToken(token)
	}
}
