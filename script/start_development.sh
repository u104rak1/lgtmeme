#!/bin/bash

source .env.local

if [ ! -f .env.local.secret ]; then
  echo ".env.local.secret does not exist. Creating keys..."
  openssl genpkey -algorithm RSA -out private_key.pem -pkeyopt rsa_keygen_bits:2048
  openssl rsa -pubout -in private_key.pem -out public_key.pem
  
  echo "JWT_PRIVATE_KEY_BASE64=$(base64 private_key.pem | tr -d '\n')" > .env.local.secret
  echo "JWT_PUBLIC_KEY_BASE64=$(base64 public_key.pem | tr -d '\n')" >> .env.local.secret
fi

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