package services

import (
	"context"
	"errors"
	"testing"

	"github.com/andriystech/lgc/mocks"
	"github.com/andriystech/lgc/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGenerateTokenSuccess(t *testing.T) {
	ctx := context.TODO()
	tr := new(mocks.TokensRepository)
	usr := &models.User{}
	tr.On("SaveToken", ctx, mock.Anything, usr).Return(nil)
	svc := NewTokenService(tr)

	_, gotErr := svc.GenerateToken(ctx, usr)

	assert.Equal(t, nil, gotErr, "GenerateToken returned unexpected result: got error %v want %v", gotErr, nil)
}

func TestGenerateTokenFail(t *testing.T) {
	ctx := context.TODO()
	tr := new(mocks.TokensRepository)
	usr := &models.User{}
	wantErr := errors.New("Some error")
	tr.On("SaveToken", ctx, mock.Anything, usr).Return(wantErr)
	svc := NewTokenService(tr)

	_, gotErr := svc.GenerateToken(ctx, usr)

	assert.Equal(t, wantErr, gotErr, "GenerateToken returned unexpected result: got error %v want %v", gotErr, wantErr)
}

func TestGetUserByTokenSuccess(t *testing.T) {
	fakeUuid := "14ef71b2-5d7c-11ec-a0f3-c46516a4fa45"
	ctx := context.TODO()
	tr := new(mocks.TokensRepository)
	usr := &models.User{UserName: "foo"}
	tr.On("GetUserByToken", ctx, fakeUuid).Return(usr, nil)
	svc := NewTokenService(tr)

	gotUsr, gotErr := svc.GetUserByToken(ctx, fakeUuid)

	assert.Nil(t, gotErr, "GetUserByToken returned unexpected result: got error %v want %v", gotErr, nil)
	assert.Equal(t, usr, gotUsr, "GetUserByToken returned unexpected result: got user %v want %v", gotUsr, usr)
}

func TestGetUserByTokenFail(t *testing.T) {
	fakeUuid := "14ef71b2-5d7c-11ec-a0f3-c46516a4fa45"
	ctx := context.TODO()
	tr := new(mocks.TokensRepository)
	err := errors.New("Some error")
	tr.On("GetUserByToken", ctx, fakeUuid).Return(nil, err)
	svc := NewTokenService(tr)

	gotUsr, gotErr := svc.GetUserByToken(ctx, fakeUuid)

	assert.Equal(t, err, gotErr, "GetUserByToken returned unexpected result: got error %v want %v", gotErr, err)
	assert.Nil(t, gotUsr, "GetUserByToken returned unexpected result: got user %v want %v", gotUsr, nil)
}
