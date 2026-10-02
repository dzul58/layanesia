CREATE TABLE IF NOT EXISTS skill_postings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category VARCHAR(50) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    country VARCHAR(100) DEFAULT 'Indonesia',
    province VARCHAR(100) NOT NULL,
    city VARCHAR(100) NOT NULL,
    district VARCHAR(100) NOT NULL,
    address_detail TEXT,
    latitude DECIMAL(10,8),
    longitude DECIMAL(11,8),
    formatted_address TEXT,
    osm_place_id VARCHAR(255),
    rate_type VARCHAR(30) NOT NULL,
    rate_amount NUMERIC(12,2) NOT NULL,
    availability VARCHAR(30) DEFAULT 'AVAILABLE',
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_skill_postings_user_id ON skill_postings(user_id);
