package services

import (
	"context"
	"testing"

	"github.com/andriystech/lgc/mocks"
	"github.com/andriystech/lgc/models"
	"github.com/stretchr/testify/assert"
)

func TestNewUserSuccess(t *testing.T) {
	ur := new(mocks.UsersRepository)
	svc := NewUserService(ur)

	_, gotErr := svc.NewUser("foo", "bar")

	assert.Nil(t, gotErr, "NewUser returned unexpected result: got error %v want %v", gotErr, nil)
}

func TestFindUserByNameSuccess(t *testing.T) {
	ctx := context.TODO()
	ur := new(mocks.UsersRepository)
	usr := &models.User{UserName: "foo"}
	ur.On("FindUserByName", ctx, "foo").Return(usr, nil)
	svc := NewUserService(ur)

	gotUsr, gotErr := svc.FindUserByName(ctx, "foo")

	assert.Nil(t, gotErr, "GenerateToken returned unexpected result: got error %v want %v", gotErr, nil)
	assert.Equal(t, usr, gotUsr, "GenerateToken returned unexpected result: got user %v want %v", gotUsr, usr)
}

func TestSaveUserSuccess(t *testing.T) {
	wantId := "1"
	ctx := context.TODO()
	ur := new(mocks.UsersRepository)
	usr := &models.User{UserName: "foo"}
	ur.On("SaveUser", ctx, usr).Return(wantId, nil)
	svc := NewUserService(ur)

	gotUsrId, gotErr := svc.SaveUser(ctx, usr)

	assert.Nil(t, gotErr, "SaveUser returned unexpected result: got error %v want %v", gotErr, nil)
	assert.Equal(t, wantId, gotUsrId, "SaveUser returned unexpected result: got user %v want %v", gotUsrId, wantId)
}
