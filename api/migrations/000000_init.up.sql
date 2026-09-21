-- Migration: 000000_init.up.sql
-- Description: Enable core PostgreSQL extensions (PostGIS for spatial data & uuid-ossp for UUID generation)

CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
