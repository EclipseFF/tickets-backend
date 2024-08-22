package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (app *Application) BuyTicketNoShah(c echo.Context) error {
	req := struct {
		UserToken    string `json:"userToken"`
		TicketTypeId *int   `json:"ticketTypeId"`
		Count        *int   `json:"count"`
	}{}

	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid request")
	}

	if req.UserToken == "" || req.TicketTypeId == nil || req.Count == nil {
		return c.JSON(http.StatusBadRequest, "invalid request")
	}

	u, err := app.models.user.GetUserBySession(&req.UserToken)
	if err != nil {

		return c.JSON(http.StatusInternalServerError, "internal server error")
	}

	result, err := app.models.tickets.BuyTicketNoShah(req.TicketTypeId, u.Id, req.Count)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "internal server error")
	}
	return c.JSON(http.StatusOK, result)
}

func (app *Application) ReadDatesForEventVenue(c echo.Context) error {
	req := struct {
		EventID *int `json:"eventId"`
		VenueID *int `json:"venueId"`
	}{}

	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid JSON")
	}
	result, err := app.models.tickets.GetDatesForEventVenue(req.EventID, req.VenueID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "internal server error")
	}
	return c.JSON(http.StatusOK, result)
}

func (app *Application) ReadDatesForEventVenueShah(c echo.Context) error {
	req := struct {
		EventID *int `json:"eventId"`
		VenueID *int `json:"venueId"`
	}{}

	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid JSON")
	}
	result, err := app.models.tickets.GetDatesForEventVenueShah(req.EventID, req.VenueID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "internal server error")
	}
	return c.JSON(http.StatusOK, result)
}
