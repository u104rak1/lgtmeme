-- health_checks
INSERT INTO health_checks (key, value)
VALUES ('healthCheckKey', 'postgresValue');

-- users
INSERT INTO users (id, name, password, role) VALUES
('2e664489-8362-4c1b-9379-2ab6f80beace', 'username', '$2a$10$/cGQ8OB.jgqENIpD4BcNZObIreR6vJPKvMXBhNv0P5XApERwKDpOW', 'admin');

-- oauth_clients
INSERT INTO oauth_clients (id, name, client_id, client_secret, redirect_uri, client_type) VALUES
('b2124953-31a8-4c16-b5cf-fdd1e40edc14', 'General Public Client', '0411a9bb-b450-4951-8d95-dfbf19dd925b', 'public_client_secret', 'http://localhost:3000/api/callback', 'confidential'),
('4533d234-5f04-4a03-8171-f1f952736373', 'Owner Private Client', 'a74983c2-c578-41fd-993b-9e4716d244ac', 'owner_client_secret', 'http://localhost:3000/api/owner/callback', 'confidential');

-- scopes
INSERT INTO scopes (code, description) VALUES
('images_read', 'Permission to read images'),
('images_create', 'Permission to create images'),
('images_update', 'Permission to update images'),
('images_delete', 'Permission to delete images');

INSERT INTO oauth_clients_scopes (client_id, scope_code) VALUES
('b2124953-31a8-4c16-b5cf-fdd1e40edc14', 'images_read'),
('b2124953-31a8-4c16-b5cf-fdd1e40edc14', 'images_create'),
('b2124953-31a8-4c16-b5cf-fdd1e40edc14', 'images_update'),
('4533d234-5f04-4a03-8171-f1f952736373', 'images_read'),
('4533d234-5f04-4a03-8171-f1f952736373', 'images_create'),
('4533d234-5f04-4a03-8171-f1f952736373', 'images_update'),
('4533d234-5f04-4a03-8171-f1f952736373', 'images_delete');

-- application_types
INSERT INTO application_types (type) VALUES
('web');

-- oauth_clients_application_types
INSERT INTO oauth_clients_application_types (client_id, application_type) VALUES
('b2124953-31a8-4c16-b5cf-fdd1e40edc14', 'web'),
('4533d234-5f04-4a03-8171-f1f952736373', 'web');
