package services

import (
	"context"

	"github.com/andriystech/lgc/db/repositories"
	"github.com/andriystech/lgc/models"
	"github.com/google/uuid"
)

type TokenService struct {
	storage *repositories.TokensRepository
}

func NewTokenService(storage *repositories.TokensRepository) *TokenService {
	return &TokenService{
		storage: storage,
	}
}

func (svc *TokenService) GenerateToken(ctx context.Context, user *models.User) (*models.Token, error) {
	token := models.NewToken(uuid.NewString())
	err := svc.storage.SaveToken(ctx, token.Payload, user)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (svc *TokenService) GetUserByToken(ctx context.Context, token string) (*models.User, error) {
	user, err := svc.storage.GetUserByToken(ctx, token)
	if err != nil {
		return nil, err
	}
	return user, nil
}
