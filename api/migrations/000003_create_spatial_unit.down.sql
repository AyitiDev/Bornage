-- Migration: 000003_create_spatial_unit.down.sql
-- Description: Drop the spatial_unit table and its associated trigger/function

DROP TRIGGER IF EXISTS trg_spatial_unit_area ON spatial_unit;
DROP FUNCTION IF EXISTS fn_calculate_spatial_unit_area();
DROP TABLE IF EXISTS spatial_unit;
