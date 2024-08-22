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
	"time"
)

func (app *Application) CreateEvent(c echo.Context) error {
	req := struct {
		ID             *int                   `json:"id"`
		Title          *string                `json:"title"`
		Type           []*internal.EventType  `json:"eventType"`
		Description    map[string]interface{} `json:"description"`
		BriefDesc      *string                `json:"briefDesc"`
		Genres         []*string              `json:"genres"`
		Venues         []*internal.Venue      `json:"venues"`
		StartTime      *time.Time             `json:"startTime"`
		EndTime        *time.Time             `json:"endTime"`
		Price          *float64               `json:"price"`
		AgeRestriction *int                   `json:"ageRestriction"`
		Rating         *float64               `json:"rating"`
		CreatedAt      *time.Time             `json:"createdAt"`
		UpdatedAt      *time.Time             `json:"updatedAt"`
	}{}
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid JSON")
	}

	jsonBytes, err := json.Marshal(req.Description)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid JSON")
	}
	jsonString := string(jsonBytes)

	for i, eventType := range req.Type {
		if eventType.ID == nil {
			t, err := app.models.event.CreateEventType(eventType.Name)
			if err != nil {
				return c.JSON(http.StatusBadRequest, "invalid JSON")
			}
			req.Type[i] = t
		}
	}

	for i, venue := range req.Venues {
		if venue.ID == nil {

			id, err := app.models.venue.CreateVenue(venue)
			if err != nil {
				fmt.Println(err.Error())
				return c.JSON(http.StatusInternalServerError, "internal server error")
			}
			req.Venues[i].ID = id
		}
	}
	timestamp := time.Now()
	event := internal.Event{
		ID:             req.ID,
		Title:          req.Title,
		Type:           req.Type,
		Description:    &jsonString,
		BriefDesc:      req.BriefDesc,
		Genres:         req.Genres,
		Venues:         req.Venues,
		StartTime:      req.StartTime,
		EndTime:        req.EndTime,
		Price:          req.Price,
		AgeRestriction: req.AgeRestriction,
		Rating:         req.Rating,
		CreatedAt:      &timestamp,
		UpdatedAt:      &timestamp,
	}

	id, err := app.models.event.CreateEvent(&event)
	if err != nil {
		fmt.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, "internal server error")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"id": id})
}

func (app *Application) GetEventsByFilter(c echo.Context) error {
	eventType := c.Param("type")
	if eventType == "" {
		return c.JSON(http.StatusBadRequest, "bad type")
	}
	pageNumber, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		pageNumber = 1
	}

	events, err := app.models.event.GetEventsByType(&eventType, &pageNumber)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "internal server error")
	}
	return c.JSON(http.StatusOK, events)
}

func (app *Application) GetEventTypes(c echo.Context) error {

	events, err := app.models.event.GetEventTypes()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "internal server error")
	}
	return c.JSON(http.StatusOK, events)
}

func (app *Application) GetEventTypeByName(c echo.Context) error {
	name := c.Param("name")
	if name == "" {
		return c.JSON(http.StatusBadRequest, "invalid name")
	}
	events, err := app.models.event.GetEventType(name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "internal server error")
	}
	return c.JSON(http.StatusOK, events)
}

func (app *Application) GetEventPagination(c echo.Context) error {
	pageParam := c.QueryParam("page")
	if pageParam == "" {
		return c.JSON(http.StatusBadRequest, "invalid page")
	}
	page, err := strconv.Atoi(pageParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid page")
	}
	events, totalPages, err := app.models.event.GetEventsPage(&page)
	if err != nil {
		fmt.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, "internal server error")
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"events": events, "totalPages": totalPages})
}

func (app *Application) GetEventImages(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid id")
	}
	images, err := app.models.event.GetImages(&id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "internal server error")
	}
	return c.JSON(http.StatusOK, images)
}

func (app *Application) GetGenres(c echo.Context) error {
	genres, err := app.models.event.GetAllGenres()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "internal server error")
	}
	return c.JSON(http.StatusOK, genres)
}

func (app *Application) GetEventById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid id")
	}
	event, err := app.models.event.GetEventById(&id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "internal server error")
	}
	return c.JSON(http.StatusOK, event)
}

