package app

import (
	"strconv"
	"time"

	"github.com/yasinsaee/go-user-service/internal/app/config"
	"github.com/yasinsaee/go-user-service/pkg/redis"
)

func InitRedis() {
	enable, err := strconv.ParseBool(config.GetEnv("REDIS_ENABLE", ""))
	if enable && err == nil {
		db, _ := strconv.Atoi(config.GetEnv("REDIS_DB", ""))
		ttl, _ := strconv.Atoi(config.GetEnv("REDIS_TTL", ""))
		redis.Init(redis.Config{
			Addr:     config.GetEnv("REDIS_ADDR", ""),
			Password: config.GetEnv("REDIS_PASSWORD", ""),
			DB:       db,
			TTL:      time.Duration(ttl) * time.Second,
		})
	}

}
