-- Migration: 000008_create_audit_log.down.sql
-- Description: Drop the audit_log table (reverse of up migration)
-- WARNING: This permanently destroys all audit records. Only for development/test rollbacks

DROP TABLE IF EXISTS audit_log;
