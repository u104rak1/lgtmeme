#!/bin/bash

source .env.local

rm -rf ./docker/volumes/storage/stub

make dependencies_start

make view_build

make migrate_up

make clear_data

make insert_data

docker exec -it lgtmeme_redis redis-cli SET healthCheckKey redisValue

echo "Create a bucket for images in Supabase Storage..."

SUPABASE_STORAGE_BASE_URL=$SUPABASE_STORAGE_BASE_URL
SUPABASE_SERVICE_ROLE_KEY=$SUPABASE_SERVICE_ROLE_KEY

response=$(curl -i -X GET "${SUPABASE_STORAGE_BASE_URL}/bucket/images" \
  -H "Authorization: Bearer ${SUPABASE_SERVICE_ROLE_KEY}")

status=$(echo "$response" | grep HTTP | awk '{print $2}')

if [ "$status" = "404" ]; then
  echo "Bucket 'images' does not exist, creating..."
  curl -i -X POST "${SUPABASE_STORAGE_BASE_URL}/bucket" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer ${SUPABASE_SERVICE_ROLE_KEY}" \
    -d '{
        "name": "images",
        "id": "images",
        "public": true
      }'
fi