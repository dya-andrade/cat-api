CREATE TABLE IF NOT EXISTS cat_thumbnails (
    id          BIGSERIAL PRIMARY KEY,
    cat_id      BIGINT NOT NULL REFERENCES cats(id) ON DELETE CASCADE,
    path        TEXT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_thumbs_cat_id ON cat_thumbnails(cat_id);