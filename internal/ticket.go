package internal

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type TicketRepo struct {
	DB *pgxpool.Pool
}

type TicketTypeNoShah struct {
	ID         *int    `json:"id"`
	EventDayId *int    `json:"event_day_id"`
	Name       *string `json:"name"`
	Price      *int    `json:"price"`
	Amount     *int    `json:"amount"`
	SoldCount  *int    `json:"sold_count"`
	Version    *int    `json:"version"`
}

type DateWithTicketsNoShah struct {
	ID      *int                `json:"id"`
	EventId *int                `json:"event_id"`
	Date    *time.Time          `json:"date"`
	Types   []*TicketTypeNoShah `json:"types"`
}

type Seat struct {
	Id        *int          `json:"id"`
	VenueId   *int          `json:"venue_id"`
	Num       *int          `json:"num"`
	Left      *int          `json:"left"`
	Top       *int          `json:"top"`
	Price     *int          `json:"price"`
	BgColor   *string       `json:"bgColor"`
	TextColor *string       `json:"textColor"`
	Types     []*TicketType `json:"types"`
	Date      *time.Time    `json:"date"`
}

type TicketType struct {
	ID     *int    `json:"id"`
	Name   *string `json:"name"`
	Price  *int    `json:"price"`
	Amount *int    `json:"amount"`
}

func (r *TicketRepo) CreateTicketsNoSham(eventId, venueId *int, days []*DateWithTicketsNoShah) error {
	tx, err := r.DB.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	for _, day := range days {
		var id int
		err := tx.QueryRow(context.Background(), `INSERT INTO event_days_no_shah(id, event_id, venue_id, date) VALUES(default, $1, $2, $3) RETURNING id`, eventId, venueId, day.Date).Scan(&id)
		if err != nil {
			return err
		}
		for _, t := range day.Types {
			_, err := tx.Exec(context.Background(), `INSERT INTO ticket_types_no_shah(id, event_day_id, name, price, amount, sold_count, version) VALUES(default, $1, $2, $3, $4, $5, $6)`, id, t.Name, t.Price, t.Amount, 0, 1)
			if err != nil {
				return err
			}
		}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}
	return nil

}

type TicketPurchaseResult struct {
	TicketID         int
	PurchaseTime     time.Time
	RemainingTickets int
}

