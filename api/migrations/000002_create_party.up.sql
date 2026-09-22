-- Migration: 000002_create_party.up.sql
-- Description: Create the party table (LADM LA_Party)
-- Stores persons or organizations that hold rights over land
-- IMPORTANT: Does NOT duplicate civil registry data. Links externally via external_id + external_system

CREATE TABLE party (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    party_type      party_type NOT NULL,
    first_name      VARCHAR(150),
    last_name       VARCHAR(150),
    -- External identity reference (Civil Registry, Tax Authority, etc)
    external_id     VARCHAR(100),
    external_system VARCHAR(100),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    -- Prevent duplicate linking of the same external identity
    CONSTRAINT uq_party_external UNIQUE (external_id, external_system)
);

COMMENT ON TABLE party IS 'LADM LA_Party: Persons or entities holding land rights. Links to external identity systems without duplicating sensitive data.';
COMMENT ON COLUMN party.external_id IS 'Identifier in the referenced external system (e.g. national civil registry ID).';
COMMENT ON COLUMN party.external_system IS 'Name of the external system providing the identity (e.g. NATIONAL_CIVIL_REGISTRY).';
