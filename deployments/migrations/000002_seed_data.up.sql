-- Insert default super tenant
INSERT INTO tenants (id, name, slug, description, is_active) VALUES
(1, 'Super Tenant', 'super', 'Default system tenant for super administrators', TRUE);

-- Insert default OIDC client for development
INSERT INTO oidc_clients (tenant_id, client_id, client_secret_hash, name, redirect_uris, scopes, grant_types, is_public, is_active) VALUES
(1, 'nexus-id-dev-client', '$2a$10$.placeholder.hash.for.development', 'Development Client',
 '["http://localhost:3000/callback", "http://localhost:8080/callback"]',
 '["openid", "profile", "email"]',
 '["authorization_code", "refresh_token"]',
 FALSE, TRUE);

-- Insert system roles for the super tenant
INSERT INTO roles (tenant_id, name, description, is_system_role) VALUES
(1, 'super_admin', 'Super administrator with full system access', TRUE),
(1, 'tenant_admin', 'Tenant administrator with tenant-level access', TRUE),
(1, 'user', 'Standard user role', TRUE);
