package auth

import (
	"context"
	"fmt"
	"time"

	"gin-layout/internal/common"
	"gin-layout/internal/domain"
)

type Service struct {
	common.BaseService
	users    UserProvider
	password common.PasswordHasher
	tokens   common.TokenManager
}

func NewService(base common.BaseService, users UserProvider, password common.PasswordHasher, tokens common.TokenManager) *Service {
	return &Service{BaseService: base, users: users, password: password, tokens: tokens}
}

func (s *Service) Login(ctx context.Context, req LoginReq) (*LoginRes, error) {
	logger := s.Log(ctx)
	logger.Debug().Any("input", req).Msg("login attempt")

	u, err := s.users.FindByAccount(ctx, req.Account)
	if err != nil {
		return nil, fmt.Errorf("login: %w", domain.ErrAccountOrPassword)
	}

	if !u.Enabled {
		return nil, domain.ErrUserDisabled
	}

	if !s.password.Compare(u.PasswordHash, req.Password) {
		return nil, domain.ErrAccountOrPassword
	}

	tokens, err := s.tokens.Issue(u.ID, u.Account)
	if err != nil {
		return nil, fmt.Errorf("issue token (userID=%d): %w", u.ID, err)
	}

	_ = s.users.UpdateLastLogin(ctx, u.ID, time.Now())

	return &LoginRes{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (s *Service) Logout(ctx context.Context, _ LogoutReq) (*LogoutRes, error) {
	if _, err := s.tokens.RevokeCurrent(ctx); err != nil {
		return nil, err
	}
	return &LogoutRes{}, nil
}

func (s *Service) RefreshToken(ctx context.Context, req RefreshTokenReq) (*RefreshTokenRes, error) {
	logger := s.Log(ctx)
	logger.Debug().Any("input", req).Msg("refresh token")

	claims, err := s.tokens.Parse(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	revoked, err := s.tokens.IsRevoked(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}
	if revoked {
		return nil, domain.ErrTokenRevoked
	}

	if err := s.tokens.Revoke(ctx, req.RefreshToken, claims.UserID, claims.ExpiresAt); err != nil {
		return nil, err
	}

	tokens, err := s.tokens.Issue(claims.UserID, claims.Subject)
	if err != nil {
		return nil, err
	}
	return &RefreshTokenRes{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}
