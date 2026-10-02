CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    active_mode VARCHAR(20) NOT NULL DEFAULT 'SEEKER',
    is_ktp_verified BOOLEAN DEFAULT false,
    ktp_number VARCHAR(50),
    ktp_image_url VARCHAR(500),
    resume_url VARCHAR(500),
    country VARCHAR(100) DEFAULT 'Indonesia',
    province VARCHAR(100) NOT NULL,
    city VARCHAR(100) NOT NULL,
    district VARCHAR(100) NOT NULL,
    address_detail TEXT,
    latitude DECIMAL(10,8),
    longitude DECIMAL(11,8),
    formatted_address TEXT,
    osm_place_id VARCHAR(255),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_phone ON users(phone);
