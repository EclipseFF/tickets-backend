package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func (app *Application) GetUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid id")
	}
	user, err := app.models.user.GetUser(id)
	if err != nil {
		switch err.Error() {
		case "no rows in result set":
			return c.JSON(http.StatusNotFound, "user not found")
		}
		return c.JSON(http.StatusInternalServerError, "internal server error")
	}

	return c.JSON(http.StatusOK, user)
}
