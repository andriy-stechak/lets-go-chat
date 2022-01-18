package repositories

import (
	"context"
	"fmt"
	"testing"

	"github.com/andriystech/lgc/config"
	"github.com/andriystech/lgc/models"
	"github.com/stretchr/testify/assert"
)

type testGetUserByTokenData struct {
	token       string
	wantErr     error
	wantUsr     *models.User
	composeRepo func() TokensRepository
}

func TestSaveToken(t *testing.T) {

	ctx := context.Background()
	repo := NewTokensRepository(&config.ServerConfig{TokenTTLInSeconds: 10})

	gotErr := repo.SaveToken(ctx, "token", &models.User{})

	assert.Nil(t, gotErr, "SaveToken returned unexpected result: got %v want %v", gotErr, nil)
}

func TestGetUserByToken(t *testing.T) {
	fakeToken := "51cdef8e-283b-4333-ade9-74ab2c7ca8fa"
	fakeUsr := &models.User{Id: "someid"}
	testConditions := []testGetUserByTokenData{
		{
			token:   fakeToken,
			wantErr: nil,
			wantUsr: fakeUsr,
			composeRepo: func() TokensRepository {
				tr := NewTokensRepository(&config.ServerConfig{TokenTTLInSeconds: 10})
				tr.SaveToken(context.Background(), fakeToken, fakeUsr)
				return tr
			},
		},
		{
			token:   fakeToken,
			wantErr: ErrTokenNotFound,
			wantUsr: nil,
			composeRepo: func() TokensRepository {
				return NewTokensRepository(&config.ServerConfig{TokenTTLInSeconds: 10})
			},
		},
		{
			token:   fakeToken,
			wantErr: ErrTokenExpired,
			wantUsr: nil,
			composeRepo: func() TokensRepository {
				tr := NewTokensRepository(&config.ServerConfig{TokenTTLInSeconds: -1})
				tr.SaveToken(context.Background(), fakeToken, fakeUsr)
				return tr
			},
		},
	}
	for _, testCond := range testConditions {
		tName := fmt.Sprintf("GetUserByToken(%v, %v) == %v, %v", context.Background(), testCond.token, testCond.wantUsr, testCond.wantErr)
		t.Run(tName, func(t *testing.T) {
			ctx := context.Background()
			repo := testCond.composeRepo()

			gotUsr, gotErr := repo.GetUserByToken(ctx, testCond.token)

			assert.Equal(t, testCond.wantErr, gotErr, "GetUserByToken returned unexpected result: got error %v want %v", gotErr, testCond.wantErr)
			assert.Equal(t, testCond.wantUsr, gotUsr, "Invalid user returned: got %v want %v", gotUsr, testCond.wantUsr)
		})
	}
}
