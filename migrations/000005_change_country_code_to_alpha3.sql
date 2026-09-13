-- +goose Up

ALTER TABLE passports
    DROP CONSTRAINT passports_nationality_code_fkey;

ALTER TABLE countries
    ALTER COLUMN code TYPE CHAR(3);

ALTER TABLE passports
    ALTER COLUMN nationality_code TYPE CHAR(3);

ALTER TABLE passports
    ADD CONSTRAINT passports_nationality_code_fkey
    FOREIGN KEY (nationality_code)
    REFERENCES countries(code);

-- +goose Down

ALTER TABLE passports
    DROP CONSTRAINT passports_nationality_code_fkey;

ALTER TABLE passports
    ALTER COLUMN nationality_code TYPE CHAR(2);

ALTER TABLE countries
    ALTER COLUMN code TYPE CHAR(2);

ALTER TABLE passports
    ADD CONSTRAINT passports_nationality_code_fkey
    FOREIGN KEY (nationality_code)
    REFERENCES countries(code);