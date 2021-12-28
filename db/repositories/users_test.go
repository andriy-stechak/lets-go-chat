package repositories

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/andriystech/lgc/facilities/mongo"
	"github.com/andriystech/lgc/mocks"
	"github.com/andriystech/lgc/models"
	"github.com/stretchr/testify/assert"
)

type testSaveUsersData struct {
	usr          *models.User
	wantId       string
	wantErr      error
	prepareMocks func(*mocks.CollectionHelper, *mocks.SingleResultHelper)
}

func TestSaveUser(t *testing.T) {
	fakeUsr := &models.User{UserName: "foo"}
	unknownErr := errors.New("Unable to save")
	testConditions := []testSaveUsersData{
		{
			usr:     fakeUsr,
			wantId:  "1",
			wantErr: nil,
			prepareMocks: func(ch *mocks.CollectionHelper, srh *mocks.SingleResultHelper) {
				ctx := context.Background()
				ch.On("FindOne", ctx, map[string]string{"userName": fakeUsr.UserName}).Return(srh, nil)
				ch.On("InsertOne", ctx, fakeUsr).Return("1", nil)
				srh.On("Decode", &models.User{}).Return(ErrUserNotFound)
			},
		},
		{
			usr:     fakeUsr,
			wantId:  "",
			wantErr: ErrUserWithNameAlreadyExists,
			prepareMocks: func(ch *mocks.CollectionHelper, srh *mocks.SingleResultHelper) {
				ctx := context.Background()
				ch.On("FindOne", ctx, map[string]string{"userName": fakeUsr.UserName}).Return(srh, nil)
				srh.On("Decode", &models.User{}).Return(nil)
			},
		},
		{
			usr:     fakeUsr,
			wantId:  "",
			wantErr: unknownErr,
			prepareMocks: func(ch *mocks.CollectionHelper, srh *mocks.SingleResultHelper) {
				ctx := context.Background()
				ch.On("FindOne", ctx, map[string]string{"userName": fakeUsr.UserName}).Return(srh, nil)
				ch.On("InsertOne", ctx, fakeUsr).Return(nil, unknownErr)
				srh.On("Decode", &models.User{}).Return(ErrUserNotFound)
			},
		},
	}

	for _, testCond := range testConditions {
		tName := fmt.Sprintf("SaveUser(%v, %v) == %v, %v", context.Background(), testCond.usr, testCond.wantId, testCond.wantErr)
		t.Run(tName, func(t *testing.T) {
			ctx := context.Background()
			ch := new(mocks.CollectionHelper)
			srh := new(mocks.SingleResultHelper)
			testCond.prepareMocks(ch, srh)
			repo := NewUsersRepository(ch)

			gotId, gotErr := repo.SaveUser(ctx, testCond.usr)

			assert.Equal(t, testCond.wantErr, gotErr, "SaveUser returned unexpected result: got error %v want %v", gotErr, testCond.wantErr)
			assert.Equal(t, testCond.wantId, gotId, "SaveUser returned unexpected result: got Id %v want %v", gotId, testCond.wantId)

			ch.AssertExpectations(t)
			srh.AssertExpectations(t)
		})
	}
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
