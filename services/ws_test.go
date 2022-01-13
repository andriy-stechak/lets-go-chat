package services

import (
	"context"
	"errors"
	"testing"

	"github.com/andriystech/lgc/facilities/ws"
	"github.com/andriystech/lgc/mocks"
	"github.com/andriystech/lgc/models"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetActiveConnectionsCount(t *testing.T) {
	ctx := context.Background()
	cr := new(mocks.ConnectionsRepository)
	mr := new(mocks.MessagesRepository)
	ur := new(mocks.UsersRepository)
	wu := new(mocks.UpgraderHelper)
	count := 1
	cr.On("CountConnections", ctx).Return(count, nil)
	svc := NewWebSocketService(cr, mr, ur, wu)

	gotCount, gotErr := svc.GetActiveConnectionsCount(ctx)

	assert.Nil(t, gotErr, "GetActiveConnectionsCount returned unexpected result: got error %v want %v", gotErr, nil)
	assert.Equal(t, count, gotCount, "GetActiveConnectionsCount returned unexpected result: got count %v want %v", gotCount, count)

	cr.AssertExpectations(t)
}

func TestGetActiveUsers(t *testing.T) {
	ctx := context.Background()
	cr := new(mocks.ConnectionsRepository)
	mr := new(mocks.MessagesRepository)
	ur := new(mocks.UsersRepository)
	wu := new(mocks.UpgraderHelper)
	clients := []string{"1-user", "2-user2"}
	cr.On("ConnectedClients", ctx).Return(clients, nil)
	svc := NewWebSocketService(cr, mr, ur, wu)

	gotClients, gotErr := svc.GetActiveUsers(ctx)

	assert.Nil(t, gotErr, "GetActiveUsers returned unexpected result: got error %v want %v", gotErr, nil)
	assert.Equal(t, clients, gotClients, "GetActiveUsers returned unexpected result: got clients %v want %v", gotClients, clients)

	cr.AssertExpectations(t)
}

func TestSendMessageToAllConnections(t *testing.T) {
	fakeMessageUuid := "14ef71b2-5d7c-11ec-a0f3-c46516a4fa44"
	sender := &models.User{Id: "14ef71b2-5d7c-11ec-a0f3-c46516a4fa45", UserName: "foo"}
	recipient := &models.User{Id: "14ef71b2-5d7c-11ec-a0f3-c46516a4fa46", UserName: "bar"}
	errorUnableToGetConnections := errors.New("Unable to get connections list")
	errorUnableToSaveMessage := errors.New("Unable to save message in database")
	errorUnableToSendMessage := errors.New("Unable to send message into websocket")
	testConditions := []struct {
		tName        string
		payload      string
		sender       *models.User
		expected     error
		prepareMocks func(
			*mocks.ConnectionsRepository,
			*mocks.MessagesRepository,
			*mocks.UsersRepository,
			*mocks.ConnHelper,
		)
	}{
		{
			tName:    "should fail with unable to get connections list error",
			payload:  "hello",
			sender:   sender,
			expected: errorUnableToGetConnections,
			prepareMocks: func(cr *mocks.ConnectionsRepository, mr *mocks.MessagesRepository, ur *mocks.UsersRepository, wc *mocks.ConnHelper) {
				cr.On("GetAllConnections", mock.Anything).Return(nil, errorUnableToGetConnections)
			},
		},
		{
			tName:    "should fail with error when unable to save message",
			payload:  "hello",
			sender:   sender,
			expected: errorUnableToSaveMessage,
			prepareMocks: func(cr *mocks.ConnectionsRepository, mr *mocks.MessagesRepository, ur *mocks.UsersRepository, wc *mocks.ConnHelper) {
				cr.On("GetAllConnections", mock.Anything).Return(map[string]ws.ConnHelper{
					sender.Id:    wc,
					recipient.Id: wc,
				}, nil)
				mr.On("SaveMessage", mock.Anything, mock.Anything).Return("", errorUnableToSaveMessage)
			},
		},
		{
			tName:    "should fail with error when unable to write message to socket",
			payload:  "hello",
			sender:   sender,
			expected: errorUnableToSendMessage,
			prepareMocks: func(cr *mocks.ConnectionsRepository, mr *mocks.MessagesRepository, ur *mocks.UsersRepository, wc *mocks.ConnHelper) {
				cr.On("GetAllConnections", mock.Anything).Return(map[string]ws.ConnHelper{
					sender.Id:    wc,
					recipient.Id: wc,
				}, nil)
				mr.On("SaveMessage", mock.Anything, mock.Anything).Return(fakeMessageUuid, nil)
				wc.On("WriteMessage", websocket.TextMessage, []byte("hello")).Return(errorUnableToSendMessage)
			},
		},
		{
			tName:    "should successfully send message to all connections",
			payload:  "hello",
			sender:   sender,
			expected: nil,
			prepareMocks: func(cr *mocks.ConnectionsRepository, mr *mocks.MessagesRepository, ur *mocks.UsersRepository, wc *mocks.ConnHelper) {
				cr.On("GetAllConnections", mock.Anything).Return(map[string]ws.ConnHelper{
					sender.Id:    wc,
					recipient.Id: wc,
				}, nil)
				mr.On("SaveMessage", mock.Anything, mock.Anything).Return(fakeMessageUuid, nil)
				wc.On("WriteMessage", websocket.TextMessage, []byte("hello")).Return(nil)
			},
		},
	}

	for _, testCond := range testConditions {
		t.Run(testCond.tName, func(t *testing.T) {
			ctx := context.Background()
			cr := new(mocks.ConnectionsRepository)
			mr := new(mocks.MessagesRepository)
			ur := new(mocks.UsersRepository)
			wu := new(mocks.UpgraderHelper)
			wc := new(mocks.ConnHelper)

			testCond.prepareMocks(cr, mr, ur, wc)
			svc := NewWebSocketService(cr, mr, ur, wu)

			gotErr := svc.SendMessageToAllConnections(ctx, testCond.payload, testCond.sender)

			assert.Equal(t, testCond.expected, gotErr, "SendMessageToAllConnections returned unexpected result: got error %v want %v", gotErr, testCond.expected)

			cr.AssertExpectations(t)
			mr.AssertExpectations(t)
			ur.AssertExpectations(t)
			wu.AssertExpectations(t)
			wc.AssertExpectations(t)
		})
	}
}

func TestLoadUserMessages(t *testing.T) {
	usr := &models.User{Id: "14ef71b2-5d7c-11ec-a0f3-c46516a4fa45", UserName: "foo"}
	msgs := []*models.Message{
		{
			Id:      "14ef71b2-5d7c-11ec-a0f3-c46516a4fa46",
			Payload: "hello",
		},
		{
			Id:      "14ef71b2-5d7c-11ec-a0f3-c46516a4fa47",
			Payload: "world",
		},
	}
	errorUnableToFindUsrMsgs := errors.New("Unable user messages")
	errorUnableToSendMessage := errors.New("Unable to send message into websocket")
	testConditions := []struct {
		tName        string
		msgs         []*models.Message
		usr          *models.User
		expected     error
		prepareMocks func(
			*mocks.MessagesRepository,
			*mocks.ConnHelper,
		)
	}{
		{
			tName:    "should fail with unable to get user messages error",
			msgs:     msgs,
			usr:      usr,
			expected: errorUnableToFindUsrMsgs,
			prepareMocks: func(mr *mocks.MessagesRepository, wc *mocks.ConnHelper) {
				mr.On("FindUserMessages", mock.Anything, usr.Id).Return(nil, errorUnableToFindUsrMsgs)
			},
		},
		{
			tName:    "should fail with unable to send messages to web socket",
			msgs:     msgs,
			usr:      usr,
			expected: errorUnableToSendMessage,
			prepareMocks: func(mr *mocks.MessagesRepository, wc *mocks.ConnHelper) {
				mr.On("FindUserMessages", mock.Anything, usr.Id).Return(msgs, nil)
				wc.On("WriteMessage", websocket.TextMessage, mock.Anything).Return(errorUnableToSendMessage)
			},
		},
		{
			tName:    "should successfully load all messages into web socket",
			msgs:     msgs,
			usr:      usr,
			expected: nil,
			prepareMocks: func(mr *mocks.MessagesRepository, wc *mocks.ConnHelper) {
				mr.On("FindUserMessages", mock.Anything, usr.Id).Return(msgs, nil)
				wc.On("WriteMessage", websocket.TextMessage, mock.Anything).Return(nil)
			},
		},
	}

	for _, testCond := range testConditions {
		t.Run(testCond.tName, func(t *testing.T) {
			ctx := context.Background()
			cr := new(mocks.ConnectionsRepository)
			mr := new(mocks.MessagesRepository)
			ur := new(mocks.UsersRepository)
			wu := new(mocks.UpgraderHelper)
			wc := new(mocks.ConnHelper)

			testCond.prepareMocks(mr, wc)
			svc := NewWebSocketService(cr, mr, ur, wu)

			gotErr := svc.LoadUserMessages(ctx, testCond.usr, wc)

			assert.Equal(t, testCond.expected, gotErr, "LoadUserMessages returned unexpected result: got error %v want %v", gotErr, testCond.expected)

			cr.AssertExpectations(t)
			mr.AssertExpectations(t)
			ur.AssertExpectations(t)
			wu.AssertExpectations(t)
			wc.AssertExpectations(t)
		})
	}
}

func TestSaveUnreadMessages(t *testing.T) {
	sender := &models.User{Id: "14ef71b2-5d7c-11ec-a0f3-c46516a4fa45", UserName: "foo"}
	recipient := &models.User{Id: "14ef71b2-5d7c-11ec-a0f3-c46516a4fa46", UserName: "bar"}
	errorUnableToGetConnections := errors.New("Unable to get connections list")
	errorUnableToFindUsrs := errors.New("Unable to find users")
	errorUnableToSaveMessage := errors.New("Unable to save message into database")
	testConditions := []struct {
		tName        string
		msg          string
		sender       *models.User
		expectedErr  error
		prepareMocks func(
			*mocks.ConnectionsRepository,
			*mocks.MessagesRepository,
			*mocks.UsersRepository,
			*mocks.ConnHelper,
		)
	}{
		{
			tName:       "should fail with unable to get connections list error",
			msg:         "hello",
			sender:      sender,
			expectedErr: errorUnableToGetConnections,
			prepareMocks: func(cr *mocks.ConnectionsRepository, mr *mocks.MessagesRepository, ur *mocks.UsersRepository, wc *mocks.ConnHelper) {
				cr.On("GetAllConnections", mock.Anything).Return(nil, errorUnableToGetConnections)
			},
		},
		{
			tName:       "should fail with unable to get users error",
			msg:         "hello",
			sender:      sender,
			expectedErr: errorUnableToFindUsrs,
			prepareMocks: func(cr *mocks.ConnectionsRepository, mr *mocks.MessagesRepository, ur *mocks.UsersRepository, wc *mocks.ConnHelper) {
				cr.On("GetAllConnections", mock.Anything).Return(map[string]ws.ConnHelper{
					sender.Id:    wc,
					recipient.Id: wc,
				}, nil)
				ur.On("FindUsersNotInIdList", mock.Anything, []string{sender.Id, recipient.Id}).Return(nil, errorUnableToFindUsrs)
			},
		},
		{
			tName:       "should fail with unable to save message into database",
			msg:         "hello",
			sender:      sender,
			expectedErr: errorUnableToSaveMessage,
			prepareMocks: func(cr *mocks.ConnectionsRepository, mr *mocks.MessagesRepository, ur *mocks.UsersRepository, wc *mocks.ConnHelper) {
				cr.On("GetAllConnections", mock.Anything).Return(map[string]ws.ConnHelper{
					sender.Id:    wc,
					recipient.Id: wc,
				}, nil)
				ur.On("FindUsersNotInIdList", mock.Anything, []string{sender.Id, recipient.Id}).Return([]*models.User{sender, recipient}, nil)
				mr.On("SaveMessage", mock.Anything, mock.Anything).Return("", errorUnableToSaveMessage)
			},
		},
		{
			tName:       "should successfully save message into database",
			msg:         "hello",
			sender:      sender,
			expectedErr: nil,
			prepareMocks: func(cr *mocks.ConnectionsRepository, mr *mocks.MessagesRepository, ur *mocks.UsersRepository, wc *mocks.ConnHelper) {
				cr.On("GetAllConnections", mock.Anything).Return(map[string]ws.ConnHelper{
					sender.Id:    wc,
					recipient.Id: wc,
				}, nil)
				ur.On("FindUsersNotInIdList", mock.Anything, []string{sender.Id, recipient.Id}).Return([]*models.User{sender, recipient}, nil)
				mr.On("SaveMessage", mock.Anything, mock.Anything).Return("14ef71b2-5d7c-11ec-a0f3-c46516a4fa45", nil)
			},
		},
	}

	for _, testCond := range testConditions {
		t.Run(testCond.tName, func(t *testing.T) {
			ctx := context.Background()
			cr := new(mocks.ConnectionsRepository)
			mr := new(mocks.MessagesRepository)
			ur := new(mocks.UsersRepository)
			wu := new(mocks.UpgraderHelper)
			wc := new(mocks.ConnHelper)

			testCond.prepareMocks(cr, mr, ur, wc)
			svc := NewWebSocketService(cr, mr, ur, wu)

			gotErr := svc.SaveUnreadMessages(ctx, testCond.sender, testCond.msg)

			assert.Equal(t, testCond.expectedErr, gotErr, "SaveUnreadMessages returned unexpected result: got error %v want %v", gotErr, testCond.expectedErr)

			cr.AssertExpectations(t)
			mr.AssertExpectations(t)
			ur.AssertExpectations(t)
			wu.AssertExpectations(t)
			wc.AssertExpectations(t)
		})
	}
}
