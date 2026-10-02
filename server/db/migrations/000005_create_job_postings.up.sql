CREATE TABLE IF NOT EXISTS job_postings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employer_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
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
    work_date DATE NOT NULL,
    duration_type VARCHAR(30) NOT NULL,
    duration_value INT NOT NULL,
    budget NUMERIC(12,2) NOT NULL,
    status VARCHAR(30) DEFAULT 'OPEN',
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_job_postings_employer_id ON job_postings(employer_id);
