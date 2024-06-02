package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func (app *Application) GetSectorsByVenueId(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid id")
	}
	sectors, err := app.models.sector.GetSectorsByVenue(id)
	if err != nil {
		switch err.Error() {
		case "no rows in result set":
			return c.JSON(http.StatusNotFound, "user not found")
		}
		return c.JSON(http.StatusInternalServerError, "internal server error")
	}
	for i, s := range sectors {
		seats, err := app.models.seat.GetSeatsBySectorID(s.ID)
		if err != nil {
			app.server.Logger.Error(err)
			return c.JSON(http.StatusInternalServerError, "internal server error")
		}
		sectors[i].Layout = seats
	}
	return c.JSON(http.StatusOK, sectors)
}
