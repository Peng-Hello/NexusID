-- Rollback seed data
DELETE FROM roles WHERE tenant_id = 1;
DELETE FROM oidc_clients WHERE tenant_id = 1;
DELETE FROM tenants WHERE id = 1;
