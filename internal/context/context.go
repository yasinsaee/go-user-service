package context

import "github.com/labstack/echo/v4"

type GlobalContext struct {
	echo.Context
}

func InitContext(h echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		g := &GlobalContext{Context: c}
		return h(g)
	}
}

