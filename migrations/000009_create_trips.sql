-- +goose Up

CREATE TABLE trips (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id UUID NOT NULL
        REFERENCES users(id)
        ON DELETE CASCADE,

    passport_id UUID NOT NULL
        REFERENCES passports(id)
        ON DELETE CASCADE,

    destination_code CHAR(3) NOT NULL
        REFERENCES countries(code),

    departure_date DATE NOT NULL,
    return_date DATE NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CHECK (return_date >= departure_date)
);

CREATE INDEX idx_trips_user_id
    ON trips (user_id);

CREATE INDEX idx_trips_passport_id
    ON trips (passport_id);

CREATE INDEX idx_trips_destination
    ON trips (destination_code);

-- +goose Down

DROP TABLE trips;