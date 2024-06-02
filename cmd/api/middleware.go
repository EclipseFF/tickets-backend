package main

import (
	"github.com/labstack/echo/v4/middleware"
)

func (app *Application) AddMiddleware() {
	DefaultLoggerConfig := middleware.LoggerConfig{
		Skipper: middleware.DefaultSkipper,
		Format: `{"time":"${time_rfc3339_nano}",` +
			`"method":"${method}","uri":"${uri}",` +
			`"status":${status},"error":"${error}","latency_human":"${latency_human}"` +
			`,"bytes_in":${bytes_in},"bytes_out":${bytes_out}}` + "\n",
		CustomTimeFormat: "2006-01-02 15:04:05.00000",
	}
	app.server.Use(middleware.LoggerWithConfig(DefaultLoggerConfig))
	app.server.Use(middleware.Recover())
	app.server.Use(middleware.CORS())
	//app.server.Use(middleware.CSRF())
	app.server.Use(middleware.Secure())

}
