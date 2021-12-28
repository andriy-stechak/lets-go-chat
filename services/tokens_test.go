package services

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/andriystech/lgc/mocks"
	"github.com/andriystech/lgc/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type generateTokenTestData struct {
	usr          *models.User
	wantErr      error
	prepareMocks func(*mocks.TokensRepository)
}

type getUserByTokenTestData struct {
	uuid         string
	wantUsr      *models.User
	wantErr      error
	prepareMocks func(*mocks.TokensRepository)
}

func TestGenerateToken(t *testing.T) {
	fakeUsr := &models.User{}
	fakeErr := errors.New("Unable to generate token")
	testConditions := []generateTokenTestData{
		{
			usr:     fakeUsr,
			wantErr: nil,
			prepareMocks: func(tr *mocks.TokensRepository) {
				tr.On("SaveToken", context.Background(), mock.Anything, fakeUsr).Return(nil)
			},
		},
		{
			usr:     fakeUsr,
			wantErr: fakeErr,
			prepareMocks: func(tr *mocks.TokensRepository) {
				tr.On("SaveToken", context.Background(), mock.Anything, fakeUsr).Return(fakeErr)
			},
		},
	}

	for _, testCond := range testConditions {
		tName := fmt.Sprintf("GenerateToken(%v, %v) == _, %v", context.Background(), testCond.usr, testCond.wantErr)
		t.Run(tName, func(t *testing.T) {
			ctx := context.Background()
			tr := new(mocks.TokensRepository)
			testCond.prepareMocks(tr)
			svc := NewTokenService(tr)

			_, gotErr := svc.GenerateToken(ctx, testCond.usr)

			assert.Equal(t, testCond.wantErr, gotErr, "GenerateToken returned unexpected result: got error %v want %v", gotErr, testCond.wantErr)
			tr.AssertExpectations(t)
		})
	}
}

func TestGetUserByToken(t *testing.T) {
	fakeUuid := "14ef71b2-5d7c-11ec-a0f3-c46516a4fa45"
	fakeUsr := &models.User{UserName: "foo"}
	fakeErr := errors.New("Unable to find user by token")

	testConditions := []getUserByTokenTestData{
		{
			uuid:    fakeUuid,
			wantUsr: fakeUsr,
			wantErr: nil,
			prepareMocks: func(tr *mocks.TokensRepository) {
				tr.On("GetUserByToken", context.Background(), fakeUuid).Return(fakeUsr, nil)
			},
		},
		{
			uuid:    fakeUuid,
			wantUsr: nil,
			wantErr: fakeErr,
			prepareMocks: func(tr *mocks.TokensRepository) {
				tr.On("GetUserByToken", context.Background(), fakeUuid).Return(nil, fakeErr)
			},
		},
	}

	for _, testCond := range testConditions {
		tName := fmt.Sprintf("GetUserByToken(%v, %v) == %v, %v", context.Background(), testCond.uuid, testCond.wantUsr, testCond.wantErr)
		t.Run(tName, func(t *testing.T) {
			ctx := context.Background()
			tr := new(mocks.TokensRepository)
			testCond.prepareMocks(tr)
			svc := NewTokenService(tr)

			gotUsr, gotErr := svc.GetUserByToken(ctx, testCond.uuid)

			assert.Equal(t, testCond.wantErr, gotErr, "GetUserByToken returned unexpected result: got error %v want %v", gotErr, testCond.wantErr)
			assert.Equal(t, testCond.wantUsr, gotUsr, "GetUserByToken returned unexpected result: got user %v want %v", gotUsr, testCond.wantUsr)
			tr.AssertExpectations(t)
		})
	}
}
