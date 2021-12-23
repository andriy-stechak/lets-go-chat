package repositories

import (
	"context"
	"testing"

	"github.com/andriystech/lgc/models"
	"github.com/stretchr/testify/assert"
)

func TestSaveTokenSuccess(t *testing.T) {
	ctx := context.Background()
	repo := NewTokensRepository(10)

	gotErr := repo.SaveToken(ctx, "token", &models.User{})

	assert.Nil(t, gotErr, "SaveToken returned unexpected result: got %v want %v", gotErr, nil)
}

func TestGetUserByTokenSuccess(t *testing.T) {
	wantUsr := &models.User{Id: "someid"}
	ctx := context.Background()
	repo := NewTokensRepository(10)
	repo.SaveToken(ctx, "token", &models.User{Id: "someid"})

	gotUsr, gotErr := repo.GetUserByToken(ctx, "token")

	assert.Nil(t, gotErr, "GetUserByToken returned unexpected result: got error %v want %v", gotErr, nil)
	assert.Equal(t, gotUsr, wantUsr, "Invalid user returned: got %v want %v", gotUsr, wantUsr)
}

func TestGetUserByTokenNotFound(t *testing.T) {
	ctx := context.Background()
	repo := NewTokensRepository(10)
	repo.SaveToken(ctx, "token", &models.User{Id: "someid"})

	gotUsr, gotErr := repo.GetUserByToken(ctx, "token2")

	assert.Nil(t, gotUsr, "GetUserByToken should throw error %v instead of success", ErrTokenNotFound)
	assert.Equal(t, ErrTokenNotFound, gotErr, "Should throw not found error: got %v want %v", gotErr, ErrTokenNotFound)
}

func TestGetUserByTokenExpired(t *testing.T) {
	ctx := context.Background()
	repo := NewTokensRepository(-1)
	repo.SaveToken(ctx, "token", &models.User{Id: "someid"})

	gotUsr, gotErr := repo.GetUserByToken(ctx, "token")

	assert.Nil(t, gotUsr, "GetUserByToken should throw error %v instead of success", ErrTokenExpired)
	assert.Equal(t, ErrTokenExpired, gotErr, "Should throw token expired error: got %v want %v", gotErr, ErrTokenExpired)
}
