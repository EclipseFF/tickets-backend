package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"tap2go/internal"
)

func (app *Application) CreateAdmin(c echo.Context) error {
	req := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid JSON")
	}
	admin := internal.Admin{
		Email: &req.Email,
		Password: internal.Password{
			Plaintext: req.Password,
			Hash:      "",
		},
	}
	err = admin.Password.SetPassword()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "invalid password")
	}
	session, err := app.models.admin.CreateAdmin(&admin)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "internal server error")
	}

	return c.JSON(http.StatusOK, session)
}

func (app *Application) LoginAdmin(c echo.Context) error {
	req := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid request")
	}
	user, err := app.models.admin.GetAdminByEmail(&req.Email)
	if err != nil {
		switch err.Error() {
		case "no rows in result set":
			return c.JSON(http.StatusNotFound, "user not found")
		default:
			return c.JSON(http.StatusInternalServerError, "internal server error")
		}
	}
	matches, err := user.Password.Matches(req.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "internal server error")
	}
	if !matches {
		return c.JSON(http.StatusUnauthorized, "unauthorized")
	}
	token, err := app.models.admin.CreateSession(user.Id)
	if err != nil {
		fmt.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, "internal server error")
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"token": token, "user": user})
}

func (app *Application) AdminLogout(c echo.Context) error {
	req := struct {
		Token string `json:"token"`
	}{}
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid request")
	}
	err = app.models.admin.DeleteSession(&req.Token)
	if err != nil {
		fmt.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, "internal server error")
	}
	return c.JSON(http.StatusOK, "ok")
}

func (app *Application) GetAdmin(c echo.Context) error {
	token := c.Param("token")
	if len(token) == 0 {
		return c.JSON(http.StatusBadRequest, "invalid token")
	}
	user, err := app.models.admin.GetAdminBySession(&token)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "internal server error")
	}
	return c.JSON(http.StatusOK, user)
}

func (app *Application) CreateEventType(c echo.Context) error {
	token := c.QueryParam("token")
	if len(token) == 0 {
		return c.JSON(http.StatusBadRequest, "invalid token")
	}
	_, err := app.models.admin.EnsureSession(&token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "not authorized")
	}
	req := struct {
		Name string `json:"name"`
	}{}
	err = c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "bad request")
	}

	eventType, err := app.models.event.CreateEventType(&req.Name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "internal server error")
	}
	return c.JSON(http.StatusOK, eventType)
}
