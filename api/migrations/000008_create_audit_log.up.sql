-- Migration: 000008_create_audit_log.up.sql
-- Description: Create the immutable append-only audit_log table
-- SECURITY RULE: UPDATE and DELETE are revoked at the database permission level
-- No application code, user, or role may modify or erase audit records

CREATE TABLE audit_log (
    id            BIGSERIAL PRIMARY KEY,
    -- Actor performing the action
    actor_id      UUID,
    actor_role    VARCHAR(100),
    actor_ip      INET,
    -- Action performed (e.g. CLAIM_SUBMITTED, STATUS_CHANGED, SOURCE_UPLOADED)
    action        VARCHAR(100) NOT NULL,
    -- Target entity
    target_table  VARCHAR(100) NOT NULL,
    target_id     UUID,
    -- JSON payloads for before/after state (NULL for creates)
    payload_before JSONB,
    payload_after  JSONB,
    -- Immutable timestamp
    occurred_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE audit_log IS 'Immutable append-only audit trail. UPDATE and DELETE are revoked at DB level. No record may ever be altered or erased.';
COMMENT ON COLUMN audit_log.payload_before IS 'JSON snapshot of the entity state before the action. NULL for INSERT operations.';
COMMENT ON COLUMN audit_log.payload_after IS 'JSON snapshot of the entity state after the action. NULL for DELETE operations.';

-- SECURITY: Revoke destructive permissions from ALL roles including PUBLIC and api_user
REVOKE UPDATE, DELETE ON audit_log FROM PUBLIC;
