-- Install uuid-ossp extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS health_checks (
    key VARCHAR(20) PRIMARY KEY,
    value VARCHAR(20) NOT NULL
);

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(20) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    role VARCHAR(20) NOT NULL
);

CREATE TABLE IF NOT EXISTS oauth_clients (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(30) NOT NULL,
    client_id UUID UNIQUE NOT NULL,
    client_secret VARCHAR(255) UNIQUE NOT NULL,
    redirect_uri TEXT NOT NULL,
    application_url TEXT NOT NULL,
    client_type VARCHAR(50) NOT NULL
);

CREATE TABLE IF NOT EXISTS master_scopes (
    code VARCHAR(20) PRIMARY KEY,
    description TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS oauth_clients_scopes (
    client_id UUID,
    scope_code VARCHAR(20),
    PRIMARY KEY (client_id, scope_code),
    FOREIGN KEY (client_id) REFERENCES oauth_clients(id),
    FOREIGN KEY (scope_code) REFERENCES master_scopes(code)
);

CREATE TABLE IF NOT EXISTS master_application_types (
    type VARCHAR(20) PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS oauth_clients_application_types (
    client_id UUID,
    application_type VARCHAR(20),
    PRIMARY KEY (client_id, application_type),
    FOREIGN KEY (client_id) REFERENCES oauth_clients(id),
    FOREIGN KEY (application_type) REFERENCES master_application_types(type)
);

CREATE TABLE IF NOT EXISTS refresh_tokens (
    token VARCHAR(255) PRIMARY KEY,
    user_id UUID NOT NULL,
    client_id UUID NOT NULL,
    scopes TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (client_id) REFERENCES oauth_clients(client_id)
);

CREATE TABLE IF NOT EXISTS images (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    url TEXT NOT NULL,
    keyword VARCHAR(50) NOT NULL DEFAULT '',
    used_count INTEGER NOT NULL DEFAULT 0,
    reported BOOLEAN NOT NULL DEFAULT false,
    confirmed BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP
);