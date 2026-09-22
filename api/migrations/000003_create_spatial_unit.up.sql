-- Migration: 000003_create_spatial_unit.up.sql
-- Description: Create the spatial_unit table (LADM LA_SpatialUnit) with PostGIS geometry
-- Stores land parcels as Polygons in EPSG:4326 (WGS84)
-- Fit for Purpose (FFP): Accepts 3m accuracy from satellite visual delineation

CREATE TABLE spatial_unit (
    id                   UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    geom                 GEOMETRY(Polygon, 4326) NOT NULL,
    calculated_area_sqm  NUMERIC(14, 4),
    capture_method       capture_method NOT NULL DEFAULT 'SATELLITE_VISUAL_DELINEATION',
    accuracy_meters      NUMERIC(6, 2) NOT NULL DEFAULT 3.0,
    administrative_unit  VARCHAR(200),
    notes                TEXT,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE spatial_unit IS 'LADM LA_SpatialUnit: Land parcel with PostGIS polygon geometry. Area is auto-calculated by trigger.';
COMMENT ON COLUMN spatial_unit.accuracy_meters IS 'Estimated capture precision in meters. FFP standard: ~3m is acceptable.';
COMMENT ON COLUMN spatial_unit.calculated_area_sqm IS 'Area in square metres, computed automatically via ST_Area on the geography cast.';

-- Trigger function: auto-calculate area on INSERT or UPDATE using PostGIS geography cast
CREATE OR REPLACE FUNCTION fn_calculate_spatial_unit_area()
RETURNS TRIGGER AS $$
BEGIN
    NEW.calculated_area_sqm := ROUND(CAST(ST_Area(NEW.geom::geography) AS NUMERIC), 4);
    NEW.updated_at := NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_spatial_unit_area
BEFORE INSERT OR UPDATE ON spatial_unit
FOR EACH ROW EXECUTE FUNCTION fn_calculate_spatial_unit_area();
