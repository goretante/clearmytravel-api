-- +goose Up

CREATE TABLE visa_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    nationality_code CHAR(3) NOT NULL
        REFERENCES countries(code),

    destination_code CHAR(3) NOT NULL
        REFERENCES countries(code),

    visa_required BOOLEAN NOT NULL,

    effective_from DATE NOT NULL,
    effective_to DATE,

    source_name TEXT NOT NULL,
    source_url TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CHECK (
        effective_to IS NULL
        OR effective_to >= effective_from
    )
);

CREATE INDEX idx_visa_rules_lookup
    ON visa_rules (
        nationality_code,
        destination_code,
        effective_from
    );

-- +goose Down

DROP TABLE visa_rules;