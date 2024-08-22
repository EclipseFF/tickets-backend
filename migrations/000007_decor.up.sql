create table decors(
    id serial primary key,
    name text,
    images varchar,
    created_at timestamp,
    venue_id int references venues(id)
)