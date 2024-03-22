#!/bin/bash

if ! command -v jq &> /dev/null; then
    echo "jq command is not installed. Please install jq to run this script."
    exit 1
fi

if ! curl --output /dev/null --silent --head --fail http://localhost:8080/api/health; then
    echo "The service at localhost:8080 is not responding to health check. Please start the service and try again."
    exit 1
else
    echo "The service at localhost:8080 is up and running."
fi

COOKIE_FILE="cookies.txt"

echo "Executes the login process..."
curl -c $COOKIE_FILE -i -X POST http://localhost:8080/api/login \
 -H "Content-Type: application/x-www-form-urlencoded" \
 -d "username=username" \
 -d "password=password"
echo -e "\n"

echo "Perform authorization processing..."
RESPONSE_TYPE="code"
CLIENT_ID="a74983c2-c578-41fd-993b-9e4716d244ac"
REDIRECT_URI="http://localhost:3000/api/owner/callback"
SCOPE="images.read%20images.create%20images.update%20images.delete"
STATE="xyz"
NONCE="abc123"
QUERY_PARAMS="response_type=${RESPONSE_TYPE}&client_id=${CLIENT_ID}&redirect_uri=${REDIRECT_URI}&scope=${SCOPE}&state=${STATE}&nonce=${NONCE}"
AUTH_RESPONSE=$(curl -L -b $COOKIE_FILE -c $COOKIE_FILE -i "http://localhost:8080/api/connect/authorize?${QUERY_PARAMS}")
echo "$AUTH_RESPONSE"
echo -e "\n"

rm -f $COOKIE_FILE

echo "Obtain Access Token using authorization_code grant..."
AUTH_CODE=$(echo "$AUTH_RESPONSE" | grep -oE 'code=[A-Za-z0-9\-]+' | cut -d'=' -f2 | cut -d'&' -f1)
CLIENT_SECRET="owner_client_secret"
TOKEN_RESPONSE=$(curl -X POST http://localhost:8080/api/connect/token \
 -H "Content-Type: application/x-www-form-urlencoded" \
 -d "grant_type=authorization_code" \
 -d "code=${AUTH_CODE}" \
 -d "redirect_uri=${REDIRECT_URI}" \
 -d "client_id=${CLIENT_ID}" \
 -d "client_secret=${CLIENT_SECRET}")
echo -e "\n"
ACCESS_TOKEN=$(echo "$TOKEN_RESPONSE" | jq -r '.accessToken')
echo "Access Token obtained with authorization_code grant: $ACCESS_TOKEN"
echo -e "\n"

echo "Obtain Access Token using refresh_token grant..."
REFRESH_TOKEN=$(echo "$TOKEN_RESPONSE" | jq -r '.refreshToken')
REFLESH_TOKEN_RESPONSE=$(curl -X POST http://localhost:8080/api/connect/token \
 -H "Content-Type: application/x-www-form-urlencoded" \
 -d "grant_type=refresh_token" \
 -d "refresh_token=${REFRESH_TOKEN}" \
 -d "client_id=${CLIENT_ID}" \
 -d "client_secret=${CLIENT_SECRET}")
echo -e "\n"
ACCESS_TOKEN_REFLESH_TOKEN=$(echo "$REFLESH_TOKEN_RESPONSE" | jq -r '.accessToken')
echo "Access Token obtained with reflesh_token grant: $ACCESS_TOKEN_REFLESH_TOKEN"
echo -e "\n"

echo "Obtain Access Token using client_credentials grant..."
CLIENT_ID_2="0411a9bb-b450-4951-8d95-dfbf19dd925b"
CLIENT_SECRET_2="public_client_secret"
CLIENT_CREDENTIALS_RESPONSE=$(curl -X POST http://localhost:8080/api/connect/token \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=client_credentials" \
  -d "client_id=${CLIENT_ID_2}" \
  -d "client_secret=${CLIENT_SECRET_2}")
ACCESS_TOKEN_CLIENT_CREDENTIALS=$(echo $CLIENT_CREDENTIALS_RESPONSE | jq -r '.accessToken')
echo "Access Token obtained with client_credentials grant: $ACCESS_TOKEN_CLIENT_CREDENTIALS"
echo -e "\n"
