package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/4epyx/testtask/model"
	"github.com/4epyx/testtask/repository"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// TokenService may be used for work with tokens (get pair of tokens or update tokens)
type TokenService struct {
	// repo is an instance of refresh token repository
	repo repository.RefreshTokenRepository
	// accessTokenTTL is a time to live  of access token
	accessTokenTTL time.Duration
	// refershTokenTTL is a time to live of refresh token
	refreshTokenTTL time.Duration
	// jwtSecretSign is a secret string, which jwt will sign
	jwtSecretKey []byte
}

// NewTokenService creates a new instance of token service
func NewTokenService(repo repository.RefreshTokenRepository) *TokenService {
	return &TokenService{repo: repo}
}

func (s *TokenService) GenerateAccessToken(ctx context.Context, userGuid string) (string, error) {
	claims := jwt.MapClaims{
		"user_guid": userGuid,
		"exp":       time.Now().Add(s.accessTokenTTL).Unix(),
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS512, claims).SignedString(s.jwtSecretKey)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (s *TokenService) GenerateRefreshToken(ctx context.Context, userGuid string) (string, error) {
	tokenUuid := uuid.NewString()
	tokenExpiration := time.Now().Add(s.refreshTokenTTL).Unix()

	claims, err := json.Marshal(map[string]interface{}{
		"id":        tokenUuid,
		"user_guid": userGuid,
		"exp":       tokenExpiration,
	})
	if err != nil {
		return "", err
	}

	encoded := base64.StdEncoding.EncodeToString(claims)
	hash, err := bcrypt.GenerateFromPassword([]byte(encoded), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	tokenData := model.RefreshToken{
		Id:         tokenUuid,
		UserGuid:   userGuid,
		Expiration: tokenExpiration,
		Token:      string(hash),
	}

	if err := s.repo.CreateToken(ctx, tokenData); err != nil {
		return "", err
	}

	return encoded, nil
}
