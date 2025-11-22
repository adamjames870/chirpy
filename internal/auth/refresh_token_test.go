package auth

import (
	"encoding/hex"
	"testing"
)

func TestMakeRefreshToken_BasicProperties(t *testing.T) {
	token, err := MakeRefreshToken()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(token) != 64 {
		t.Fatalf("expected hex-encoded token length 64, got %d", len(token))
	}

	// Check valid hex
	decoded, err := hex.DecodeString(token)
	if err != nil {
		t.Fatalf("token is not valid hex: %v", err)
	}

	if len(decoded) != 32 {
		t.Fatalf("expected decoded length 32 bytes, got %d", len(decoded))
	}
}

func TestMakeRefreshToken_Uniqueness(t *testing.T) {
	const n = 1000
	seen := make(map[string]struct{})

	for i := 0; i < n; i++ {
		tok, err := MakeRefreshToken()
		if err != nil {
			t.Fatalf("unexpected error at iteration %d: %v", i, err)
		}

		if _, exists := seen[tok]; exists {
			t.Fatalf("duplicate token generated at iteration %d: %s", i, tok)
		}
		seen[tok] = struct{}{}
	}
}
