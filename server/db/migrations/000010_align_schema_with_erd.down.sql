ALTER TABLE ratings          DROP CONSTRAINT IF EXISTS chk_ratings_not_self, DROP CONSTRAINT IF EXISTS chk_ratings_reviewer_role, DROP CONSTRAINT IF EXISTS chk_ratings_stars;
ALTER TABLE connections      DROP CONSTRAINT IF EXISTS chk_connections_not_self, DROP CONSTRAINT IF EXISTS chk_connections_status, DROP CONSTRAINT IF EXISTS chk_connections_source_type;
ALTER TABLE job_offers       DROP CONSTRAINT IF EXISTS chk_job_offers_not_self, DROP CONSTRAINT IF EXISTS chk_job_offers_budget, DROP CONSTRAINT IF EXISTS chk_job_offers_status;
ALTER TABLE job_applications DROP CONSTRAINT IF EXISTS chk_job_applications_status;
ALTER TABLE job_postings     DROP CONSTRAINT IF EXISTS chk_job_postings_budget, DROP CONSTRAINT IF EXISTS chk_job_postings_duration_value, DROP CONSTRAINT IF EXISTS chk_job_postings_status, DROP CONSTRAINT IF EXISTS chk_job_postings_duration_type, DROP CONSTRAINT IF EXISTS chk_job_postings_category;
ALTER TABLE skill_postings   DROP CONSTRAINT IF EXISTS chk_skill_postings_rate_amount, DROP CONSTRAINT IF EXISTS chk_skill_postings_availability, DROP CONSTRAINT IF EXISTS chk_skill_postings_rate_type, DROP CONSTRAINT IF EXISTS chk_skill_postings_category;
ALTER TABLE subscriptions    DROP CONSTRAINT IF EXISTS chk_subscriptions_period, DROP CONSTRAINT IF EXISTS chk_subscriptions_amount, DROP CONSTRAINT IF EXISTS chk_subscriptions_plan_type;
ALTER TABLE users            DROP CONSTRAINT IF EXISTS chk_users_active_mode;

DROP INDEX IF EXISTS uq_locations_province_city_district;
DROP INDEX IF EXISTS uq_connections_source;
DROP INDEX IF EXISTS uq_job_offers_pending_per_employer;
CREATE UNIQUE INDEX IF NOT EXISTS idx_offer_skill_employer ON job_offers (skill_posting_id, employer_id);
