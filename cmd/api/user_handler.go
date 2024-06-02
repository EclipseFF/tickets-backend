package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"tap2go/internal"
)

func (app *Application) GetUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid id")
	}
	user, err := app.models.user.GetUserById(id)
	if err != nil {
		switch err.Error() {
		case "no rows in result set":
			return c.JSON(http.StatusNotFound, "user not found")
		}
		return c.JSON(http.StatusInternalServerError, "internal server error")
	}

	return c.JSON(http.StatusOK, user)
}

func (app *Application) GetUserBySession(c echo.Context) error {
	req := struct {
		Token string `json:"token"`
	}{}
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid token")
	}
	user, err := app.models.user.GetUserBySession(&req.Token)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "internal server error")
	}
	return c.JSON(http.StatusOK, user)
}

func (app *Application) CreateUser(c echo.Context) error {
	req := struct {
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Password string `json:"password"`
	}{}
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid request")
	}
	var u internal.User
	u.Email = &req.Email
	u.Phone = &req.Phone
	u.Password.Plaintext = req.Password
	err = u.Password.SetPassword()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "internal server error")
	}
	session, err := app.models.user.CreateUser(&u)
	if err != nil {
		switch err.Error() {
		case "ERROR: duplicate key value violates unique constraint \"users_email_key\" (SQLSTATE 23505)":
			return c.JSON(http.StatusBadRequest, "email is used")
		case "ERROR: duplicate key value violates unique constraint \"users_phone_key\" (SQLSTATE 23505)":
			return c.JSON(http.StatusBadRequest, "phone is used")
		}
		return c.JSON(http.StatusInternalServerError, "internal server error")
	}
	return c.JSON(http.StatusOK, session)
}

func (app *Application) AuthenticateUser(c echo.Context) error {
	req := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid request")
	}
	user, err := app.models.user.GetUserByEmail(&req.Email)
	if err != nil {
		switch err.Error() {
		case "no rows in result set":
			return c.JSON(http.StatusNotFound, "user not found")
		default:
			fmt.Println(err.Error())
			return c.JSON(http.StatusInternalServerError, "internal server error")
		}
	}
	matches, err := user.Password.Matches(req.Password)
	if err != nil {
		fmt.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, "internal server error")
	}
	if !matches {
		return c.JSON(http.StatusUnauthorized, "unauthorized")
	}
	token, err := app.models.user.CreateSession(user.Id)
	if err != nil {
		fmt.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, "internal server error")
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"token": token, "user": user})
}

func (app *Application) Logout(c echo.Context) error {
	req := struct {
		Token string `json:"token"`
	}{}
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid request")
	}
	err = app.models.user.DeleteSession(&req.Token)
	if err != nil {
		fmt.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, "internal server error")
	}
	return c.JSON(http.StatusOK, "ok")
}

func (app *Application) GetAdditionalUserData(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid id")
	}
	user, err := app.models.user.GetUserAdditional(id)
	if err != nil {
		switch err.Error() {
		case "no rows in result set":
			return c.JSON(http.StatusNotFound, "user not found")
		}
		fmt.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, "internal server error")
	}

	return c.JSON(http.StatusOK, user)
}
