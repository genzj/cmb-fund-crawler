package api

import (
	"expvar"
	"github.com/dgrijalva/jwt-go"
	"github.com/genzj/cmb-fund-crawler/db"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

// TODO read signing key from argument or environment variable
const JwtSigningKey = "insecure"

func StartAPIServer() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/login", login)

	api := e.Group("/api")
	api.Use(middleware.JWT([]byte(JwtSigningKey)))
	registerFundRoutes(api)
	// Routes
	e.GET("/_stats", echo.WrapHandler(expvar.Handler()))

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

func login(c echo.Context) error {
	info := &struct {
		Username string
		Password string
	}{}

	if err := c.Bind(info); err != nil {
		return echo.ErrBadRequest
	}

	dbi := db.GetDatabaseInstance()
	if user, err := db.GetUser(dbi, info.Username); err != nil || user == nil {
		return echo.ErrBadRequest
	} else if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(info.Password)) == nil {
		// Create token
		token := jwt.New(jwt.SigningMethodHS256)

		// Set claims
		claims := token.Claims.(jwt.MapClaims)
		claims["name"] = user.Name
		// TODO shorten valid time duration and add refresh mechanism
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte(JwtSigningKey))
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, map[string]string{
			"token": t,
		})
	}
	return echo.ErrUnauthorized
}
