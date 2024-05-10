package main

import (
	"github.com/labstack/echo/v4"
)

func (app *Application) AddRoutes() {
	app.server.GET("/healthcheck", func(c echo.Context) error {
		return c.JSON(200, "OK")
	})

	api := app.server.Group("/api")
	version := api.Group("/v1")
	usersRoutes := version.Group("/user")
	usersRoutes.GET("/:id", app.GetUser)
}
