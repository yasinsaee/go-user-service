package user_token_store

import (
	"time"

	"github.com/yasinsaee/go-user-service/pkg/redis"
)

type refreshTokenStoreImpl struct {
	refreshTokenExpire int64
}

// NewRefreshTokenStore returns a new instance of RefreshTokenStore
func NewRefreshTokenStore(refreshTokenExpire int64) *refreshTokenStoreImpl {
	return &refreshTokenStoreImpl{
		refreshTokenExpire: refreshTokenExpire,
	}
}

func (s *refreshTokenStoreImpl) Set(userID, refreshToken string) error {
	key := "refresh_token:" + userID + ":" + refreshToken
	ttl := time.Duration(s.refreshTokenExpire) * 24 * time.Hour
	return redis.Set(key, refreshToken, ttl)
}

func (s *refreshTokenStoreImpl) Exists(userID, refreshToken string) (bool, error) {
	key := "refresh_token:" + userID + ":" + refreshToken
	stored, err := redis.Get(key)
	if err == redis.ErrKeyNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return stored == refreshToken, nil
}

func (s *refreshTokenStoreImpl) Delete(userID, refreshToken string) error {
	key := "refresh_token:" + userID + ":" + refreshToken
	return redis.Remove(key)
}
