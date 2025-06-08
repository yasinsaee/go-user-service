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
	exp, err := strconv.Atoi(config.GetEnv("JWT_EXP", "24"))
	if err != nil {
		logger.Error("youre expire date jwt is not ok")
	}
	jwt.Init(jwt.JWTConfig{
		PrivateKey: loadKey("private.key"),
		PublicKey:  loadKey("public.key"),
		Exp:        exp,
	})
}
