CREATE TABLE event_days_no_shah (
                            id SERIAL PRIMARY KEY,
                            event_id INT REFERENCES events(id),
                            venue_id INT REFERENCES venues(id),
                            date TIMESTAMPTZ NOT NULL
);

CREATE TABLE ticket_types_no_shah (
                              id SERIAL PRIMARY KEY,
                              event_day_id INT REFERENCES event_days_no_shah(id),
                              name VARCHAR(255) NOT NULL,
                              price DECIMAL(10, 2) NOT NULL,
                              amount INT NOT NULL, -- Total number of tickets available
                              sold_count INT DEFAULT 0, -- Number of tickets sold
                              version INT DEFAULT 1     -- For optimistic locking
);

CREATE TABLE tickets_no_shah (
                         id SERIAL PRIMARY KEY,
                         ticket_type_id INT REFERENCES ticket_types_no_shah(id),
                         purchase_time TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                         is_reserved BOOLEAN DEFAULT FALSE, -- To handle reservation before final purchase
                         user_id INT -- Assuming users table exists
);