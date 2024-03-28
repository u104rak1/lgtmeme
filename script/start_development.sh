#!/bin/bash

source .env.local

make dependencies_start

make view_build

make migrate_up

make clear_data

make insert_data

docker exec -it lgtmeme_redis redis-cli SET healthCheckKey redisValue
