CREATE TABLE shah_seats (
                       id SERIAL PRIMARY KEY,
                        venue_id INT REFERENCES venues(id) ON DELETE CASCADE,
                       num INT,
                       "left" INT,
                       top INT,
                       price INT,
                       bg_color VARCHAR(255),
                       text_color VARCHAR(255),
                       date timestamp
);
CREATE TABLE shah_ticket_types (
                              id SERIAL PRIMARY KEY,
                              name VARCHAR(255),
                              price INT,
                              amount INT
);

CREATE TABLE shah_seat_ticket_types (
                                   seat_id INT REFERENCES shah_seats(id) ON DELETE CASCADE,
                                   ticket_type_id INT REFERENCES shah_ticket_types(id) ON DELETE CASCADE,
                                   PRIMARY KEY (seat_id, ticket_type_id)
);