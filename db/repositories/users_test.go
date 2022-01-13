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
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
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

func TestFindUsersNotInIdList(t *testing.T) {
	errUnableToFind := errors.New("Unable to run find query")
	errUnableToParse := errors.New("Unable to parse result")
	ids := []string{
		"14ef71b2-5d7c-11ec-a0f3-c46516a4fa45",
		"14ef71b2-5d7c-11ec-a0f3-c46516a4fa46",
	}
	testConditions := []struct {
		tName        string
		ids          []string
		expectedErr  error
		expectedRes  []*models.User
		prepareMocks func(*mocks.CollectionHelper, *mocks.MultiResultHelper)
	}{
		{
			tName:       "should fail with unable to find error",
			ids:         ids,
			expectedErr: errUnableToFind,
			expectedRes: nil,
			prepareMocks: func(ch *mocks.CollectionHelper, mrh *mocks.MultiResultHelper) {
				ch.On("Find", mock.Anything, bson.M{"_id": bson.M{"$nin": ids}}).Return(nil, errUnableToFind)
			},
		},
		{
			tName:       "should return empty list when no documents found",
			ids:         ids,
			expectedErr: nil,
			expectedRes: []*models.User(nil),
			prepareMocks: func(ch *mocks.CollectionHelper, mrh *mocks.MultiResultHelper) {
				mrh.On("All", mock.Anything, mock.Anything).Return(mongo.ErrNoDocuments)
				ch.On("Find", mock.Anything, bson.M{"_id": bson.M{"$nin": ids}}).Return(mrh, nil)
			},
		},
		{
			tName:       "should fail with unable to parse result error",
			ids:         ids,
			expectedErr: errUnableToParse,
			expectedRes: nil,
			prepareMocks: func(ch *mocks.CollectionHelper, mrh *mocks.MultiResultHelper) {
				mrh.On("All", mock.Anything, mock.Anything).Return(errUnableToParse)
				ch.On("Find", mock.Anything, bson.M{"_id": bson.M{"$nin": ids}}).Return(mrh, nil)
			},
		},
		{
			tName:       "should succesfully return list of users",
			ids:         ids,
			expectedErr: nil,
			expectedRes: []*models.User(nil),
			prepareMocks: func(ch *mocks.CollectionHelper, mrh *mocks.MultiResultHelper) {
				mrh.On("All", mock.Anything, mock.Anything).Return(nil)
				ch.On("Find", mock.Anything, bson.M{"_id": bson.M{"$nin": ids}}).Return(mrh, nil)
			},
		},
	}
	for _, testCond := range testConditions {
		t.Run(testCond.tName, func(t *testing.T) {
			ctx := context.Background()
			ch := new(mocks.CollectionHelper)
			mrh := new(mocks.MultiResultHelper)

			testCond.prepareMocks(ch, mrh)
			repo := NewUsersRepository(ch)

			gotRes, gotErr := repo.FindUsersNotInIdList(ctx, testCond.ids)

			assert.Equal(t, testCond.expectedErr, gotErr, "FindUsersNotInIdList returned unexpected error: got error %v want %v", gotErr, testCond.expectedErr)
			assert.Equal(t, testCond.expectedRes, gotRes, "FindUsersNotInIdList returned unexpected result: got %v want %v", gotRes, testCond.expectedRes)

			ch.AssertExpectations(t)
			mrh.AssertExpectations(t)
		})
	}
}
