-- +goose Up

CREATE TABLE countries (
    code CHAR(2) PRIMARY KEY,
    name TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down

DROP TABLE countries;