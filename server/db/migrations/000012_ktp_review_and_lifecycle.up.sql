-- ============================================================================
-- 000012 - Alur review verifikasi KTP & lifecycle status job/connection.
-- ============================================================================

-- Verifikasi KTP tidak lagi klaim mandiri: UNVERIFIED -> PENDING_REVIEW -> VERIFIED | REJECTED
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS ktp_status           VARCHAR(30) NOT NULL DEFAULT 'UNVERIFIED',
    ADD COLUMN IF NOT EXISTS ktp_submitted_at     TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS ktp_verified_at      TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS ktp_rejection_reason TEXT;

UPDATE users SET is_ktp_verified = FALSE WHERE is_ktp_verified IS NULL;
ALTER TABLE users ALTER COLUMN is_ktp_verified SET NOT NULL, ALTER COLUMN is_ktp_verified SET DEFAULT FALSE;

UPDATE users
SET ktp_status      = 'VERIFIED',
    ktp_verified_at = COALESCE(ktp_verified_at, updated_at, CURRENT_TIMESTAMP)
WHERE is_ktp_verified = TRUE AND ktp_status <> 'VERIFIED';

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'chk_users_ktp_status') THEN
        ALTER TABLE users ADD CONSTRAINT chk_users_ktp_status
            CHECK (ktp_status IN ('UNVERIFIED', 'PENDING_REVIEW', 'VERIFIED', 'REJECTED'));
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'chk_users_ktp_flag_consistency') THEN
        ALTER TABLE users ADD CONSTRAINT chk_users_ktp_flag_consistency
            CHECK (is_ktp_verified = (ktp_status = 'VERIFIED'));
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_users_ktp_status ON users (ktp_status) WHERE ktp_status = 'PENDING_REVIEW';

-- Lifecycle
ALTER TABLE job_postings ADD COLUMN IF NOT EXISTS closed_at    TIMESTAMPTZ;
ALTER TABLE connections  ADD COLUMN IF NOT EXISTS completed_at TIMESTAMPTZ;
ALTER TABLE connections  ADD COLUMN IF NOT EXISTS updated_at   TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP;
