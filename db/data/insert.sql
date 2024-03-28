-- master_scopes
INSERT INTO master_scopes (code, description) VALUES
('images.read', 'Permission to read images'),
('images.create', 'Permission to create images'),
('images.update', 'Permission to update images'),
('images.delete', 'Permission to delete images');

-- master_application_types
INSERT INTO master_application_types (type) VALUES
('web');

-- health_checks
INSERT INTO health_checks (key, value)
VALUES ('healthCheckKey', 'postgresValue');

-- users
INSERT INTO users (id, name, password, role) VALUES
('2e664489-8362-4c1b-9379-2ab6f80beace', 'username', '$2a$10$/cGQ8OB.jgqENIpD4BcNZObIreR6vJPKvMXBhNv0P5XApERwKDpOW', 'admin');

-- oauth_clients
INSERT INTO oauth_clients (id, name, client_id, client_secret, redirect_uri, application_url, client_type) VALUES
('b2124953-31a8-4c16-b5cf-fdd1e40edc14', 'General Public Client', '0411a9bb-b450-4951-8d95-dfbf19dd925b', 'public_client_secret', null, 'http://localhost:3000', 'confidential'),
('4533d234-5f04-4a03-8171-f1f952736373', 'Owner Private Client', 'a74983c2-c578-41fd-993b-9e4716d244ac', 'owner_client_secret', 'http://localhost:3000/api/auth/callback', 'http://localhost:3000', 'confidential');

INSERT INTO oauth_clients_scopes (client_id, scope_code) VALUES
('b2124953-31a8-4c16-b5cf-fdd1e40edc14', 'images.read'),
('b2124953-31a8-4c16-b5cf-fdd1e40edc14', 'images.create'),
('b2124953-31a8-4c16-b5cf-fdd1e40edc14', 'images.update'),
('4533d234-5f04-4a03-8171-f1f952736373', 'images.read'),
('4533d234-5f04-4a03-8171-f1f952736373', 'images.create'),
('4533d234-5f04-4a03-8171-f1f952736373', 'images.update'),
('4533d234-5f04-4a03-8171-f1f952736373', 'images.delete');

-- oauth_clients_application_types
INSERT INTO oauth_clients_application_types (client_id, application_type) VALUES
('b2124953-31a8-4c16-b5cf-fdd1e40edc14', 'web'),
('4533d234-5f04-4a03-8171-f1f952736373', 'web');
