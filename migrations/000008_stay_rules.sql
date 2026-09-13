-- +goose Up

CREATE TABLE stay_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    nationality_code CHAR(3) NOT NULL
        REFERENCES countries(code),

    destination_code CHAR(3) NOT NULL
        REFERENCES countries(code),

    max_stay_days INTEGER NOT NULL,

    effective_from DATE NOT NULL,
    effective_to DATE,

    source_name TEXT NOT NULL,
    source_url TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CHECK (max_stay_days > 0),

    CHECK (
        effective_to IS NULL
        OR effective_to >= effective_from
    )
);

CREATE INDEX idx_stay_rules_lookup
    ON stay_rules (
        nationality_code,
        destination_code,
        effective_from
    );

-- +goose Down

DROP TABLE stay_rules;