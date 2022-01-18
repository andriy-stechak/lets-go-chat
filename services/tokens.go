package services

import (
	"context"

	"github.com/andriystech/lgc/db/repositories"
	"github.com/andriystech/lgc/models"
	"github.com/google/uuid"
)

type TokenService interface {
	GenerateToken(context.Context, *models.User) (*models.Token, error)
	GetUserByToken(context.Context, string) (*models.User, error)
}

type TokenServiceContainer struct {
	storage repositories.TokensRepository
}

func NewTokenService(storage repositories.TokensRepository) TokenService {
	return &TokenServiceContainer{
		storage: storage,
	}
}

func (svc *TokenServiceContainer) GenerateToken(ctx context.Context, user *models.User) (*models.Token, error) {
	token := models.NewToken(uuid.NewString())
	err := svc.storage.SaveToken(ctx, token.Payload, user)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (svc *TokenServiceContainer) GetUserByToken(ctx context.Context, token string) (*models.User, error) {
	user, err := svc.storage.GetUserByToken(ctx, token)
	if err != nil {
		return nil, err
	}
	return user, nil
}
