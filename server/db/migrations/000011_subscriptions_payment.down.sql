ALTER TABLE subscriptions DROP CONSTRAINT IF EXISTS excl_subscriptions_active_overlap;
DROP INDEX IF EXISTS uq_subscriptions_pending_per_plan;
ALTER TABLE subscriptions
    DROP CONSTRAINT IF EXISTS chk_subscriptions_paid_consistency,
    DROP CONSTRAINT IF EXISTS chk_subscriptions_payment_status,
    DROP CONSTRAINT IF EXISTS chk_subscriptions_status;
ALTER TABLE subscriptions ALTER COLUMN status SET DEFAULT 'ACTIVE';
DROP INDEX IF EXISTS uq_subscriptions_order_id;
ALTER TABLE subscriptions
    DROP COLUMN IF EXISTS updated_at,
    DROP COLUMN IF EXISTS paid_at,
    DROP COLUMN IF EXISTS redirect_url,
    DROP COLUMN IF EXISTS snap_token,
    DROP COLUMN IF EXISTS payment_type,
    DROP COLUMN IF EXISTS payment_provider,
    DROP COLUMN IF EXISTS payment_status,
    DROP COLUMN IF EXISTS order_id;
