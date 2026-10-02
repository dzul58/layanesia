-- ============================================================================
-- 000013 - Index untuk query terpanas (cek langganan, pencarian marketplace,
--          listing per user) + trigram untuk pencarian teks judul/deskripsi.
-- ============================================================================

-- Cek langganan aktif: WHERE user_id=? AND plan_type=? AND status=? AND expires_at>now()
CREATE INDEX IF NOT EXISTS idx_subscriptions_active_lookup
    ON subscriptions (user_id, plan_type, status, expires_at DESC);

-- Filter lokasi memakai equality case-insensitive (nilai berasal dari master data).
CREATE INDEX IF NOT EXISTS idx_job_postings_search
    ON job_postings (status, LOWER(province), LOWER(city), LOWER(district), created_at DESC);
CREATE INDEX IF NOT EXISTS idx_skill_postings_search
    ON skill_postings (availability, LOWER(province), LOWER(city), LOWER(district), created_at DESC);

CREATE INDEX IF NOT EXISTS idx_job_postings_created_at   ON job_postings (created_at DESC);
CREATE INDEX IF NOT EXISTS idx_skill_postings_created_at ON skill_postings (created_at DESC);

-- Listing per user
CREATE INDEX IF NOT EXISTS idx_job_applications_applicant ON job_applications (applicant_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_job_applications_job       ON job_applications (job_posting_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_job_offers_employer        ON job_offers (employer_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_job_offers_worker          ON job_offers (worker_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_connections_employer       ON connections (employer_id, status);
CREATE INDEX IF NOT EXISTS idx_connections_worker         ON connections (worker_id, status);
CREATE INDEX IF NOT EXISTS idx_ratings_reviewer           ON ratings (reviewer_id);

-- Pencarian teks ILIKE '%kata%' pada judul/deskripsi.
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE INDEX IF NOT EXISTS idx_job_postings_title_trgm   ON job_postings   USING gin (title gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_job_postings_desc_trgm    ON job_postings   USING gin (description gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_skill_postings_title_trgm ON skill_postings USING gin (title gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_skill_postings_desc_trgm  ON skill_postings USING gin (description gin_trgm_ops);
