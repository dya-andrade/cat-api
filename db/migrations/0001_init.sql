CREATE TABLE IF NOT EXISTS cats (
    id          BIGSERIAL PRIMARY KEY,
    name        TEXT NOT NULL,
    age_years   INT  NOT NULL CHECK (age_years >= 0),
    breed       TEXT,
    coat_color  TEXT,
    weight_kg   NUMERIC(5,2) CHECK (weight_kg >= 0),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_cats_created_at ON cats(created_at DESC);

CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON cats
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();