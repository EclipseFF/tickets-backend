package internal

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type VenueRepo struct {
	DB *pgxpool.Pool
}

type Venue struct {
	ID       *int    `json:"id"`
	Name     *string `json:"name"`
	Location *string `json:"location"`
}

func (m *VenueRepo) CreateVenue(venue *Venue) (*int, error) {
	tx, err := m.DB.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())
	stmt := `INSERT INTO venues (id, name, location) VALUES (default, $1, $2) RETURNING id`
	err = tx.QueryRow(context.Background(), stmt, venue.Name, venue.Location).Scan(&venue.ID)
	if err != nil {
		return nil, err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}
	return venue.ID, nil
}

func (m *VenueRepo) GetVenuesByEvent(eventId *int) ([]*Venue, error) {
	tx, err := m.DB.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())
	rows, err := tx.Query(context.Background(), "SELECT venue_id FROM event_venues where event_id = $1", eventId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	venues := make([]*Venue, 0)
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		v, err := m.GetVenueById(&id)
		if err != nil {
			return nil, err
		}
		venues = append(venues, v)
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}
	return venues, nil
}

func (m *VenueRepo) GetVenueById(id *int) (*Venue, error) {
	tx, err := m.DB.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())
	var v Venue
	v.ID = id
	err = tx.QueryRow(context.Background(), `SELECT name, location FROM venues where id = $1`, *id).Scan(&v.Name, &v.Location)
	if err != nil {
		return nil, err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (m *VenueRepo) GetAll() ([]*Venue, error) {
	tx, err := m.DB.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())
	rows, err := tx.Query(context.Background(), "SELECT id, name, location FROM venues")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	venues := make([]*Venue, 0)
	for rows.Next() {
		var v Venue
		err = rows.Scan(&v.ID, &v.Name, &v.Location)
		if err != nil {
			return nil, err
		}
		venues = append(venues, &v)
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}
	return venues, nil
}
