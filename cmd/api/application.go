package main

import (
	"github.com/labstack/echo/v4"
	"tap2go/internal"
)

type Models struct {
	user internal.UserModel
}

type Config struct {
	dsn  *string
	port *string
}

type Application struct {
	server *echo.Echo
	models Models
	config Config
}

func NewApplication(dsn, port *string) (*Application, error) {
	app := Application{}
	app.server = echo.New()
	app.AddMiddleware()
	app.AddRoutes()
	app.config.port = port
	app.config.dsn = dsn
	pool, err := ConnectPgPoolConfigured(app.config.dsn)
	if err != nil {
		return nil, err
	}
	app.models.user = internal.UserModel{DB: pool}
	return &app, nil
}
