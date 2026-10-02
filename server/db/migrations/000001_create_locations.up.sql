CREATE TABLE IF NOT EXISTS locations (
    id SERIAL PRIMARY KEY,
    country VARCHAR(100) DEFAULT 'Indonesia' NOT NULL,
    province VARCHAR(100) NOT NULL,
    city_or_regency VARCHAR(100) NOT NULL,
    district VARCHAR(100) NOT NULL,
    postal_code VARCHAR(20)
);

CREATE INDEX IF NOT EXISTS idx_locations_province ON locations(province);
CREATE INDEX IF NOT EXISTS idx_locations_city_or_regency ON locations(city_or_regency);
CREATE INDEX IF NOT EXISTS idx_locations_district ON locations(district);
