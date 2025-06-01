package app

import (
	"strconv"

	"github.com/yasinsaee/go-user-service/internal/app/config"
	"github.com/yasinsaee/go-user-service/pkg/jwt"
	"github.com/yasinsaee/go-user-service/pkg/logger"
)

func InitJWT() {
	exp, err := strconv.Atoi(config.GetEnv("JWT_EXP", "24"))
	if err != nil {
		logger.Error("youre expire date jwt is not ok")
	}
	jwt.Init(jwt.JWTConfig{
		Secret: config.GetEnv("JWT_SECRET", "defaultsecret"),
		Exp:    exp,
	})
}
