package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
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

func (app *Application) GetVenuesByEvent(c echo.Context) error {
	param := c.Param("id")
	if param == "" {
		return c.JSON(http.StatusBadRequest, "bad id")
	}

	id, err := strconv.Atoi(param)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "bad id")
	}
	venues, err := app.models.venue.GetVenuesByEvent(&id)
	return c.JSON(http.StatusOK, venues)
}

func (app *Application) GetVenueById(c echo.Context) error {
	param := c.Param("id")
	if param == "" {
		return c.JSON(http.StatusBadRequest, "bad id")
	}
	id, err := strconv.Atoi(param)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "bad id")
	}
	venue, err := app.models.venue.GetVenueById(&id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "internal server error")
	}
	return c.JSON(http.StatusOK, venue)

}
