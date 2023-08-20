package repository

import (
	"context"

	"github.com/4epyx/testtask/model"
)

// RefreshTokenRepository defines methods for storing refresh tokens in database
type RefreshTokenRepository interface {
	GetTokenById(context.Context, string) (model.RefreshToken, error)
	CreateToken(context.Context, model.RefreshToken) error
	DeleteToken(context.Context, string) error
}
