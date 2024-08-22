package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"tap2go/internal"
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
	/*for i, s := range sectors {
		_, err := app.models.seat.GetSeatsBySectorID(*s.ID)
		if err != nil {
			app.server.Logger.Error(err)
			return c.JSON(http.StatusInternalServerError, "internal server error")
		}
	}*/
	return c.JSON(http.StatusOK, sectors)
}

func (app *Application) CreateSector(c echo.Context) error {
	req := struct {
		EventId *int `form:"eventId"`
		VenueId *int `form:"venueId"`
		Sectors []*internal.Sector
	}{}

	form, err := c.MultipartForm()
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid form")
	}
	tmp, err := strconv.Atoi(form.Value["eventId"][0])
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid id")
	}
	req.EventId = &tmp
	tmp, err = strconv.Atoi(form.Value["venueId"][0])
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid id")
	}
	req.VenueId = &tmp

	sectors := form.Value["sectors"]
	for _, sector := range sectors {
		temp := struct {
			Name   string `form:"name"`
			Height int    `form:"height"`
			Width  int    `form:"width"`
			IsLink bool   `form:"isLink"`
			Left   int    `form:"left"`
			Top    int    `form:"top"`
			Uuid   string `form:"uuidTemp"`
			Image  string `form:"image"`
		}{}
		err = json.Unmarshal([]byte(sector), &temp)
		if err != nil {
			return c.JSON(http.StatusBadRequest, "invalid JSON")
		}
		sector := internal.Sector{
			VenueID: req.VenueId,
			Name:    &temp.Name,
			Height:  &temp.Height,
			Width:   &temp.Width,
			IsLink:  &temp.IsLink,
			Left:    &temp.Left,
			Top:     &temp.Top,
			Image:   &temp.Image,
		}
		req.Sectors = append(req.Sectors, &sector)
	}

	s, err := app.models.sector.CreateSectors(req.VenueId, req.Sectors)
	if err != nil {
		fmt.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, "internal server error")
	}
	req.Sectors = s
	files := form.File["images"]
	for _, sector := range req.Sectors {
		for _, file := range files {
			if sector.Image != nil && file.Filename == *sector.Image {
				src, err := file.Open()
				if err != nil {
					return err
				}
				defer src.Close()
				err = os.MkdirAll("static/sectors/"+strconv.Itoa(*sector.ID), os.ModePerm)
				if err != nil {
					return err
				}
				tempUUID, err := uuid.NewV7()
				if err != nil {
					return c.JSON(http.StatusInternalServerError, "internal server error")
				}
				file.Filename = tempUUID.String() + filepath.Ext(file.Filename)
				dst, err := os.Create("./static/sectors/" + strconv.Itoa(*sector.ID) + "/" + file.Filename)
				if err != nil {
					return err
				}
				defer dst.Close()

				// Copy
				if _, err = io.Copy(dst, src); err != nil {
					return err
				}

				err = app.models.sector.UpdateImage(&file.Filename, sector.ID)
				if err != nil {
					return err
				}
			}
		}
	}
	return c.JSON(http.StatusOK, req)
}
