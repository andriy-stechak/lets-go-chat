package repositories

import (
	"context"
	"errors"
	"testing"

	"github.com/andriystech/lgc/facilities/mongo"
	"github.com/andriystech/lgc/mocks"
	"github.com/andriystech/lgc/models"
	"github.com/stretchr/testify/assert"
)

func TestSaveUserSuccess(t *testing.T) {
	ctx := context.Background()
	usr := &models.User{UserName: "foo"}
	c := new(mocks.CollectionHelper)
	srh := new(mocks.SingleResultHelper)
	c.On("FindOne", ctx, map[string]string{"userName": usr.UserName}).Return(srh, nil)
	c.On("InsertOne", ctx, usr).Return("1", nil)
	srh.On("Decode", &models.User{}).Return(ErrUserNotFound)
	repo := NewUsersRepository(c)

	gotId, gotErr := repo.SaveUser(ctx, usr)

	assert.Nil(t, gotErr, "SaveUser returned unexpected result: got error %v", gotErr)
	assert.Equal(t, "1", gotId, "SaveUser returned unexpected result: got %v want 1", gotId)

	c.AssertExpectations(t)
	srh.AssertExpectations(t)
}

func TestSaveUserAlreadyExist(t *testing.T) {
	ctx := context.Background()
	usr := &models.User{UserName: "foo"}
	c := new(mocks.CollectionHelper)
	srh := new(mocks.SingleResultHelper)
	c.On("FindOne", ctx, map[string]string{"userName": usr.UserName}).Return(srh, nil)
	srh.On("Decode", &models.User{}).Return(nil)
	repo := NewUsersRepository(c)

	gotId, gotErr := repo.SaveUser(ctx, usr)

	assert.Equal(t, "", gotId, "SaveUser returned unexpected result: got user id %v instead of %v", gotId, "")
	assert.Equal(t, ErrUserWithNameAlreadyExists, gotErr, "SaveUser returned unexpected result: got success instead of %v", ErrUserWithNameAlreadyExists)

	c.AssertExpectations(t)
	srh.AssertExpectations(t)
}

func TestSaveUserFailToInsert(t *testing.T) {
	ctx := context.Background()
	usr := &models.User{UserName: "foo"}
	c := new(mocks.CollectionHelper)
	srh := new(mocks.SingleResultHelper)
	wantError := errors.New("Unable to save")
	c.On("FindOne", ctx, map[string]string{"userName": usr.UserName}).Return(srh, nil)
	c.On("InsertOne", ctx, usr).Return(nil, wantError)
	srh.On("Decode", &models.User{}).Return(ErrUserNotFound)
	repo := NewUsersRepository(c)

	gotId, gotErr := repo.SaveUser(ctx, usr)

	assert.Equal(t, "", gotId, "SaveUser returned unexpected result: got user id %v instead of %v", gotId, "")
	assert.Equal(t, wantError, gotErr, "SaveUser returned unexpected result: got success instead of %v", ErrUserWithNameAlreadyExists)

	c.AssertExpectations(t)
	srh.AssertExpectations(t)
}

func TestFindUserByNameNoItems(t *testing.T) {
	ctx := context.Background()
	usr := &models.User{UserName: "foo"}
	c := new(mocks.CollectionHelper)
	srh := new(mocks.SingleResultHelper)
	c.On("FindOne", ctx, map[string]string{"userName": usr.UserName}).Return(srh, nil)
	srh.On("Decode", &models.User{}).Return(mongo.ErrNoDocuments)
	repo := NewUsersRepository(c)

	gotUsr, gotErr := repo.FindUserByName(ctx, usr.UserName)

	assert.Nil(t, gotUsr, "FindUserByName returned unexpected result: got user %v want %v", gotUsr, nil)
	assert.Equal(t, ErrUserNotFound, gotErr, "FindUserByName returned unexpected result: got success instead of %v", ErrUserNotFound)

	c.AssertExpectations(t)
	srh.AssertExpectations(t)
}
