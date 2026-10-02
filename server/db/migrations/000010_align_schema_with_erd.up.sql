-- ============================================================================
-- 000010 - Menyelaraskan skema fisik dengan ERD_v0.md
-- Migrasi ini idempoten dan aman dijalankan pada:
--   (a) database baru yang dibuat dari 000001..000009, maupun
--   (b) database lama yang dibuat oleh GORM AutoMigrate (tanpa CASCADE/CHECK,
--       dengan kolom legacy users.password_hash).
-- ============================================================================

-- 1. Kolom legacy yang membuat registrasi gagal (NOT NULL tanpa default).
ALTER TABLE users DROP COLUMN IF EXISTS password_hash;

-- 2. Default timestamp yang hilang pada DB hasil AutoMigrate.
ALTER TABLE users          ALTER COLUMN created_at SET DEFAULT CURRENT_TIMESTAMP, ALTER COLUMN updated_at SET DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE subscriptions  ALTER COLUMN created_at SET DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE skill_postings ALTER COLUMN created_at SET DEFAULT CURRENT_TIMESTAMP, ALTER COLUMN updated_at SET DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE job_postings   ALTER COLUMN created_at SET DEFAULT CURRENT_TIMESTAMP, ALTER COLUMN updated_at SET DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE job_applications ALTER COLUMN created_at SET DEFAULT CURRENT_TIMESTAMP, ALTER COLUMN updated_at SET DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE job_offers     ALTER COLUMN created_at SET DEFAULT CURRENT_TIMESTAMP, ALTER COLUMN updated_at SET DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE connections    ALTER COLUMN created_at SET DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE ratings        ALTER COLUMN created_at SET DEFAULT CURRENT_TIMESTAMP;

-- 3. Semua FOREIGN KEY harus ON DELETE CASCADE (sesuai ERD).
DO $$
DECLARE
    r RECORD;
    base_def TEXT;
BEGIN
    FOR r IN
        SELECT c.conname,
               c.conrelid::regclass AS tbl,
               pg_get_constraintdef(c.oid) AS def
        FROM pg_constraint c
        JOIN pg_namespace n ON n.oid = c.connamespace
        WHERE n.nspname = 'public'
          AND c.contype = 'f'
          AND c.confdeltype <> 'c'
    LOOP
        base_def := regexp_replace(r.def, '\s+ON DELETE\s+(NO ACTION|RESTRICT|CASCADE|SET NULL|SET DEFAULT)', '', 'gi');
        EXECUTE format('ALTER TABLE %s DROP CONSTRAINT %I', r.tbl, r.conname);
        EXECUTE format('ALTER TABLE %s ADD CONSTRAINT %I %s ON DELETE CASCADE', r.tbl, r.conname, base_def);
    END LOOP;
END $$;

-- 4. Unique job_offers harus PARTIAL (hanya untuk status PENDING) - ERD §2.6.
ALTER TABLE job_offers DROP CONSTRAINT IF EXISTS idx_offer_skill_employer;
DROP INDEX IF EXISTS idx_offer_skill_employer;
CREATE UNIQUE INDEX IF NOT EXISTS uq_job_offers_pending_per_employer
    ON job_offers (skill_posting_id, employer_id)
    WHERE status = 'PENDING';

-- 5. Satu connection per sumber transaksi (polymorphic guard tambahan).
CREATE UNIQUE INDEX IF NOT EXISTS uq_connections_source
    ON connections (source_type, source_id);

-- 6. Master data lokasi harus unik. Bersihkan duplikat sebelum memasang constraint.
DELETE FROM locations a
USING locations b
WHERE a.id > b.id
  AND a.province = b.province
  AND a.city_or_regency = b.city_or_regency
  AND a.district = b.district;
CREATE UNIQUE INDEX IF NOT EXISTS uq_locations_province_city_district
    ON locations (province, city_or_regency, district);

