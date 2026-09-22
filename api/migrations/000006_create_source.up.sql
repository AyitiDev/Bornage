-- Migration: 000006_create_source.up.sql
-- Description: Create the source table (LADM LA_AdministrativeSource)
-- Stores documentary or photographic evidence linked to a claim
-- SHA-256 hash ensures tamper detection after upload

CREATE TABLE source (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    claim_id     UUID NOT NULL REFERENCES claim(id) ON DELETE CASCADE,
    -- S3 compatible object key (stored in MinIO)
    s3_key       VARCHAR(500) NOT NULL,
    file_type    VARCHAR(100) NOT NULL,
    original_filename VARCHAR(255),
    -- SHA-256 hex digest (64 characters) - mandatory for tamper detection
    sha256_hash  CHAR(64) NOT NULL,
    file_size_bytes BIGINT,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    -- Enforce SHA-256 hex format: exactly 64 lowercase hex characters
    CONSTRAINT chk_sha256_format
        CHECK (sha256_hash ~ '^[a-f0-9]{64}$')
);

COMMENT ON TABLE source IS 'LADM LA_AdministrativeSource: Documentary and photographic evidence for a claim with SHA-256 tamper detection.';
COMMENT ON COLUMN source.sha256_hash IS 'SHA-256 hex digest of the file content. Used to detect tampering after upload. Exactly 64 lowercase hex characters.';
COMMENT ON COLUMN source.s3_key IS 'MinIO/S3 object storage key for secure evidence retrieval via pre-signed URL.';
