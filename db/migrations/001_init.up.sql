-- Install uuid-ossp extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Up migration users
CREATE TABLE IF NOT EXISTS "users" (
    "id" UUID NOT NULL DEFAULT uuid_generate_v4(),
    "name" VARCHAR(10) NOT NULL UNIQUE,
    "password" TEXT NOT NULL,
    PRIMARY KEY ("id")
);

-- Up migration health_checks
CREATE TABLE IF NOT EXISTS "health_checks" (
    "id" UUID NOT NULL DEFAULT uuid_generate_v4(),
    "key" VARCHAR(20) NOT NULL,
    "value" VARCHAR(20) NOT NULL,
    PRIMARY KEY ("id")
);