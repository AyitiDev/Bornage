-- Migration: 000001_create_enums.down.sql
-- Description: Drop all LADM domain enum types (reverse of up migration)

DROP TYPE IF EXISTS capture_method;
DROP TYPE IF EXISTS party_type;
DROP TYPE IF EXISTS rrr_type;
DROP TYPE IF EXISTS claim_status;
