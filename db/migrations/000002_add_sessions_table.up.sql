
create table if not exists sessions(
    id text primary key not null,
    user_id text not null,
    refresh_token text not null,
    created_at timestamp default now(),
    expires_at timestamp
)