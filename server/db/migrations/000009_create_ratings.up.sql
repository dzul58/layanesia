CREATE TABLE IF NOT EXISTS ratings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    connection_id UUID NOT NULL REFERENCES connections(id) ON DELETE CASCADE,
    reviewer_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    reviewee_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    reviewer_role VARCHAR(30) NOT NULL,
    rating_stars INT NOT NULL CHECK (rating_stars >= 1 AND rating_stars <= 5),
    comment TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT idx_rating_conn_rev UNIQUE(connection_id, reviewer_id)
);

CREATE INDEX IF NOT EXISTS idx_ratings_reviewee_id ON ratings(reviewee_id);
