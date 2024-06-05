package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"tap2go/internal"
	"time"
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

func (app *Application) UpdateAdditionalUserData(c echo.Context) error {
	req := struct {
		UserId     *int    `json:"user_id"`
		Email      *string `json:"email"`
		Password   *string `json:"password"`
		Phone      *string `json:"phone"`
		Surname    *string `json:"surname"`
		Name       *string `json:"name"`
		Patronymic *string `json:"patronymic"`
		Birthdate  *string `json:"birthdate"`
	}{}
	err := c.Bind(&req)
	if err != nil {

		return c.JSON(http.StatusBadRequest, "invalid request")
	}

	oldUser, err := app.models.user.GetUserById(*req.UserId)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println(21)
		return c.JSON(http.StatusInternalServerError, "internal server error")
	}
	createNewRecod := false
	oldUserData, err := app.models.user.GetUserAdditional(*req.UserId)
	if err != nil {
		switch err.Error() {
		case "no rows in result set":
			createNewRecod = true
		default:
			return c.JSON(http.StatusInternalServerError, "internal server error")
		}

	}
	u := new(internal.User)
	d := new(internal.AdditionalUserData)
	u.Id = req.UserId
	d.UserId = *req.UserId

	if *req.Email == "" {
		u.Email = oldUser.Email
	} else {
		u.Email = req.Email
	}

	if *req.Password == "" {
		u.Password.Hash = oldUser.Password.Hash
	} else {
		u.Password.Plaintext = *req.Password
		err = u.Password.SetPassword()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "internal server error")
		}
	}

	if *req.Phone == "" {
		u.Phone = oldUser.Phone
	} else {
		u.Phone = req.Phone
	}

	if *req.Surname == "" {
		d.Surname = oldUserData.Surname
	} else {
		d.Surname = req.Surname
	}

	if *req.Name == "" {
		d.Name = oldUserData.Name
	} else {
		d.Name = req.Name
	}

	if *req.Patronymic == "" {
		d.Patronymic = oldUserData.Patronymic
	} else {
		d.Patronymic = req.Patronymic
	}

	if *req.Birthdate == "" {
		d.DateOfBirth = oldUserData.DateOfBirth
	} else {
		bdate, err := time.Parse(time.RFC3339, *req.Birthdate)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "internal server error")
		}
		d.DateOfBirth = &bdate
	}

	err = app.models.user.UpdateUser(u, d, createNewRecod)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "internal server error")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"user": u, "additional": d})
}
