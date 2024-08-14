package main

import (
	"github.com/labstack/echo/v4"
	"tap2go/internal"
)

type Models struct {
	user    *internal.UserRepo
	sector  *internal.SectorRepo
	seat    *internal.SeatRepo
	venue   *internal.VenueRepo
	event   *internal.EventRepo
	tickets *internal.TicketRepo
	admin   *internal.AdminRepo
	news    *internal.NewsRepo
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

func NewApp(dsn, port *string) (*Application, error) {
	app := Application{}
	app.server = echo.New()
	app.server.HideBanner = true
	app.AddMiddleware()
	app.AddRoutes()
	app.config.port = port
	app.config.dsn = dsn
	pool, err := ConnectPgPoolConfigured(app.config.dsn)
	if err != nil {
		return nil, err
	}
	app.models.user = &internal.UserRepo{DB: pool}
	app.models.sector = &internal.SectorRepo{DB: pool}
	app.models.seat = &internal.SeatRepo{DB: pool}
	app.models.venue = &internal.VenueRepo{DB: pool}
	app.models.event = &internal.EventRepo{DB: pool}
	app.models.tickets = &internal.TicketRepo{DB: pool}
	app.models.admin = &internal.AdminRepo{DB: pool}
	app.models.news = &internal.NewsRepo{DB: pool}
	return &app, nil
}
