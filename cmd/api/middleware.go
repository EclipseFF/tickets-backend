package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

func (app *Application) AddMiddleware() {
	log := logrus.New()
	app.server.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
			log.WithFields(logrus.Fields{
				"URI":            values.URI,
				"status":         values.Status,
				"method":         c.Request().Method,
				"content-length": c.Request().ContentLength,
				"response-size":  values.ResponseSize,
				"error":          values.Error,
			}).Info("request")

			return nil
		},
	}))
	app.server.Use(middleware.Recover())
	app.server.Use(middleware.CORS())
	app.server.Use(middleware.CSRF())
	app.server.Use(middleware.Secure())

}
