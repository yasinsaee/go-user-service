package app

import (
	"os"
	"strconv"

	"github.com/yasinsaee/go-user-service/internal/app/config"
	"github.com/yasinsaee/go-user-service/pkg/jwt"
	"github.com/yasinsaee/go-user-service/pkg/logger"
)

func loadKey(path string) []byte {

	data, err := os.ReadFile(path)
	if err != nil {
		panic("cannot read key file: " + err.Error())
	}
	return data
}

func InitJWT() {
	accessExp, err := strconv.Atoi(config.GetEnv("JWT_ACCESS_TOKEN_EXP", "1"))
	if err != nil {
		logger.Error("youre expire date jwt is not ok")
	}

	refreshExp, err := strconv.Atoi(config.GetEnv("JWT_REFRESH_TOKEN_EXP", "30"))
	if err != nil {
		logger.Error("youre expire date jwt is not ok")
	}
	jwt.Init(jwt.JWTConfig{
		PrivateKey:      loadKey(config.GetEnv("PRIVATE_KEY_PATH", "../keys/private.key")),
		PublicKey:       loadKey(config.GetEnv("PUBLIC_KEY_PATH", "../keys/public.key")),
		AccessTokenExp:  accessExp,
		RefreshTokenExp: refreshExp,
	})
}
