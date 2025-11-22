package ratelimiter

import (
	"errors"
	"fmt"
	"time"

	"github.com/yasinsaee/go-user-service/pkg/redis"
)

type RedisOTPRateLimiter struct {
	limit int
}

func NewRedisOTPRateLimiter(limit int) *RedisOTPRateLimiter {
	return &RedisOTPRateLimiter{
		limit: limit,
	}
}

func (r *RedisOTPRateLimiter) CanSend(receiver string) (bool, error) {
	key := "otp:" + receiver

	val, err := redis.Get(key)
	if err != nil {
		if errors.Is(err, redis.ErrKeyNotFound) {
			return true, nil
		}
		return false, err
	}

	count, err := redisClientStrToInt(val)
	if err != nil {
		return false, err
	}

	return count < r.limit, nil
}

func (r *RedisOTPRateLimiter) MarkSend(receiver string, ttl time.Duration) error {
	key := "otp:" + receiver

	if _, err := redis.Incr(key); err != nil {
		return err
	}

	if err := redis.Expire(key, ttl); err != nil {
		return err
	}

	return nil
}

func redisClientStrToInt(val string) (int, error) {
	var count int
	_, err := fmt.Sscanf(val, "%d", &count)
	return count, err
}
