package internal

import (
	"context"
	"errors"
	"github.com/essentialkaos/translit/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type EventRepo struct {
	DB *pgxpool.Pool
}

type EventDescription struct {
	Type *string `json:"doc"`
}

type Event struct {
	ID             *int         `json:"id"`
	Title          *string      `json:"title"`
	Type           []*EventType `json:"eventType"`
	Description    *string      `json:"description"`
	BriefDesc      *string      `json:"brief_desc"`
	Genres         []*string    `json:"genres"`
	Venues         []*Venue     `json:"venues"`
	StartTime      *time.Time   `json:"startTime"`
	EndTime        *time.Time   `json:"endTime"`
	Price          *float64     `json:"price"`
	AgeRestriction *int         `json:"ageRestriction"`
	Rating         *float64     `json:"rating"`
	CreatedAt      *time.Time   `json:"createdAt"`
	UpdatedAt      *time.Time   `json:"updatedAt"`
}

type EventImages struct {
	EventId    int       `json:"event_id"`
	Posters    []*string `json:"posters"`
	MainImages []*string `json:"main_images"`
}

type EventType struct {
	ID             *int    `json:"id"`
	Name           *string `json:"name"`
	TranslatedName *string `json:"translatedName"`
}

func (m *EventRepo) CreateEvent(event *Event) (*int, error) {
	tx, err := m.DB.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())
	var id int
	row := tx.QueryRow(context.Background(), `INSERT INTO events VALUES(default, $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id`,
		event.Title, event.Description, event.BriefDesc, event.Genres, event.StartTime, event.EndTime, event.Price, event.AgeRestriction, event.Rating, event.CreatedAt, event.UpdatedAt)

	err = row.Scan(&id)
	if err != nil {
		return nil, err
	}

	for _, venue := range event.Venues {
		_, err := tx.Exec(context.Background(), `INSERT INTO event_venues VALUES($1, $2)`, id, venue.ID)
		if err != nil {
			return nil, err
		}
	}

	for _, t := range event.Type {
		_, err := tx.Exec(context.Background(), `INSERT INTO event_types VALUES($1, $2)`, id, t.ID)
		if err != nil {
			return nil, err
		}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}

	return &id, nil
}

func (m *EventRepo) GetEventsByType(tip *string, pageNumber *int) ([]*Event, error) {
	tx, err := m.DB.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())

	var typeId int
	err = tx.QueryRow(context.Background(), `SELECT id FROM types WHERE translated_name = $1`, *tip).Scan(&typeId)
	if err != nil {
		return nil, err
	}
	limit := 20 * *pageNumber
	offset := limit * (*pageNumber - 1)

	rows, err := tx.Query(context.Background(), `SELECT event_id FROM event_types WHERE type_id = $1 ORDER BY event_id desc LIMIT $2 OFFSET $3`, typeId, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := make([]*Event, 0)
	venueRepo := &VenueRepo{DB: m.DB}
	for rows.Next() {
		var eId int
		err := rows.Scan(&eId)
		if err != nil {
			return nil, err
		}
		e, err := m.GetEventById(&eId)
		if err != nil {
			return nil, err
		}
		venues, err := venueRepo.GetVenuesByEvent(&eId)
		if err != nil {
			return nil, err
		}
		e.Venues = venues
		events = append(events, e)
	}

	rows.Close()
	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (m *EventRepo) CreateEventType(name *string) (*EventType, error) {
	translatedName := translit.ICAO(*name)
	tx, err := m.DB.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())
	var id int
	err = tx.QueryRow(context.Background(), `INSERT INTO types(id, name, translated_name) VALUES(default, $1, $2) RETURNING id`, *name, translatedName).Scan(&id)
	if err != nil {
		return nil, err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}
	t := EventType{
		ID:             &id,
		Name:           name,
		TranslatedName: &translatedName,
	}

	return &t, nil
}

