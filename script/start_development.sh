#!/bin/bash

source .env.local

make dependencies_start

make view_build

make migrate_up

make insert_data

docker exec -it my_authn_authz_redis redis-cli SET healthCheckKey redisValue
