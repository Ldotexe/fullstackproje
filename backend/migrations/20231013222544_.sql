-- +goose Up
-- +goose StatementBegin
create table users(
    id BIGSERIAL PRIMARY KEY NOT NULL,
    login text NOT NULL,
    password_hash text not null,
    is_admin bool not null
);

create table messages(
    from_login text,
    to_login text,
    txt text,
    ts text
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table users;
drop table messages;
-- +goose StatementEnd

