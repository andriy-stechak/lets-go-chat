package repositories

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/andriystech/lgc/models"
)

var ErrTokenNotFound = errors.New("token not found")
var ErrTokenExpired = errors.New("token has been expired")

type TokensRepository interface {
	SaveToken(context.Context, string, *models.User) error
	GetUserByToken(context.Context, string) (*models.User, error)
}

type inMemoryRecord struct {
	user      *models.User
	expiresAt int
}

type tokensStorage struct {
	db  map[string]*inMemoryRecord
	ttl int
	mu  *sync.Mutex
}

func NewTokensRepository(ttl int) TokensRepository {
	return &tokensStorage{
		db:  map[string]*inMemoryRecord{},
		ttl: ttl,
		mu:  &sync.Mutex{},
	}
}

func (r *tokensStorage) SaveToken(ctx context.Context, token string, user *models.User) error {
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

func (r *tokensStorage) GetUserByToken(ctx context.Context, token string) (*models.User, error) {
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
