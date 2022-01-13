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

func TestSaveMessage(t *testing.T) {
	fakeMsg := &models.Message{Payload: "hello"}
	unknownErr := errors.New("Unable to save")
	testConditions := []struct {
		msg          *models.Message
		wantId       string
		wantErr      error
		prepareMocks func(*mocks.CollectionHelper)
	}{
		{
			msg:     fakeMsg,
			wantId:  "1",
			wantErr: nil,
			prepareMocks: func(ch *mocks.CollectionHelper) {
				ctx := context.Background()
				ch.On("InsertOne", ctx, fakeMsg).Return("1", nil)
			},
		},
		{
			msg:     fakeMsg,
			wantId:  "",
			wantErr: unknownErr,
			prepareMocks: func(ch *mocks.CollectionHelper) {
				ctx := context.Background()
				ch.On("InsertOne", ctx, fakeMsg).Return(nil, unknownErr)
			},
		},
	}

	for _, testCond := range testConditions {
		tName := fmt.Sprintf("SaveMessage(%v, %v) == %v, %v", context.Background(), testCond.msg, testCond.wantId, testCond.wantErr)
		t.Run(tName, func(t *testing.T) {
			ctx := context.Background()
			ch := new(mocks.CollectionHelper)
			testCond.prepareMocks(ch)
			repo := NewMessagesRepository(ch)

			gotId, gotErr := repo.SaveMessage(ctx, testCond.msg)

			assert.Equal(t, testCond.wantErr, gotErr, "SaveMessage returned unexpected result: got error %v want %v", gotErr, testCond.wantErr)
			assert.Equal(t, testCond.wantId, gotId, "SaveMessage returned unexpected result: got Id %v want %v", gotId, testCond.wantId)

			ch.AssertExpectations(t)
		})
	}
}

func TestFindUserMessages(t *testing.T) {
	errUnableToFind := errors.New("Unable to run find query")
	errUnableToParse := errors.New("Unable to parse result")
	id := "14ef71b2-5d7c-11ec-a0f3-c46516a4fa45"
	testConditions := []struct {
		tName        string
		id           string
		expectedErr  error
		expectedRes  []*models.Message
		prepareMocks func(*mocks.CollectionHelper, *mocks.MultiResultHelper)
	}{
		{
			tName:       "should fail with unable to find error",
			id:          id,
			expectedErr: errUnableToFind,
			expectedRes: nil,
			prepareMocks: func(ch *mocks.CollectionHelper, mrh *mocks.MultiResultHelper) {
				ch.On("Find", mock.Anything, bson.M{"recipientId": id}).Return(nil, errUnableToFind)
			},
		},
		{
			tName:       "should return empty list when no documents found",
			id:          id,
			expectedErr: nil,
			expectedRes: []*models.Message(nil),
			prepareMocks: func(ch *mocks.CollectionHelper, mrh *mocks.MultiResultHelper) {
				mrh.On("All", mock.Anything, mock.Anything).Return(mongo.ErrNoDocuments)
				ch.On("Find", mock.Anything, bson.M{"recipientId": id}).Return(mrh, nil)
			},
		},
		{
			tName:       "should fail with unable to parse result error",
			id:          id,
			expectedErr: errUnableToParse,
			expectedRes: nil,
			prepareMocks: func(ch *mocks.CollectionHelper, mrh *mocks.MultiResultHelper) {
				mrh.On("All", mock.Anything, mock.Anything).Return(errUnableToParse)
				ch.On("Find", mock.Anything, bson.M{"recipientId": id}).Return(mrh, nil)
			},
		},
		{
			tName:       "should succesfully return list of messages",
			id:          id,
			expectedErr: nil,
			expectedRes: []*models.Message(nil),
			prepareMocks: func(ch *mocks.CollectionHelper, mrh *mocks.MultiResultHelper) {
				mrh.On("All", mock.Anything, mock.Anything).Return(nil)
				ch.On("Find", mock.Anything, bson.M{"recipientId": id}).Return(mrh, nil)
			},
		},
	}
	for _, testCond := range testConditions {
		t.Run(testCond.tName, func(t *testing.T) {
			ctx := context.Background()
			ch := new(mocks.CollectionHelper)
			mrh := new(mocks.MultiResultHelper)

			testCond.prepareMocks(ch, mrh)
			repo := NewMessagesRepository(ch)

			gotRes, gotErr := repo.FindUserMessages(ctx, testCond.id)

			assert.Equal(t, testCond.expectedErr, gotErr, "FindUserMessages returned unexpected error: got error %v want %v", gotErr, testCond.expectedErr)
			assert.Equal(t, testCond.expectedRes, gotRes, "FindUserMessages returned unexpected result: got %v want %v", gotRes, testCond.expectedRes)

			ch.AssertExpectations(t)
			mrh.AssertExpectations(t)
		})
	}
}
