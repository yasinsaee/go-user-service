package app

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/yasinsaee/go-user-service/internal/app/config"
)

func StartApp() {
	config.LoadEnv()

	InitMongo()

	InitJWT()

	e := echo.New()

	InitRoutes(e)

	port := config.GetEnv("PORT", "8080")
	fmt.Println("✅ Server running on port:", port)
	e.Logger.Fatal(e.Start(":" + port))
}
