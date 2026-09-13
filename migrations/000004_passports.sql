-- +goose Up

CREATE TABLE passports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id UUID NOT NULL
        REFERENCES users(id)
        ON DELETE CASCADE,

    nationality_code CHAR(2) NOT NULL
        REFERENCES countries(code),

    expires_at DATE NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_passports_user_id
    ON passports(user_id);

-- +goose Down

DROP TABLE passports;