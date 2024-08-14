create table news(
    id serial primary key,
    name text,
    images varchar[],
    description text,
    created_at timestamp
);