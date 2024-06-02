package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (app *Application) GetAllVenues(c echo.Context) error {
	venues, err := app.models.venue.GetAll()
	if err != nil {
		switch err.Error() {
		case "no rows in result set":
			return c.JSON(http.StatusNotFound, "no venues found")
		}
		return c.JSON(http.StatusInternalServerError, "internal server error")
	}

	return c.JSON(http.StatusOK, venues)
}
