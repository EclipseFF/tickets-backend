create table types(
                            id SERIAL primary key,
                            name text,
                            translated_name text
);

CREATE TABLE venues (
                        id SERIAL PRIMARY KEY,
                        name VARCHAR(255) NOT NULL,
                        location VARCHAR(255)
);

CREATE TABLE sectors (
                         id SERIAL PRIMARY KEY,
                         venue_id INT NOT NULL,
                         name VARCHAR(255) NOT NULL,
                         layout JSONB NOT NULL,
                         FOREIGN KEY (venue_id) REFERENCES venues (id) ON DELETE CASCADE
);

CREATE TABLE seats (
                       id SERIAL PRIMARY KEY,
                       sector_id INT NOT NULL,
                       x INT NOT NULL,
                       y INT NOT NULL,
                       is_available BOOLEAN DEFAULT TRUE,
                       FOREIGN KEY (sector_id) REFERENCES sectors (id) ON DELETE CASCADE
);

CREATE TABLE tickets (
                         id SERIAL PRIMARY KEY,
                         seat_id INT NOT NULL,
                         user_id INT NOT NULL,
                         purchase_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                         FOREIGN KEY (seat_id) REFERENCES seats (id) ON DELETE CASCADE
);

CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       email text unique ,
                       password text,
                       phone text unique
);

create table sessions(
                         token text primary key,
                         user_id int references users(id)
);