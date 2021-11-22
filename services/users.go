package services

import (
	"context"

	"github.com/andriystech/lgc/db/repositories"
	"github.com/andriystech/lgc/models"
	"github.com/andriystech/lgc/pkg/hasher"
	"github.com/google/uuid"
)

type UserService struct {
	storage *repositories.UsersRepository
}

func NewUserService(storage *repositories.UsersRepository) *UserService {
	return &UserService{
		storage: storage,
	}
}

func (svc *UserService) NewUser(name, password string) (*models.User, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	userId := id.String()
	passwordHash, err := hasher.HashPassword(password)
	if err != nil {
		return nil, err
	}
	return models.NewUser(userId, name, passwordHash), nil
}

func (svc *UserService) FindUserByName(ctx context.Context, name string) (*models.User, error) {
	return svc.storage.FindUserByName(ctx, name)
}

func (svc *UserService) SaveUser(ctx context.Context, user *models.User) (string, error) {
	return svc.storage.SaveUser(ctx, user)
}
