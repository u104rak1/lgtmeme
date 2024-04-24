INSERT INTO master_scopes (code, description) VALUES
('images.read', 'Permission to read images'),
('images.create', 'Permission to create images'),
('images.update', 'Permission to update images'),
('images.delete', 'Permission to delete images')
ON CONFLICT (code) DO NOTHING;

INSERT INTO master_application_types (type) VALUES
('web')
ON CONFLICT (type) DO NOTHING;

INSERT INTO health_checks (key, value)
VALUES ('healthCheckKey', 'postgresValue')
ON CONFLICT (key) DO NOTHING;

INSERT INTO users (id, name, password, role) VALUES
('2e664489-8362-4c1b-9379-2ab6f80beace', 'username', '$2a$10$/cGQ8OB.jgqENIpD4BcNZObIreR6vJPKvMXBhNv0P5XApERwKDpOW', 'admin')
ON CONFLICT (id) DO NOTHING;

INSERT INTO oauth_clients (id, name, client_id, client_secret, redirect_uri, application_url, client_type) VALUES
('b2124953-31a8-4c16-b5cf-fdd1e40edc14', 'General Public Client', '0411a9bb-b450-4951-8d95-dfbf19dd925b', 'public_client_secret', '', 'http://localhost:8080', 'confidential'),
('4533d234-5f04-4a03-8171-f1f952736373', 'Admin Private Client', 'a74983c2-c578-41fd-993b-9e4716d244ac', 'admin_client_secret', 'http://localhost:8080/client-api/admin/callback', 'http://localhost:8080', 'confidential')
ON CONFLICT (id) DO NOTHING;

INSERT INTO oauth_clients_scopes (oauth_client_id, scope_code) VALUES
('b2124953-31a8-4c16-b5cf-fdd1e40edc14', 'images.read'),
('b2124953-31a8-4c16-b5cf-fdd1e40edc14', 'images.create'),
('b2124953-31a8-4c16-b5cf-fdd1e40edc14', 'images.update'),
('4533d234-5f04-4a03-8171-f1f952736373', 'images.read'),
('4533d234-5f04-4a03-8171-f1f952736373', 'images.create'),
('4533d234-5f04-4a03-8171-f1f952736373', 'images.update'),
('4533d234-5f04-4a03-8171-f1f952736373', 'images.delete')
ON CONFLICT (oauth_client_id, scope_code) DO NOTHING;

INSERT INTO oauth_clients_application_types (oauth_client_id, application_type) VALUES
('b2124953-31a8-4c16-b5cf-fdd1e40edc14', 'web'),
('4533d234-5f04-4a03-8171-f1f952736373', 'web')
ON CONFLICT (oauth_client_id, application_type) DO NOTHING;

INSERT INTO images (id, url, keyword, used_count, reported, confirmed, created_at) VALUES
('a2128761-21a8-53c6-b6cd-1578eaf12c14', 'https://placehold.jp/300x300.png', '300 * 300 sample A', 2, false, false, '2024-01-01 00:00:00'),
('b2128761-21a8-53c6-b6cd-1578eaf12c14', 'https://placehold.jp/400x400.png', '400 * 400 sample B', 1, false, false, '2024-02-01 00:00:00'),
('c2128761-21a8-53c6-b6cd-1578eaf12c14', 'https://placehold.jp/500x500.png', '500 * 500 sample C', 0, false, false, '2024-03-01 00:00:00')
ON CONFLICT (id) DO NOTHING;
