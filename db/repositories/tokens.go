package repositories

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/andriystech/lgc/config"
	"github.com/andriystech/lgc/models"
)

var ErrTokenNotFound = errors.New("token not found")
var ErrTokenExpired = errors.New("token has been expired")

type inMemoryRecord struct {
	user      *models.User
	expiresAt int
}

type TokensRepository struct {
	db  map[string]*inMemoryRecord
	ttl int
	mu  *sync.Mutex
}

func NewTokensRepository() *TokensRepository {
	serviceConfig := config.GetServerConfig()
	return &TokensRepository{
		db:  map[string]*inMemoryRecord{},
		ttl: serviceConfig.TokenTTLInSeconds,
		mu:  &sync.Mutex{},
	}
}

func (r *TokensRepository) SaveToken(ctx context.Context, token string, user *models.User) error {
	// TODO: Use real storage for tokens
	r.mu.Lock()
	defer r.mu.Unlock()
	now := int(time.Now().Unix())
	r.db[token] = &inMemoryRecord{
		user:      user,
		expiresAt: int(now + r.ttl),
	}
	return nil
}

func (r *TokensRepository) GetUserByToken(ctx context.Context, token string) (*models.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	record := r.db[token]
	if record == nil {
		return nil, ErrTokenNotFound
	}
	now := int(time.Now().Unix())
	if record.expiresAt <= now {
		delete(r.db, token)
		return nil, ErrTokenExpired
	}
	delete(r.db, token)
	return record.user, nil
}
