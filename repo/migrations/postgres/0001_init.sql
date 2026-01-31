-- +goose Up

create table if not exists credentials (
    user_id       uuid primary key default uuidv7(),
    email         text not null,
    password_hash text not null,
    created_at    timestamptz not null default now(),

    constraint credentials_email_unique unique (email)
);

-- +goose Down

drop table if exists credentials;
