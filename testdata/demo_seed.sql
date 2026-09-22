-- testdata/demo_seed.sql
-- Purely synthetic and mock demonstration seed dataset for local development
-- NO REAL PERSONAL OR PARCEL DATA IS CONTAINED IN THIS FILE

-- 1. Insert Mock Parties (Identities linked via external_id)
-- Note: Real systems reference external Civil Registry / Tax IDs without duplicating sensitive records
INSERT INTO party (id, first_name, last_name, external_id, external_system, party_type, created_at)
VALUES 
  ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Jean-Luc', 'Valcin', 'HT-ID-9901', 'NATIONAL_CIVIL_REGISTRY', 'INDIVIDUAL', NOW()),
  ('b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'Marie-Rose', 'Desrosiers', 'HT-ID-9902', 'NATIONAL_CIVIL_REGISTRY', 'INDIVIDUAL', NOW()),
  ('c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'Coopérative Agricole Nord', 'TAX-ORG-5001', 'TAX_AUTHORITY', 'ORGANIZATION', NOW())
ON CONFLICT DO NOTHING;

-- 2. Insert Mock Spatial Units (Parcels with PostGIS Polygons in EPSG:4326)
-- Estimated precision ~3m Fit for Purpose (FFP) visual boundary capture method
INSERT INTO spatial_unit (id, geom, calculated_area_sqm, capture_method, accuracy_meters, created_at)
VALUES 
  (
    'd0eebc99-9c0b-4ef8-bb6d-6bb9bd380a44',
    ST_GeomFromText('POLYGON((-72.338 18.539, -72.337 18.539, -72.337 18.538, -72.338 18.538, -72.338 18.539))', 4326),
    12350.50,
    'SATELLITE_VISUAL_DELINEATION',
    3.0,
    NOW()
  ),
  (
    'e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a55',
    ST_GeomFromText('POLYGON((-72.336 18.539, -72.335 18.539, -72.335 18.538, -72.336 18.538, -72.336 18.539))', 4326),
    11800.00,
    'GPS_WALKED_PERIMETER',
    2.5,
    NOW()
  )
ON CONFLICT DO NOTHING;

-- 3. Insert Mock Claims (Registration applications in state machine lifecycle)
INSERT INTO claim (id, tracking_code, status, claimed_geom, claimant_id, created_at)
VALUES 
  (
    'f0eebc99-9c0b-4ef8-bb6d-6bb9bd380a66',
    'CLAIM-2026-HT-001',
    'PUBLISHED_FOR_OBJECTION',
    ST_GeomFromText('POLYGON((-72.338 18.539, -72.337 18.539, -72.337 18.538, -72.338 18.538, -72.338 18.539))', 4326),
    'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
    NOW()
  )
ON CONFLICT DO NOTHING;
