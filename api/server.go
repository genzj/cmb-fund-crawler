package api

import (
	"expvar"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func StartAPIServer() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", hello)
	e.GET("/_stats", echo.WrapHandler(expvar.Handler()))

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