-- 7. CHECK constraints: enum, nilai positif, dan anti self-reference.
DO $$
BEGIN
    -- users
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'chk_users_active_mode') THEN
        ALTER TABLE users ADD CONSTRAINT chk_users_active_mode CHECK (active_mode IN ('SEEKER', 'EMPLOYER'));
    END IF;

    -- subscriptions
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'chk_subscriptions_plan_type') THEN
        ALTER TABLE subscriptions ADD CONSTRAINT chk_subscriptions_plan_type CHECK (plan_type IN ('APPLY_JOB_5K', 'POST_SKILL_10K'));
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'chk_subscriptions_amount') THEN
        ALTER TABLE subscriptions ADD CONSTRAINT chk_subscriptions_amount CHECK (amount >= 0);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'chk_subscriptions_period') THEN
        ALTER TABLE subscriptions ADD CONSTRAINT chk_subscriptions_period CHECK (expires_at > starts_at);
    END IF;

    -- skill_postings
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'chk_skill_postings_category') THEN
        ALTER TABLE skill_postings ADD CONSTRAINT chk_skill_postings_category CHECK (category IN ('SERABUTAN', 'PROFESIONAL'));
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'chk_skill_postings_rate_type') THEN
        ALTER TABLE skill_postings ADD CONSTRAINT chk_skill_postings_rate_type CHECK (rate_type IN ('PER_HOUR', 'PER_DAY'));
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'chk_skill_postings_availability') THEN
        ALTER TABLE skill_postings ADD CONSTRAINT chk_skill_postings_availability CHECK (availability IN ('AVAILABLE', 'BUSY'));
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'chk_skill_postings_rate_amount') THEN
        ALTER TABLE skill_postings ADD CONSTRAINT chk_skill_postings_rate_amount CHECK (rate_amount > 0);
    END IF;

    -- job_postings
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'chk_job_postings_category') THEN
        ALTER TABLE job_postings ADD CONSTRAINT chk_job_postings_category CHECK (category IN ('SERABUTAN', 'PROFESIONAL'));
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'chk_job_postings_duration_type') THEN
        ALTER TABLE job_postings ADD CONSTRAINT chk_job_postings_duration_type CHECK (duration_type IN ('HOURS', 'DAYS'));
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'chk_job_postings_status') THEN
        ALTER TABLE job_postings ADD CONSTRAINT chk_job_postings_status CHECK (status IN ('OPEN', 'IN_PROGRESS', 'DONE', 'CANCELLED'));
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'chk_job_postings_duration_value') THEN
        ALTER TABLE job_postings ADD CONSTRAINT chk_job_postings_duration_value CHECK (duration_value > 0);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'chk_job_postings_budget') THEN
        ALTER TABLE job_postings ADD CONSTRAINT chk_job_postings_budget CHECK (budget > 0);
    END IF;

    -- job_applications
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'chk_job_applications_status') THEN
        ALTER TABLE job_applications ADD CONSTRAINT chk_job_applications_status CHECK (status IN ('PENDING', 'ACCEPTED', 'REJECTED'));
    END IF;

    -- job_offers
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'chk_job_offers_status') THEN
        ALTER TABLE job_offers ADD CONSTRAINT chk_job_offers_status CHECK (status IN ('PENDING', 'ACCEPTED', 'REJECTED'));
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'chk_job_offers_budget') THEN
        ALTER TABLE job_offers ADD CONSTRAINT chk_job_offers_budget CHECK (offered_budget > 0);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'chk_job_offers_not_self') THEN
        ALTER TABLE job_offers ADD CONSTRAINT chk_job_offers_not_self CHECK (employer_id <> worker_id);
    END IF;

    -- connections
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'chk_connections_source_type') THEN
        ALTER TABLE connections ADD CONSTRAINT chk_connections_source_type CHECK (source_type IN ('JOB_APPLICATION', 'JOB_OFFER'));
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'chk_connections_status') THEN
        ALTER TABLE connections ADD CONSTRAINT chk_connections_status CHECK (status IN ('ACTIVE', 'COMPLETED'));
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'chk_connections_not_self') THEN
        ALTER TABLE connections ADD CONSTRAINT chk_connections_not_self CHECK (employer_id <> worker_id);
    END IF;

    -- ratings
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'chk_ratings_stars') THEN
        ALTER TABLE ratings ADD CONSTRAINT chk_ratings_stars CHECK (rating_stars >= 1 AND rating_stars <= 5);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'chk_ratings_reviewer_role') THEN
        ALTER TABLE ratings ADD CONSTRAINT chk_ratings_reviewer_role CHECK (reviewer_role IN ('EMPLOYER', 'WORKER'));
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'chk_ratings_not_self') THEN
        ALTER TABLE ratings ADD CONSTRAINT chk_ratings_not_self CHECK (reviewer_id <> reviewee_id);
    END IF;
END $$;
