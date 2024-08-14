package main

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"tap2go/internal"
)

func (app *Application) GetAllNews(c echo.Context) error {
	news, err := app.models.news.GetAllNews()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "internal server error")
	}
	return c.JSON(http.StatusOK, news)
}

func (app *Application) GetNewsById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid id")
	}
	news, err := app.models.news.GetNewsById(&id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "internal server error")
	}
	return c.JSON(http.StatusOK, news)
}

func (app *Application) GetNewsPagination(c echo.Context) error {
	pageParam := c.QueryParam("page")
	if pageParam == "" {
		return c.JSON(http.StatusBadRequest, "invalid page")
	}
	page, err := strconv.Atoi(pageParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid page")
	}

	news, totalPages, err := app.models.news.GetPaginatedNews(10, (page-1)*10)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "internal server error")
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"news": news, "totalPages": totalPages})
}

func (app *Application) GetLatestNews(c echo.Context) error {

	news, err := app.models.news.GetLatestNews(1)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "internal server error")
	}
	return c.JSON(http.StatusOK, news)
}

func (app *Application) CreateNews(c echo.Context) error {

	req := internal.News{}

	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid JSON")
	}
	n, err := app.models.news.CreateNews(&req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "internal server error")
	}
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	files := form.File["images"]

	if len(files) > 0 {
		err = os.MkdirAll("static/news/"+strconv.Itoa(*n.Id), os.ModePerm)
		if err != nil {
			return err
		}
		filenames := make([]*string, 0)
		for _, file := range files {

			src, err := file.Open()
			if err != nil {
				return err
			}
			defer src.Close()
			temp := uuid.New().String() + filepath.Ext(file.Filename)
			dst, err := os.Create("static/news/" + strconv.Itoa(*n.Id) + "/" + temp)
			if err != nil {
				return err
			}
			defer dst.Close()

			if _, err = io.Copy(dst, src); err != nil {
				return err
			}

			filenames = append(filenames, &temp)

		}
		err = app.models.news.SetNewsImages(filenames, n.Id)
		if err != nil {
			go app.models.news.DeleteNews(n.Id)
			return c.JSON(http.StatusInternalServerError, "internal server error")
		}
	}

	return c.JSON(http.StatusOK, n)
}
