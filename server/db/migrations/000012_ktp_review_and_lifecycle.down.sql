ALTER TABLE connections  DROP COLUMN IF EXISTS updated_at, DROP COLUMN IF EXISTS completed_at;
ALTER TABLE job_postings DROP COLUMN IF EXISTS closed_at;
DROP INDEX IF EXISTS idx_users_ktp_status;
ALTER TABLE users
    DROP CONSTRAINT IF EXISTS chk_users_ktp_flag_consistency,
    DROP CONSTRAINT IF EXISTS chk_users_ktp_status;
ALTER TABLE users
    DROP COLUMN IF EXISTS ktp_rejection_reason,
    DROP COLUMN IF EXISTS ktp_verified_at,
    DROP COLUMN IF EXISTS ktp_submitted_at,
    DROP COLUMN IF EXISTS ktp_status;