func (app *Application) UploadImages(c echo.Context) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "internal server error")
	}
	param := form.Value["eventId"]
	eventId, err := strconv.Atoi(param[0])
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid id")
	}
	mainImagesNames := make([]*string, 0)
	postersNames := make([]*string, 0)

	mainImages := form.File["main_images"]
	for _, file := range mainImages {
		temp, err := uuid.NewV7()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "internal server error")
		}
		file.Filename = temp.String() + filepath.Ext(file.Filename)
		mainImagesNames = append(mainImagesNames, &file.Filename)
	}

	posters := form.File["posters"]
	for _, file := range posters {
		temp, err := uuid.NewV7()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "internal server error")
		}
		file.Filename = temp.String() + filepath.Ext(file.Filename)
		postersNames = append(postersNames, &file.Filename)
	}
	if len(mainImagesNames) == 0 && len(postersNames) == 0 {
		return c.JSON(http.StatusBadRequest, "invalid upload images")
	}
	err = app.models.event.CreateEventImage(eventId, mainImagesNames, postersNames)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "internal server error")
	}
	folder := "./static/" + strconv.Itoa(eventId)
	err = os.Mkdir(folder, 0705)
	if err != nil && !os.IsExist(err) {
		return c.JSON(http.StatusInternalServerError, "internal server error")
	}

	for _, file := range mainImages {
		// Source
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		// Destination
		dst, err := os.Create(folder + "/" + file.Filename)
		if err != nil {
			return err
		}
		defer dst.Close()

		// Copy
		if _, err = io.Copy(dst, src); err != nil {
			return err
		}

	}

	for _, file := range posters {
		// Source
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		// Destination
		dst, err := os.Create(folder + "/" + file.Filename)
		if err != nil {
			return err
		}
		defer dst.Close()

		// Copy
		if _, err = io.Copy(dst, src); err != nil {
			return err
		}

	}
	return c.JSON(http.StatusOK, "success")
}

func (app *Application) UploadTicketsNoShah(c echo.Context) error {
	req := struct {
		EventId *int `json:"eventId"`
		VenueId *int `json:"venueId"`
		Days    []*internal.DateWithTicketsNoShah
	}{}

	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid JSON")
	}
	err = app.models.tickets.CreateTicketsNoSham(req.EventId, req.VenueId, req.Days)
	if err != nil {
		fmt.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, "internal server error")
	}
	return c.JSON(http.StatusOK, "success")
}

func (app *Application) UploadTicketsWithShah(c echo.Context) error {
	req := struct {
		VenueId *int               `json:"venueId"`
		EventId *int               `json:"eventId"`
		Seats   [][]*internal.Seat `json:"seats"`
	}{}

	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid JSON")
	}

	err = app.models.tickets.CreateTicketsWithSham(req.EventId, req.VenueId, req.Seats)
	if err != nil {
		fmt.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, "internal server error")
	}
	return c.JSON(http.StatusOK, req)
}

func (app *Application) UploadDecorWithShah(c echo.Context) error {
	req := struct {
		VenueId *int `form:"venueId"`
		Items   []*internal.Decor
	}{}
	form, err := c.MultipartForm()
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid form")
	}
	tmp, err := strconv.Atoi(form.Value["venueId"][0])
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid id")
	}
	req.VenueId = &tmp

	items := form.Value["items"]
	for _, item := range items {
		temp := struct {
			Name     string `form:"name"`
			Left     int    `form:"left"`
			Top      int    `form:"top"`
			Uuid     string `form:"uuidTemp"`
			Image    string `form:"image"`
			Filename string `form:"filename"`
		}{}
		err = json.Unmarshal([]byte(item), &temp)
		if err != nil {
			return c.JSON(http.StatusBadRequest, "invalid JSON")
		}
		dec := internal.Decor{
			Id:    nil,
			Name:  &temp.Name,
			Image: &temp.Filename,
			Uuid:  &temp.Uuid,
			Left:  &temp.Left,
			Top:   &temp.Top,
		}
		req.Items = append(req.Items, &dec)
	}

	temp, err := app.models.event.CreateDecors(req.VenueId, req.Items)
	if err != nil {
		fmt.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, "internal server error")
	}

	req.Items = temp
	files := form.File["images"]
	for _, item := range req.Items {
		for _, file := range files {
			fmt.Println(file.Filename)
			fmt.Println(*item.Image)
			if item.Image != nil && file.Filename == *item.Image {
				src, err := file.Open()
				if err != nil {
					return err
				}
				defer src.Close()
				err = os.MkdirAll("static/decors/"+strconv.Itoa(*item.Id), os.ModePerm)
				if err != nil {
					return err
				}
				tempUUID, err := uuid.NewV7()
				if err != nil {
					return c.JSON(http.StatusInternalServerError, "internal server error")
				}
				file.Filename = tempUUID.String() + filepath.Ext(file.Filename)
				dst, err := os.Create("./static/decors/" + strconv.Itoa(*item.Id) + "/" + file.Filename)
				if err != nil {
					return err
				}
				defer dst.Close()

				// Copy
				if _, err = io.Copy(dst, src); err != nil {
					return err
				}

				err = app.models.event.UpdateImageDecor(&file.Filename, item.Id)
				if err != nil {
					return err
				}
			}
		}
	}

	return c.JSON(http.StatusOK, req)
}
