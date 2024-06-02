create table admin_users(
    id serial primary key,
    email text not null unique,
    password text not null
);

create table admin_session(
                         token text primary key,
                         admin_id int references admin_users(id)
);