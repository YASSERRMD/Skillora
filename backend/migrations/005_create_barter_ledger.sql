-- Migration: 005_create_barter_ledger
-- Implements a double-entry accounting system for the barter economy.

-- Barter transactions record the agreement between two users.
CREATE TABLE IF NOT EXISTS barter_transactions (
    id             UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    initiator_id   UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    receiver_id    UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    initiator_skill_id UUID   NOT NULL REFERENCES skills(id) ON DELETE CASCADE,
    receiver_skill_id  UUID   NOT NULL REFERENCES skills(id) ON DELETE CASCADE,
    credit_amount  INT         NOT NULL CHECK (credit_amount > 0),
    status         VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'accepted', 'completed', 'cancelled')),
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_barter_initiator ON barter_transactions(initiator_id);
CREATE INDEX IF NOT EXISTS idx_barter_receiver  ON barter_transactions(receiver_id);
CREATE INDEX IF NOT EXISTS idx_barter_status    ON barter_transactions(status);

-- Ledger entries: every completed barter produces exactly two entries.
-- Debit: credit_amount subtracted from payer
-- Credit: credit_amount added to payee
CREATE TABLE IF NOT EXISTS ledger_entries (
    id             UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_id UUID        NOT NULL REFERENCES barter_transactions(id) ON DELETE CASCADE,
    user_id        UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    entry_type     VARCHAR(6)  NOT NULL CHECK (entry_type IN ('debit', 'credit')),
    amount         INT         NOT NULL CHECK (amount > 0),
    balance_after  INT         NOT NULL,   -- Running balance snapshot
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_ledger_user_id        ON ledger_entries(user_id);
CREATE INDEX IF NOT EXISTS idx_ledger_transaction_id ON ledger_entries(transaction_id);

-- Apply the updated_at trigger on barter_transactions.
CREATE TRIGGER set_barter_transactions_updated_at
    BEFORE UPDATE ON barter_transactions
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
