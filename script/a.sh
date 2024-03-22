#!/bin/bash

if ! command -v jq &> /dev/null; then
    echo "jq command is not installed. Please install jq to run this script."
    exit 1
fi

cookie_jar="cookies.txt"

echo "ログイン処理を実行します..."
curl -c $cookie_jar -i -X POST http://localhost:8080/api/login \
 -H "Content-Type: application/x-www-form-urlencoded" \
 -d "username=username" \
 -d "password=password"
echo -e "\n"

# 認可処理
echo "認可処理を実行します..."
response_type="code"
client_id="a74983c2-c578-41fd-993b-9e4716d244ac"
redirect_uri="http://localhost:3000/api/owner/callback"
scope="images.read%20images.create%20images.update%20images.delete"
state="xyz"
nonce="abc123"
query_params="response_type=${response_type}&client_id=${client_id}&redirect_uri=${redirect_uri}&scope=${scope}&state=${state}&nonce=${nonce}"

auth_response=$(curl -L -b $cookie_jar -c $cookie_jar -i "http://localhost:8080/api/connect/authorize?${query_params}")
echo "$auth_response"
echo -e "\n"

auth_code=$(echo "$auth_response" | grep -oE 'code=[A-Za-z0-9\-]+' | cut -d'=' -f2 | cut -d'&' -f1)

# アクセストークン取得
echo "アクセストークンを取得します..."
client_secret="owner_client_secret"
access_token_response=$(curl -X POST http://localhost:8080/api/connect/token \
 -H "Content-Type: application/x-www-form-urlencoded" \
 -d "grant_type=authorization_code" \
 -d "code=${auth_code}" \
 -d "redirect_uri=http://localhost:3000/api/owner/callback" \
 -d "client_id=${client_id}" \
 -d "client_secret=${client_secret}")
echo -e "\n"

access_token=$(echo "$access_token_response" | jq -r '.accessToken')
echo "Access Token: $access_token"