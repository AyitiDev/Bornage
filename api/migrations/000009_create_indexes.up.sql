-- Migration: 000009_create_indexes.up.sql
-- Description: Create GiST spatial indexes for PostGIS queries and B Tree indexes on foreign keys
-- GiST indexes enable fast ST_Intersects / ST_Contains / ST_DWithin queries on parcels

-- GiST spatial index on spatial_unit geometry (core overlap detection engine)
CREATE INDEX idx_spatial_unit_geom ON spatial_unit USING GIST(geom);

-- GiST spatial index on claim geometry (overlap analysis against pending claims)
CREATE INDEX idx_claim_geom ON claim USING GIST(claimed_geom);

-- B-Tree index on claim.tracking_code (public anonymous lookups - must be fast)
CREATE INDEX idx_claim_tracking_code ON claim(tracking_code);

-- B-Tree index on claim.status (dashboard filtering by state)
CREATE INDEX idx_claim_status ON claim(status);

-- B-Tree index on claim.objection_end_date (query for expired objection windows)
CREATE INDEX idx_claim_objection_end_date ON claim(objection_end_date);

-- B-Tree indexes on foreign key columns (prevent sequential scans on JOINs)
CREATE INDEX idx_rrr_party_id ON rrr(party_id);
CREATE INDEX idx_rrr_spatial_unit_id ON rrr(spatial_unit_id);
CREATE INDEX idx_source_claim_id ON source(claim_id);
CREATE INDEX idx_witness_claim_id ON witness(claim_id);
CREATE INDEX idx_audit_log_target_id ON audit_log(target_id);
CREATE INDEX idx_audit_log_actor_id ON audit_log(actor_id);
CREATE INDEX idx_audit_log_occurred_at ON audit_log(occurred_at DESC);
