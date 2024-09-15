create table if not exists users
(
    id   uuid default gen_random_uuid() primary key,
    username text not null,
    email text not null unique
)