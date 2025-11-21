package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

func TestMakeAndValidateJWT(t *testing.T) {
	secret := "supersecret"
	userID := uuid.New()

	t.Run("valid token", func(t *testing.T) {
		tokenStr, err := MakeJWT(userID, secret, time.Minute)
		if err != nil {
			t.Fatalf("makeJWT returned error: %v", err)
		}

		gotID, err := ValidateJWT(tokenStr, secret)
		if err != nil {
			t.Fatalf("ValidateJWT returned unexpected error: %v", err)
		}

		if gotID != userID {
			t.Fatalf("expected userID %s, got %s", userID, gotID)
		}
	})

	t.Run("expired token", func(t *testing.T) {
		tokenStr, err := MakeJWT(userID, secret, -1*time.Second)
		if err != nil {
			t.Fatalf("makeJWT returned error: %v", err)
		}

		_, err = ValidateJWT(tokenStr, secret)
		if err == nil {
			t.Fatal("expected error for expired token, got nil")
		}
	})

	t.Run("wrong signing secret", func(t *testing.T) {
		tokenStr, err := MakeJWT(userID, secret, time.Minute)
		if err != nil {
			t.Fatalf("makeJWT returned error: %v", err)
		}

		_, err = ValidateJWT(tokenStr, "wrong-secret")
		if err == nil {
			t.Fatal("expected error for wrong secret, got nil")
		}
	})

	t.Run("bad subject (non-UUID)", func(t *testing.T) {
		claims := jwt.MapClaims{
			"sub": "not-a-uuid",
			"exp": time.Now().Add(time.Minute).Unix(),
		}

		badToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenStr, err := badToken.SignedString([]byte(secret))
		if err != nil {
			t.Fatalf("SignedString failed: %v", err)
		}

		_, err = ValidateJWT(tokenStr, secret)
		if err == nil {
			t.Fatal("expected error for invalid subject UUID, got nil")
		}
	})

	t.Run("malformed token string", func(t *testing.T) {
		_, err := ValidateJWT("this-is-not-a-jwt", secret)
		if err == nil {
			t.Fatal("expected error for malformed token, got nil")
		}
	})

	t.Run("missing subject claim", func(t *testing.T) {
		claims := jwt.MapClaims{
			"exp": time.Now().Add(time.Minute).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenStr, err := token.SignedString([]byte(secret))
		if err != nil {
			t.Fatalf("SignedString failed: %v", err)
		}

		_, err = ValidateJWT(tokenStr, secret)
		if err == nil {
			t.Fatal("expected error for missing subject, got nil")
		}
	})

	t.Run("unsigned token", func(t *testing.T) {
		claims := jwt.MapClaims{
			"sub": userID.String(),
			"exp": time.Now().Add(time.Minute).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
		tokenStr, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
		if err != nil {
			t.Fatalf("SignedString failed: %v", err)
		}

		_, err = ValidateJWT(tokenStr, secret)
		if err == nil {
			t.Fatal("expected error for unsigned token, got nil")
		}
	})
}

func TestGetBearerToken(t *testing.T) {

	t.Run("valid bearer token", func(t *testing.T) {
		h := http.Header{}
		h.Set("Authorization", "Bearer abc123")

		token, err := GetBearerToken(h)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if token != "abc123" {
			t.Fatalf("expected token 'abc123', got '%s'", token)
		}
	})

	t.Run("missing Authorization header", func(t *testing.T) {
		h := http.Header{}

		_, err := GetBearerToken(h)
		if err == nil {
			t.Fatal("expected error for missing Authorization header, got nil")
		}
	})

	t.Run("empty Authorization header", func(t *testing.T) {
		h := http.Header{}
		h.Set("Authorization", "")

		_, err := GetBearerToken(h)
		if err == nil {
			t.Fatal("expected error for empty Authorization header, got nil")
		}
	})

	t.Run("missing token after Bearer", func(t *testing.T) {
		h := http.Header{}
		h.Set("Authorization", "Bearer")

		_, err := GetBearerToken(h)
		if err == nil {
			t.Fatal("expected error for missing token, got nil")
		}
	})

	t.Run("Bearer with extra spaces but valid token", func(t *testing.T) {
		h := http.Header{}
		h.Set("Authorization", "Bearer    xyz")

		token, err := GetBearerToken(h)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if token != "xyz" {
			t.Fatalf("expected 'xyz', got '%s'", token)
		}
	})

	t.Run("multiple Authorization headers (HTTP allows this)", func(t *testing.T) {
		h := http.Header{}
		h.Add("Authorization", "Bearer abc")
		h.Add("Authorization", "Bearer def")

		// expected behavior:
		// Most robust implementations use headers.Get() → first value.
		// If your implementation concatenates multiple headers, adjust this test.
		token, err := GetBearerToken(h)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if token != "abc" {
			t.Fatalf("expected first header token 'abc', got '%s'", token)
		}
	})

	t.Run("Bearer with leading/trailing spaces", func(t *testing.T) {
		h := http.Header{}
		h.Set("Authorization", "   Bearer tok999   ")

		// Depends on your implementation.
		// This test expects trimming of the entire header string, which is typical.
		token, err := GetBearerToken(h)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if token != "tok999" {
			t.Fatalf("expected 'tok999', got '%s'", token)
		}
	})
}
