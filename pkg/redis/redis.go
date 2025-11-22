package redis

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/yasinsaee/go-user-service/pkg/logger"
)

type Config struct {
	Addr     string        `json:"addr" yaml:"addr"`
	Password string        `json:"password" yaml:"password"`
	DB       int           `json:"db" yaml:"db"`
	TTL      time.Duration `json:"ttl" yaml:"ttl"`
}

var (
	DB            *RedisDB
	DefaultConfig = Config{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		TTL:      24 * time.Hour,
	}
	ErrKeyNotFound = errors.New("redis: key not found")
)

type RedisDB struct {
	Client *redis.Client
	TTL    time.Duration
}

// Init connects to Redis
func Init(cfg Config) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		logger.Error("Error connecting to Redis, " + err.Error())
		for i := 0; i < 3; i++ {
			_, err = client.Ping(ctx).Result()
			if err == nil {
				logger.Info("Connected to Redis")
				break
			} else {
				logger.Error("Error connecting to Redis, " + err.Error())
			}
		}
	} else {
		logger.Info("Connected to Redis")
	}

	DB = &RedisDB{
		Client: client,
		TTL:    cfg.TTL,
	}
}

// Set sets a key with optional TTL
func Set(key string, value interface{}, ttl ...time.Duration) error {
	ctx := context.Background()
	exp := DB.TTL
	if len(ttl) > 0 {
		exp = ttl[0]
	}
	return DB.Client.Set(ctx, key, value, exp).Err()
}

// Get returns value by key
func Get(key string) (string, error) {
	ctx := context.Background()
	val, err := DB.Client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", ErrKeyNotFound
	}
	return val, err
}

// Remove deletes a key
func Remove(key string) error {
	ctx := context.Background()
	return DB.Client.Del(ctx, key).Err()
}

// Exists checks if key exists
func Exists(key string) (bool, error) {
	ctx := context.Background()
	val, err := DB.Client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return val > 0, nil
}

// Increment a key
func Incr(key string) (int64, error) {
	ctx := context.Background()
	return DB.Client.Incr(ctx, key).Result()
}

// Decrement a key
func Decr(key string) (int64, error) {
	ctx := context.Background()
	return DB.Client.Decr(ctx, key).Result()
}

// TTL returns remaining time for a key
func TTL(key string) (time.Duration, error) {
	ctx := context.Background()
	return DB.Client.TTL(ctx, key).Result()
}
