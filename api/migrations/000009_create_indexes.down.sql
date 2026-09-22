-- Migration: 000009_create_indexes.down.sql
-- Description: Drop all GiST and B-Tree indexes (reverse of up migration)

DROP INDEX IF EXISTS idx_audit_log_occurred_at;
DROP INDEX IF EXISTS idx_audit_log_actor_id;
DROP INDEX IF EXISTS idx_audit_log_target_id;
DROP INDEX IF EXISTS idx_witness_claim_id;
DROP INDEX IF EXISTS idx_source_claim_id;
DROP INDEX IF EXISTS idx_rrr_spatial_unit_id;
DROP INDEX IF EXISTS idx_rrr_party_id;
DROP INDEX IF EXISTS idx_claim_objection_end_date;
DROP INDEX IF EXISTS idx_claim_status;
DROP INDEX IF EXISTS idx_claim_tracking_code;
DROP INDEX IF EXISTS idx_claim_geom;
DROP INDEX IF EXISTS idx_spatial_unit_geom;
