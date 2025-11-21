-- +goose Up

ALTER TABLE users
ADD hashed_password VARCHAR NOT NULL DEFAULT 'unset';

-- +goose Down

ALTER TABLE users
DROP hashed_password;