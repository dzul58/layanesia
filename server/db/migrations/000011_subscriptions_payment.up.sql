-- ============================================================================
-- 000011 - Subscriptions: kolom pembayaran Midtrans & siklus status.
--   status         : PENDING -> ACTIVE -> EXPIRED  (atau CANCELLED)
--   payment_status : PENDING -> PAID | FAILED | EXPIRED
-- ============================================================================

ALTER TABLE subscriptions
    ADD COLUMN IF NOT EXISTS order_id         VARCHAR(64),
    ADD COLUMN IF NOT EXISTS payment_status   VARCHAR(30) NOT NULL DEFAULT 'PENDING',
    ADD COLUMN IF NOT EXISTS payment_provider VARCHAR(30) NOT NULL DEFAULT 'MIDTRANS',
    ADD COLUMN IF NOT EXISTS payment_type     VARCHAR(50),
    ADD COLUMN IF NOT EXISTS snap_token       VARCHAR(255),
    ADD COLUMN IF NOT EXISTS redirect_url     TEXT,
    ADD COLUMN IF NOT EXISTS paid_at          TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS updated_at       TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP;

-- Backfill data lama: semua langganan yang sudah ACTIVE dianggap sudah dibayar.
UPDATE subscriptions
SET order_id       = COALESCE(order_id, 'LEGACY-' || id::text),
    payment_status = CASE WHEN status = 'ACTIVE' THEN 'PAID' ELSE payment_status END,
    paid_at        = CASE WHEN status = 'ACTIVE' THEN COALESCE(paid_at, created_at) ELSE paid_at END
WHERE order_id IS NULL OR (status = 'ACTIVE' AND payment_status <> 'PAID');

ALTER TABLE subscriptions ALTER COLUMN order_id SET NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS uq_subscriptions_order_id ON subscriptions (order_id);

-- Default status baru: PENDING sampai pembayaran terkonfirmasi lewat webhook.
ALTER TABLE subscriptions ALTER COLUMN status SET DEFAULT 'PENDING';

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'chk_subscriptions_status') THEN
        ALTER TABLE subscriptions ADD CONSTRAINT chk_subscriptions_status
            CHECK (status IN ('PENDING', 'ACTIVE', 'EXPIRED', 'CANCELLED'));
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'chk_subscriptions_payment_status') THEN
        ALTER TABLE subscriptions ADD CONSTRAINT chk_subscriptions_payment_status
            CHECK (payment_status IN ('PENDING', 'PAID', 'FAILED', 'EXPIRED', 'REFUNDED'));
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'chk_subscriptions_paid_consistency') THEN
        ALTER TABLE subscriptions ADD CONSTRAINT chk_subscriptions_paid_consistency
            CHECK (status <> 'ACTIVE' OR payment_status = 'PAID');
    END IF;
END $$;

-- Hanya satu checkout PENDING per user per paket pada satu waktu.
CREATE UNIQUE INDEX IF NOT EXISTS uq_subscriptions_pending_per_plan
    ON subscriptions (user_id, plan_type)
    WHERE status = 'PENDING';

-- Periode ACTIVE untuk user+paket yang sama tidak boleh saling tumpang tindih
-- (mengizinkan perpanjangan berurutan, menolak duplikasi aktif secara bersamaan).
CREATE EXTENSION IF NOT EXISTS btree_gist;
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'excl_subscriptions_active_overlap') THEN
        ALTER TABLE subscriptions ADD CONSTRAINT excl_subscriptions_active_overlap
            EXCLUDE USING gist (
                user_id   WITH =,
                plan_type WITH =,
                tstzrange(starts_at, expires_at, '[)') WITH &&
            ) WHERE (status = 'ACTIVE');
    END IF;
END $$;
