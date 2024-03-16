-- Install uuid-ossp extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Up migration
CREATE TABLE IF NOT EXISTS "users" (
    "id" UUID NOT NULL DEFAULT uuid_generate_v4(),
    "name" VARCHAR(10) NOT NULL,
    "password" TEXT NOT NULL,
    PRIMARY KEY ("id")
);