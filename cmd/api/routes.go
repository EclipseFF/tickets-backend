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

	version.Static("/static", "static")

	adminRoutes := version.Group("/admin")
	adminRoutes.POST("/user-admin", app.CreateAdmin)
	adminRoutes.POST("/login", app.LoginAdmin)
	adminRoutes.DELETE("/logout", app.AdminLogout)
	adminRoutes.GET("/user-admin/:token", app.GetAdmin)
	adminRoutes.POST("/type/create", app.CreateEventType)

	usersRoutes := version.Group("/user")
	usersRoutes.POST("/register", app.CreateUser)
	usersRoutes.POST("/login", app.AuthenticateUser)
	usersRoutes.DELETE("/logout", app.Logout)
	usersRoutes.GET("/:id", app.GetAdditionalUserData)
	usersRoutes.POST("", app.GetUserBySession)
	usersRoutes.GET("/additional/:id", app.GetAdditionalUserData)

	eventRoutes := version.Group("/event")
	eventRoutes.GET("/page/:page", app.GetEventPagination)
	eventRoutes.POST("/create", app.CreateEvent)
	eventRoutes.GET("/:id", app.GetEventById)
	eventRoutes.GET("/images/:id", app.GetEventImages)
	eventRoutes.GET("/genres", app.GetGenres)
	eventRoutes.GET("/type/:type", app.GetEventsByFilter)

	typeRoutes := version.Group("/type")
	typeRoutes.GET("/all", app.GetEventTypes)
	typeRoutes.GET("/:name", app.GetEventTypeByName)
	//typeRoutes.POST("", app.CreateEventType)

	venueRoutes := version.Group("/venue")
	venueRoutes.GET("/all", app.GetAllVenues)

	sectorRoutes := version.Group("/sector")
	sectorRoutes.GET("/:id", app.GetSectorsByVenueId)
}
