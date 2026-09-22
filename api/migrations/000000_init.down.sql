-- Migration: 000000_init.down.sql
-- Description: Revert initial extensions setup

DROP EXTENSION IF EXISTS "uuid-ossp";
DROP EXTENSION IF EXISTS postgis;
