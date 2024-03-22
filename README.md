# my_authn_authz

## Redis
```
docker exec -it my_authn_authz_redis redis-cli
keys *
get [key]
```

## /api/connect/authorize
[login](http://localhost:8080/api/connect/authorize?response_type=code&client_id=a74983c2-c578-41fd-993b-9e4716d244ac&redirect_uri=http://localhost:3000/api/owner/callback&scope=images_read%20images_create%20images_update%20images_delete&state=xyz&nonce=abc123)

## /api/connect/token
```
curl -X POST http://localhost:8080/api/connect/token \
    -H "Content-Type: application/x-www-form-urlencoded" \
    -d "grant_type=authorization_code" \
    -d "code=YOUR_AUTHORIZATION_CODE" \
    -d "redirect_uri=http://localhost:3000/api/owner/callback" \
    -d "client_id=a74983c2-c578-41fd-993b-9e4716d244ac" \
    -d "client_secret=owner_client_secret"
```