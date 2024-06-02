CREATE TABLE events (
                        id SERIAL PRIMARY KEY,
                        title TEXT,
                        description TEXT,
                        brief_desc TEXT,
                        genre TEXT[], -- genre stored as an array of text
                        start_time TIMESTAMP,
                        end_time TIMESTAMP,
                        price NUMERIC,
                        age_restriction INTEGER,
                        rating NUMERIC,
                        created_at TIMESTAMP,
                        updated_at TIMESTAMP
);

create table event_types (
                             event_id INTEGER REFERENCES events(id) ON DELETE CASCADE,
                             type_id INTEGER REFERENCES types(id) ON DELETE CASCADE,
                             PRIMARY KEY (event_id, type_id)
);

CREATE TABLE event_venues (
                              event_id INTEGER REFERENCES events(id) ON DELETE CASCADE,
                              venue_id INTEGER REFERENCES venues(id) ON DELETE CASCADE,
                              PRIMARY KEY (event_id, venue_id)
);

create TABLE additional_user_data (
    user_id int primary key references users(id),
    surname text,
    name text,
    patronymic text,
    date_of_birth date
);

create table event_images(
    event_id int primary key references events(id),
    posters text[],
    main_images text[]
);