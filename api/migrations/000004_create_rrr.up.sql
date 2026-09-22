-- Migration: 000004_create_rrr.up.sql
-- Description: Create the rrr table (LADM LA_RRR - Rights, Restrictions, Responsibilities)
-- Links a party to a spatial_unit with a specific tenure type

CREATE TABLE rrr (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    rrr_type         rrr_type NOT NULL,
    party_id         UUID NOT NULL REFERENCES party(id) ON DELETE RESTRICT,
    spatial_unit_id  UUID NOT NULL REFERENCES spatial_unit(id) ON DELETE RESTRICT,
    share_fraction   NUMERIC(5, 4) CHECK (share_fraction > 0 AND share_fraction <= 1),
    start_date       DATE,
    end_date         DATE,
    notes            TEXT,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    -- A party can hold only one RRR of each type per spatial unit
    CONSTRAINT uq_rrr_party_unit_type UNIQUE (party_id, spatial_unit_id, rrr_type)
);

COMMENT ON TABLE rrr IS 'LADM LA_RRR: Rights, Restrictions and Responsibilities linking a party to a spatial unit.';
COMMENT ON COLUMN rrr.share_fraction IS 'Fractional share of the right (e.g. 0.5 for 50% co-ownership). Null means full right.';
