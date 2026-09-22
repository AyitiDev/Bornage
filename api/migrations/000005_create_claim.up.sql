-- Migration: 000005_create_claim.up.sql
-- Description: Create the claim table (registration application in state machine lifecycle)
-- DESIGN RULE: PUBLISHED_FOR_OBJECTION is mandatory. No direct path from UNDER_REVIEW to REGISTERED

CREATE TABLE claim (
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    -- Anonymous public tracking code (safe to share with claimant)
    tracking_code       VARCHAR(50) NOT NULL UNIQUE,
    status              claim_status NOT NULL DEFAULT 'DRAFT',
    claimed_geom        GEOMETRY(Polygon, 4326) NOT NULL,
    claimant_id         UUID NOT NULL REFERENCES party(id) ON DELETE RESTRICT,
    -- Set when transitioning to PUBLISHED_FOR_OBJECTION (30-60 days mandatory window)
    objection_end_date  DATE,
    -- Administrative jurisdiction of the parcel
    administrative_unit VARCHAR(200),
    notes               TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    -- Enforce: objection_end_date must be set before the claim can be REGISTERED
    CONSTRAINT chk_registered_requires_objection_date
        CHECK (status != 'REGISTERED' OR objection_end_date IS NOT NULL)
);

COMMENT ON TABLE claim IS 'Land registration application subject to the LADM state machine. PUBLISHED_FOR_OBJECTION is a mandatory non-skippable step.';
COMMENT ON COLUMN claim.tracking_code IS 'Anonymous public-facing code given to the claimant for status tracking. Never exposes personal data.';
COMMENT ON COLUMN claim.objection_end_date IS 'End of 30-60 day mandatory public objection window. REGISTERED requires this date to be in the past.';