func (r *TicketRepo) BuyTicketNoShah(ticketTypeID, userID, count *int) (*TicketPurchaseResult, error) {
	var result TicketPurchaseResult
	ctx := context.Background()
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// Lock the ticket type row to prevent concurrent updates
	row := tx.QueryRow(ctx, `
		SELECT id, amount - sold_count AS remaining_tickets
		FROM ticket_types_no_shah	
		WHERE id = $1
		FOR UPDATE
	`, ticketTypeID)

	var remainingTickets int
	err = row.Scan(&result.TicketID, &remainingTickets)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("ticket type not found")
		}
		return nil, err
	}

	if remainingTickets <= 0 {
		return nil, fmt.Errorf("no tickets available")
	}

	// Insert the purchased ticket and get its ID and purchase time
	err = tx.QueryRow(ctx, `
		WITH inserted_ticket AS (
			INSERT INTO tickets (ticket_type_id, user_id)
			VALUES ($1, $2)
			RETURNING id, purchase_time
		)
		UPDATE ticket_types
		SET sold_count = sold_count + 1
		WHERE id = $1
		RETURNING (SELECT id FROM inserted_ticket), 
				  (SELECT purchase_time FROM inserted_ticket), 
				  amount - sold_count AS remaining_tickets
	`, ticketTypeID, userID).Scan(&result.TicketID, &result.PurchaseTime, &result.RemainingTickets)

	if err != nil {
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *TicketRepo) GetDatesForEventVenue(eventId, venueId *int) ([]*DateWithTicketsNoShah, error) {

	tx, err := r.DB.Begin(context.Background())
	if err != nil {
		return nil, err
	}

	defer tx.Rollback(context.Background())
	rows, err := tx.Query(context.Background(), `SELECT id, date FROM event_days_no_shah WHERE event_id = $1 AND venue_id = $2`, eventId, venueId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	dates := make([]*DateWithTicketsNoShah, 0)
	for rows.Next() {
		var d DateWithTicketsNoShah
		err = rows.Scan(&d.ID, &d.Date)
		if err != nil {
			return nil, err
		}
		temp, err := r.GetTypesForDate(d.ID)
		if err != nil {
			return nil, err
		}
		d.Types = temp
		d.EventId = eventId
		dates = append(dates, &d)
	}
	return dates, nil
}

func (r *TicketRepo) GetTypesForDate(dateId *int) ([]*TicketTypeNoShah, error) {
	tx, err := r.DB.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())
	rows, err := tx.Query(context.Background(), `SELECT id, name, price, amount, sold_count, version FROM ticket_types_no_shah WHERE event_day_id = $1`, dateId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]*TicketTypeNoShah, 0)
	for rows.Next() {
		var t TicketTypeNoShah
		t.EventDayId = dateId
		err = rows.Scan(&t.ID, &t.Name, &t.Price, &t.Amount, &t.SoldCount, &t.Version)
		if err != nil {
			return nil, err
		}
		result = append(result, &t)
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *TicketRepo) CreateTicketsWithSham(eventId, venueId *int, seats [][]*Seat) error {
	tx, err := r.DB.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(), `DELETE FROM shah_seats WHERE venue_id = $1`, venueId)
	if err != nil {
		return err
	}

	uniqueTypes := GetUniqueTicketTypes(seats)

	for i, ticketType := range uniqueTypes {
		if ticketType.Amount == nil {
			temp := 1
			uniqueTypes[i].Amount = &temp
		}
		err := tx.QueryRow(context.Background(), `INSERT INTO shah_ticket_types(id, name, price, amount) values (default, $1, $2, $3) returning id`, ticketType.Name, ticketType.Price, ticketType.Amount).Scan(&uniqueTypes[i].ID)
		if err != nil {
			return err
		}
	}

	for _, uniqueType := range uniqueTypes {
		for rowIndex, row := range seats {
			for seatIndex, seat := range row {
				for typeIndex, seatType := range seat.Types {
					if seatType.Amount == nil {
						seats[rowIndex][seatIndex].Types[typeIndex].Amount = uniqueType.Amount
					}
					if *uniqueType.Name == *seatType.Name && *uniqueType.Price == *seatType.Price && *uniqueType.Amount == *seatType.Amount {
						seats[rowIndex][seatIndex].Types[typeIndex].ID = uniqueType.ID
					}
				}
			}
		}
	}

	for _, row := range seats {
		for _, seat := range row {
			var id int
			err := tx.QueryRow(context.Background(), `INSERT INTO shah_seats(id, venue_id, num, "left", top, price, bg_color, text_color) values (default, $1, $2, $3, $4, $5, $6, $7) returning id`, venueId, seat.Num, seat.Left, seat.Top, seat.Price, seat.BgColor, seat.TextColor).Scan(&id)
			if err != nil {
				return err
			}

			for _, stype := range seat.Types {
				_, err = tx.Exec(context.Background(), `INSERT INTO shah_seat_ticket_types(seat_id, ticket_type_id) values ($1, $2)`, id, stype.ID)
				if err != nil {
					return err
				}
			}
		}
	}
	return tx.Commit(context.Background())

}

func GetUniqueTicketTypes(seats [][]*Seat) []TicketType {
	uniqueTypes := make(map[int]TicketType)

	for _, row := range seats {
		for _, seat := range row {
			if seat != nil {
				for _, ticketType := range seat.Types {
					if ticketType != nil && ticketType.ID != nil {
						if _, exists := uniqueTypes[*ticketType.ID]; !exists {
							uniqueTypes[*ticketType.ID] = *ticketType
						}
					}
				}
			}
		}
	}

	result := make([]TicketType, 0, len(uniqueTypes))
	for _, ticketType := range uniqueTypes {
		result = append(result, ticketType)
	}

	return result
}

func (r *TicketRepo) AddType(seatId, typeId *int) error {
	tx, err := r.DB.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())
	_, err = tx.Exec(context.Background(), `INSERT INTO shah_seat_ticket_types(seat_id, ticket_type_id) values ($1, $2)`, seatId, typeId)
	if err != nil {
		return err
	}
	return tx.Commit(context.Background())
}

func (r *TicketRepo) GetDatesForEventVenueShah(eventId, venueId *int) ([]*Seat, error) {

	tx, err := r.DB.Begin(context.Background())
	if err != nil {
		return nil, err
	}

	defer tx.Rollback(context.Background())
	rows, err := tx.Query(context.Background(), `SELECT id, venue_id, num, "left", top, price, bg_color, text_color, date FROM shah_seats WHERE venue_id = $1`, venueId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	dates := make([]*Seat, 0)
	for rows.Next() {
		var d Seat
		err = rows.Scan(&d.Id, &d.VenueId, &d.Num, &d.Left, &d.Top, &d.Price, &d.BgColor, &d.TextColor, &d.Date)
		if err != nil {
			return nil, err
		}

		dates = append(dates, &d)
	}
	return dates, nil
}
