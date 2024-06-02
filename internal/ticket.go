package internal

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type TicketRepo struct {
	DB *pgxpool.Pool
}

type Ticket struct {
	ID           int       `json:"id"`
	SeatID       int       `json:"seat_id"`
	UserID       int       `json:"user_id"`
	PurchaseDate time.Time `json:"purchase_date"`
}