func (m *EventRepo) GetEventTypes() ([]*EventType, error) {
	tx, err := m.DB.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())

	rows, err := tx.Query(context.Background(), `SELECT id, name, translated_name FROM types order by id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var events []*EventType
	for rows.Next() {
		var e EventType
		err = rows.Scan(&e.ID, &e.Name, &e.TranslatedName)
		if err != nil {
			return nil, err
		}
		events = append(events, &e)
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (m *EventRepo) GetEventTypeById(id *int) (*EventType, error) {
	tx, err := m.DB.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())
	var e EventType
	err = tx.QueryRow(context.Background(), `SELECT name, translated_name FROM types where id = $1`, id).Scan(&e.Name, &e.TranslatedName)
	if err != nil {
		return nil, err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}
	e.ID = id
	return &e, nil
}

func (m *EventRepo) GetEventType(name string) (*EventType, error) {
	tx, err := m.DB.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())
	var e EventType
	err = tx.QueryRow(context.Background(), `SELECT id, name, translated_name FROM types WHERE translated_name = $1`, name).Scan(&e.ID, &e.Name, &e.TranslatedName)
	if err != nil {

		return nil, err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}

	return &e, nil
}

func (m *EventRepo) GetEventTypeByEvent(eventId *int) ([]*EventType, error) {
	tx, err := m.DB.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())

	rows, err := tx.Query(context.Background(), `SELECT type_id FROM event_types WHERE event_id = $1`, *eventId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	eTypes := make([]*EventType, 0)
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		e, err := m.GetEventTypeById(&id)
		if err != nil {
			return nil, err
		}
		eTypes = append(eTypes, e)
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}

	return eTypes, nil
}

func (m *EventRepo) GetEventsPage(page *int) ([]*Event, *int, error) {
	tx, err := m.DB.Begin(context.Background())
	if err != nil {
		return nil, nil, err
	}
	defer tx.Rollback(context.Background())
	var totalPages int
	err = tx.QueryRow(context.Background(), `SELECT count(*) FROM events`).Scan(&totalPages)
	if err != nil {
		return nil, nil, err
	}
	stmt := `SELECT id, title, description, brief_desc, genre, start_time, end_time, price, age_restriction, rating, created_at, updated_at FROM events order by start_time desc limit $1 OFFSET $2`
	rows, err := tx.Query(context.Background(), stmt, 10, *page*10)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	var events []*Event
	venueRepo := VenueRepo{DB: m.DB}
	for rows.Next() {
		var e Event
		err := rows.Scan(&e.ID, &e.Title, &e.Description, &e.BriefDesc, &e.Genres, &e.StartTime, &e.EndTime, &e.Price, &e.AgeRestriction, &e.Rating, &e.CreatedAt, &e.UpdatedAt)
		if err != nil {
			return nil, nil, err
		}
		venues, err := venueRepo.GetVenuesByEvent(e.ID)
		eventTypes, err := m.GetEventTypeByEvent(e.ID)
		if err != nil {
			return nil, nil, err
		}
		e.Venues = venues
		e.Type = eventTypes
		events = append(events, &e)
	}

	rows.Close()
	err = tx.Commit(context.Background())
	if err != nil {
		return nil, nil, err
	}
	return events, &totalPages, nil
}

func (m *EventRepo) GetImages(id *int) (*EventImages, error) {
	tx, err := m.DB.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())
	var images EventImages
	err = tx.QueryRow(context.Background(), `SELECT posters, main_images FROM event_images where event_id = $1`, *id).Scan(&images.Posters, &images.MainImages)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, nil
		}
		return nil, err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}
	images.EventId = *id
	return &images, nil
}

func (m *EventRepo) GetAllGenres() ([]*string, error) {
	tx, err := m.DB.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())
	genres := make([]*string, 0)
	rows, err := tx.Query(context.Background(), `SELECT DISTINCT UNNEST(genre) AS genre FROM events;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var genre string
		err = rows.Scan(&genre)
		if err != nil {
			return nil, err
		}
		genres = append(genres, &genre)
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}
	return genres, nil
}

func (m *EventRepo) GetEventById(id *int) (*Event, error) {
	tx, err := m.DB.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())
	var e Event
	err = tx.QueryRow(context.Background(), `SELECT * FROM events where id = $1`, *id).Scan(&e.ID, &e.Title, &e.Description, &e.BriefDesc, &e.Genres, &e.StartTime, &e.EndTime, &e.Price, &e.AgeRestriction, &e.Rating, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		return nil, err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}
	return &e, nil
}
