package internal

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SectorRepo struct {
	DB *pgxpool.Pool
}

type Sector struct {
	ID      int    `json:"id"`
	VenueID int    `json:"venue_id"`
	Name    string `json:"name"`
	//Layout  pgtype.JSONB `json:"layout"`
	Layout [][]*Seat `json:"layout"`
}

func (m *SectorRepo) GetSectorsByVenue(venueId int) ([]*Sector, error) {
	tx, err := m.DB.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())

	rows, err := tx.Query(context.Background(), `SELECT id, venue_id, name FROM sectors WHERE venue_id = $1`, venueId)
	tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}
	sectors := make([]*Sector, 0)
	for rows.Next() {
		var s Sector
		err = rows.Scan(&s.ID, &s.VenueID, &s.Name)
		if err != nil {
			return nil, err
		}
		sectors = append(sectors, &s)
	}
	return sectors, nil
}
