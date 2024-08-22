package internal

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"strconv"
)

type SectorRepo struct {
	DB *pgxpool.Pool
}

type Sector struct {
	ID      *int    `json:"id"`
	VenueID *int    `json:"venue_id"`
	Name    *string `json:"name"`
	Height  *int    `json:"height"`
	Width   *int    `json:"width"`
	IsLink  *bool   `json:"isLink"`
	Left    *int    `json:"left"`
	Top     *int    `json:"top"`
	Image   *string `json:"image"`
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

func (m *SectorRepo) CreateSectors(venueId *int, sectors []*Sector) ([]*Sector, error) {
	tx, err := m.DB.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())
	rows, _ := tx.Query(context.Background(), `DELETE FROM sectors WHERE venue_id = $1 returning id`, venueId)
	for rows.Next() {

		var id int
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		os.RemoveAll("./static/sectors/" + strconv.Itoa(id))
	}

	for i, s := range sectors {
		err := tx.QueryRow(context.Background(), `INSERT INTO sectors
		(venue_id, name, height, width, is_link, "left", top, image)
		VALUES
		($1, $2, $3, $4, $5, $6, $7, $8) returning id`, venueId, s.Name, s.Height, s.Width, s.IsLink, s.Left, s.Top, s.Image).Scan(&s.ID)
		if err != nil {
			return nil, err
		}
		sectors[i] = s
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}
	return sectors, nil
}

func (m *SectorRepo) UpdateImage(filename *string, id *int) error {

	tx, err := m.DB.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())
	_, err = tx.Exec(context.Background(), `UPDATE sectors SET image = $1 WHERE id = $2`, filename, id)
	if err != nil {
		return err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}
	return nil
}
