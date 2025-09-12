-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id         INTEGER      PRIMARY KEY,
    email      VARCHAR(255) NOT NULL UNIQUE,
    password   VARCHAR(255) NOT NULL,
    created_at TIMESTAMP    WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP    WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_person FOREIGN KEY(id) REFERENCES people(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd