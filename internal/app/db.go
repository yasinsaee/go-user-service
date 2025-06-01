package app

import (
	"github.com/yasinsaee/go-user-service/internal/app/config"
	"github.com/yasinsaee/go-user-service/pkg/mongo"
)

func InitMongo() {
	mongo.Init(mongo.Config{
		URI:           config.GetEnv("MONGO_URI", "mongodb://localhost:27017"),
		DB:            config.GetEnv("MONGO_DB", "user_service"),
		Username:      config.GetEnv("MONGO_USERNAME", ""),
		Password:      config.GetEnv("MONGO_PASSWORD", ""),
		AuthMechanism: config.GetEnv("MONGO_AUTH_MECH", "SCRAM-SHA-1"),
		AuthSource:    config.GetEnv("MONGO_AUTH_SOURCE", ""),
	})
}
