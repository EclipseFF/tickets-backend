package internal

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SeatRepo struct {
	DB *pgxpool.Pool
}

func (m *SeatRepo) GetSeatsBySectorID(sectorId int) ([][]*Seat, error) {
	tx, err := m.DB.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())
	rows, err := tx.Query(context.Background(), `SELECT id FROM seats WHERE sector_id = $1`, sectorId)
	tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}

	seats := make([]*Seat, 0)
	for rows.Next() {
		var s Seat
		err := rows.Scan(&s.Price)
		if err != nil {
			return nil, err
		}
		seats = append(seats, &s)
	}
	resp := make([][]*Seat, 0)
	for _, seat := range seats {
		resp = append(resp, []*Seat{seat})
	}
	return resp, nil
}
