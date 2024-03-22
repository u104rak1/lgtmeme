-- Down migration
DROP TABLE IF EXISTS "users";
DROP TABLE IF EXISTS "health_checks";
DROP TABLE IF EXISTS "oauth_clients_application_types";
DROP TABLE IF EXISTS "oauth_clients_scopes";
DROP TABLE IF EXISTS "oauth_clients";
DROP TABLE IF EXISTS "master_application_types";
DROP TABLE IF EXISTS "master_scopes";