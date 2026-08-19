package infra

import (
	"testing"

	"gin-layout/config"
	"gin-layout/internal/token"
)

func TestJWTIssuer_Issue_SetsTokenTypes(t *testing.T) {
	issuer := NewJWTIssuer(&config.JWTConfig{
		Secret:         "test-secret",
		AccessExpired:  5,
		RefreshExpired: 60,
	})

	pair, err := issuer.Issue(1, "admin")
	if err != nil {
		t.Fatalf("Issue() error = %v", err)
	}

	access, err := issuer.Parse(pair.AccessToken)
	if err != nil {
		t.Fatalf("Parse(access) error = %v", err)
	}
	if access.Type != token.TypeAccess {
		t.Fatalf("access token type = %q, want %q", access.Type, token.TypeAccess)
	}

	refresh, err := issuer.Parse(pair.RefreshToken)
	if err != nil {
		t.Fatalf("Parse(refresh) error = %v", err)
	}
	if refresh.Type != token.TypeRefresh {
		t.Fatalf("refresh token type = %q, want %q", refresh.Type, token.TypeRefresh)
	}
}
