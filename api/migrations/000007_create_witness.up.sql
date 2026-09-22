-- Migration: 000007_create_witness.up.sql
-- Description: Create the witness table for on site boundary witnesses
-- Neighboring occupants who confirm the parcel boundary in the field

CREATE TABLE witness (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    claim_id        UUID NOT NULL REFERENCES claim(id) ON DELETE CASCADE,
    full_name       VARCHAR(300) NOT NULL,
    -- S3 key for the witness signature image (touch/canvas capture or scanned)
    signature_s3_key VARCHAR(500),
    -- S3 key for witness photo or ID document photo
    photo_s3_key    VARCHAR(500),
    -- GPS coordinates where witness was present during field verification
    witness_lat     NUMERIC(10, 7),
    witness_lng     NUMERIC(10, 7),
    confirmed_at    TIMESTAMPTZ,
    notes           TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE witness IS 'On-site boundary witnesses (neighboring occupants) confirming parcel limits during field data collection.';
COMMENT ON COLUMN witness.signature_s3_key IS 'S3 key for the digital touch/canvas signature image captured in the field PWA.';
COMMENT ON COLUMN witness.confirmed_at IS 'Timestamp when the witness confirmed the boundary on-site.';
