package models

import (
	"lgc/pkg/hasher"

	"github.com/google/uuid"
)

const NameMinLength = 3

const PasswordMinLength = 6

type User struct {
	Id       string
	UserName string
	Password string
}

func NewUser(name, password string) (*User, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	userId := id.String()
	userPassword, err := hasher.HashPassword(password)
	if err != nil {
		return nil, err
	}
	return &User{
		Id:       userId,
		UserName: name,
		Password: userPassword,
	}, nil
}
