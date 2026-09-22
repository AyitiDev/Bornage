-- Migration: 000001_create_enums.up.sql
-- Description: Create all LADM domain enum types used across core tables

-- Claim lifecycle states (strict state machine -- no skipping PUBLISHED_FOR_OBJECTION)
CREATE TYPE claim_status AS ENUM (
    'DRAFT',
    'SUBMITTED',
    'UNDER_REVIEW',
    'PUBLISHED_FOR_OBJECTION',
    'REGISTERED',
    'RETURNED_TO_FIELD',
    'REJECTED',
    'DISPUTED'
);

-- Right / Restriction / Responsibility type (LADM LA_RRR)
CREATE TYPE rrr_type AS ENUM (
    'OWNERSHIP',
    'LEASE',
    'USUFRUCT',
    'INFORMAL_OCCUPATION',
    'EASEMENT',
    'MORTGAGE',
    'RESTRICTION'
);

-- Party type (LADM LA_Party)
CREATE TYPE party_type AS ENUM (
    'INDIVIDUAL',
    'ORGANIZATION',
    'STATE'
);

-- Capture method used to delineate the spatial unit boundary
CREATE TYPE capture_method AS ENUM (
    'SATELLITE_VISUAL_DELINEATION',
    'GPS_WALKED_PERIMETER',
    'SURVEY_TOTAL_STATION',
    'DRONE_PHOTOGRAMMETRY'
);
