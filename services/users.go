package services

import (
	"context"

	"github.com/andriystech/lgc/db/repositories"
	"github.com/andriystech/lgc/models"
	"github.com/andriystech/lgc/pkg/hasher"
	"github.com/google/uuid"
)

type UserService interface {
	NewUser(string, string) (*models.User, error)
	FindUserByName(context.Context, string) (*models.User, error)
	SaveUser(context.Context, *models.User) (string, error)
}

type userService struct {
	storage repositories.UsersRepository
}

func NewUserService(storage repositories.UsersRepository) UserService {
	return &userService{
		storage: storage,
	}
}

func (svc *userService) NewUser(name, password string) (*models.User, error) {
	userId := uuid.NewString()
	passwordHash, err := hasher.HashPassword(password)
	if err != nil {
		return nil, err
	}
	return models.NewUser(userId, name, passwordHash), nil
}

func (svc *userService) FindUserByName(ctx context.Context, name string) (*models.User, error) {
	return svc.storage.FindUserByName(ctx, name)
}

func (svc *userService) SaveUser(ctx context.Context, user *models.User) (string, error) {
	return svc.storage.SaveUser(ctx, user)
}
