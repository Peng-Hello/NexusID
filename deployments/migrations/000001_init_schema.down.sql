-- Rollback migration
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS authorization_codes;
DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS user_roles;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS oidc_clients;
DROP TABLE IF EXISTS tenants;
